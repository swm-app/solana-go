package mpltokenauthrules

import (
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

func testOwner(t *testing.T) solanago.PublicKey {
	return mustPubkey(t, "7qCcNjYqadaN1w2xqYqnPE4q7Jgoxsc7eSpKVxAq2VWh")
}

// TestV2SyntheticAccountParsesAndEvaluates wraps a synthetic V2 revision in a
// full account - including alignment padding between the header and the first
// revision - and checks it decodes and evaluates. This mirrors the on-chain
// layout where revision offsets, not a fixed offset 9, locate each revision.
func TestV2SyntheticAccountParsesAndEvaluates(t *testing.T) {
	owner := testOwner(t)
	rs := v2RuleSet(owner, "synthetic", []string{"Transfer:Owner"},
		[][]byte{v2AmountRule(1, 2 /* Eq */, PayloadKeyAmount)})

	for _, pad := range []int{0, 7} {
		data := wrapAccount([][]byte{rs}, pad)
		acct, err := ParseAccount(data)
		require.NoErrorf(t, err, "pad=%d", pad)
		require.Len(t, acct.Revisions, 1)
		require.Equal(t, uint8(libVersionV2), acct.Revisions[0].LibVersion)
		got := acct.Revisions[0].RuleSet
		require.NotNil(t, got)
		require.Equal(t, owner, got.Owner)
		require.Equal(t, "synthetic", got.Name)
		require.Equal(t, []string{"Transfer:Owner"}, operationNames(got))

		require.Equal(t, Allow, EvaluateOperation(got, OperationTransferOwner, EvalContext{
			Payload: Payload{PayloadKeyAmount: NewNumberPayload(1)},
		}).Verdict)
		require.Equal(t, Deny, EvaluateOperation(got, OperationTransferOwner, EvalContext{
			Payload: Payload{PayloadKeyAmount: NewNumberPayload(2)},
		}).Verdict)
	}
}

// TestV2EveryConstraintType builds a rule set with one operation per constraint
// type and confirms each round-parses into the expected unified Rule kind with
// the correct decoded fields.
func TestV2EveryConstraintType(t *testing.T) {
	owner := testOwner(t)
	pk := mustPubkey(t, "So11111111111111111111111111111111111111112")
	var root [32]byte
	root[0] = 0xab

	type tc struct {
		op    string
		rule  []byte
		check func(t *testing.T, r Rule)
	}
	cases := []tc{
		{"AdditionalSigner", v2AdditionalSignerRule(pk), func(t *testing.T, r Rule) {
			require.Equal(t, RuleAdditionalSigner, r.Kind)
			require.Equal(t, pk, r.PublicKey)
		}},
		{"All", v2AllRule(v2PassRule()), func(t *testing.T, r Rule) {
			require.Equal(t, RuleAll, r.Kind)
			require.Len(t, r.Rules, 1)
			require.Equal(t, RulePass, r.Rules[0].Kind)
		}},
		{"Amount", v2AmountRule(7, 3 /* GtEq */, PayloadKeyAmount), func(t *testing.T, r Rule) {
			require.Equal(t, RuleAmount, r.Kind)
			require.Equal(t, uint64(7), r.Amount)
			require.Equal(t, CompareOpGtEq, r.Operator)
			require.Equal(t, PayloadKeyAmount, r.Field)
		}},
		{"Any", v2AnyRule(v2PassRule(), v2NamespaceRule()), func(t *testing.T, r Rule) {
			require.Equal(t, RuleAny, r.Kind)
			require.Len(t, r.Rules, 2)
		}},
		{"Frequency", v2FrequencyRule(pk), func(t *testing.T, r Rule) {
			require.Equal(t, RuleFrequency, r.Kind)
			require.Equal(t, pk, r.PublicKey)
		}},
		{"IsWallet", v2IsWalletRule(PayloadKeySource), func(t *testing.T, r Rule) {
			require.Equal(t, RuleIsWallet, r.Kind)
			require.Equal(t, PayloadKeySource, r.Field)
		}},
		{"Namespace", v2NamespaceRule(), func(t *testing.T, r Rule) {
			require.Equal(t, RuleNamespace, r.Kind)
		}},
		{"Not", v2NotRule(v2PassRule()), func(t *testing.T, r Rule) {
			require.Equal(t, RuleNot, r.Kind)
			require.NotNil(t, r.Rule)
			require.Equal(t, RulePass, r.Rule.Kind)
		}},
		{"Pass", v2PassRule(), func(t *testing.T, r Rule) {
			require.Equal(t, RulePass, r.Kind)
		}},
		{"PDAMatch", v2PDAMatchRule(&pk, PayloadKeyDestination, PayloadKeySource), func(t *testing.T, r Rule) {
			require.Equal(t, RulePDAMatch, r.Kind)
			require.NotNil(t, r.Program)
			require.Equal(t, pk, *r.Program)
			require.Equal(t, PayloadKeyDestination, r.Field)
			require.Equal(t, PayloadKeySource, r.Field2)
		}},
		{"PDAMatchDefaultProgram", v2PDAMatchRule(nil, PayloadKeyDestination, PayloadKeySource), func(t *testing.T, r Rule) {
			require.Equal(t, RulePDAMatch, r.Kind)
			require.Nil(t, r.Program) // all-zero sentinel decodes to nil
		}},
		{"ProgramOwned", v2ProgramOwnedRule(pk, PayloadKeySource), func(t *testing.T, r Rule) {
			require.Equal(t, RuleProgramOwned, r.Kind)
			require.Equal(t, pk, r.PublicKey)
			require.Equal(t, PayloadKeySource, r.Field)
		}},
		{"ProgramOwnedList", v2ListRule(v2ProgramOwnedList, PayloadKeySource, []solanago.PublicKey{pk, owner}), func(t *testing.T, r Rule) {
			require.Equal(t, RuleProgramOwnedList, r.Kind)
			require.Equal(t, []solanago.PublicKey{pk, owner}, r.PublicKeys)
		}},
		{"ProgramOwnedTree", v2TreeRule(v2ProgramOwnedTree, PayloadKeySource, PayloadKeyDestination, root), func(t *testing.T, r Rule) {
			require.Equal(t, RuleProgramOwnedTree, r.Kind)
			require.Equal(t, PayloadKeySource, r.Field)
			require.Equal(t, PayloadKeyDestination, r.Field2)
			require.Equal(t, root, r.Root)
		}},
		{"PubkeyListMatch", v2ListRule(v2PubkeyListMatch, PayloadKeyDestination, []solanago.PublicKey{pk}), func(t *testing.T, r Rule) {
			require.Equal(t, RulePubkeyListMatch, r.Kind)
			require.Equal(t, []solanago.PublicKey{pk}, r.PublicKeys)
		}},
		{"PubkeyMatch", v2PubkeyMatchRule(pk, PayloadKeyDestination), func(t *testing.T, r Rule) {
			require.Equal(t, RulePubkeyMatch, r.Kind)
			require.Equal(t, pk, r.PublicKey)
			require.Equal(t, PayloadKeyDestination, r.Field)
		}},
		{"PubkeyTreeMatch", v2TreeRule(v2PubkeyTreeMatch, PayloadKeySource, PayloadKeyDestination, root), func(t *testing.T, r Rule) {
			require.Equal(t, RulePubkeyTreeMatch, r.Kind)
			require.Equal(t, root, r.Root)
		}},
	}

	ops := make([]string, len(cases))
	rules := make([][]byte, len(cases))
	for i, c := range cases {
		ops[i] = c.op
		rules[i] = c.rule
	}
	acct, err := ParseAccount(wrapAccount([][]byte{v2RuleSet(owner, "kinds", ops, rules)}, 0))
	require.NoError(t, err)
	rs := acct.Revisions[0].RuleSet
	require.NotNil(t, rs)
	require.Len(t, operationNames(rs), len(cases))
	for _, c := range cases {
		t.Run(c.op, func(t *testing.T) {
			c.check(t, rs.Operations[c.op])
		})
	}
}

// TestV2UnknownDiscriminantPreserved confirms an unrecognized constraint type is
// preserved as an Unknown node that still consumes its declared length, so
// following rules keep parsing.
func TestV2UnknownDiscriminantPreserved(t *testing.T) {
	owner := testOwner(t)
	unknown := v2Rule(9999, []byte{0x01, 0x02, 0x03, 0x04})
	rs := v2RuleSet(owner, "mixedkinds",
		[]string{"first", "second", "third"},
		[][]byte{v2PassRule(), unknown, v2NamespaceRule()})

	acct, err := ParseAccount(wrapAccount([][]byte{rs}, 0))
	require.NoError(t, err)
	got := acct.Revisions[0].RuleSet
	require.NotNil(t, got)

	require.Equal(t, RulePass, got.Operations["first"].Kind)
	u := got.Operations["second"]
	require.Equal(t, RuleUnknown, u.Kind)
	require.Equal(t, "ConstraintType(9999)", u.RawKind)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, u.Raw)
	require.Equal(t, RuleNamespace, got.Operations["third"].Kind)

	// An unknown rule evaluates to Indeterminate.
	require.Equal(t, Indeterminate, EvaluateRule(u, EvalContext{}).Verdict)
}

