// Package mpltokenauthrules is a hand-written reader and static evaluator for
// the Metaplex Token Auth Rules program account state
// (program id auth9SigNpDKz4sJJ1DfCTuZrZNSAgh9sFD3rboVmgg).
//
// This package never builds or sends instructions. It does two things:
//
//  1. Deserialization of a rule-set PDA account. A rule-set account is a
//     9-byte Borsh header, followed by any number of revisions laid out
//     back-to-back, followed by a revision map that records the byte offset of
//     every revision. Each revision is one of two on-wire formats: V1, an
//     rmp_serde (msgpack) encoding of a RuleSetV1 struct, or V2, a flat binary
//     TLV encoding. Both formats describe the same rule algebra, so they are
//     parsed into one unified Rule tree (see rule.go).
//
//  2. A three-valued static evaluator of an operation (primarily
//     "Transfer:Owner") over the parsed rule tree. Because a client cannot see
//     everything the on-chain program sees at validation time, the evaluator
//     returns one of three verdicts:
//
//     Allow  - the on-chain program would certainly pass this operation for
//     any transaction consistent with the provided EvalContext.
//     Deny   - the on-chain program would certainly fail it.
//     Indeterminate - it cannot be proven either way from the provided context.
//
// The zero value of Verdict is Indeterminate. The evaluator never silently
// answers Allow or Deny for something it cannot prove. Merkle-proof rules
// (PubkeyTreeMatch, ProgramOwnedTree), rules whose inputs are absent from the
// context, and unknown rule kinds all degrade to Indeterminate.
//
// A nil EvalContext.Payload means "the payload is unknown" and forces every
// payload-dependent rule to Indeterminate. A non-nil but empty Payload means
// "the payload is known and contains no keys"; a rule that reads a missing key
// then resolves to Deny, because the on-chain program errors on a missing
// payload value and that failure is certain once the payload is known. The
// same known-versus-unknown distinction applies to EvalContext.Signers and
// EvalContext.AccountOwners: a nil map is unknown, a non-nil map is
// authoritative.
//
// Parsing philosophy: an unknown but well-delimited construct (an unrecognized
// V1 enum variant name, an unrecognized V2 constraint discriminant, or a
// revision whose lib version is not V1 or V2) is preserved rather than
// rejected, so that a rule set carrying one exotic rule still parses and its
// other branches still evaluate. Unknown rules become explicit Unknown nodes
// that evaluate to Indeterminate; an unknown revision lib version keeps the raw
// revision bytes with a nil RuleSet. Only genuinely malformed input -
// truncation, out-of-bounds offsets, a known construct with a corrupt body -
// is a parse error, always reported with byte-offset context.
//
// Two on-chain rule kinds, Frequency and IsWallet, have no code path that
// succeeds in the current on-chain program; every evaluation of them fails.
// This package therefore evaluates them to Deny unconditionally.
package mpltokenauthrules
