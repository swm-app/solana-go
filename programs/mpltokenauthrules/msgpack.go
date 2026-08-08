package mpltokenauthrules

import (
	"encoding/binary"
	"fmt"
)

// Minimal msgpack reader covering exactly the subset that rmp_serde emits for a
// RuleSetV1 with the default (compact) configuration: nil, bool, every integer
// width, str (fixstr / str8 / str16 / str32), bin (bin8 / bin16 / bin32), array
// (fixarray / array16 / array32) and map (fixmap / map16 / map32). Floats and
// extension types are rejected: they never appear in this format, so their
// presence signals corrupt data. Every length is validated against the
// remaining input before any slice or allocation, so arbitrary bytes cannot
// cause an out-of-bounds access or an oversized preallocation.

type mpType uint8

const (
	mpNil mpType = iota
	mpBool
	mpUint
	mpInt
	mpStr
	mpBin
	mpArray
	mpMap
)

type mpValue struct {
	typ  mpType
	b    bool
	u    uint64
	i    int64
	str  string
	bin  []byte
	arr  []mpValue
	pair []mpKV
}

type mpKV struct {
	key mpValue
	val mpValue
}

// maxMsgpackDepth bounds the nesting of arrays and maps the reader will descend
// into. Every nested container consumes one level; a legitimate RuleSetV1 nests
// only a handful of levels (the top struct array, the operations map, each rule
// map, and any All/Any/Not sub-rules), so this ceiling is far above any real
// rule set while preventing a maliciously deep input from exhausting the Go
// stack (an unrecoverable fatal error, not a catchable panic).
const maxMsgpackDepth = 1024

type mpReader struct {
	data  []byte
	pos   int
	depth int
}

func newMsgpackReader(data []byte) *mpReader {
	return &mpReader{data: data}
}

func (r *mpReader) errorf(format string, args ...any) error {
	return fmt.Errorf("msgpack at offset %d: %s", r.pos, fmt.Sprintf(format, args...))
}

// need reports an error unless at least n bytes remain from the current
// position.
func (r *mpReader) need(n int) error {
	if n < 0 || r.pos+n > len(r.data) {
		return r.errorf("unexpected end of input: need %d byte(s), have %d", n, len(r.data)-r.pos)
	}
	return nil
}

