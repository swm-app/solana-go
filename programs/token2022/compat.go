package token2022

import solanago "github.com/gagliardetto/solana-go"

// NewTransferFeeTransferCheckedWithFeeInstruction is a compatibility helper
// matching the README example signature.
func NewTransferFeeTransferCheckedWithFeeInstruction(
	sourceTokenAccount solanago.PublicKey,
	mint solanago.PublicKey,
	destinationTokenAccount solanago.PublicKey,
	authority solanago.PublicKey,
	multiSigners []solanago.PublicKey,
	amount uint64,
	decimals uint8,
	fee uint64,
) (solanago.Instruction, error) {
	return NewTransferCheckedWithFeeInstruction(
		TransferCheckedWithFeeParams{
			Amount:   amount,
			Decimals: decimals,
			Fee:      fee,
		},
		TransferCheckedWithFeeAccounts{
			Source:       sourceTokenAccount,
			Mint:         mint,
			Destination:  destinationTokenAccount,
			Authority:    authority,
			MultiSigners: multiSigners,
		},
	)
}
