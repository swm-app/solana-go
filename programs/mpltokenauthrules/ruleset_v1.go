package mpltokenauthrules

import (
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
)

// parseRuleSetV1 decodes a V1 (msgpack) rule-set payload. The input is the
// revision bytes with the external 1-byte lib-version marker already stripped,
// i.e. exactly what rmp_serde::from_slice receives on-chain.
//
// RuleSetV1 is serialized with rmp_serde's default configuration, which encodes
// a struct as a plain array of its field values in declared order (field names
// dropped). The array is therefore [lib_version, owner, rule_set_name,
// operations].
func parseRuleSetV1(payload []byte) (*RuleSet, error) {
	r := newMsgpackReader(payload)
	top, err := r.readValue()
	if err != nil {
		return nil, fmt.Errorf("v1 rule set: %w", err)
	}
	if top.typ != mpArray || len(top.arr) != 4 {
		return nil, fmt.Errorf("v1 rule set: expected 4-element struct array, got %s", describeValue(top))
	}

	libVersion, ok := top.arr[0].asUint()
	if !ok {
		return nil, fmt.Errorf("v1 rule set: lib_version is not an integer")
	}

	ownerBytes, ok := top.arr[1].asPubkey()
	if !ok {
		return nil, fmt.Errorf("v1 rule set: owner is not a 32-byte pubkey")
	}

	if top.arr[2].typ != mpStr {
		return nil, fmt.Errorf("v1 rule set: rule_set_name is not a string")
	}
	name := top.arr[2].str

	opsVal := top.arr[3]
	if opsVal.typ != mpMap {
		return nil, fmt.Errorf("v1 rule set: operations is not a map")
	}

	rs := &RuleSet{
		LibVersion: uint8(libVersion),
		Owner:      solanago.PublicKeyFromBytes(ownerBytes[:]),
		Name:       name,
		Operations: make(map[string]Rule, len(opsVal.pair)),
	}
	for _, kv := range opsVal.pair {
		if kv.key.typ != mpStr {
			return nil, fmt.Errorf("v1 rule set: operation key is not a string")
		}
		rule, err := parseRuleV1(kv.val)
		if err != nil {
			return nil, fmt.Errorf("v1 rule set: operation %q: %w", kv.key.str, err)
		}
		if _, exists := rs.Operations[kv.key.str]; !exists {
			rs.OperationNames = append(rs.OperationNames, kv.key.str)
		}
		rs.Operations[kv.key.str] = rule
	}
	return rs, nil
}

// parseRuleV1 decodes one msgpack-encoded Rule. Unit variants (Pass, Namespace)
// are a bare string; every other variant is a single-entry map whose key is the
// variant name and whose value is an array of the variant's field values. An
// unrecognized variant name is preserved as an Unknown rule; a recognized
// variant with the wrong body shape is an error, since the format is fixed.
func parseRuleV1(v mpValue) (Rule, error) {
	switch v.typ {
	case mpStr:
		switch v.str {
		case "Pass":
			return Rule{Kind: RulePass}, nil
		case "Namespace":
			return Rule{Kind: RuleNamespace}, nil
		default:
			return Rule{Kind: RuleUnknown, RawKind: v.str}, nil
		}
	case mpMap:
		if len(v.pair) != 1 {
			return Rule{}, fmt.Errorf("rule map must have exactly one entry, got %d", len(v.pair))
		}
		entry := v.pair[0]
		if entry.key.typ != mpStr {
			return Rule{}, fmt.Errorf("rule variant key is not a string")
		}
		name := entry.key.str
		if entry.val.typ != mpArray {
			return Rule{}, fmt.Errorf("rule variant %q body is not an array", name)
		}
		return parseRuleV1Variant(name, entry.val.arr)
	default:
		return Rule{}, fmt.Errorf("rule is neither a string nor a map: %s", describeValue(v))
	}
}

