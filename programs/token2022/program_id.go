package token2022

import solanago "github.com/gagliardetto/solana-go"

var ProgramID = solanago.Token2022ProgramID

const ProgramName = "Token2022"

func SetProgramID(pubkey solanago.PublicKey) {
	ProgramID = pubkey
	solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

func init() {
	if !ProgramID.IsZero() {
		solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
	}
}
