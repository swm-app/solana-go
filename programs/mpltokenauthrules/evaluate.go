package mpltokenauthrules

import (
	"fmt"
	"strings"

	solanago "github.com/gagliardetto/solana-go"
)

// Verdict is the three-valued result of a static evaluation. The zero value is
// Indeterminate, so an unset or defaulted Verdict never reads as Allow.
type Verdict int

const (
	// Indeterminate: the outcome cannot be proven either way from the
	// provided context.
	Indeterminate Verdict = iota
	// Allow: the on-chain program would certainly pass the operation.
	Allow
	// Deny: the on-chain program would certainly fail the operation.
	Deny
)

func (v Verdict) String() string {
	switch v {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	default:
		return "Indeterminate"
	}
}

// Result is a verdict with a human-readable reason. For a composite rule the
// reason is that of the deciding child.
type Result struct {
	Verdict Verdict
	Reason  string
}

func allow(reason string) Result { return Result{Verdict: Allow, Reason: reason} }
func deny(reason string) Result  { return Result{Verdict: Deny, Reason: reason} }
func indet(reason string) Result { return Result{Verdict: Indeterminate, Reason: reason} }

// Transfer operation name strings, as Token Metadata sends them to the auth
// rules program. MigrationDelegate is included for completeness, but Token
// Metadata skips auth-rule validation for it entirely.
const (
	OperationTransferOwner             = "Transfer:Owner"
	OperationTransferTransferDelegate  = "Transfer:TransferDelegate"
	OperationTransferSaleDelegate      = "Transfer:SaleDelegate"
	OperationTransferMigrationDelegate = "Transfer:MigrationDelegate"
)

// Payload key strings actually populated for a transfer by Token Metadata.
const (
	PayloadKeyAmount      = "Amount"
	PayloadKeyAuthority   = "Authority"
	PayloadKeySource      = "Source"
	PayloadKeyDestination = "Destination"
)

// PayloadKind identifies which arm of a PayloadValue is set. It mirrors the
// on-chain PayloadType enum.
type PayloadKind int

const (
	PayloadPubkey PayloadKind = iota
	PayloadSeeds
	PayloadMerkleProof
	PayloadNumber
)

// PayloadValue is one entry of a Payload: exactly one arm is meaningful,
// selected by Kind.
type PayloadValue struct {
	Kind        PayloadKind
	Pubkey      solanago.PublicKey
	Seeds       [][]byte
	MerkleProof [][32]byte
	Number      uint64
}

// NewPubkeyPayload builds a Pubkey payload value.
func NewPubkeyPayload(pk solanago.PublicKey) PayloadValue {
	return PayloadValue{Kind: PayloadPubkey, Pubkey: pk}
}

// NewNumberPayload builds a Number payload value.
func NewNumberPayload(n uint64) PayloadValue {
	return PayloadValue{Kind: PayloadNumber, Number: n}
}

// NewSeedsPayload builds a Seeds payload value.
func NewSeedsPayload(seeds [][]byte) PayloadValue {
	return PayloadValue{Kind: PayloadSeeds, Seeds: seeds}
}

// NewMerkleProofPayload builds a MerkleProof payload value.
func NewMerkleProofPayload(proof [][32]byte) PayloadValue {
	return PayloadValue{Kind: PayloadMerkleProof, MerkleProof: proof}
}

func (v PayloadValue) asPubkey() (solanago.PublicKey, bool) {
	if v.Kind == PayloadPubkey {
		return v.Pubkey, true
	}
	return solanago.PublicKey{}, false
}

func (v PayloadValue) asNumber() (uint64, bool) {
	if v.Kind == PayloadNumber {
		return v.Number, true
	}
	return 0, false
}

func (v PayloadValue) asSeeds() ([][]byte, bool) {
	if v.Kind == PayloadSeeds {
		return v.Seeds, true
	}
	return nil, false
}

// Payload maps a payload key to its value. A nil Payload means the payload is
// unknown (every payload-dependent rule becomes Indeterminate); a non-nil but
// empty Payload means the payload is known to contain no keys.
type Payload map[string]PayloadValue

// EvalContext carries the optional knowledge available to the evaluator. Every
// field is optional. For each map, nil means "unknown" and a non-nil map is
// authoritative: a key's absence from a non-nil map is a known fact, not a gap.
type EvalContext struct {
	// Payload is the operation payload keyed by payload-key string.
	Payload Payload
	// Signers records which accounts are known to have signed. A nil map is
	// unknown; in a non-nil map, absence means the account did not sign.
	Signers map[solanago.PublicKey]bool
	// AccountOwners maps an account pubkey to its owner program. A nil map is
	// unknown; a missing entry in a non-nil map is treated as unknown for that
	// account and yields Indeterminate for rules that need it.
	AccountOwners map[solanago.PublicKey]solanago.PublicKey
}

