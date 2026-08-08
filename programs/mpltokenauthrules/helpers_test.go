package mpltokenauthrules

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"strings"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

// Fixture file base names (without the testdata/ prefix or .b64 suffix). Every
// fixture is a raw mainnet rule-set PDA account, base64-encoded, captured
// 2026-07-17 via Helius getAccountInfo / getMultipleAccounts. See
// testdata/README.md for provenance of each.
const (
	// V1 (msgpack) fixtures.
	fixtureFoundation    = "eBJLFYPxJmMGKuFwpDWkzxZeUrad92kZRC5BJLpzyT9"  // "Metaplex Foundation Rule Set", 9 revisions
	fixtureCompatibility = "AdH2Utn6Fus15ZhtenW4hZBQnvtLgM1YCW2MfVp7pYS5" // "Compatibility Rule Set", 1 revision
	fixtureElementerra   = "7QU6FKMp6gFD24DHWB8AayQgArGZHtYHFr1rc94dhwga" // "elementerra", 2 revisions

	// V2 (binary TLV) fixtures.
	fixtureRHO           = "DJw6T8nxKf4V4b15mPfSZuZQDG2vaXXaBKvNJH2jCsWk" // "RHO", 9 revisions all V2 (largest)
	fixtureNoPermissions = "Hjf4FzbS9zgZV5rda9ZLoiTqCKNrDujkbEVamp3b5fM6" // "no_permissions", 1 V2 revision, 0 operations (smallest, padded)

	// Mixed accounts carrying both V1 and V2 revisions.
	fixtureMixedMe     = "8icKNaVtc6TGF8KPg2KHkm39JPTwtfSBc6DgZqx8UpTi" // "me", rev0 V1 then rev1 V2
	fixtureMixedAmigos = "HThAqeQQvRQQaPvsPaDZ4h4C94GoDWEPvqymxNL7WBrn" // "Amigos Odyssey", rev0/1 V2 then rev2 V1
)

// readFixtureRaw reads and base64-decodes a testdata fixture into raw account
// bytes without a testing handle, for use from fuzz seeding.
func readFixtureRaw(name string) ([]byte, error) {
	raw, err := os.ReadFile("testdata/" + name + ".b64")
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(strings.TrimSpace(string(raw)))
}

// loadFixtureBytes reads and base64-decodes a testdata fixture into raw account
// bytes.
func loadFixtureBytes(t *testing.T, name string) []byte {
	t.Helper()
	data, err := readFixtureRaw(name)
	require.NoError(t, err)
	return data
}

// loadFixtureAccount reads a fixture and parses it, requiring success.
func loadFixtureAccount(t *testing.T, name string) *RuleSetAccount {
	t.Helper()
	acct, err := ParseAccount(loadFixtureBytes(t, name))
	require.NoError(t, err)
	return acct
}

func mustPubkey(t *testing.T, s string) solanago.PublicKey {
	t.Helper()
	pk, err := solanago.PublicKeyFromBase58(s)
	require.NoError(t, err)
	return pk
}

// operationNames returns the ordered operation names of a rule set.
func operationNames(rs *RuleSet) []string { return rs.OperationNames }

// ---------------------------------------------------------------------------
// Minimal msgpack encoder for synthesizing V1 rule-set payloads in tests. It
// mirrors what rmp_serde emits: structs as arrays, enums as single-entry maps,
// and Solana pubkeys as a 32-element array of byte integers.
// ---------------------------------------------------------------------------