func parseRuleV1Variant(name string, args []mpValue) (Rule, error) {
	expectArgs := func(n int) error {
		if len(args) != n {
			return fmt.Errorf("rule %q expects %d field(s), got %d", name, n, len(args))
		}
		return nil
	}

	switch name {
	case "All", "Any":
		if err := expectArgs(1); err != nil {
			return Rule{}, err
		}
		if args[0].typ != mpArray {
			return Rule{}, fmt.Errorf("rule %q rules field is not an array", name)
		}
		children := make([]Rule, 0, capHint(len(args[0].arr)))
		for i, c := range args[0].arr {
			child, err := parseRuleV1(c)
			if err != nil {
				return Rule{}, fmt.Errorf("rule %q child %d: %w", name, i, err)
			}
			children = append(children, child)
		}
		kind := RuleAll
		if name == "Any" {
			kind = RuleAny
		}
		return Rule{Kind: kind, Rules: children}, nil

	case "Not":
		if err := expectArgs(1); err != nil {
			return Rule{}, err
		}
		child, err := parseRuleV1(args[0])
		if err != nil {
			return Rule{}, fmt.Errorf("rule Not child: %w", err)
		}
		return Rule{Kind: RuleNot, Rule: &child}, nil

	case "AdditionalSigner":
		if err := expectArgs(1); err != nil {
			return Rule{}, err
		}
		pk, err := readV1Pubkey(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleAdditionalSigner, PublicKey: pk}, nil

	case "PubkeyMatch":
		if err := expectArgs(2); err != nil {
			return Rule{}, err
		}
		pk, err := readV1Pubkey(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		field, err := readV1String(name, args[1])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RulePubkeyMatch, PublicKey: pk, Field: field}, nil

	case "PubkeyListMatch":
		if err := expectArgs(2); err != nil {
			return Rule{}, err
		}
		list, err := readV1PubkeyList(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		field, err := readV1String(name, args[1])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RulePubkeyListMatch, PublicKeys: list, Field: field}, nil

	case "PubkeyTreeMatch", "ProgramOwnedTree":
		if err := expectArgs(3); err != nil {
			return Rule{}, err
		}
		root, err := readV1Pubkey(name, args[0]) // [u8;32], same array-of-32 encoding
		if err != nil {
			return Rule{}, err
		}
		pubkeyField, err := readV1String(name, args[1])
		if err != nil {
			return Rule{}, err
		}
		proofField, err := readV1String(name, args[2])
		if err != nil {
			return Rule{}, err
		}
		kind := RulePubkeyTreeMatch
		if name == "ProgramOwnedTree" {
			kind = RuleProgramOwnedTree
		}
		return Rule{Kind: kind, Root: root, Field: pubkeyField, Field2: proofField}, nil

	case "PDAMatch":
		if err := expectArgs(3); err != nil {
			return Rule{}, err
		}
		var program *solanago.PublicKey
		if args[0].typ != mpNil {
			pk, err := readV1Pubkey(name, args[0])
			if err != nil {
				return Rule{}, err
			}
			key := solanago.PublicKeyFromBytes(pk[:])
			program = &key
		}
		pdaField, err := readV1String(name, args[1])
		if err != nil {
			return Rule{}, err
		}
		seedsField, err := readV1String(name, args[2])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RulePDAMatch, Program: program, Field: pdaField, Field2: seedsField}, nil

	case "ProgramOwned":
		if err := expectArgs(2); err != nil {
			return Rule{}, err
		}
		pk, err := readV1Pubkey(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		field, err := readV1String(name, args[1])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleProgramOwned, PublicKey: pk, Field: field}, nil

	case "ProgramOwnedList", "ProgramOwnedSet":
		if err := expectArgs(2); err != nil {
			return Rule{}, err
		}
		list, err := readV1PubkeyList(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		field, err := readV1String(name, args[1])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleProgramOwnedList, PublicKeys: list, Field: field}, nil

	case "Amount":
		if err := expectArgs(3); err != nil {
			return Rule{}, err
		}
		amount, ok := args[0].asUint()
		if !ok {
			return Rule{}, fmt.Errorf("rule Amount amount is not an unsigned integer")
		}
		if args[1].typ != mpStr {
			return Rule{}, fmt.Errorf("rule Amount operator is not a string")
		}
		op, err := compareOpFromString(args[1].str)
		if err != nil {
			return Rule{}, err
		}
		field, err := readV1String(name, args[2])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleAmount, Amount: amount, Operator: op, Field: field}, nil

	case "Frequency":
		if err := expectArgs(1); err != nil {
			return Rule{}, err
		}
		pk, err := readV1Pubkey(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleFrequency, PublicKey: pk}, nil

	case "IsWallet":
		if err := expectArgs(1); err != nil {
			return Rule{}, err
		}
		field, err := readV1String(name, args[0])
		if err != nil {
			return Rule{}, err
		}
		return Rule{Kind: RuleIsWallet, Field: field}, nil

	default:
		return Rule{Kind: RuleUnknown, RawKind: name}, nil
	}
}

func readV1Pubkey(variant string, v mpValue) (solanago.PublicKey, error) {
	raw, ok := v.asPubkey()
	if !ok {
		return solanago.PublicKey{}, fmt.Errorf("rule %q: expected 32-byte pubkey, got %s", variant, describeValue(v))
	}
	return solanago.PublicKeyFromBytes(raw[:]), nil
}

func readV1PubkeyList(variant string, v mpValue) ([]solanago.PublicKey, error) {
	if v.typ != mpArray {
		return nil, fmt.Errorf("rule %q: expected pubkey array, got %s", variant, describeValue(v))
	}
	out := make([]solanago.PublicKey, 0, capHint(len(v.arr)))
	for i, e := range v.arr {
		raw, ok := e.asPubkey()
		if !ok {
			return nil, fmt.Errorf("rule %q: element %d is not a 32-byte pubkey", variant, i)
		}
		out = append(out, solanago.PublicKeyFromBytes(raw[:]))
	}
	return out, nil
}

func readV1String(variant string, v mpValue) (string, error) {
	if v.typ != mpStr {
		return "", fmt.Errorf("rule %q: expected string field, got %s", variant, describeValue(v))
	}
	return v.str, nil
}

func compareOpFromString(s string) (CompareOp, error) {
	switch s {
	case "Lt":
		return CompareOpLt, nil
	case "LtEq":
		return CompareOpLtEq, nil
	case "Eq":
		return CompareOpEq, nil
	case "GtEq":
		return CompareOpGtEq, nil
	case "Gt":
		return CompareOpGt, nil
	default:
		return 0, fmt.Errorf("unknown compare operator %q", s)
	}
}

func describeValue(v mpValue) string {
	switch v.typ {
	case mpNil:
		return "nil"
	case mpBool:
		return "bool"
	case mpUint:
		return "uint"
	case mpInt:
		return "int"
	case mpStr:
		return "string"
	case mpBin:
		return "bin"
	case mpArray:
		return fmt.Sprintf("array(%d)", len(v.arr))
	case mpMap:
		return fmt.Sprintf("map(%d)", len(v.pair))
	default:
		return "unknown"
	}
}
