// Hand-written companion to the generated program client: PDA derivations and
// well-known program ids the IDL does not carry.

package mpltokenmetadata

import solanago "github.com/gagliardetto/solana-go"

// TokenAuthRulesProgramID is the Metaplex Token Auth Rules program, which
// validates rule sets referenced by programmable NFTs. The Transfer instruction
// expects it in the authorization_rules_program slot whenever a rule set
// account is supplied.
var TokenAuthRulesProgramID = solanago.MustPublicKeyFromBase58("auth9SigNpDKz4sJJ1DfCTuZrZNSAgh9sFD3rboVmgg")

var metadataSeed = []byte("metadata")

// MetadataPDA derives the metadata account for a mint:
// seeds ["metadata", program id, mint].
func MetadataPDA(mint solanago.PublicKey) (solanago.PublicKey, error) {
	addr, _, err := solanago.FindProgramAddress(
		[][]byte{metadataSeed, ProgramID.Bytes(), mint.Bytes()},
		ProgramID,
	)
	return addr, err
}

// MasterEditionPDA derives the master edition account for a mint:
// seeds ["metadata", program id, mint, "edition"].
func MasterEditionPDA(mint solanago.PublicKey) (solanago.PublicKey, error) {
	addr, _, err := solanago.FindProgramAddress(
		[][]byte{metadataSeed, ProgramID.Bytes(), mint.Bytes(), []byte("edition")},
		ProgramID,
	)
	return addr, err
}

// TokenRecordPDA derives the token record account for a (mint, token account)
// pair: seeds ["metadata", program id, mint, "token_record", token account].
// Programmable NFTs keep per-token-account state (lock state, delegate, rule
// set revision) in this record.
func TokenRecordPDA(mint, tokenAccount solanago.PublicKey) (solanago.PublicKey, error) {
	addr, _, err := solanago.FindProgramAddress(
		[][]byte{metadataSeed, ProgramID.Bytes(), mint.Bytes(), []byte("token_record"), tokenAccount.Bytes()},
		ProgramID,
	)
	return addr, err
}
