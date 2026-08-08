package mpltokenauthrules

import (
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
)

// V2 binary constraint discriminants (u32). Uninitialized (0) is a zeroed-memory
// sentinel that never appears as a valid on-disk constraint.
const (
	v2AdditionalSigner uint32 = 1
	v2All              uint32 = 2
	v2Amount           uint32 = 3
	v2Any              uint32 = 4
	v2Frequency        uint32 = 5
	v2IsWallet         uint32 = 6
	v2Namespace        uint32 = 7
	v2Not              uint32 = 8
	v2Pass             uint32 = 9
	v2PDAMatch         uint32 = 10
	v2ProgramOwned     uint32 = 11
	v2ProgramOwnedList uint32 = 12
	v2ProgramOwnedTree uint32 = 13
	v2PubkeyListMatch  uint32 = 14
	v2PubkeyMatch      uint32 = 15
	v2PubkeyTreeMatch  uint32 = 16
)

const (
	v2HeaderSize = 8  // RuleV2 TLV header: [constraint_type u32, length u32]
	v2Str32Size  = 32 // Str32 fixed field width
	v2PubkeySize = 32
	v2RuleSetHdr = 72 // [u32;2] header + owner(32) + rule_set_name Str32(32)
)

// maxV2RuleDepth bounds how deeply nested V2 rules (All/Any/Not) may be. A
// legitimate rule tree nests only a handful of levels; this ceiling prevents a
// maliciously deep TLV chain from exhausting the Go stack during parsing (and,
// since the parsed tree is no deeper than this, during evaluation as well),
// which would be an unrecoverable fatal error rather than a catchable panic.
const maxV2RuleDepth = 256

// parseRuleSetV2 decodes a V2 (binary) rule-set payload. The input is the full
// revision slice; unlike V1, the lib-version marker is embedded as the low byte
// of the first u32 header word and is not stripped before this call.
func parseRuleSetV2(payload []byte) (*RuleSet, error) {
	if len(payload) < v2RuleSetHdr {
		return nil, fmt.Errorf("v2 rule set: truncated header, need %d bytes, have %d", v2RuleSetHdr, len(payload))
	}
	libVersion := payload[0]
	size := le32(payload[4:8])

	owner := solanago.PublicKeyFromBytes(payload[8:40])
	name := decodeStr32(payload[40:72])

	// Each operation is a Str32; the rules follow the operations block.
	opsEnd := v2RuleSetHdr + int(size)*v2Str32Size
	if opsEnd < v2RuleSetHdr || opsEnd > len(payload) {
		return nil, fmt.Errorf("v2 rule set: operations block of %d entries overruns %d bytes", size, len(payload))
	}

	rs := &RuleSet{
		LibVersion: libVersion,
		Owner:      owner,
		Name:       name,
		Operations: make(map[string]Rule, size),
	}

	names := make([]string, 0, capHint(int(size)))
	for i := 0; i < int(size); i++ {
		off := v2RuleSetHdr + i*v2Str32Size
		names = append(names, decodeStr32(payload[off:off+v2Str32Size]))
	}

	cursor := opsEnd
	for i := 0; i < int(size); i++ {
		rule, consumed, err := parseRuleV2(payload[cursor:], 0)
		if err != nil {
			return nil, fmt.Errorf("v2 rule set: operation %d (%q): %w", i, names[i], err)
		}
		cursor += consumed
		opName := names[i]
		if _, exists := rs.Operations[opName]; !exists {
			rs.OperationNames = append(rs.OperationNames, opName)
		}
		rs.Operations[opName] = rule
	}
	return rs, nil
}