// TestV2MalformedBodies covers structural errors that must be reported, not
// panicked or silently accepted.
func TestV2MalformedBodies(t *testing.T) {
	t.Run("pubkey list not a multiple of 32", func(t *testing.T) {
		bad := v2Rule(v2ProgramOwnedList, append(str32(PayloadKeySource), 0x01, 0x02, 0x03))
		_, _, err := parseRuleV2(bad, 0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not a multiple of")
	})

	t.Run("fixed-length body wrong size", func(t *testing.T) {
		bad := v2Rule(v2Amount, make([]byte, 47)) // Amount body must be 48
		_, _, err := parseRuleV2(bad, 0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be 48 bytes")
	})

	t.Run("declared length overruns available", func(t *testing.T) {
		// A header claiming a 100-byte body but only 4 body bytes present.
		bad := []byte{byte(v2Pass), 0, 0, 0, 100, 0, 0, 0, 1, 2, 3, 4}
		_, _, err := parseRuleV2(bad, 0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "overruns")
	})
}

// TestV2ZeroLengthTLVCannotLoop confirms an All whose declared child count
// cannot be satisfied by the body terminates with an error rather than spinning:
// every child consumes at least the 8-byte header, so the parser always makes
// progress or fails.
func TestV2ZeroLengthTLVCannotLoop(t *testing.T) {
	// An All declaring 5 children but a body holding only the 8-byte count.
	body := make([]byte, 8)
	body[0] = 5 // size = 5
	all := v2Rule(v2All, body)
	_, _, err := parseRuleV2(all, 0)
	require.Error(t, err)

	// An All declaring 1 child but with a truncated child header also fails.
	body2 := make([]byte, 8+4)
	body2[0] = 1
	_, _, err = parseRuleV2(v2Rule(v2All, body2), 0)
	require.Error(t, err)
}

// TestV2DepthCap pins the V2 nesting ceiling: maxV2RuleDepth nested Not rules
// parse, one more errors.
func TestV2DepthCap(t *testing.T) {
	nest := func(depth int) []byte {
		rule := v2PassRule()
		for i := 0; i < depth; i++ {
			rule = v2NotRule(rule)
		}
		return rule
	}

	_, _, err := parseRuleV2(nest(maxV2RuleDepth), 0)
	require.NoError(t, err)

	_, _, err = parseRuleV2(nest(maxV2RuleDepth+1), 0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "maximum depth")
}

// TestSyntheticMixedAccount builds an account whose revision 0 is a synthetic V1
// (msgpack) revision and revision 1 is a synthetic V2 (binary) revision, the
// same shape observed on mainnet.
func TestSyntheticMixedAccount(t *testing.T) {
	owner := mustPubkey(t, "So11111111111111111111111111111111111111112")

	v1 := v1RuleSet(owner, "mix-v1", []encKV{
		{key: encStr(OperationTransferOwner), val: v1Pass()},
		{key: encStr("Transfer:Base"), val: v1Amount(1, "Eq", PayloadKeyAmount)},
	})
	v2 := v2RuleSet(owner, "mix-v2", []string{OperationTransferOwner},
		[][]byte{v2AmountRule(1, 2 /* Eq */, PayloadKeyAmount)})

	acct, err := ParseAccount(wrapAccount([][]byte{v1, v2}, 0))
	require.NoError(t, err)
	require.Len(t, acct.Revisions, 2)

	require.Equal(t, uint8(libVersionV1), acct.Revisions[0].LibVersion)
	require.Equal(t, uint8(libVersionV2), acct.Revisions[1].LibVersion)
	require.NotNil(t, acct.Revisions[0].RuleSet)
	require.NotNil(t, acct.Revisions[1].RuleSet)

	require.Equal(t, "mix-v1", acct.Revisions[0].RuleSet.Name)
	require.Equal(t, "mix-v2", acct.Revisions[1].RuleSet.Name)

	// V1 revision: Transfer:Owner is Pass => Allow.
	require.Equal(t, Allow, EvaluateOperation(acct.Revisions[0].RuleSet, OperationTransferOwner, EvalContext{}).Verdict)
	// V2 latest: Transfer:Owner is Amount==1.
	require.Equal(t, Allow, EvaluateOperation(acct.Latest().RuleSet, OperationTransferOwner, EvalContext{
		Payload: Payload{PayloadKeyAmount: NewNumberPayload(1)},
	}).Verdict)
}

// TestV1SyntheticVariants exercises the V1 msgpack decoder on hand-built rules,
// including an unrecognized enum variant that is preserved as Unknown.
func TestV1SyntheticVariants(t *testing.T) {
	owner := mustPubkey(t, "So11111111111111111111111111111111111111112")
	rev := v1RuleSet(owner, "v1syn", []encKV{
		{key: encStr("Transfer:Owner"), val: v1Amount(2, "GtEq", PayloadKeyAmount)},
		{key: encStr("Transfer:Base"), val: v1Namespace()},
		{key: encStr("Exotic"), val: v1Variant("SomeFutureRule", encUint(1))},
	})

	rs, err := ParseRuleSet(rev)
	require.NoError(t, err)
	require.Equal(t, "v1syn", rs.Name)
	require.Equal(t, owner, rs.Owner)

	amount := rs.Operations["Transfer:Owner"]
	require.Equal(t, RuleAmount, amount.Kind)
	require.Equal(t, uint64(2), amount.Amount)
	require.Equal(t, CompareOpGtEq, amount.Operator)

	require.Equal(t, RuleNamespace, rs.Operations["Transfer:Base"].Kind)

	exotic := rs.Operations["Exotic"]
	require.Equal(t, RuleUnknown, exotic.Kind)
	require.Equal(t, "SomeFutureRule", exotic.RawKind)
	require.Equal(t, Indeterminate, EvaluateRule(exotic, EvalContext{}).Verdict)
}
