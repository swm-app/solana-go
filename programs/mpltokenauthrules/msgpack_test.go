package mpltokenauthrules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func readOne(t *testing.T, b []byte) mpValue {
	t.Helper()
	v, err := newMsgpackReader(b).readValue()
	require.NoError(t, err)
	return v
}

func TestMsgpackScalars(t *testing.T) {
	t.Run("positive fixint", func(t *testing.T) {
		v := readOne(t, []byte{0x00})
		require.Equal(t, mpUint, v.typ)
		require.Equal(t, uint64(0), v.u)
		v = readOne(t, []byte{0x7f})
		require.Equal(t, uint64(127), v.u)
	})

	t.Run("negative fixint", func(t *testing.T) {
		v := readOne(t, []byte{0xff})
		require.Equal(t, mpInt, v.typ)
		require.Equal(t, int64(-1), v.i)
		v = readOne(t, []byte{0xe0})
		require.Equal(t, int64(-32), v.i)
	})

	t.Run("uint widths", func(t *testing.T) {
		require.Equal(t, uint64(0x80), readOne(t, []byte{0xcc, 0x80}).u)
		require.Equal(t, uint64(0x0102), readOne(t, []byte{0xcd, 0x01, 0x02}).u)
		require.Equal(t, uint64(0x01020304), readOne(t, []byte{0xce, 0x01, 0x02, 0x03, 0x04}).u)
		require.Equal(t, uint64(0x0102030405060708),
			readOne(t, []byte{0xcf, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}).u)
	})

	t.Run("int widths sign extend", func(t *testing.T) {
		require.Equal(t, int64(-1), readOne(t, []byte{0xd0, 0xff}).i)
		require.Equal(t, int64(-2), readOne(t, []byte{0xd1, 0xff, 0xfe}).i)
		require.Equal(t, int64(-1), readOne(t, []byte{0xd2, 0xff, 0xff, 0xff, 0xff}).i)
		require.Equal(t, int64(1),
			readOne(t, []byte{0xd3, 0, 0, 0, 0, 0, 0, 0, 1}).i)
	})

	t.Run("nil and bool", func(t *testing.T) {
		require.Equal(t, mpNil, readOne(t, []byte{0xc0}).typ)
		require.False(t, readOne(t, []byte{0xc2}).b)
		vt := readOne(t, []byte{0xc3})
		require.Equal(t, mpBool, vt.typ)
		require.True(t, vt.b)
	})
}

func TestMsgpackStrings(t *testing.T) {
	t.Run("fixstr", func(t *testing.T) {
		v := readOne(t, append([]byte{0xa3}, "abc"...))
		require.Equal(t, mpStr, v.typ)
		require.Equal(t, "abc", v.str)
	})
	t.Run("str8", func(t *testing.T) {
		v := readOne(t, append([]byte{0xd9, 0x03}, "xyz"...))
		require.Equal(t, "xyz", v.str)
	})
	t.Run("str16", func(t *testing.T) {
		v := readOne(t, append([]byte{0xda, 0x00, 0x02}, "hi"...))
		require.Equal(t, "hi", v.str)
	})
	t.Run("str32", func(t *testing.T) {
		v := readOne(t, append([]byte{0xdb, 0x00, 0x00, 0x00, 0x01}, "z"...))
		require.Equal(t, "z", v.str)
	})
}

func TestMsgpackBin(t *testing.T) {
	t.Run("bin8", func(t *testing.T) {
		v := readOne(t, []byte{0xc4, 0x02, 0xaa, 0xbb})
		require.Equal(t, mpBin, v.typ)
		require.Equal(t, []byte{0xaa, 0xbb}, v.bin)
	})
	t.Run("bin16", func(t *testing.T) {
		v := readOne(t, []byte{0xc5, 0x00, 0x01, 0xcc})
		require.Equal(t, []byte{0xcc}, v.bin)
	})
	t.Run("bin32", func(t *testing.T) {
		v := readOne(t, []byte{0xc6, 0x00, 0x00, 0x00, 0x01, 0xdd})
		require.Equal(t, []byte{0xdd}, v.bin)
	})
	t.Run("bin copies the input", func(t *testing.T) {
		in := []byte{0xc4, 0x01, 0x11}
		v := readOne(t, in)
		in[2] = 0x22 // mutate source
		require.Equal(t, []byte{0x11}, v.bin)
	})
}

func TestMsgpackArrays(t *testing.T) {
	t.Run("fixarray", func(t *testing.T) {
		v := readOne(t, []byte{0x92, 0x01, 0x02})
		require.Equal(t, mpArray, v.typ)
		require.Len(t, v.arr, 2)
		require.Equal(t, uint64(1), v.arr[0].u)
		require.Equal(t, uint64(2), v.arr[1].u)
	})
	t.Run("array16", func(t *testing.T) {
		v := readOne(t, []byte{0xdc, 0x00, 0x01, 0x05})
		require.Len(t, v.arr, 1)
		require.Equal(t, uint64(5), v.arr[0].u)
	})
	t.Run("array32", func(t *testing.T) {
		v := readOne(t, []byte{0xdd, 0x00, 0x00, 0x00, 0x01, 0x06})
		require.Len(t, v.arr, 1)
		require.Equal(t, uint64(6), v.arr[0].u)
	})
}

