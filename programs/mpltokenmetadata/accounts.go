// Hand-written companion to the generated program client: account decoding
// that the generator gets wrong.

package mpltokenmetadata

import (
	binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

// UnmarshalMetadataAccount decodes a Metadata account from raw on-chain bytes,
// including the trailing ProgrammableConfig enum that the generated
// UnmarshalMetadata silently loses: the generated decoder calls
// decoder.Decode on a pointer-to-interface field, which consumes no bytes and
// leaves the field pointing at a nil interface. This helper re-runs the
// generated decoder (whose position stops exactly at the enum payload because
// of that zero-consumption) and then decodes the enum explicitly, repairing
// the ProgrammableConfig field in place.
//
// On-chain Metadata accounts are fixed-size with zero padding after the borsh
// payload; the decoder tolerates trailing bytes by construction (it never
// requires EOF).
func UnmarshalMetadataAccount(data []byte) (*Metadata, error) {
	decoder := binary.NewBorshDecoder(data)
	obj := new(Metadata)
	if err := obj.UnmarshalWithDecoder(decoder); err != nil {
		return nil, err
	}
	// Non-nil means the option byte read true — the enum payload follows at
	// the decoder's current position.
	if obj.ProgrammableConfig != nil {
		pc, err := DecodeProgrammableConfig(decoder)
		if err != nil {
			return nil, err
		}
		*obj.ProgrammableConfig = pc
	}
	return obj, nil
}

// RuleSet returns the metadata's programmable rule set address, or nil when
// the account carries no programmable config or the config has no rule set.
func (obj *Metadata) RuleSet() *solanago.PublicKey {
	if obj.ProgrammableConfig == nil {
		return nil
	}
	v1, ok := (*obj.ProgrammableConfig).(*ProgrammableConfig_V1)
	if !ok || v1 == nil {
		return nil
	}
	return v1.RuleSet
}