// parseRuleV2 decodes one RuleV2 TLV record and returns the decoded rule plus
// the total number of bytes consumed (header + body). An unrecognized
// constraint discriminant is preserved as an Unknown rule that still consumes
// its declared length, keeping the surrounding rule set parseable; a recognized
// constraint whose body length is wrong is an error.
func parseRuleV2(data []byte, depth int) (Rule, int, error) {
	if depth > maxV2RuleDepth {
		return Rule{}, 0, fmt.Errorf("rule nesting exceeds maximum depth %d", maxV2RuleDepth)
	}
	if len(data) < v2HeaderSize {
		return Rule{}, 0, fmt.Errorf("rule header truncated: need %d bytes, have %d", v2HeaderSize, len(data))
	}
	constraintType := le32(data[0:4])
	length := le32(data[4:8])

	total := v2HeaderSize + int(length)
	if total < v2HeaderSize || total > len(data) {
		return Rule{}, 0, fmt.Errorf("rule body of %d bytes overruns %d available", length, len(data)-v2HeaderSize)
	}
	body := data[v2HeaderSize:total]

	rule, err := decodeV2Body(constraintType, body, depth)
	if err != nil {
		return Rule{}, 0, err
	}
	return rule, total, nil
}

func decodeV2Body(constraintType uint32, body []byte, depth int) (Rule, error) {
	switch constraintType {
	case v2AdditionalSigner:
		if err := expectV2Len("AdditionalSigner", body, 32); err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleAdditionalSigner, PublicKey: solanago.PublicKeyFromBytes(body[0:32])}, nil

	case v2All, v2Any:
		children, err := decodeV2Composite(body, depth)
		if err != nil {
			return Rule{}, err
		}
		kind := RuleAll
		if constraintType == v2Any {
			kind = RuleAny
		}
		return Rule{Kind: kind, Rules: children}, nil

	case v2Not:
		child, consumed, err := parseRuleV2(body, depth+1)
		if err != nil {
			return Rule{}, fmt.Errorf("Not nested rule: %w", err)
		}
		if consumed != len(body) {
			return Rule{}, fmt.Errorf("Not body has %d trailing bytes", len(body)-consumed)
		}
		return Rule{Kind: RuleNot, Rule: &child}, nil

	case v2Amount:
		if err := expectV2Len("Amount", body, 48); err != nil {
			return Rule{}, err
		}
		amount := le64(body[0:8])
		operatorRaw := le64(body[8:16])
		op, err := compareOpFromU64(operatorRaw)
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleAmount, Amount: amount, Operator: op, Field: decodeStr32(body[16:48])}, nil

	case v2Frequency:
		if err := expectV2Len("Frequency", body, 32); err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleFrequency, PublicKey: solanago.PublicKeyFromBytes(body[0:32])}, nil

	case v2IsWallet:
		if err := expectV2Len("IsWallet", body, 32); err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleIsWallet, Field: decodeStr32(body[0:32])}, nil

	case v2Namespace:
		if err := expectV2Len("Namespace", body, 0); err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleNamespace}, nil

	case v2Pass:
		if err := expectV2Len("Pass", body, 0); err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RulePass}, nil

	case v2PDAMatch:
		if err := expectV2Len("PDAMatch", body, 96); err != nil {
			return Rule{}, err
		}
		var program *solanago.PublicKey
		if !isZeroPubkey(body[0:32]) {
			key := solanago.PublicKeyFromBytes(body[0:32])
			program = &key
		}
		return Rule{
			Kind:    RulePDAMatch,
			Program: program,
			Field:   decodeStr32(body[32:64]),
			Field2:  decodeStr32(body[64:96]),
		}, nil

	case v2ProgramOwned:
		if err := expectV2Len("ProgramOwned", body, 64); err != nil {
			return Rule{}, err
		}
		return Rule{
			Kind:      RuleProgramOwned,
			PublicKey: solanago.PublicKeyFromBytes(body[0:32]),
			Field:     decodeStr32(body[32:64]),
		}, nil

	case v2ProgramOwnedList:
		field, keys, err := decodeV2FieldThenPubkeys("ProgramOwnedList", body)
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleProgramOwnedList, Field: field, PublicKeys: keys}, nil

	case v2PubkeyListMatch:
		field, keys, err := decodeV2FieldThenPubkeys("PubkeyListMatch", body)
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RulePubkeyListMatch, Field: field, PublicKeys: keys}, nil

	case v2PubkeyMatch:
		if err := expectV2Len("PubkeyMatch", body, 64); err != nil {
			return Rule{}, err
		}
		return Rule{
			Kind:      RulePubkeyMatch,
			PublicKey: solanago.PublicKeyFromBytes(body[0:32]),
			Field:     decodeStr32(body[32:64]),
		}, nil

	case v2PubkeyTreeMatch, v2ProgramOwnedTree:
		if err := expectV2Len("TreeMatch", body, 96); err != nil {
			return Rule{}, err
		}
		var root [32]byte
		copy(root[:], body[64:96])
		kind := RulePubkeyTreeMatch
		if constraintType == v2ProgramOwnedTree {
			kind = RuleProgramOwnedTree
		}
		return Rule{
			Kind:   kind,
			Field:  decodeStr32(body[0:32]),
			Field2: decodeStr32(body[32:64]),
			Root:   root,
		}, nil

	default:
		raw := make([]byte, len(body))
		copy(raw, body)
		return Rule{
			Kind:    RuleUnknown,
			RawKind: fmt.Sprintf("ConstraintType(%d)", constraintType),
			Raw:     raw,
		}, nil
	}
}