func (r *mpReader) readByte() (byte, error) {
	if err := r.need(1); err != nil {
		return 0, err
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

func (r *mpReader) readN(n int) ([]byte, error) {
	if err := r.need(n); err != nil {
		return nil, err
	}
	b := r.data[r.pos : r.pos+n]
	r.pos += n
	return b, nil
}

// readLen reads an unsigned length of the given byte width and validates that
// it cannot exceed the number of bytes left in the input (every array/map/str
// element consumes at least one byte, so a length larger than the remaining
// input is necessarily truncated). This bounds preallocation.
func (r *mpReader) readLen(width int) (int, error) {
	raw, err := r.readN(width)
	if err != nil {
		return 0, err
	}
	var n uint64
	for _, b := range raw {
		n = (n << 8) | uint64(b)
	}
	remaining := len(r.data) - r.pos
	if n > uint64(remaining) {
		return 0, r.errorf("declared length %d exceeds remaining %d bytes", n, remaining)
	}
	return int(n), nil
}

func (r *mpReader) readValue() (mpValue, error) {
	tag, err := r.readByte()
	if err != nil {
		return mpValue{}, err
	}

	switch {
	case tag <= 0x7f: // positive fixint
		return mpValue{typ: mpUint, u: uint64(tag)}, nil
	case tag >= 0xe0: // negative fixint
		return mpValue{typ: mpInt, i: int64(int8(tag))}, nil
	case tag >= 0x80 && tag <= 0x8f: // fixmap
		return r.readMap(int(tag & 0x0f))
	case tag >= 0x90 && tag <= 0x9f: // fixarray
		return r.readArray(int(tag & 0x0f))
	case tag >= 0xa0 && tag <= 0xbf: // fixstr
		return r.readStr(int(tag & 0x1f))
	}

	switch tag {
	case 0xc0:
		return mpValue{typ: mpNil}, nil
	case 0xc2:
		return mpValue{typ: mpBool, b: false}, nil
	case 0xc3:
		return mpValue{typ: mpBool, b: true}, nil
	case 0xc4:
		return r.readBin(1)
	case 0xc5:
		return r.readBin(2)
	case 0xc6:
		return r.readBin(4)
	case 0xcc:
		return r.readUint(1)
	case 0xcd:
		return r.readUint(2)
	case 0xce:
		return r.readUint(4)
	case 0xcf:
		return r.readUint(8)
	case 0xd0:
		return r.readIntSigned(1)
	case 0xd1:
		return r.readIntSigned(2)
	case 0xd2:
		return r.readIntSigned(4)
	case 0xd3:
		return r.readIntSigned(8)
	case 0xd9:
		return r.readStrLen(1)
	case 0xda:
		return r.readStrLen(2)
	case 0xdb:
		return r.readStrLen(4)
	case 0xdc:
		n, err := r.readLen(2)
		if err != nil {
			return mpValue{}, err
		}
		return r.readArray(n)
	case 0xdd:
		n, err := r.readLen(4)
		if err != nil {
			return mpValue{}, err
		}
		return r.readArray(n)
	case 0xde:
		n, err := r.readLen(2)
		if err != nil {
			return mpValue{}, err
		}
		return r.readMap(n)
	case 0xdf:
		n, err := r.readLen(4)
		if err != nil {
			return mpValue{}, err
		}
		return r.readMap(n)
	default:
		return mpValue{}, r.errorf("unsupported msgpack tag 0x%02x", tag)
	}
}

func (r *mpReader) readUint(width int) (mpValue, error) {
	raw, err := r.readN(width)
	if err != nil {
		return mpValue{}, err
	}
	var n uint64
	for _, b := range raw {
		n = (n << 8) | uint64(b)
	}
	return mpValue{typ: mpUint, u: n}, nil
}

func (r *mpReader) readIntSigned(width int) (mpValue, error) {
	raw, err := r.readN(width)
	if err != nil {
		return mpValue{}, err
	}
	var n uint64
	for _, b := range raw {
		n = (n << 8) | uint64(b)
	}
	// Sign-extend from the given width.
	shift := uint(64 - 8*width)
	return mpValue{typ: mpInt, i: int64(n<<shift) >> shift}, nil
}

func (r *mpReader) readStr(n int) (mpValue, error) {
	raw, err := r.readN(n)
	if err != nil {
		return mpValue{}, err
	}
	return mpValue{typ: mpStr, str: string(raw)}, nil
}

func (r *mpReader) readStrLen(width int) (mpValue, error) {
	n, err := r.readLen(width)
	if err != nil {
		return mpValue{}, err
	}
	return r.readStr(n)
}

func (r *mpReader) readBin(width int) (mpValue, error) {
	n, err := r.readLen(width)
	if err != nil {
		return mpValue{}, err
	}
	raw, err := r.readN(n)
	if err != nil {
		return mpValue{}, err
	}
	cp := make([]byte, n)
	copy(cp, raw)
	return mpValue{typ: mpBin, bin: cp}, nil
}

func (r *mpReader) readArray(n int) (mpValue, error) {
	if r.depth >= maxMsgpackDepth {
		return mpValue{}, r.errorf("nesting exceeds maximum depth %d", maxMsgpackDepth)
	}
	r.depth++
	defer func() { r.depth-- }()
	out := make([]mpValue, 0, capHint(n))
	for range n {
		v, err := r.readValue()
		if err != nil {
			return mpValue{}, err
		}
		out = append(out, v)
	}
	return mpValue{typ: mpArray, arr: out}, nil
}

func (r *mpReader) readMap(n int) (mpValue, error) {
	if r.depth >= maxMsgpackDepth {
		return mpValue{}, r.errorf("nesting exceeds maximum depth %d", maxMsgpackDepth)
	}
	r.depth++
	defer func() { r.depth-- }()
	out := make([]mpKV, 0, capHint(n))
	for range n {
		k, err := r.readValue()
		if err != nil {
			return mpValue{}, err
		}
		v, err := r.readValue()
		if err != nil {
			return mpValue{}, err
		}
		out = append(out, mpKV{key: k, val: v})
	}
	return mpValue{typ: mpMap, pair: out}, nil
}

// capHint caps the initial slice capacity for a declared element count so that
// a large-but-in-bounds length does not force an oversized allocation up front;
// the slice still grows to the real count via append.
func capHint(n int) int {
	const maxHint = 1024
	if n < 0 {
		return 0
	}
	if n > maxHint {
		return maxHint
	}
	return n
}

// asUint returns the value as an unsigned integer if it is a non-negative
// integer.
func (v mpValue) asUint() (uint64, bool) {
	switch v.typ {
	case mpUint:
		return v.u, true
	case mpInt:
		if v.i >= 0 {
			return uint64(v.i), true
		}
	}
	return 0, false
}

// asByte returns the value as a single byte if it is an integer in [0, 255].
func (v mpValue) asByte() (byte, bool) {
	n, ok := v.asUint()
	if !ok || n > 0xff {
		return 0, false
	}
	return byte(n), true
}

// asPubkey interprets the value as a 32-element array of byte-sized integers,
// the wire encoding rmp_serde produces for a Solana Pubkey (a transparent
// [u8; 32], not a msgpack bin).
func (v mpValue) asPubkey() ([32]byte, bool) {
	var out [32]byte
	if v.typ != mpArray || len(v.arr) != 32 {
		return out, false
	}
	for idx, e := range v.arr {
		b, ok := e.asByte()
		if !ok {
			return out, false
		}
		out[idx] = b
	}
	return out, true
}

// le32 reads a little-endian u32; used by the account-level Borsh sections
// rather than by the msgpack body.
func le32(b []byte) uint32 { return binary.LittleEndian.Uint32(b) }

// le64 reads a little-endian u64.
func le64(b []byte) uint64 { return binary.LittleEndian.Uint64(b) }
