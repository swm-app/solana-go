package mpltokenauthrules

import solanago "github.com/gagliardetto/solana-go"

// ProgramID is the Metaplex Token Auth Rules program.
var ProgramID = solanago.MustPublicKeyFromBase58("auth9SigNpDKz4sJJ1DfCTuZrZNSAgh9sFD3rboVmgg")

// ruleSetSeedPrefix is the first seed of a rule-set PDA. It mirrors the
// on-chain PREFIX constant ("rule_set").
var ruleSetSeedPrefix = []byte("rule_set")

// RuleSetPDA derives the rule-set account address for a given creator and rule
// set name: seeds ["rule_set", creator, name].
func RuleSetPDA(creator solanago.PublicKey, name string) (solanago.PublicKey, error) {
	addr, _, err := solanago.FindProgramAddress(
		[][]byte{ruleSetSeedPrefix, creator.Bytes(), []byte(name)},
		ProgramID,
	)
	return addr, err
}