// decodeV2Composite decodes an All/Any body: a u64 element count followed by
// that many nested RuleV2 records.
func decodeV2Composite(body []byte, depth int) ([]Rule, error) {
	if len(body) < 8 {
		return nil, fmt.Errorf("composite body truncated: need at least 8 bytes, have %d", len(body))
	}
	size := le64(body[0:8])
	children := make([]Rule, 0, capHint(clampToInt(size)))
	offset := 8
	for i := range size {
		child, consumed, err := parseRuleV2(body[offset:], depth+1)
		if err != nil {
			return nil, fmt.Errorf("nested rule %d: %w", i, err)
		}
		children = append(children, child)
		offset += consumed
	}
	return children, nil
}

// decodeV2FieldThenPubkeys decodes the shared layout of ProgramOwnedList and
// PubkeyListMatch: a Str32 field followed by a raw array of pubkeys whose count
// is implied by the remaining body length.
func decodeV2FieldThenPubkeys(variant string, body []byte) (string, []solanago.PublicKey, error) {
	if len(body) < v2Str32Size {
		return "", nil, fmt.Errorf("%s body truncated: need at least %d bytes, have %d", variant, v2Str32Size, len(body))
	}
	field := decodeStr32(body[0:v2Str32Size])
	rest := body[v2Str32Size:]
	if len(rest)%v2PubkeySize != 0 {
		return "", nil, fmt.Errorf("%s pubkey list of %d bytes is not a multiple of %d", variant, len(rest), v2PubkeySize)
	}
	n := len(rest) / v2PubkeySize
	keys := make([]solanago.PublicKey, 0, capHint(n))
	for i := range n {
		off := i * v2PubkeySize
		keys = append(keys, solanago.PublicKeyFromBytes(rest[off:off+v2PubkeySize]))
	}
	return field, keys, nil
}

func expectV2Len(variant string, body []byte, want int) error {
	if len(body) != want {
		return fmt.Errorf("%s body must be %d bytes, got %d", variant, want, len(body))
	}
	return nil
}

func compareOpFromU64(v uint64) (CompareOp, error) {
	switch v {
	case 0:
		return CompareOpLt, nil
	case 1:
		return CompareOpLtEq, nil
	case 2:
		return CompareOpEq, nil
	case 3:
		return CompareOpGtEq, nil
	case 4:
		return CompareOpGt, nil
	default:
		return 0, fmt.Errorf("unknown compare operator %d", v)
	}
}

// decodeStr32 reads a 32-byte NUL-padded UTF-8 field, truncated at the first
// NUL byte.
func decodeStr32(b []byte) string {
	end := len(b)
	for i, c := range b {
		if c == 0 {
			end = i
			break
		}
	}
	return string(b[:end])
}

func isZeroPubkey(b []byte) bool {
	for _, c := range b {
		if c != 0 {
			return false
		}
	}
	return true
}

// clampToInt caps a u64 element count to a non-negative int; the subsequent
// per-element bounds checks reject any count that would actually overrun the
// body, so this only prevents an overflow when computing a capacity hint.
func clampToInt(n uint64) int {
	const maxInt = int(^uint(0) >> 1)
	if n > uint64(maxInt) {
		return maxInt
	}
	return int(n)
}
