package mpltokenauthrules

import (
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

// leaf verdict shorthands for composing hand-built fold tables. RulePass is a
// certain Allow; RuleFrequency is a certain Deny; RuleUnknown is Indeterminate.
var (
	leafAllow = Rule{Kind: RulePass}
	leafDeny  = Rule{Kind: RuleFrequency}
	leafIndet = Rule{Kind: RuleUnknown, RawKind: "synthetic"}
)

func TestVerdictZeroValueIsIndeterminate(t *testing.T) {
	var v Verdict
	require.Equal(t, Indeterminate, v)
	require.Equal(t, "Indeterminate", v.String())
	// A defaulted Result likewise reads as Indeterminate, never Allow.
	require.Equal(t, Indeterminate, Result{}.Verdict)
}

func TestEvalLeafKinds(t *testing.T) {
	require.Equal(t, Allow, EvaluateRule(leafAllow, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateRule(leafDeny, EvalContext{}).Verdict)
	require.Equal(t, Indeterminate, EvaluateRule(leafIndet, EvalContext{}).Verdict)
}

// TestEvalAllFold covers the AND fold, including the empty-All success that
// matches the on-chain All semantics.
func TestEvalAllFold(t *testing.T) {
	tests := []struct {
		name     string
		children []Rule
		want     Verdict
	}{
		{"empty all allows", nil, Allow},
		{"all allow", []Rule{leafAllow, leafAllow}, Allow},
		{"any deny denies", []Rule{leafAllow, leafDeny}, Deny},
		{"deny dominates indet", []Rule{leafIndet, leafDeny}, Deny},
		{"indet without deny stays indet", []Rule{leafAllow, leafIndet}, Indeterminate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rule{Kind: RuleAll, Rules: tt.children}
			require.Equal(t, tt.want, EvaluateRule(r, EvalContext{}).Verdict)
		})
	}
}

// TestEvalAnyFold covers the OR fold, including the empty-Any failure that
// matches the on-chain Any semantics.
func TestEvalAnyFold(t *testing.T) {
	tests := []struct {
		name     string
		children []Rule
		want     Verdict
	}{
		{"empty any denies", nil, Deny},
		{"any allow allows", []Rule{leafDeny, leafAllow}, Allow},
		{"allow dominates indet", []Rule{leafIndet, leafAllow}, Allow},
		{"indet without allow stays indet", []Rule{leafDeny, leafIndet}, Indeterminate},
		{"all deny denies", []Rule{leafDeny, leafDeny}, Deny},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rule{Kind: RuleAny, Rules: tt.children}
			require.Equal(t, tt.want, EvaluateRule(r, EvalContext{}).Verdict)
		})
	}
}

// TestEvalNotFold pins the three-valued Not semantics. On-chain Not maps
// Success->Failure and Failure->Success but leaves an aborting Error as Error;
// the three-valued verdict collapses Failure and Error into Deny, so Not(Deny)
// cannot soundly become Allow and degrades to Indeterminate. This matches the
// Rust Not::validate in programs/token-auth-rules/src/state/v2/constraint/not.rs.
func TestEvalNotFold(t *testing.T) {
	tests := []struct {
		name  string
		child Rule
		want  Verdict
	}{
		{"not allow is deny", leafAllow, Deny},
		{"not deny is indeterminate", leafDeny, Indeterminate},
		{"not indeterminate is indeterminate", leafIndet, Indeterminate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			child := tt.child
			r := Rule{Kind: RuleNot, Rule: &child}
			require.Equal(t, tt.want, EvaluateRule(r, EvalContext{}).Verdict)
		})
	}
	// A Not with no child cannot be decided.
	require.Equal(t, Indeterminate, EvaluateRule(Rule{Kind: RuleNot}, EvalContext{}).Verdict)
}

// TestEvalAmount checks every comparison operator in both directions plus the
// missing / non-numeric / unknown-payload degradations.
func TestEvalAmount(t *testing.T) {
	const rule = uint64(5)
	ops := []struct {
		op        CompareOp
		less      Verdict // payload 4 (< 5)
		equal     Verdict // payload 5 (== 5)
		greater   Verdict // payload 6 (> 5)
		opDisplay string
	}{
		{CompareOpLt, Allow, Deny, Deny, "Lt"},
		{CompareOpLtEq, Allow, Allow, Deny, "LtEq"},
		{CompareOpEq, Deny, Allow, Deny, "Eq"},
		{CompareOpGtEq, Deny, Allow, Allow, "GtEq"},
		{CompareOpGt, Deny, Deny, Allow, "Gt"},
	}
	for _, tc := range ops {
		t.Run(tc.opDisplay, func(t *testing.T) {
			require.Equal(t, tc.opDisplay, tc.op.String())
			r := Rule{Kind: RuleAmount, Field: PayloadKeyAmount, Operator: tc.op, Amount: rule}
			eval := func(n uint64) Verdict {
				return EvaluateRule(r, EvalContext{Payload: Payload{PayloadKeyAmount: NewNumberPayload(n)}}).Verdict
			}
			require.Equal(t, tc.less, eval(4))
			require.Equal(t, tc.equal, eval(5))
			require.Equal(t, tc.greater, eval(6))
		})
	}

	r := Rule{Kind: RuleAmount, Field: PayloadKeyAmount, Operator: CompareOpEq, Amount: 1}
	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{Payload: Payload{}}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload: Payload{PayloadKeyAmount: NewPubkeyPayload(solanago.PublicKey{})},
	}).Verdict)
}

