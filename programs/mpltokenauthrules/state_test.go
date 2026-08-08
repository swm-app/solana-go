package mpltokenauthrules

import (
	"encoding/binary"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

// Golden values in this file are pinned against real mainnet rule-set PDA
// accounts captured 2026-07-17 (see testdata/README.md). They assert concrete
// decoded structures and evaluator verdicts, not merely absence of errors.

// TestFoundationFixture pins the "Metaplex Foundation Rule Set", the royalty
// enforcement rule set used by Mad Lads, y00ts, Okay Bears and others. It has 9
// V1 (msgpack) revisions whose rule algebra evolved from a namespaced royalty
// tree into an unconditional Pass.
func TestFoundationFixture(t *testing.T) {
	acct := loadFixtureAccount(t, fixtureFoundation)

	require.Equal(t, uint8(1), acct.Header.Key)
	require.Equal(t, uint64(18924), acct.Header.RevMapVersionLocation)
	require.Len(t, acct.Revisions, 9)

	owner := mustPubkey(t, "ELskdHjzTQ6F4bBibhk4iqy63gSPs8ELec9HbfAaSDJk")
	for i, rev := range acct.Revisions {
		require.Equalf(t, uint8(libVersionV1), rev.LibVersion, "revision %d lib version", i)
		require.NotNilf(t, rev.RuleSet, "revision %d rule set", i)
		require.Equal(t, owner, rev.RuleSet.Owner)
		require.Equal(t, "Metaplex Foundation Rule Set", rev.RuleSet.Name)
	}

	// Latest revision (rev 8): every operation resolves to Pass, so a transfer
	// is unconditionally allowed regardless of payload knowledge.
	latest := acct.Latest().RuleSet
	require.Equal(t, []string{
		"Transfer:WalletToWallet", "Transfer:Owner", "Transfer:MigrationDelegate",
		"Transfer:SaleDelegate", "Transfer:TransferDelegate", "Delegate:LockedTransfer",
		"Delegate:Update", "Delegate:Transfer", "Delegate:Utility", "Delegate:Staking",
		"Delegate:Authority", "Delegate:Collection", "Delegate:Use", "Delegate:Sale",
	}, operationNames(latest))
	require.Equal(t, RulePass, latest.Operations[OperationTransferOwner].Kind)
	for _, ctx := range []EvalContext{
		{},
		{Payload: Payload{}},
		{Payload: Payload{PayloadKeyAmount: NewNumberPayload(1)}},
	} {
		require.Equal(t, Allow, EvaluateOperation(latest, OperationTransferOwner, ctx).Verdict)
	}

	// Middle revision (rev 2): Transfer:Owner is a Namespace that redirects to
	// the "Transfer" operation, an All[Amount==1, Any[ProgramOwnedList]] royalty
	// tree.
	rev2 := acct.Revisions[2].RuleSet
	require.Equal(t, []string{
		"Transfer", "Transfer:WalletToWallet", "Transfer:Owner", "Transfer:MigrationDelegate",
		"Transfer:SaleDelegate", "Transfer:TransferDelegate", "Delegate", "Delegate:Update",
		"Delegate:Transfer", "Delegate:Utility", "Delegate:Staking", "Delegate:Authority",
		"Delegate:Collection", "Delegate:Use", "Delegate:Sale",
	}, operationNames(rev2))
	require.Equal(t, RuleNamespace, rev2.Operations[OperationTransferOwner].Kind)

	transfer := rev2.Operations["Transfer"]
	require.Equal(t, RuleAll, transfer.Kind)
	require.Len(t, transfer.Rules, 2)
	amount := transfer.Rules[0]
	require.Equal(t, RuleAmount, amount.Kind)
	require.Equal(t, PayloadKeyAmount, amount.Field)
	require.Equal(t, CompareOpEq, amount.Operator)
	require.Equal(t, uint64(1), amount.Amount)
	anyRule := transfer.Rules[1]
	require.Equal(t, RuleAny, anyRule.Kind)
	require.Len(t, anyRule.Rules, 1)
	pol := anyRule.Rules[0]
	require.Equal(t, RuleProgramOwnedList, pol.Kind)
	require.Equal(t, "Source|Destination|Authority", pol.Field)
	require.Len(t, pol.PublicKeys, 13)

	// Verdict matrix for Transfer:Owner via the Namespace redirect on rev 2.
	require.Equal(t, Indeterminate, EvaluateOperation(rev2, OperationTransferOwner, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateOperation(rev2, OperationTransferOwner, EvalContext{Payload: Payload{}}).Verdict)
	// Amount present and equal to 1 (Amount passes) but no owned-program account
	// in the payload -> the whole tree still certainly fails.
	require.Equal(t, Deny, EvaluateOperation(rev2, OperationTransferOwner, EvalContext{
		Payload: Payload{PayloadKeyAmount: NewNumberPayload(1)},
	}).Verdict)
	// Amount present but not equal to 1 (Amount fails) -> All short-circuits Deny.
	require.Equal(t, Deny, EvaluateOperation(rev2, OperationTransferOwner, EvalContext{
		Payload: Payload{PayloadKeyAmount: NewNumberPayload(2)},
	}).Verdict)
	// Amount passes and a Source account is owned by a listed program: the owner
	// matches but the on-chain non-empty-data check cannot be verified
	// statically, so the ProgramOwnedList (and the whole tree) is Indeterminate.
	source := solanago.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
	require.Equal(t, Indeterminate, EvaluateOperation(rev2, OperationTransferOwner, EvalContext{
		Payload: Payload{
			PayloadKeyAmount: NewNumberPayload(1),
			PayloadKeySource: NewPubkeyPayload(source),
		},
		AccountOwners: map[solanago.PublicKey]solanago.PublicKey{source: pol.PublicKeys[0]},
	}).Verdict)

	// Earliest revision (rev 0): Transfer:Owner is a Namespace, but its redirect
	// target "Transfer" does not exist (rev 0 uses "Transfer:Base"), so the
	// operation resolves to OperationNotFound => Deny.
	rev0 := acct.Revisions[0].RuleSet
	require.Contains(t, operationNames(rev0), "Transfer:Base")
	require.NotContains(t, operationNames(rev0), "Transfer")
	require.Equal(t, RuleNamespace, rev0.Operations[OperationTransferOwner].Kind)
	require.Equal(t, Deny, EvaluateOperation(rev0, OperationTransferOwner, EvalContext{}).Verdict)
}

// TestCompatibilityFixture pins the Metaplex "Compatibility Rule Set", a
// no-restrictions rule set: a single V1 revision whose every operation is Pass.
func TestCompatibilityFixture(t *testing.T) {
	acct := loadFixtureAccount(t, fixtureCompatibility)

	require.Equal(t, uint8(1), acct.Header.Key)
	require.Len(t, acct.Revisions, 1)
	rs := acct.Revisions[0].RuleSet
	require.NotNil(t, rs)
	require.Equal(t, uint8(libVersionV1), acct.Revisions[0].LibVersion)
	require.Equal(t, mustPubkey(t, "ELskdHjzTQ6F4bBibhk4iqy63gSPs8ELec9HbfAaSDJk"), rs.Owner)
	require.Equal(t, "Compatibility Rule Set", rs.Name)
	require.Len(t, operationNames(rs), 14)

	for _, op := range operationNames(rs) {
		require.Equal(t, RulePass, rs.Operations[op].Kind)
	}
	require.Equal(t, Allow, EvaluateOperation(rs, OperationTransferOwner, EvalContext{}).Verdict)
}

// TestElementerraFixture pins the Elementerra collection's custom rule set. Its
// first revision has zero operations and the latest maps everything to Pass.
func TestElementerraFixture(t *testing.T) {
	acct := loadFixtureAccount(t, fixtureElementerra)

	require.Len(t, acct.Revisions, 2)
	owner := mustPubkey(t, "3qnGpWSmtdVDm2iw8LZUQ7wHSKKyYrtnyNuXuv2Na8Nf")
	for i, rev := range acct.Revisions {
		require.Equalf(t, uint8(libVersionV1), rev.LibVersion, "revision %d", i)
		require.NotNil(t, rev.RuleSet)
		require.Equal(t, owner, rev.RuleSet.Owner)
		require.Equal(t, "elementerra", rev.RuleSet.Name)
	}
	// Revision 0 carries no operations; any operation is OperationNotFound.
	require.Empty(t, operationNames(acct.Revisions[0].RuleSet))
	require.Equal(t, Deny, EvaluateOperation(acct.Revisions[0].RuleSet, OperationTransferOwner, EvalContext{}).Verdict)

	latest := acct.Latest().RuleSet
	require.Len(t, operationNames(latest), 14)
	require.Equal(t, RulePass, latest.Operations[OperationTransferOwner].Kind)
	require.Equal(t, Allow, EvaluateOperation(latest, OperationTransferOwner, EvalContext{}).Verdict)
}

// TestRHOFixture pins the "RHO" rule set: 9 revisions, every one a V2 (binary
// TLV) revision. Its latest Transfer tree nests a Not over an empty
// ProgramOwnedList, exercising the V2 composite and Not decoders on real data.
func TestRHOFixture(t *testing.T) {
	acct := loadFixtureAccount(t, fixtureRHO)

	require.Equal(t, uint8(1), acct.Header.Key)
	require.Equal(t, uint64(8808), acct.Header.RevMapVersionLocation)
	require.Len(t, acct.Revisions, 9)

	owner := mustPubkey(t, "7qCcNjYqadaN1w2xqYqnPE4q7Jgoxsc7eSpKVxAq2VWh")
	for i, rev := range acct.Revisions {
		require.Equalf(t, uint8(libVersionV2), rev.LibVersion, "revision %d lib version", i)
		require.NotNilf(t, rev.RuleSet, "revision %d rule set", i)
		require.Equal(t, owner, rev.RuleSet.Owner)
		require.Equal(t, "RHO", rev.RuleSet.Name)
	}

	latest := acct.Latest().RuleSet
	require.Len(t, operationNames(latest), 16)
	require.Equal(t, "Transfer", operationNames(latest)[0])
	require.Equal(t, RulePass, latest.Operations["Transfer:WalletToWallet"].Kind)
	require.Equal(t, RuleNamespace, latest.Operations[OperationTransferOwner].Kind)

	transfer := latest.Operations["Transfer"]
	require.Equal(t, RuleAll, transfer.Kind)
	require.Len(t, transfer.Rules, 2)
	require.Equal(t, RuleAmount, transfer.Rules[0].Kind)
	require.Equal(t, uint64(1), transfer.Rules[0].Amount)
	require.Equal(t, CompareOpEq, transfer.Rules[0].Operator)
	inner := transfer.Rules[1]
	require.Equal(t, RuleAll, inner.Kind)
	require.Len(t, inner.Rules, 1)
	not := inner.Rules[0]
	require.Equal(t, RuleNot, not.Kind)
	require.NotNil(t, not.Rule)
	require.Equal(t, RuleProgramOwnedList, not.Rule.Kind)
	require.Equal(t, "Source|Destination|Authority", not.Rule.Field)
	require.Empty(t, not.Rule.PublicKeys)

	// Transfer:WalletToWallet is Pass => Allow even with no payload knowledge.
	require.Equal(t, Allow, EvaluateOperation(latest, "Transfer:WalletToWallet", EvalContext{}).Verdict)
	// Transfer:Owner redirects to the Transfer tree: unknown payload is
	// Indeterminate, a known-empty payload makes the Amount rule certainly fail.
	require.Equal(t, Indeterminate, EvaluateOperation(latest, OperationTransferOwner, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateOperation(latest, OperationTransferOwner, EvalContext{Payload: Payload{}}).Verdict)
}

// TestNoPermissionsFixture pins the smallest V2 fixture: a single revision with
// zero operations. It also proves the parser honours the revision-map offsets
// rather than assuming the first revision starts immediately after the 9-byte
// header - this account has 7 bytes of alignment padding, so revision 0 begins
// at byte offset 16.
func TestNoPermissionsFixture(t *testing.T) {
	data := loadFixtureBytes(t, fixtureNoPermissions)
	acct, err := ParseAccount(data)
	require.NoError(t, err)

	require.Len(t, acct.Revisions, 1)
	require.Equal(t, uint8(libVersionV2), acct.Revisions[0].LibVersion)
	rs := acct.Revisions[0].RuleSet
	require.NotNil(t, rs)
	require.Equal(t, mustPubkey(t, "CSiscYYM1N1k6MtuX39QP6a5mBHHxxzsnc53bZyB13jo"), rs.Owner)
	require.Equal(t, "no_permissions", rs.Name)
	require.Empty(t, operationNames(rs))
	require.Equal(t, Deny, EvaluateOperation(rs, OperationTransferOwner, EvalContext{}).Verdict)

	// The revision has no operations, so its raw payload is exactly the 72-byte
	// V2 rule-set header.
	require.Len(t, acct.Revisions[0].Raw, 72)

	// Read the on-disk revision map directly and confirm the alignment padding.
	loc := acct.Header.RevMapVersionLocation
	require.Equal(t, uint64(88), loc)
	require.Equal(t, uint8(revMapVersion), data[loc])
	require.Equal(t, uint32(1), binary.LittleEndian.Uint32(data[loc+1:loc+5]))
	firstOffset := binary.LittleEndian.Uint64(data[loc+5 : loc+13])
	require.Equal(t, uint64(16), firstOffset, "revision 0 starts after 7 bytes of alignment padding")
}

// TestMixedMeFixture pins a mainnet account whose two revisions use different
// on-wire formats: revision 0 is V1 (msgpack) and revision 1 is V2 (binary).
// Both decode and evaluate.
func TestMixedMeFixture(t *testing.T) {
	acct := loadFixtureAccount(t, fixtureMixedMe)

	require.Len(t, acct.Revisions, 2)
	require.Equal(t, uint8(libVersionV1), acct.Revisions[0].LibVersion)
	require.Equal(t, uint8(libVersionV2), acct.Revisions[1].LibVersion)

	owner := mustPubkey(t, "B1gSiVjGmpzaupQgffCtRdsgAChboVtBQbTeurMbBVxw")
	for _, rev := range acct.Revisions {
		require.NotNil(t, rev.RuleSet)
		require.Equal(t, owner, rev.RuleSet.Owner)
		require.Equal(t, "me", rev.RuleSet.Name)
	}

	// Latest (V2) revision: the Transfer royalty tree.
	latest := acct.Latest().RuleSet
	transfer := latest.Operations["Transfer"]
	require.Equal(t, RuleAll, transfer.Kind)
	require.Len(t, transfer.Rules, 2)
	require.Equal(t, RuleAmount, transfer.Rules[0].Kind)
	require.Equal(t, RuleAny, transfer.Rules[1].Kind)
	require.Len(t, transfer.Rules[1].Rules, 1)
	require.Equal(t, RuleProgramOwnedList, transfer.Rules[1].Rules[0].Kind)
	require.Len(t, transfer.Rules[1].Rules[0].PublicKeys, 2)

	require.Equal(t, Allow, EvaluateOperation(latest, "Transfer:WalletToWallet", EvalContext{}).Verdict)
	require.Equal(t, Indeterminate, EvaluateOperation(latest, OperationTransferOwner, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateOperation(latest, OperationTransferOwner, EvalContext{Payload: Payload{}}).Verdict)
}

// TestMixedAmigosFixture pins a mainnet account whose newest revision is V1 even
// though earlier revisions are V2 (revisions 0 and 1 are V2, revision 2 is V1),
// confirming the parser does not assume a single format per account.
func TestMixedAmigosFixture(t *testing.T) {
	acct := loadFixtureAccount(t, fixtureMixedAmigos)

	require.Len(t, acct.Revisions, 3)
	require.Equal(t, uint8(libVersionV2), acct.Revisions[0].LibVersion)
	require.Equal(t, uint8(libVersionV2), acct.Revisions[1].LibVersion)
	require.Equal(t, uint8(libVersionV1), acct.Revisions[2].LibVersion)

	latest := acct.Latest().RuleSet
	require.Equal(t, uint8(libVersionV1), acct.Latest().LibVersion)
	require.Equal(t, mustPubkey(t, "GK5oBTZuRaCdHbDE2XEWhLs8dXwmGLJL8yp3SRsReGjw"), latest.Owner)
	require.Equal(t, "Amigos Odyssey", latest.Name)

	// The V1 latest revision's Transfer:WalletToWallet is All[Amount, IsWallet,
	// IsWallet]; IsWallet always fails on-chain, so the operation is Deny even
	// with an unknown payload.
	wtw := latest.Operations["Transfer:WalletToWallet"]
	require.Equal(t, RuleAll, wtw.Kind)
	require.Len(t, wtw.Rules, 3)
	require.Equal(t, RuleAmount, wtw.Rules[0].Kind)
	require.Equal(t, RuleIsWallet, wtw.Rules[1].Kind)
	require.Equal(t, RuleIsWallet, wtw.Rules[2].Kind)
	require.Equal(t, Deny, EvaluateOperation(latest, "Transfer:WalletToWallet", EvalContext{}).Verdict)

	require.Equal(t, Indeterminate, EvaluateOperation(latest, OperationTransferOwner, EvalContext{}).Verdict)
	require.Equal(t, Deny, EvaluateOperation(latest, OperationTransferOwner, EvalContext{Payload: Payload{}}).Verdict)
}
