package mplcore

import (
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

// Expected instruction bytes and account layout are pinned against two live
// mainnet TransferV1 transactions so the generated builder can be checked
// against real on-chain data, not just its own IDL:
//   - 4SvwWadbEVpwFKKRqnmQNcYLppFMz5TXpuXQnFUxHchRbfSt9e5cZ8gbNd7NygHH9pKFcfMBhMpRabsXt4qf1Uad
//   - LPDHpVTN2G6GhEtHA42hSWseNYu2Hj4vckxwnQzJBpH1ZF8DCxhF7A6eY59kKGFvXe9DzDZSjsX1Cv4FY4JU4pF
//
// In both, the system_program and log_wrapper accounts carry the mpl-core
// program ID itself: the Metaplex convention for encoding an omitted
// optional positional account.
func TestNewTransferV1Instruction_MatchesMainnetTransaction(t *testing.T) {
	tests := []struct {
		name       string
		asset      string
		collection string
		payer      string
		authority  string
		newOwner   string
	}{
		{
			name:       "4SvwWadbEVpwFKKRqnmQNcYLppFMz5TXpuXQnFUxHchRbfSt9e5cZ8gbNd7NygHH9pKFcfMBhMpRabsXt4qf1Uad",
			asset:      "9kszY5urtVkHCrDiMY5S3bCck3YtChBYbrTLp5qDzSvL",
			collection: "phygZDQZJZVHvJGYPGoKPYUtXw7mstSYtTtcuh8LJcC",
			payer:      "mGrEwYbMo98b6koTfy9pSfFJA1FZoxFxY7kjcsKf9C5",
			authority:  "DAHMmedRWVNVY9TSno81B9tciTbm7Jqof1z621WcQmq3",
			newOwner:   "62Q9eeDY3eM8A5CnprBGYMPShdBjAzdpBdr71QHsS8dS",
		},
		{
			name:       "LPDHpVTN2G6GhEtHA42hSWseNYu2Hj4vckxwnQzJBpH1ZF8DCxhF7A6eY59kKGFvXe9DzDZSjsX1Cv4FY4JU4pF",
			asset:      "FFDSgGYDMGFZRYstPzkcPBoCUSjJSy7obS1MDHj1yPeb",
			collection: "CCryptUfeFSZ3Fgc9FLeKrhLVAP67FSqi1GuVoj9CRac",
			payer:      "GachaNgyXTU3zFogQ8Z5jR2BLXs8215X2AtEH18VxJq3",
			authority:  "Low6UekJP3QrFVMfNRTL8CPK2SiGFhvp57sgF2pkmVu",
			newOwner:   "GvSLpfy7dzBcnCU3uo4tqzFXmwp2ecxK276gV15fsgDH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := solanago.MustPublicKeyFromBase58(tt.asset)
			collection := solanago.MustPublicKeyFromBase58(tt.collection)
			payer := solanago.MustPublicKeyFromBase58(tt.payer)
			authority := solanago.MustPublicKeyFromBase58(tt.authority)
			newOwner := solanago.MustPublicKeyFromBase58(tt.newOwner)

			inst, err := NewTransferV1Instruction(
				TransferV1Args{CompressionProof: nil},
				asset,
				collection,
				payer,
				authority,
				newOwner,
				ProgramID,
				ProgramID,
			)
			require.NoError(t, err)

			data, err := inst.Data()
			require.NoError(t, err)
			require.Equal(t, []byte{14, 0}, data)

			require.Equal(t, ProgramID, inst.ProgramID())

			expected := []*solanago.AccountMeta{
				solanago.NewAccountMeta(asset, true, false),
				solanago.NewAccountMeta(collection, false, false),
				solanago.NewAccountMeta(payer, true, true),
				solanago.NewAccountMeta(authority, false, true),
				solanago.NewAccountMeta(newOwner, false, false),
				solanago.NewAccountMeta(ProgramID, false, false),
				solanago.NewAccountMeta(ProgramID, false, false),
			}
			require.Equal(t, expected, inst.Accounts())
		})
	}
}
