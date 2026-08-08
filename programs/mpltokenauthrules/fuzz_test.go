package mpltokenauthrules

import (
	"encoding/binary"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

// evaluateAll walks every revision of a parsed account and evaluates a spread of
// operations under several contexts. It exists to drive the full decode ->
// evaluate pipeline so a fuzzing or adversarial input that decodes cannot panic
// during evaluation either.
func evaluateAll(acct *RuleSetAccount) {
	if acct == nil {
		return
	}
	pk := solanago.PublicKey{1}
	contexts := []EvalContext{
		{},
		{Payload: Payload{}},
		{Payload: Payload{
			PayloadKeyAmount:      NewNumberPayload(1),
			PayloadKeySource:      NewPubkeyPayload(pk),
			PayloadKeyDestination: NewSeedsPayload([][]byte{{1, 2}}),
		}},
	}
	ops := []string{
		OperationTransferOwner, "Transfer", "Transfer:WalletToWallet",
		"Delegate", "Delegate:Sale", "Nonexistent:Operation",
	}
	for i := range acct.Revisions {
		rs := acct.Revisions[i].RuleSet
		if rs == nil {
			continue
		}
		for _, op := range ops {
			for _, ctx := range contexts {
				_ = EvaluateOperation(rs, op, ctx)
			}
		}
	}
}

// malformedAccounts returns a table of intentionally broken account byte-slices
// derived from a real fixture and from synthetic builders. None may panic when
// parsed; each is expected to either error or parse cleanly.
func malformedAccounts(t *testing.T) map[string][]byte {
	t.Helper()
	real := loadFixtureBytes(t, fixtureFoundation)
	out := map[string][]byte{
		"nil":             nil,
		"empty":           {},
		"one byte":        {1},
		"header only":     real[:ruleSetHeaderLen],
		"header plus one": real[:ruleSetHeaderLen+1],
	}

	// Truncations at a spread of prefixes of the real account.
	for _, frac := range []int{4, 3, 2} {
		n := len(real) / frac
		cp := make([]byte, n)
		copy(cp, real)
		out["truncated 1/"+string(rune('0'+frac))] = cp
	}
	// Truncated one byte short of the full account.
	short := make([]byte, len(real)-1)
	copy(short, real)
	out["truncated minus one"] = short

	// Rev-map version location out of range (below the header and past the end).
	locLow := make([]byte, len(real))
	copy(locLow, real)
	binary.LittleEndian.PutUint64(locLow[1:9], 3)
	out["rev map loc below header"] = locLow

	locHigh := make([]byte, len(real))
	copy(locHigh, real)
	binary.LittleEndian.PutUint64(locHigh[1:9], uint64(len(real))+1000)
	out["rev map loc past end"] = locHigh

	// A synthetic well-formed account, then corruptions of its revision map.
	owner := mustPubkey(t, "So11111111111111111111111111111111111111112")
	rs := v2RuleSet(owner, "adv", []string{"Transfer:Owner"},
		[][]byte{v2AmountRule(1, 2, PayloadKeyAmount)})
	good := wrapAccount([][]byte{rs}, 0)
	loc := int(binary.LittleEndian.Uint64(good[1:9]))

	// Huge Borsh vec length in the revision map.
	hugeCount := make([]byte, len(good))
	copy(hugeCount, good)
	binary.LittleEndian.PutUint32(hugeCount[loc+1:loc+5], 0xffffffff)
	out["huge rev map count"] = hugeCount

	// Descending / overlapping revision offset (points before the header).
	badOffset := make([]byte, len(good))
	copy(badOffset, good)
	binary.LittleEndian.PutUint64(badOffset[loc+5:loc+13], 0)
	out["revision offset in header"] = badOffset

	// Revision offset past the account end.
	pastEnd := make([]byte, len(good))
	copy(pastEnd, good)
	binary.LittleEndian.PutUint64(pastEnd[loc+5:loc+13], uint64(len(good))+50)
	out["revision offset past end"] = pastEnd

	// A V2 revision nested far beyond the depth cap, wrapped in an account.
	deep := v2PassRule()
	for i := 0; i < maxV2RuleDepth+20; i++ {
		deep = v2NotRule(deep)
	}
	deepRS := v2RuleSet(owner, "deep", []string{"Transfer:Owner"}, [][]byte{deep})
	out["deep nesting"] = wrapAccount([][]byte{deepRS}, 0)

	// A revision map claiming two revisions but offsets in descending order.
	two := wrapAccount([][]byte{rs, rs}, 0)
	loc2 := int(binary.LittleEndian.Uint64(two[1:9]))
	o0 := binary.LittleEndian.Uint64(two[loc2+5 : loc2+13])
	o1 := binary.LittleEndian.Uint64(two[loc2+13 : loc2+21])
	binary.LittleEndian.PutUint64(two[loc2+5:loc2+13], o1)
	binary.LittleEndian.PutUint64(two[loc2+13:loc2+21], o0)
	out["descending revision offsets"] = two

	return out
}

// TestAdversarialNoPanic feeds a table of malformed accounts through the parser
// and evaluator: every input must either parse or return an error, never panic.
func TestAdversarialNoPanic(t *testing.T) {
	for name, data := range malformedAccounts(t) {
		t.Run(name, func(t *testing.T) {
			require.NotPanics(t, func() {
				acct, err := ParseAccount(data)
				if err == nil {
					evaluateAll(acct)
				}
			})
		})
	}
}

func fuzzSeeds(f *testing.F) {
	for _, name := range []string{
		fixtureFoundation, fixtureCompatibility, fixtureElementerra,
		fixtureRHO, fixtureNoPermissions, fixtureMixedMe, fixtureMixedAmigos,
	} {
		raw, err := readFixtureRaw(name)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(raw)
	}
	// A synthetic V2 account seed.
	owner := solanago.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	rs := v2RuleSet(owner, "seed", []string{"Transfer:Owner"},
		[][]byte{v2AmountRule(1, 2, PayloadKeyAmount)})
	f.Add(wrapAccount([][]byte{rs}, 0))
}

// FuzzParseAccount asserts the parser never panics on arbitrary input and that a
// successful parse can be fully evaluated without panicking.
func FuzzParseAccount(f *testing.F) {
	fuzzSeeds(f)
	f.Fuzz(func(t *testing.T, data []byte) {
		acct, err := ParseAccount(data)
		if err != nil {
			return
		}
		evaluateAll(acct)
	})
}

// FuzzEvaluate reuses the same corpus and drives parse-then-evaluate; it guards
// the evaluator specifically against panics on decoded-but-hostile rule trees.
func FuzzEvaluate(f *testing.F) {
	fuzzSeeds(f)
	f.Fuzz(func(t *testing.T, data []byte) {
		acct, err := ParseAccount(data)
		if err != nil {
			return
		}
		evaluateAll(acct)
	})
}
