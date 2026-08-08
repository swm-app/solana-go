package mpltokenauthrules

import (
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
)

const (
	// ruleSetHeaderLen is the fixed Borsh header: 1-byte Key + 8-byte
	// rev_map_version_location.
	ruleSetHeaderLen = 9
	// revMapVersion is the only revision-map version that exists on-chain.
	revMapVersion = 1
	// libVersionV1 / libVersionV2 are the two known revision lib versions.
	libVersionV1 = 1
	libVersionV2 = 2
)

// Header is the fixed 9-byte Borsh header at the start of a rule-set PDA.
type Header struct {
	// Key is the account-type discriminant; a live rule set has Key == 1.
	Key uint8
	// RevMapVersionLocation is the absolute byte offset of the revision-map
	// version byte (the byte immediately preceding the serialized revision
	// map).
	RevMapVersionLocation uint64
}

// Revision is one on-chain rule-set revision. RuleSet is nil when the revision
// carries an unknown lib version; Raw always holds the exact revision byte
// slice (the lib-version byte plus the serialized payload).
type Revision struct {
	LibVersion uint8
	RuleSet    *RuleSet
	Raw        []byte
}

// RuleSet is a single decoded revision: its owner, name, and the map from
// operation name to top-level rule. OperationNames preserves the on-wire order
// of operations (deterministic for V2; msgpack map order for V1).
type RuleSet struct {
	LibVersion     uint8
	Owner          solanago.PublicKey
	Name           string
	Operations     map[string]Rule
	OperationNames []string
}

// RuleSetAccount is a fully parsed rule-set PDA account: its header plus every
// revision in on-chain order (index == revision number).
type RuleSetAccount struct {
	Header    Header
	Revisions []Revision
}

// Latest returns the most recent revision, which is the one Token Metadata uses
// when a token record does not pin a specific revision. It returns nil when the
// account has no revisions.
func (a *RuleSetAccount) Latest() *Revision {
	if a == nil || len(a.Revisions) == 0 {
		return nil
	}
	return &a.Revisions[len(a.Revisions)-1]
}

// ParseAccount decodes a complete rule-set PDA account: the header, the
// revision map, and every revision. A revision with an unknown lib version is
// preserved with a nil RuleSet rather than rejected. A structurally broken
// header, revision map, or known-version revision is an error with offset
// context.
func ParseAccount(data []byte) (*RuleSetAccount, error) {
	if len(data) < ruleSetHeaderLen {
		return nil, fmt.Errorf("rule set account: truncated header, need %d bytes, have %d", ruleSetHeaderLen, len(data))
	}

	header := Header{
		Key:                   data[0],
		RevMapVersionLocation: le64(data[1:9]),
	}

	loc := header.RevMapVersionLocation
	// The version byte must be readable and sit after the header.
	if loc < ruleSetHeaderLen || loc >= uint64(len(data)) {
		return nil, fmt.Errorf("rule set account: rev map version location %d out of range [%d, %d)", loc, ruleSetHeaderLen, len(data))
	}
	if got := data[loc]; got != revMapVersion {
		return nil, fmt.Errorf("rule set account: unsupported rev map version %d at offset %d", got, loc)
	}

	offsets, err := readRevisionMap(data, int(loc)+1)
	if err != nil {
		return nil, err
	}

	acct := &RuleSetAccount{Header: header}
	for i, start := range offsets {
		// The end of revision i is the start of revision i+1, or the revision
		// map version byte for the final revision.
		end := int(loc)
		if i+1 < len(offsets) {
			end = int(offsets[i+1])
		}
		if start < ruleSetHeaderLen || start > end || end > int(loc) {
			return nil, fmt.Errorf("rule set account: revision %d has invalid byte range [%d, %d)", i, start, end)
		}
		revBytes := data[start:end]
		rev, err := parseRevision(revBytes)
		if err != nil {
			return nil, fmt.Errorf("rule set account: revision %d: %w", i, err)
		}
		acct.Revisions = append(acct.Revisions, rev)
	}
	return acct, nil
}

// readRevisionMap decodes the Borsh Vec<usize> at start: a u32 little-endian
// count followed by that many u64 little-endian revision offsets.
func readRevisionMap(data []byte, start int) ([]int, error) {
	if start < 0 || start+4 > len(data) {
		return nil, fmt.Errorf("rule set account: revision map length prefix out of range at offset %d", start)
	}
	count := le32(data[start : start+4])
	pos := start + 4

	// Each entry is 8 bytes; reject a count that cannot fit in the remaining
	// input before allocating.
	remaining := len(data) - pos
	if int64(count)*8 > int64(remaining) {
		return nil, fmt.Errorf("rule set account: revision map declares %d entries but only %d bytes remain", count, remaining)
	}

	offsets := make([]int, 0, capHint(int(count)))
	for i := uint32(0); i < count; i++ {
		off := le64(data[pos : pos+8])
		if off > uint64(len(data)) {
			return nil, fmt.Errorf("rule set account: revision map entry %d offset %d exceeds account size %d", i, off, len(data))
		}
		offsets = append(offsets, int(off))
		pos += 8
	}
	return offsets, nil
}

// parseRevision decodes one revision from its raw byte slice. The slice begins
// at the lib-version byte. For V1 the msgpack payload follows the version byte;
// for V2 the version marker is embedded in the payload header and the whole
// slice is decoded. An unknown lib version is preserved (nil RuleSet) rather
// than rejected.
func parseRevision(revBytes []byte) (Revision, error) {
	if len(revBytes) == 0 {
		return Revision{}, fmt.Errorf("empty revision slice")
	}
	rev := Revision{
		LibVersion: revBytes[0],
		Raw:        revBytes,
	}
	switch revBytes[0] {
	case libVersionV1:
		rs, err := parseRuleSetV1(revBytes[1:])
		if err != nil {
			return Revision{}, err
		}
		rev.RuleSet = rs
	case libVersionV2:
		rs, err := parseRuleSetV2(revBytes)
		if err != nil {
			return Revision{}, err
		}
		rev.RuleSet = rs
	default:
		// Unknown lib version: keep the raw bytes, leave RuleSet nil.
	}
	return rev, nil
}

// ParseRuleSet decodes a single revision's bytes (starting at the lib-version
// byte) into a RuleSet. It returns an error for an unknown lib version, since a
// caller asking for a decoded RuleSet cannot proceed without one.
func ParseRuleSet(revBytes []byte) (*RuleSet, error) {
	rev, err := parseRevision(revBytes)
	if err != nil {
		return nil, err
	}
	if rev.RuleSet == nil {
		return nil, fmt.Errorf("unsupported rule set lib version %d", rev.LibVersion)
	}
	return rev.RuleSet, nil
}
