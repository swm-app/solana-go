package mpltokenauthrules

import (
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
)

// RuleKind identifies one node of the unified rule tree shared by V1 (msgpack)
// and V2 (binary) rule sets. The zero value is RuleUnknown so that an
// uninitialized Rule is treated as undecidable rather than as a known kind.
type RuleKind int

const (
	// RuleUnknown is a rule whose kind could not be recognized. It is
	// preserved with diagnostics in RawKind/Raw and always evaluates to
	// Indeterminate.
	RuleUnknown RuleKind = iota
	RuleAll
	RuleAny
	RuleNot
	RuleAdditionalSigner
	RulePubkeyMatch
	RulePubkeyListMatch
	RulePubkeyTreeMatch
	RulePDAMatch
	RuleProgramOwned
	RuleProgramOwnedList
	RuleProgramOwnedTree
	RuleAmount
	RuleFrequency
	RuleIsWallet
	RulePass
	RuleNamespace
)

var ruleKindNames = map[RuleKind]string{
	RuleUnknown:          "Unknown",
	RuleAll:              "All",
	RuleAny:              "Any",
	RuleNot:              "Not",
	RuleAdditionalSigner: "AdditionalSigner",
	RulePubkeyMatch:      "PubkeyMatch",
	RulePubkeyListMatch:  "PubkeyListMatch",
	RulePubkeyTreeMatch:  "PubkeyTreeMatch",
	RulePDAMatch:         "PDAMatch",
	RuleProgramOwned:     "ProgramOwned",
	RuleProgramOwnedList: "ProgramOwnedList",
	RuleProgramOwnedTree: "ProgramOwnedTree",
	RuleAmount:           "Amount",
	RuleFrequency:        "Frequency",
	RuleIsWallet:         "IsWallet",
	RulePass:             "Pass",
	RuleNamespace:        "Namespace",
}

func (k RuleKind) String() string {
	if s, ok := ruleKindNames[k]; ok {
		return s
	}
	return fmt.Sprintf("RuleKind(%d)", int(k))
}

// CompareOp mirrors the on-chain amount comparison operator. Its integer values
// match the V2 Operator encoding exactly (Lt=0, LtEq=1, Eq=2, GtEq=3, Gt=4).
type CompareOp uint8

const (
	CompareOpLt CompareOp = iota
	CompareOpLtEq
	CompareOpEq
	CompareOpGtEq
	CompareOpGt
)

func (o CompareOp) String() string {
	switch o {
	case CompareOpLt:
		return "Lt"
	case CompareOpLtEq:
		return "LtEq"
	case CompareOpEq:
		return "Eq"
	case CompareOpGtEq:
		return "GtEq"
	case CompareOpGt:
		return "Gt"
	default:
		return fmt.Sprintf("CompareOp(%d)", uint8(o))
	}
}

// Rule is one node of a rule tree. Only the fields relevant to Kind are
// populated; the rest keep their zero values. It carries no evaluation logic -
// it is a pure data model produced by the V1 and V2 parsers and consumed by the
// evaluator.
type Rule struct {
	Kind RuleKind

	// Composites.
	Rules []Rule // All, Any
	Rule  *Rule  // Not

	// Payload field name(s). Field is the primary field for every leaf that
	// reads the payload; Field2 is the secondary field for the two-field
	// leaves (PDAMatch seeds field, PubkeyTreeMatch/ProgramOwnedTree proof
	// field).
	Field  string
	Field2 string

	// Single pubkey: AdditionalSigner account, PubkeyMatch pubkey,
	// ProgramOwned program, Frequency authority.
	PublicKey solanago.PublicKey

	// Pubkey list: PubkeyListMatch pubkeys, ProgramOwnedList/Set programs.
	PublicKeys []solanago.PublicKey

	// PDAMatch derivation program. Nil means "use the account owner as the
	// program" (V1 None, or the V2 all-zero sentinel).
	Program *solanago.PublicKey

	// Merkle root for PubkeyTreeMatch / ProgramOwnedTree.
	Root [32]byte

	// Amount comparison.
	Amount   uint64
	Operator CompareOp

	// Diagnostics for RuleUnknown: RawKind names the unrecognized construct
	// (an enum variant name for V1, a "ConstraintType(n)" tag for V2) and Raw
	// holds the undecoded body bytes when available.
	RawKind string
	Raw     []byte
}