// EvaluateOperation resolves an operation name to its top-level rule - applying
// the on-chain namespace-fallback semantics - and evaluates it. When the
// operation is not present in the rule set, the on-chain program fails with
// OperationNotFound, so this returns Deny; that failure is certain from the
// parsed rule set. A nil RuleSet yields Indeterminate.
func EvaluateOperation(rs *RuleSet, operation string, ctx EvalContext) Result {
	if rs == nil {
		return indet("rule set is unavailable (unknown lib version)")
	}
	rule, ok := lookupOperationRule(rs, operation)
	if !ok {
		return deny(fmt.Sprintf("operation %q not found in rule set", operation))
	}
	return EvaluateRule(*rule, ctx)
}

// lookupOperationRule implements the on-chain get_rule resolution: a missing
// key is not found; a key mapped to Namespace redirects to the prefix before
// the first ':' (and only fails to resolve when there is no such prefix).
func lookupOperationRule(rs *RuleSet, operation string) (*Rule, bool) {
	seen := make(map[string]bool)
	current := operation
	for {
		if seen[current] {
			return nil, false // guard against a cyclic Namespace chain
		}
		seen[current] = true

		rule, ok := rs.Operations[current]
		if !ok {
			return nil, false
		}
		if rule.Kind != RuleNamespace {
			r := rule
			return &r, true
		}
		idx := strings.IndexByte(current, ':')
		if idx < 0 {
			return nil, false
		}
		current = current[:idx]
	}
}

// EvaluateRule walks a rule tree and returns its three-valued verdict.
func EvaluateRule(r Rule, ctx EvalContext) Result {
	switch r.Kind {
	case RulePass:
		return allow("Pass")
	case RuleNamespace:
		// A Namespace reached during evaluation (rather than resolved away by
		// lookupOperationRule) always fails on-chain.
		return deny("Namespace rule validated directly always fails")
	case RuleAll:
		return evalAll(r.Rules, ctx)
	case RuleAny:
		return evalAny(r.Rules, ctx)
	case RuleNot:
		if r.Rule == nil {
			return indet("Not rule has no child")
		}
		return negate(EvaluateRule(*r.Rule, ctx))
	case RuleAdditionalSigner:
		return evalAdditionalSigner(r, ctx)
	case RulePubkeyMatch:
		return evalPubkeyMatch(r, ctx)
	case RulePubkeyListMatch:
		return evalPubkeyListMatch(r, ctx)
	case RuleAmount:
		return evalAmount(r, ctx)
	case RulePDAMatch:
		return evalPDAMatch(r, ctx)
	case RuleProgramOwned:
		return evalProgramOwned(r, ctx)
	case RuleProgramOwnedList:
		return evalProgramOwnedList(r, ctx)
	case RulePubkeyTreeMatch, RuleProgramOwnedTree:
		// Merkle-proof verification is deliberately not performed here.
		return indet(fmt.Sprintf("%s requires merkle-proof verification", r.Kind))
	case RuleFrequency:
		return deny("Frequency is not implemented on-chain and always fails")
	case RuleIsWallet:
		return deny("IsWallet is not implemented on-chain and always fails")
	case RuleUnknown:
		return indet(fmt.Sprintf("unknown rule kind %q", r.RawKind))
	default:
		return indet(fmt.Sprintf("unhandled rule kind %s", r.Kind))
	}
}

// evalAll folds children with AND semantics: any certain Deny makes the whole
// certainly Deny; otherwise any Indeterminate keeps it Indeterminate; an empty
// All passes (matching the on-chain empty-All success).
func evalAll(children []Rule, ctx EvalContext) Result {
	indetReason := ""
	sawIndet := false
	for _, c := range children {
		res := EvaluateRule(c, ctx)
		switch res.Verdict {
		case Deny:
			return res
		case Indeterminate:
			sawIndet = true
			indetReason = res.Reason
		}
	}
	if sawIndet {
		return indet(indetReason)
	}
	return allow("all sub-rules allow")
}