func TestEvalAdditionalSigner(t *testing.T) {
	signer := solanago.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	r := Rule{Kind: RuleAdditionalSigner, PublicKey: signer}

	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{}).Verdict)
	require.Equal(t, Allow, EvaluateRule(r, EvalContext{
		Signers: map[solanago.PublicKey]bool{signer: true},
	}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Signers: map[solanago.PublicKey]bool{},
	}).Verdict)
}

func TestEvalPubkeyMatch(t *testing.T) {
	target := solanago.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	other := solanago.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
	r := Rule{Kind: RulePubkeyMatch, PublicKey: target, Field: PayloadKeyDestination}

	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{Payload: Payload{}}).Verdict)
	require.Equal(t, Allow, EvaluateRule(r, EvalContext{
		Payload: Payload{PayloadKeyDestination: NewPubkeyPayload(target)},
	}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload: Payload{PayloadKeyDestination: NewPubkeyPayload(other)},
	}).Verdict)
}

func TestEvalPubkeyListMatch(t *testing.T) {
	a := solanago.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	b := solanago.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
	other := solanago.MustPublicKeyFromBase58("Stake11111111111111111111111111111111111111")
	r := Rule{Kind: RulePubkeyListMatch, PublicKeys: []solanago.PublicKey{a, b}, Field: PayloadKeyDestination}

	require.Equal(t, Allow, EvaluateRule(r, EvalContext{
		Payload: Payload{PayloadKeyDestination: NewPubkeyPayload(b)},
	}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload: Payload{PayloadKeyDestination: NewPubkeyPayload(other)},
	}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{Payload: Payload{}}).Verdict)

	// A multi-field match is decoded identically for V1 and V2 in the unified
	// tree but validated differently on-chain, so it cannot be decided soundly.
	multi := Rule{Kind: RulePubkeyListMatch, PublicKeys: []solanago.PublicKey{a}, Field: "Source|Destination"}
	require.Equal(t, Indeterminate, EvaluateRule(multi, EvalContext{
		Payload: Payload{PayloadKeySource: NewPubkeyPayload(a)},
	}).Verdict)
}

func TestEvalProgramOwned(t *testing.T) {
	account := solanago.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	program := solanago.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
	wrong := solanago.MustPublicKeyFromBase58("Stake11111111111111111111111111111111111111")
	r := Rule{Kind: RuleProgramOwned, PublicKey: program, Field: PayloadKeySource}
	payload := Payload{PayloadKeySource: NewPubkeyPayload(account)}

	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{}).Verdict)
	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{Payload: payload}).Verdict)
	// Owner matches: on-chain also requires non-empty account data, unknowable
	// statically -> Indeterminate.
	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{
		Payload:       payload,
		AccountOwners: map[solanago.PublicKey]solanago.PublicKey{account: program},
	}).Verdict)
	// Owner mismatch is a certain failure.
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload:       payload,
		AccountOwners: map[solanago.PublicKey]solanago.PublicKey{account: wrong},
	}).Verdict)
}

func TestEvalProgramOwnedList(t *testing.T) {
	account := solanago.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	program := solanago.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
	wrong := solanago.MustPublicKeyFromBase58("Stake11111111111111111111111111111111111111")
	r := Rule{Kind: RuleProgramOwnedList, PublicKeys: []solanago.PublicKey{program}, Field: PayloadKeySource}
	payload := Payload{PayloadKeySource: NewPubkeyPayload(account)}

	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{
		Payload:       payload,
		AccountOwners: map[solanago.PublicKey]solanago.PublicKey{account: program},
	}).Verdict)
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload:       payload,
		AccountOwners: map[solanago.PublicKey]solanago.PublicKey{account: wrong},
	}).Verdict)
}