// encUint encodes a non-negative integer using the narrowest fixint/uint width.
func encUint(n uint64) []byte {
	switch {
	case n <= 0x7f:
		return []byte{byte(n)}
	case n <= 0xff:
		return []byte{0xcc, byte(n)}
	case n <= 0xffff:
		return []byte{0xcd, byte(n >> 8), byte(n)}
	case n <= 0xffffffff:
		return []byte{0xce, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	default:
		b := make([]byte, 9)
		b[0] = 0xcf
		binary.BigEndian.PutUint64(b[1:], n)
		return b
	}
}

// encStr encodes a string using the narrowest fixstr/str width.
func encStr(s string) []byte {
	n := len(s)
	var hdr []byte
	switch {
	case n <= 0x1f:
		hdr = []byte{0xa0 | byte(n)}
	case n <= 0xff:
		hdr = []byte{0xd9, byte(n)}
	case n <= 0xffff:
		hdr = []byte{0xda, byte(n >> 8), byte(n)}
	default:
		hdr = []byte{0xdb, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	}
	return append(hdr, s...)
}

// encArrHeader encodes an array header for n elements.
func encArrHeader(n int) []byte {
	switch {
	case n <= 0x0f:
		return []byte{0x90 | byte(n)}
	case n <= 0xffff:
		return []byte{0xdc, byte(n >> 8), byte(n)}
	default:
		return []byte{0xdd, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	}
}

// encArr encodes an array of already-encoded element values.
func encArr(elems ...[]byte) []byte {
	out := encArrHeader(len(elems))
	for _, e := range elems {
		out = append(out, e...)
	}
	return out
}

// encMapHeader encodes a map header for n key/value pairs.
func encMapHeader(n int) []byte {
	switch {
	case n <= 0x0f:
		return []byte{0x80 | byte(n)}
	case n <= 0xffff:
		return []byte{0xde, byte(n >> 8), byte(n)}
	default:
		return []byte{0xdf, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	}
}

// encKV is one ordered key/value entry of an encoded map.
type encKV struct {
	key []byte
	val []byte
}

// encMap encodes an ordered map from already-encoded key/value byte slices.
func encMap(pairs ...encKV) []byte {
	out := encMapHeader(len(pairs))
	for _, p := range pairs {
		out = append(out, p.key...)
		out = append(out, p.val...)
	}
	return out
}

// encPubkey encodes a 32-byte pubkey as rmp_serde does: a 32-element array of
// byte integers (not a msgpack bin).
func encPubkey(pk solanago.PublicKey) []byte {
	elems := make([][]byte, 32)
	for i, b := range pk.Bytes() {
		elems[i] = encUint(uint64(b))
	}
	return encArr(elems...)
}

// v1Pass / v1Namespace encode the two unit-variant rules (bare strings).
func v1Pass() []byte      { return encStr("Pass") }
func v1Namespace() []byte { return encStr("Namespace") }

// v1Variant encodes a data-carrying rule variant: a single-entry map whose key
// is the variant name and whose value is the array of the variant's fields.
func v1Variant(name string, fields ...[]byte) []byte {
	return encMap(encKV{key: encStr(name), val: encArr(fields...)})
}

func v1Amount(amount uint64, op, field string) []byte {
	return v1Variant("Amount", encUint(amount), encStr(op), encStr(field))
}

// v1RuleSet builds a complete V1 revision (lib-version byte 0x01 followed by the
// msgpack struct array [lib_version, owner, name, operations]).
func v1RuleSet(owner solanago.PublicKey, name string, ops []encKV) []byte {
	top := encArr(encUint(uint64(libVersionV1)), encPubkey(owner), encStr(name), encMap(ops...))
	return append([]byte{libVersionV1}, top...)
}

// ---------------------------------------------------------------------------
// V2 binary rule / rule-set builders. Layouts confirmed against the Rust source
// programs/token-auth-rules/src/state/v2/*.
// ---------------------------------------------------------------------------

// str32 encodes a NUL-padded 32-byte fixed field.
func str32(s string) []byte {
	var b [32]byte
	copy(b[:], s)
	return b[:]
}

// v2Rule assembles a RuleV2 TLV record: [constraint_type u32, length u32] + body.
func v2Rule(constraintType uint32, body []byte) []byte {
	out := make([]byte, 8)
	binary.LittleEndian.PutUint32(out[0:4], constraintType)
	binary.LittleEndian.PutUint32(out[4:8], uint32(len(body)))
	return append(out, body...)
}

func v2PassRule() []byte      { return v2Rule(v2Pass, nil) }
func v2NamespaceRule() []byte { return v2Rule(v2Namespace, nil) }

func v2AdditionalSignerRule(pk solanago.PublicKey) []byte {
	return v2Rule(v2AdditionalSigner, pk.Bytes())
}

func v2AmountRule(amount uint64, operator uint64, field string) []byte {
	body := make([]byte, 0, 48)
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], amount)
	body = append(body, buf[:]...)
	binary.LittleEndian.PutUint64(buf[:], operator)
	body = append(body, buf[:]...)
	body = append(body, str32(field)...)
	return v2Rule(v2Amount, body)
}

func v2PubkeyMatchRule(pk solanago.PublicKey, field string) []byte {
	return v2Rule(v2PubkeyMatch, append(pk.Bytes(), str32(field)...))
}

func v2ProgramOwnedRule(pk solanago.PublicKey, field string) []byte {
	return v2Rule(v2ProgramOwned, append(pk.Bytes(), str32(field)...))
}

func v2FrequencyRule(pk solanago.PublicKey) []byte {
	return v2Rule(v2Frequency, pk.Bytes())
}

func v2IsWalletRule(field string) []byte {
	return v2Rule(v2IsWallet, str32(field))
}

// v2PDAMatchRule builds a PDAMatch. A nil program serializes the all-zero
// sentinel (derive against the account owner).
func v2PDAMatchRule(program *solanago.PublicKey, pdaField, seedsField string) []byte {
	body := make([]byte, 0, 96)
	if program == nil {
		body = append(body, make([]byte, 32)...)
	} else {
		body = append(body, program.Bytes()...)
	}
	body = append(body, str32(pdaField)...)
	body = append(body, str32(seedsField)...)
	return v2Rule(v2PDAMatch, body)
}

func v2ListRule(constraintType uint32, field string, keys []solanago.PublicKey) []byte {
	body := str32(field)
	for _, k := range keys {
		body = append(body, k.Bytes()...)
	}
	return v2Rule(constraintType, body)
}

func v2TreeRule(constraintType uint32, pubkeyField, proofField string, root [32]byte) []byte {
	body := str32(pubkeyField)
	body = append(body, str32(proofField)...)
	body = append(body, root[:]...)
	return v2Rule(constraintType, body)
}

func v2NotRule(child []byte) []byte { return v2Rule(v2Not, child) }

func v2Composite(constraintType uint32, children ...[]byte) []byte {
	body := make([]byte, 8)
	binary.LittleEndian.PutUint64(body[0:8], uint64(len(children)))
	for _, c := range children {
		body = append(body, c...)
	}
	return v2Rule(constraintType, body)
}

func v2AllRule(children ...[]byte) []byte { return v2Composite(v2All, children...) }
func v2AnyRule(children ...[]byte) []byte { return v2Composite(v2Any, children...) }

// v2RuleSet builds a full V2 revision payload: [u32 lib_version|0, u32 size,
// owner(32), name Str32(32), operations(size*Str32), rules...]. The number of
// operations must equal the number of rules.
func v2RuleSet(owner solanago.PublicKey, name string, operations []string, rules [][]byte) []byte {
	out := make([]byte, 0, 72)
	out = append(out, byte(libVersionV2), 0, 0, 0)
	var sizeBuf [4]byte
	binary.LittleEndian.PutUint32(sizeBuf[:], uint32(len(operations)))
	out = append(out, sizeBuf[:]...)
	out = append(out, owner.Bytes()...)
	out = append(out, str32(name)...)
	for _, op := range operations {
		out = append(out, str32(op)...)
	}
	for _, r := range rules {
		out = append(out, r...)
	}
	return out
}

// wrapAccount frames one or more revision byte-slices into a complete rule-set
// PDA account: 9-byte header, optional zero padding, the revisions back to back,
// the revision-map version byte, then the Borsh revision map (u32 count + u64
// absolute offsets). padBefore inserts that many zero bytes between the header
// and the first revision, reproducing the on-chain alignment padding.
func wrapAccount(revisions [][]byte, padBefore int) []byte {
	out := make([]byte, 0)
	out = append(out, 1) // Key: a live rule set
	locPos := len(out)
	out = append(out, make([]byte, 8)...) // rev-map version location placeholder
	out = append(out, make([]byte, padBefore)...)

	offsets := make([]uint64, 0, len(revisions))
	for _, r := range revisions {
		offsets = append(offsets, uint64(len(out)))
		out = append(out, r...)
	}

	loc := uint64(len(out))
	binary.LittleEndian.PutUint64(out[locPos:locPos+8], loc)
	out = append(out, revMapVersion)

	var cnt [4]byte
	binary.LittleEndian.PutUint32(cnt[:], uint32(len(offsets)))
	out = append(out, cnt[:]...)
	for _, o := range offsets {
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], o)
		out = append(out, b[:]...)
	}
	return out
}