// evalAny folds children with OR semantics: any certain Allow makes the whole
// certainly Allow; otherwise any Indeterminate keeps it Indeterminate; an empty
// Any fails (matching the on-chain empty-Any failure).
func evalAny(children []Rule, ctx EvalContext) Result {
	indetReason := ""
	sawIndet := false
	for _, c := range children {
		res := EvaluateRule(c, ctx)
		switch res.Verdict {
		case Allow:
			return res
		case Indeterminate:
			sawIndet = true
			indetReason = res.Reason
		}
	}
	if sawIndet {
		return indet(indetReason)
	}
	return deny("no sub-rule allows")
}

// negate implements the on-chain Not semantics, which is not a plain boolean
// inversion. On-chain a child evaluates to one of Success, Failure, or Error,
// and Not maps Success->Failure, Failure->Success, but Error->Error: an
// aborting error (a missing payload value, a missing account, empty account
// data, an unimplemented rule) stays a failure through Not rather than becoming
// a pass. The three-valued verdict collapses both Failure and Error into Deny
// and cannot tell them apart, so a Deny child cannot be inverted to Allow
// without risking an unsound Allow for a child that actually aborts. A definite
// Allow is always an on-chain Success, so its negation is a definite Deny; a
// Deny degrades to Indeterminate.
func negate(res Result) Result {
	switch res.Verdict {
	case Allow:
		return deny("negation of an allowed rule")
	case Deny:
		return indet("negation of a denied rule is undecidable: the child may fail via an aborting error, which Not does not turn into a pass")
	default:
		return res
	}
}

func evalAdditionalSigner(r Rule, ctx EvalContext) Result {
	if ctx.Signers == nil {
		return indet("signer set is unknown")
	}
	if ctx.Signers[r.PublicKey] {
		return allow(fmt.Sprintf("%s signed", r.PublicKey))
	}
	return deny(fmt.Sprintf("required signer %s did not sign", r.PublicKey))
}

func evalPubkeyMatch(r Rule, ctx EvalContext) Result {
	if ctx.Payload == nil {
		return indet("payload is unknown")
	}
	pk, res, ok := payloadPubkey(r.Field, ctx)
	if !ok {
		return res
	}
	if pk == r.PublicKey {
		return allow(fmt.Sprintf("payload %q matches %s", r.Field, r.PublicKey))
	}
	return deny(fmt.Sprintf("payload %q %s does not match %s", r.Field, pk, r.PublicKey))
}

func evalPubkeyListMatch(r Rule, ctx EvalContext) Result {
	if ctx.Payload == nil {
		return indet("payload is unknown")
	}
	// A '|'-joined field name is interpreted differently by the two on-chain
	// formats: V2 splits the field and checks each payload pubkey against the
	// list directly, whereas V1 rewrites a multi-field PubkeyListMatch into a
	// ProgramOwnedList over the same keys, comparing account owners (not the
	// pubkeys) against the list and additionally requiring non-empty account
	// data. The unified rule tree does not record which format produced this
	// rule, so a multi-field match cannot be decided soundly and degrades to
	// Indeterminate. A single field is treated identically by both formats and
	// stays fully decidable: allow if the payload pubkey is in the list,
	// otherwise deny (a missing or non-pubkey field cannot make it pass).
	if strings.Contains(r.Field, "|") {
		return indet("multi-field PubkeyListMatch is version-dependent (V1 compares account owners, V2 compares pubkeys)")
	}
	v, present := ctx.Payload[r.Field]
	if !present {
		return deny(fmt.Sprintf("payload %q missing", r.Field))
	}
	pk, ok := v.asPubkey()
	if !ok {
		return deny(fmt.Sprintf("payload %q is not a pubkey", r.Field))
	}
	if containsPubkey(r.PublicKeys, pk) {
		return allow(fmt.Sprintf("payload %q %s is in the allow list", r.Field, pk))
	}
	return deny(fmt.Sprintf("payload %q %s is not in the allow list", r.Field, pk))
}

func evalAmount(r Rule, ctx EvalContext) Result {
	if ctx.Payload == nil {
		return indet("payload is unknown")
	}
	v, present := ctx.Payload[r.Field]
	if !present {
		return deny(fmt.Sprintf("payload %q missing", r.Field))
	}
	amount, ok := v.asNumber()
	if !ok {
		return deny(fmt.Sprintf("payload %q is not a number", r.Field))
	}
	if compareAmount(amount, r.Operator, r.Amount) {
		return allow(fmt.Sprintf("payload %q %d %s %d", r.Field, amount, r.Operator, r.Amount))
	}
	return deny(fmt.Sprintf("payload %q %d fails %s %d", r.Field, amount, r.Operator, r.Amount))
}

