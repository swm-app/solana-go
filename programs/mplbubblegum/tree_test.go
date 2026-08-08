package mplbubblegum

import (
	"encoding/base64"
	"os"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

func loadTreeAccountData(t *testing.T, path string) []byte {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	data, err := base64.StdEncoding.DecodeString(string(raw))
	require.NoError(t, err)
	return data
}

// TestParseTreeAccount_MatchesMainnetTrees pins the parser against the two
// tree accounts backing the instruction fixtures: a V1 tree owned by SPL
// Account Compression and a V2 tree owned by MPL Account Compression.
func TestParseTreeAccount_MatchesMainnetTrees(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		wantMaxDepth  uint32
		wantBuffer    uint32
		wantCanopy    uint32
		wantAuthority string
		wantSlot      uint64
	}{
		{
			name:          "GYGG71sQjYSZjxiXXfDjTLwHj4cNi4vUJvLsreANYAep (V1)",
			path:          "testdata/tree_transfer.b64",
			wantMaxDepth:  16,
			wantBuffer:    64,
			wantCanopy:    8,
			wantAuthority: "3XnzMMXQwE3BqSK1RStLbxUCtZU1J6dWZ5axnVi3kcmu",
			wantSlot:      404404323,
		},
		{
			name:          "Dw4aZw7qkgbwHyEPA55p9uN818S6inD3L1h8E4Z8s45H (V2)",
			path:          "testdata/tree_transfer_v2.b64",
			wantMaxDepth:  20,
			wantBuffer:    1024,
			wantCanopy:    13,
			wantAuthority: "7mWBfPKHipiXha3qsY8u9jJ2n9WQCmkXa861Ncejg2AU",
			wantSlot:      397531840,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadTreeAccountData(t, tt.path)

			tree, err := ParseTreeAccount(data, uint64(len(data)))
			require.NoError(t, err)
			require.Equal(t, tt.wantMaxDepth, tree.MaxDepth)
			require.Equal(t, tt.wantBuffer, tree.MaxBufferSize)
			require.Equal(t, tt.wantCanopy, tree.CanopyDepth)
			require.Equal(t, tt.wantAuthority, tree.Authority.String())
			require.Equal(t, tt.wantSlot, tree.CreationSlot)

			// The header alone, paired with the account's length, yields the
			// same tree — this is what a header-only dataSlice fetch supplies.
			sliced, err := ParseTreeAccount(data[:TreeHeaderSize], uint64(len(data)))
			require.NoError(t, err)
			require.Equal(t, tree, sliced)

			// The merkle proof carried in a Bubblegum transfer instruction
			// covers only the levels not already cached in the canopy.
			proofLen := tt.wantMaxDepth - tt.wantCanopy
			if tt.wantMaxDepth == 16 {
				require.EqualValues(t, 8, proofLen) // matches the V1 fixture's 8 trailing proof accounts
			} else {
				require.EqualValues(t, 7, proofLen) // matches the V2 fixture's 7 trailing proof accounts
			}
		})
	}
}

func TestParseTreeAccount_Errors(t *testing.T) {
	valid := loadTreeAccountData(t, "testdata/tree_transfer.b64")

	t.Run("truncated header", func(t *testing.T) {
		_, err := ParseTreeAccount(valid[:TreeHeaderSize-1], uint64(len(valid)))
		require.Error(t, err)
	})

	t.Run("account shorter than the data supplied for it", func(t *testing.T) {
		_, err := ParseTreeAccount(valid, uint64(len(valid)-1))
		require.Error(t, err)
	})

	t.Run("wrong account type", func(t *testing.T) {
		corrupted := make([]byte, len(valid))
		copy(corrupted, valid)
		corrupted[0] = 2
		_, err := ParseTreeAccount(corrupted, uint64(len(corrupted)))
		require.Error(t, err)
	})

	t.Run("canopy inconsistent with power-of-two rule", func(t *testing.T) {
		// Drop 32 bytes from the canopy: the remaining canopy length is
		// still a multiple of 32, but the resulting node count is no
		// longer 2^(c+1)-2 for any integer c.
		_, err := ParseTreeAccount(valid[:TreeHeaderSize], uint64(len(valid)-32))
		require.Error(t, err)
	})

	t.Run("account too short for declared max_depth/max_buffer_size", func(t *testing.T) {
		_, err := ParseTreeAccount(valid[:TreeHeaderSize], TreeHeaderSize+10)
		require.Error(t, err)
	})
}

func TestTransferVersionForTreeOwner(t *testing.T) {
	v1, err := TransferVersionForTreeOwner(SPLAccountCompressionProgramID)
	require.NoError(t, err)
	require.Equal(t, TransferVersionV1, v1)

	v2, err := TransferVersionForTreeOwner(MPLAccountCompressionProgramID)
	require.NoError(t, err)
	require.Equal(t, TransferVersionV2, v2)

	_, err = TransferVersionForTreeOwner(solanago.MustPublicKeyFromBase58("11111111111111111111111111111111"))
	require.Error(t, err)
}