// TestEvalPDAMatch builds a real PDA with solana-go and checks the derivable /
// wrong / missing paths.
func TestEvalPDAMatch(t *testing.T) {
	program := solanago.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
	seeds := [][]byte{[]byte("hello"), []byte("world")}
	pda, _, err := solanago.FindProgramAddress(seeds, program)
	require.NoError(t, err)

	prog := program
	r := Rule{Kind: RulePDAMatch, Program: &prog, Field: PayloadKeyDestination, Field2: PayloadKeySource}

	// Correct account + seeds derive the PDA => Allow.
	require.Equal(t, Allow, EvaluateRule(r, EvalContext{
		Payload: Payload{
			PayloadKeyDestination: NewPubkeyPayload(pda),
			PayloadKeySource:      NewSeedsPayload(seeds),
		},
	}).Verdict)
	// Wrong account for the same seeds => Deny.
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload: Payload{
			PayloadKeyDestination: NewPubkeyPayload(program),
			PayloadKeySource:      NewSeedsPayload(seeds),
		},
	}).Verdict)
	// Missing seeds field => Deny (payload known, value absent).
	require.Equal(t, Deny, EvaluateRule(r, EvalContext{
		Payload: Payload{PayloadKeyDestination: NewPubkeyPayload(pda)},
	}).Verdict)
	// Unknown payload => Indeterminate.
	require.Equal(t, Indeterminate, EvaluateRule(r, EvalContext{}).Verdict)

	// Default (nil) program derives against the account owner, so the owner must
	// be known.
	def := Rule{Kind: RulePDAMatch, Field: PayloadKeyDestination, Field2: PayloadKeySource}
	require.Equal(t, Indeterminate, EvaluateRule(def, EvalContext{
		Payload: Payload{
			PayloadKeyDestination: NewPubkeyPayload(pda),
			PayloadKeySource:      NewSeedsPayload(seeds),
		},
	}).Verdict)
	require.Equal(t, Allow, EvaluateRule(def, EvalContext{
		Payload: Payload{
			PayloadKeyDestination: NewPubkeyPayload(pda),
			PayloadKeySource:      NewSeedsPayload(seeds),
		},
		AccountOwners: map[solanago.PublicKey]solanago.PublicKey{pda: program},
	}).Verdict)
}

// TestEvalMerkleAndUnimplemented pins the rules that never decide either way and
// the rules that always deny.
func TestEvalMerkleAndUnimplemented(t *testing.T) {
	payload := EvalContext{Payload: Payload{
		PayloadKeySource:      NewPubkeyPayload(solanago.PublicKey{}),
		PayloadKeyDestination: NewMerkleProofPayload([][32]byte{{1}}),
	}}
	for _, kind := range []RuleKind{RulePubkeyTreeMatch, RuleProgramOwnedTree} {
		r := Rule{Kind: kind, Field: PayloadKeySource, Field2: PayloadKeyDestination}
		require.Equalf(t, Indeterminate, EvaluateRule(r, payload).Verdict, "%s", kind)
	}
	require.Equal(t, Deny, EvaluateRule(Rule{Kind: RuleFrequency}, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateRule(Rule{Kind: RuleIsWallet, Field: PayloadKeySource}, EvalContext{}).Verdict)
	require.Equal(t, Indeterminate, EvaluateRule(Rule{Kind: RuleUnknown, RawKind: "Exotic"}, EvalContext{}).Verdict)
	// A Namespace validated directly (not resolved away) always fails on-chain.
	require.Equal(t, Deny, EvaluateRule(Rule{Kind: RuleNamespace}, EvalContext{}).Verdict)
}

// TestNamespaceRedirect covers the operation-resolution fallback on hand-built
// rule sets: a redirect to the prefix operation, a redirect whose prefix is
// missing, and a cyclic redirect chain that must terminate without looping.
func TestNamespaceRedirect(t *testing.T) {
	t.Run("redirect to prefix", func(t *testing.T) {
		rs := &RuleSet{Operations: map[string]Rule{
			"Transfer":       {Kind: RulePass},
			"Transfer:Owner": {Kind: RuleNamespace},
		}}
		require.Equal(t, Allow, EvaluateOperation(rs, "Transfer:Owner", EvalContext{}).Verdict)
	})

	t.Run("redirect with missing prefix denies", func(t *testing.T) {
		rs := &RuleSet{Operations: map[string]Rule{
			"Transfer:Owner": {Kind: RuleNamespace},
		}}
		// Transfer:Owner -> "Transfer" which is absent -> OperationNotFound.
		require.Equal(t, Deny, EvaluateOperation(rs, "Transfer:Owner", EvalContext{}).Verdict)
	})

	t.Run("redirect with no separator denies", func(t *testing.T) {
		rs := &RuleSet{Operations: map[string]Rule{
			"Transfer": {Kind: RuleNamespace},
		}}
		// "Transfer" is a Namespace with no ':' to split on -> not resolvable.
		require.Equal(t, Deny, EvaluateOperation(rs, "Transfer", EvalContext{}).Verdict)
	})

	t.Run("chained namespace redirects terminate", func(t *testing.T) {
		// Every operation on the chain is a Namespace. Resolution strips the
		// suffix after the first ':' at each hop ("A:B:C" -> "A"), so the chain
		// always shrinks and cannot loop; the seen-set guard in
		// lookupOperationRule additionally bounds it. It must terminate with a
		// non-Allow verdict rather than spin.
		rs := &RuleSet{Operations: map[string]Rule{
			"A:B:C": {Kind: RuleNamespace},
			"A":     {Kind: RuleNamespace},
		}}
		res := EvaluateOperation(rs, "A:B:C", EvalContext{})
		require.NotEqual(t, Allow, res.Verdict)
		require.Equal(t, Deny, res.Verdict)
	})
}

func TestEvaluateOperationNilRuleSet(t *testing.T) {
	require.Equal(t, Indeterminate, EvaluateOperation(nil, OperationTransferOwner, EvalContext{}).Verdict)
}