func evalPDAMatch(r Rule, ctx EvalContext) Result {
	if ctx.Payload == nil {
		return indet("payload is unknown")
	}
	account, res, ok := payloadPubkey(r.Field, ctx)
	if !ok {
		return res
	}
	seedsVal, present := ctx.Payload[r.Field2]
	if !present {
		return deny(fmt.Sprintf("payload seeds %q missing", r.Field2))
	}
	seeds, ok := seedsVal.asSeeds()
	if !ok {
		return deny(fmt.Sprintf("payload %q is not seeds", r.Field2))
	}

	var program solanago.PublicKey
	if r.Program != nil {
		program = *r.Program
	} else {
		// The default (all-zero) program means the derivation uses the
		// account's owner as the program id.
		if ctx.AccountOwners == nil {
			return indet("account owner unknown for default-program PDA derivation")
		}
		owner, known := ctx.AccountOwners[account]
		if !known {
			return indet(fmt.Sprintf("owner of %s unknown for PDA derivation", account))
		}
		program = owner
	}

	derived, _, err := solanago.FindProgramAddress(seeds, program)
	if err != nil {
		return deny(fmt.Sprintf("PDA derivation failed: %v", err))
	}
	if derived == account {
		return allow(fmt.Sprintf("payload %q %s is the derived PDA", r.Field, account))
	}
	return deny(fmt.Sprintf("payload %q %s is not the PDA derived for %s", r.Field, account, program))
}

func evalProgramOwned(r Rule, ctx EvalContext) Result {
	if ctx.Payload == nil {
		return indet("payload is unknown")
	}
	account, res, ok := payloadPubkey(r.Field, ctx)
	if !ok {
		return res
	}
	if ctx.AccountOwners == nil {
		return indet("account owners are unknown")
	}
	owner, known := ctx.AccountOwners[account]
	if !known {
		return indet(fmt.Sprintf("owner of %s is unknown", account))
	}
	if owner != r.PublicKey {
		return deny(fmt.Sprintf("%s owner %s is not %s", account, owner, r.PublicKey))
	}
	// Owner matches, but the on-chain rule also requires the account data to be
	// non-empty, which the static context cannot confirm.
	return indet(fmt.Sprintf("%s owner matches but account data emptiness is unknown", account))
}

func evalProgramOwnedList(r Rule, ctx EvalContext) Result {
	if ctx.Payload == nil {
		return indet("payload is unknown")
	}
	// A deny is only sound when every candidate field certainly fails. A field
	// whose owner is in the list, or whose owner is unknown, could still pass
	// (subject to the account-data check the static context cannot see), so it
	// forces Indeterminate.
	for _, field := range strings.Split(r.Field, "|") {
		v, present := ctx.Payload[field]
		if !present {
			continue
		}
		account, ok := v.asPubkey()
		if !ok {
			continue
		}
		if ctx.AccountOwners == nil {
			return indet("account owners are unknown")
		}
		owner, known := ctx.AccountOwners[account]
		if !known {
			return indet(fmt.Sprintf("owner of %s is unknown", account))
		}
		if containsPubkey(r.PublicKeys, owner) {
			return indet(fmt.Sprintf("%s owner is listed but account data emptiness is unknown", account))
		}
	}
	return deny("no payload account is owned by a listed program")
}

// payloadPubkey looks up a payload pubkey field under the known-payload
// assumption (caller has already verified ctx.Payload != nil). It returns the
// pubkey, or a decided Result explaining why the field is unusable.
func payloadPubkey(field string, ctx EvalContext) (solanago.PublicKey, Result, bool) {
	v, present := ctx.Payload[field]
	if !present {
		return solanago.PublicKey{}, deny(fmt.Sprintf("payload %q missing", field)), false
	}
	pk, ok := v.asPubkey()
	if !ok {
		return solanago.PublicKey{}, deny(fmt.Sprintf("payload %q is not a pubkey", field)), false
	}
	return pk, Result{}, true
}

func compareAmount(payloadAmount uint64, op CompareOp, ruleAmount uint64) bool {
	switch op {
	case CompareOpLt:
		return payloadAmount < ruleAmount
	case CompareOpLtEq:
		return payloadAmount <= ruleAmount
	case CompareOpEq:
		return payloadAmount == ruleAmount
	case CompareOpGtEq:
		return payloadAmount >= ruleAmount
	case CompareOpGt:
		return payloadAmount > ruleAmount
	default:
		return false
	}
}

func containsPubkey(list []solanago.PublicKey, target solanago.PublicKey) bool {
	for _, pk := range list {
		if pk == target {
			return true
		}
	}
	return false
}