func TestMsgpackMaps(t *testing.T) {
	t.Run("fixmap", func(t *testing.T) {
		v := readOne(t, append(append([]byte{0x81, 0xa1}, "k"...), 0x09))
		require.Equal(t, mpMap, v.typ)
		require.Len(t, v.pair, 1)
		require.Equal(t, "k", v.pair[0].key.str)
		require.Equal(t, uint64(9), v.pair[0].val.u)
	})
	t.Run("map16", func(t *testing.T) {
		v := readOne(t, append(append([]byte{0xde, 0x00, 0x01, 0xa1}, "k"...), 0x09))
		require.Len(t, v.pair, 1)
	})
	t.Run("map32", func(t *testing.T) {
		v := readOne(t, append(append([]byte{0xdf, 0x00, 0x00, 0x00, 0x01, 0xa1}, "k"...), 0x09))
		require.Len(t, v.pair, 1)
	})
}

func TestMsgpackPubkeyHelper(t *testing.T) {
	// A 32-element array of byte integers decodes as a pubkey; other shapes do not.
	elems := make([][]byte, 32)
	for i := range elems {
		elems[i] = encUint(uint64(i))
	}
	v := readOne(t, encArr(elems...))
	pk, ok := v.asPubkey()
	require.True(t, ok)
	require.Equal(t, byte(31), pk[31])

	_, ok = readOne(t, encArr(encUint(1), encUint(2))).asPubkey()
	require.False(t, ok)
	// An element out of byte range fails the pubkey cast.
	bad := make([][]byte, 32)
	for i := range bad {
		bad[i] = encUint(300)
	}
	_, ok = readOne(t, encArr(bad...)).asPubkey()
	require.False(t, ok)
}

// TestMsgpackRejectsFloatAndExt confirms floats and extension types - which
// never appear in a RuleSetV1 - are rejected rather than silently skipped.
func TestMsgpackRejectsFloatAndExt(t *testing.T) {
	rejected := map[string][]byte{
		"float32":   {0xca, 0, 0, 0, 0},
		"float64":   {0xcb, 0, 0, 0, 0, 0, 0, 0, 0},
		"fixext1":   {0xd4, 0, 0},
		"fixext2":   {0xd5, 0, 0, 0},
		"ext8":      {0xc7, 0x01, 0x00, 0x00},
		"ext16":     {0xc8, 0x00, 0x01, 0x00, 0x00},
		"ext32":     {0xc9, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00},
		"neverused": {0xc1},
	}
	for name, b := range rejected {
		t.Run(name, func(t *testing.T) {
			_, err := newMsgpackReader(b).readValue()
			require.Error(t, err)
			require.Contains(t, err.Error(), "unsupported msgpack tag")
		})
	}
}

func TestMsgpackTruncation(t *testing.T) {
	cases := map[string][]byte{
		"empty input":          {},
		"uint16 missing byte":  {0xcd, 0x01},
		"fixstr short":         {0xa3, 'a'},
		"str8 length only":     {0xd9},
		"bin8 short body":      {0xc4, 0x04, 0x01},
		"fixarray short":       {0x92, 0x01},
		"fixmap missing value": {0x81, 0xa1, 'k'},
	}
	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := newMsgpackReader(b).readValue()
			require.Error(t, err)
		})
	}
}

// TestMsgpackNoUnboundedPrealloc confirms a declared length larger than the
// remaining input is rejected up front, before any slice of that size is
// allocated.
func TestMsgpackNoUnboundedPrealloc(t *testing.T) {
	// array32 declaring ~4 billion elements with no element bytes present.
	_, err := newMsgpackReader([]byte{0xdd, 0xff, 0xff, 0xff, 0xff}).readValue()
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds remaining")

	// map32 declaring a huge pair count with no data.
	_, err = newMsgpackReader([]byte{0xdf, 0xff, 0xff, 0xff, 0xff}).readValue()
	require.Error(t, err)

	// str32 / bin32 declaring a huge byte length.
	_, err = newMsgpackReader([]byte{0xdb, 0xff, 0xff, 0xff, 0xff}).readValue()
	require.Error(t, err)
	_, err = newMsgpackReader([]byte{0xc6, 0xff, 0xff, 0xff, 0xff}).readValue()
	require.Error(t, err)
}

// TestMsgpackDepthCap pins the recursion ceiling: 1024 nested containers parse,
// 1025 error. The reader increments depth before descending, so the boundary is
// the maxMsgpackDepth constant exactly.
func TestMsgpackDepthCap(t *testing.T) {
	// k nested fixarrays: k-1 single-element wrappers around an empty innermost.
	nest := func(k int) []byte {
		out := make([]byte, 0, k)
		for i := 0; i < k-1; i++ {
			out = append(out, 0x91)
		}
		return append(out, 0x90)
	}

	_, err := newMsgpackReader(nest(maxMsgpackDepth)).readValue()
	require.NoError(t, err)

	_, err = newMsgpackReader(nest(maxMsgpackDepth + 1)).readValue()
	require.Error(t, err)
	require.Contains(t, err.Error(), "maximum depth")
}
