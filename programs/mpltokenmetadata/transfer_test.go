package mpltokenmetadata

import (
	"encoding/base64"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

// Expected instruction bytes, account layout, and PDA derivations are pinned
// against two live mainnet Transfer transactions so the generated builder can
// be checked against real on-chain data, not just its own IDL:
//
//	21y9xtg4MyLKDiKEAmW9bUpxD9PCgnbDMPrwEmXdBpaTvk4iCKaJe1orPYkJ6DLBciDwPMnKi18rx8s8xtTvNFNp
//	  (metadata carries no rule set — the authorization_rules slot holds the
//	  token-metadata program id, the Metaplex convention for an omitted
//	  optional positional account)
//	5NUZM1XuP5x3Dzkia2mRN59eoyFzGLRqkKp5LLQRDLMTJfZA7JpvoxhjLcdakgnw84PBNxpRo5AEcNMJ5bhJBEJz
//	  (metadata carries rule set eBJLFYPxJmMGKuFwpDWkzxZeUrad92kZRC5BJLpzyT9)
//
// In both, the account keys at positions 1 (token_owner) and 9 (authority)
// coincide — a holder transfer. The on-chain message reports position 1 as a
// signer only because the message-level flag is OR-merged across every
// reference to the key; the per-instruction metas below carry the builder's
// unmerged view (token_owner read-only non-signer).
func TestNewTransferInstruction_MatchesMainnetTransaction(t *testing.T) {
	tests := []struct {
		name             string
		token            string
		tokenOwner       string
		destination      string
		destinationOwner string
		mint             string
		metadata         string
		edition          string
		ownerRecord      string
		destRecord       string
		payer            string
		authRules        string
	}{
		{
			name:             "21y9xtg4MyLKDiKEAmW9bUpxD9PCgnbDMPrwEmXdBpaTvk4iCKaJe1orPYkJ6DLBciDwPMnKi18rx8s8xtTvNFNp",
			token:            "HuD3nFg3W8qwzNp74ghv9bixDn4bHt92HrJ9UQBSLe38",
			tokenOwner:       "4tyMk97N1TQT44wwtto4T4s7UYMAFMoBkySfP8PRjLbJ",
			destination:      "Fj5Vi5koZafnif61V6uLhUyKspkCzykkM8b8aHLCkmz3",
			destinationOwner: "6JHVE4rFevAc527fKvbPHEgnRLkuVYWJTcbhYQjCbFFL",
			mint:             "95nrq4rbGKGCcoa6iPswPfTDP2jDxLs7E9Eb2mVwcVwc",
			metadata:         "GMcyCnKVLU8WWwceR2pXjuegKfptSK1WpvNHnmDp3qgx",
			edition:          "DbFJo4gLa9Umh7GDPeaZozJCD71QuuYJfJERiDX9edSd",
			ownerRecord:      "FZDsyWYKQe6qDWsvbG4qmtWkhJ7GSCsG1JEamofjsHp3",
			destRecord:       "Db95BHFvEHXGageDYLC4mtcZPNdT36zDYRNwM7qiDu2Z",
			payer:            "7EReCbXkG3pqev8DdRn7QT5NJxxpxd4ERKQBbmfgtKRc",
			authRules:        ProgramID.String(),
		},
		{
			name:             "5NUZM1XuP5x3Dzkia2mRN59eoyFzGLRqkKp5LLQRDLMTJfZA7JpvoxhjLcdakgnw84PBNxpRo5AEcNMJ5bhJBEJz",
			token:            "D2pBzT4jrP99UHYBYqAe9e8QSxxknmngU6DPYG8GD4Fo",
			tokenOwner:       "21RQUA4o9ZUYvHKmYedFva6BU2DL9GrHkXJ1Fx2ckvwZ",
			destination:      "98QURUj7TCx8QtSks18Vzfox7HLrRBrozBYvDS3nFKLf",
			destinationOwner: "Low6UekJP3QrFVMfNRTL8CPK2SiGFhvp57sgF2pkmVu",
			mint:             "3waNhdX51iNUdDKzm94Z4B4e7HgkAUPiGgYr9yT2JBmV",
			metadata:         "7izaEkfvzZtdDYAB1JgP3D5smcNQS6hnzitR3bJmEr5h",
			edition:          "23tvPhZ82kmR8DhryMUYHuDhDNh4US8233dDzYnJpM2y",
			ownerRecord:      "2m4wy9qdEJTuwkdkCdbCNjYZFN11ux4RrxkVipVMTPbi",
			destRecord:       "HJrf3Mj6sAMLHK7Lj8Ji9kTXhLkJwHGbhCpaawzPpivt",
			payer:            "GachaNgyXTU3zFogQ8Z5jR2BLXs8215X2AtEH18VxJq3",
			authRules:        "eBJLFYPxJmMGKuFwpDWkzxZeUrad92kZRC5BJLpzyT9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := solanago.MustPublicKeyFromBase58(tt.token)
			tokenOwner := solanago.MustPublicKeyFromBase58(tt.tokenOwner)
			destination := solanago.MustPublicKeyFromBase58(tt.destination)
			destinationOwner := solanago.MustPublicKeyFromBase58(tt.destinationOwner)
			mint := solanago.MustPublicKeyFromBase58(tt.mint)
			metadata := solanago.MustPublicKeyFromBase58(tt.metadata)
			edition := solanago.MustPublicKeyFromBase58(tt.edition)
			ownerRecord := solanago.MustPublicKeyFromBase58(tt.ownerRecord)
			destRecord := solanago.MustPublicKeyFromBase58(tt.destRecord)
			payer := solanago.MustPublicKeyFromBase58(tt.payer)
			authRules := solanago.MustPublicKeyFromBase58(tt.authRules)

			inst, err := NewTransferInstruction(
				&TransferArgs_V1{Amount: 1},
				token,
				tokenOwner,
				destination,
				destinationOwner,
				mint,
				metadata,
				edition,
				ownerRecord,
				destRecord,
				tokenOwner, // authority = holder
				payer,
				solanago.SystemProgramID,
				solanago.SysVarInstructionsPubkey,
				solanago.TokenProgramID,
				solanago.SPLAssociatedTokenAccountProgramID,
				TokenAuthRulesProgramID,
				authRules,
			)
			require.NoError(t, err)

			data, err := inst.Data()
			require.NoError(t, err)
			// discriminator 49, enum variant V1 (0), amount 1 as LE u64,
			// authorization_data None (0)
			require.Equal(t, []byte{49, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, data)

			require.Equal(t, ProgramID, inst.ProgramID())

			expected := []*solanago.AccountMeta{
				solanago.NewAccountMeta(token, true, false),
				solanago.NewAccountMeta(tokenOwner, false, false),
				solanago.NewAccountMeta(destination, true, false),
				solanago.NewAccountMeta(destinationOwner, false, false),
				solanago.NewAccountMeta(mint, false, false),
				solanago.NewAccountMeta(metadata, true, false),
				solanago.NewAccountMeta(edition, false, false),
				solanago.NewAccountMeta(ownerRecord, true, false),
				solanago.NewAccountMeta(destRecord, true, false),
				solanago.NewAccountMeta(tokenOwner, false, true),
				solanago.NewAccountMeta(payer, true, true),
				solanago.NewAccountMeta(solanago.SystemProgramID, false, false),
				solanago.NewAccountMeta(solanago.SysVarInstructionsPubkey, false, false),
				solanago.NewAccountMeta(solanago.TokenProgramID, false, false),
				solanago.NewAccountMeta(solanago.SPLAssociatedTokenAccountProgramID, false, false),
				solanago.NewAccountMeta(TokenAuthRulesProgramID, false, false),
				solanago.NewAccountMeta(authRules, false, false),
			}
			require.Equal(t, expected, inst.Accounts())
		})
	}
}

// PDA derivations pinned against the same two mainnet transactions: metadata
// and master edition derive from the mint; token records derive from
// (mint, token account) for both the source and destination sides.
func TestPDAs_MatchMainnetAccounts(t *testing.T) {
	tests := []struct {
		name        string
		mint        string
		token       string
		destination string
		metadata    string
		edition     string
		ownerRecord string
		destRecord  string
	}{
		{
			name:        "95nrq4rbGKGCcoa6iPswPfTDP2jDxLs7E9Eb2mVwcVwc",
			mint:        "95nrq4rbGKGCcoa6iPswPfTDP2jDxLs7E9Eb2mVwcVwc",
			token:       "HuD3nFg3W8qwzNp74ghv9bixDn4bHt92HrJ9UQBSLe38",
			destination: "Fj5Vi5koZafnif61V6uLhUyKspkCzykkM8b8aHLCkmz3",
			metadata:    "GMcyCnKVLU8WWwceR2pXjuegKfptSK1WpvNHnmDp3qgx",
			edition:     "DbFJo4gLa9Umh7GDPeaZozJCD71QuuYJfJERiDX9edSd",
			ownerRecord: "FZDsyWYKQe6qDWsvbG4qmtWkhJ7GSCsG1JEamofjsHp3",
			destRecord:  "Db95BHFvEHXGageDYLC4mtcZPNdT36zDYRNwM7qiDu2Z",
		},
		{
			name:        "3waNhdX51iNUdDKzm94Z4B4e7HgkAUPiGgYr9yT2JBmV",
			mint:        "3waNhdX51iNUdDKzm94Z4B4e7HgkAUPiGgYr9yT2JBmV",
			token:       "D2pBzT4jrP99UHYBYqAe9e8QSxxknmngU6DPYG8GD4Fo",
			destination: "98QURUj7TCx8QtSks18Vzfox7HLrRBrozBYvDS3nFKLf",
			metadata:    "7izaEkfvzZtdDYAB1JgP3D5smcNQS6hnzitR3bJmEr5h",
			edition:     "23tvPhZ82kmR8DhryMUYHuDhDNh4US8233dDzYnJpM2y",
			ownerRecord: "2m4wy9qdEJTuwkdkCdbCNjYZFN11ux4RrxkVipVMTPbi",
			destRecord:  "HJrf3Mj6sAMLHK7Lj8Ji9kTXhLkJwHGbhCpaawzPpivt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mint := solanago.MustPublicKeyFromBase58(tt.mint)

			metadata, err := MetadataPDA(mint)
			require.NoError(t, err)
			require.Equal(t, tt.metadata, metadata.String())

			edition, err := MasterEditionPDA(mint)
			require.NoError(t, err)
			require.Equal(t, tt.edition, edition.String())

			ownerRecord, err := TokenRecordPDA(mint, solanago.MustPublicKeyFromBase58(tt.token))
			require.NoError(t, err)
			require.Equal(t, tt.ownerRecord, ownerRecord.String())

			destRecord, err := TokenRecordPDA(mint, solanago.MustPublicKeyFromBase58(tt.destination))
			require.NoError(t, err)
			require.Equal(t, tt.destRecord, destRecord.String())
		})
	}
}

// Raw on-chain bytes captured 2026-07-16 via getMultipleAccounts (base64).
// Metadata accounts are fixed-size with zero padding after the borsh payload.
const (
	// Metadata GMcyCnKVLU8WWwceR2pXjuegKfptSK1WpvNHnmDp3qgx — programmable
	// config present, rule set NOT set.
	metadataNoRuleSetB64 = "BI8VunRlQ+1PjsFbTMv30XVVR+G6TdpFunkyNKGS5mGLeBer9YMsl47mkOPBd6xZFEXBctI860Odldal1UkY6CcgAAAAUHJvZ3JhbW1hYmxlTm9uRnVuZ2libGUgVG9rZW4AAAAKAAAAcE5GVAAAAAAAAMgAAABodHRwczovL2dhdGV3YXkucGluaXQuaW8vaXBmcy9RbWRwbkxycFA3djd6V0R1aVdTdmJrUWN4V0FHU3E2Tmo5UmNyNVJlUWpNMnBoLzMxMDQuanNvbgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQEAAACPFbp0ZUPtT47BW0zL99F1VUfhuk3aRbp5MjShkuZhiwFkAAAB/gEFAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="

	// Metadata FzbB8YtBeMqxuSHqutJHr2JfjWFnvcKM3DVqjaxLMNHn — rule set
	// eBJLFYPxJmMGKuFwpDWkzxZeUrad92kZRC5BJLpzyT9.
	metadataRuleSetB64 = "BLhJvSf1pWAOic+nA+axhrqDqHE78GhDY+x+MkbVM/zyMmRt3ZKJKTBdAiBbZNfV5gP0hHFP9rZ4jXQakL7BNq4gAAAAMjAxOSAjU00xODkgRnVsbCBBcnQvQmxhc3RvaXNlIEcKAAAAQ09MTEVDVE9SAMgAAABodHRwczovL2Fyd2VhdmUubmV0L2syWlk3c2d4a2FQNlpxSERuODNQTmZXMEt2dC1tWk44OXhWbDNsZXpXT0kAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMgAAQEAAAC4Sb0n9aVgDonPpwPmsYa6g6hxO/BoQ2PsfjJG1TP88gBkAAEB/wEEAQGmenUWm3sm4gcWrNZ52hLQ6Qw2PFdY1cFlpTDn2o/CKAAAAQABCYYiheNxCpDVHZ5HAt6andXp/cisgdLSrNHh3cgk/sQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="

	// TokenRecord Db95BHFvEHXGageDYLC4mtcZPNdT36zDYRNwM7qiDu2Z — freshly
	// created by the pinned transfer: state Unlocked, no delegate.
	tokenRecordB64 = "C/8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
)

func TestUnmarshalMetadataAccount(t *testing.T) {
	t.Run("no rule set", func(t *testing.T) {
		raw, err := base64.StdEncoding.DecodeString(metadataNoRuleSetB64)
		require.NoError(t, err)

		m, err := UnmarshalMetadataAccount(raw)
		require.NoError(t, err)
		require.Equal(t, Key_MetadataV1, m.Key)
		require.Equal(t, "95nrq4rbGKGCcoa6iPswPfTDP2jDxLs7E9Eb2mVwcVwc", m.Mint.String())
		require.NotNil(t, m.TokenStandard)
		require.Equal(t, TokenStandard_ProgrammableNonFungibleEdition, *m.TokenStandard)
		require.NotNil(t, m.ProgrammableConfig)
		require.Nil(t, m.RuleSet())
	})

	t.Run("rule set", func(t *testing.T) {
		raw, err := base64.StdEncoding.DecodeString(metadataRuleSetB64)
		require.NoError(t, err)

		m, err := UnmarshalMetadataAccount(raw)
		require.NoError(t, err)
		require.Equal(t, Key_MetadataV1, m.Key)
		require.Equal(t, "4PiCtSrYRLUBKsVSnv5hdavektHe52hNx8ABVmzcDVA5", m.Mint.String())
		require.NotNil(t, m.TokenStandard)
		require.Equal(t, TokenStandard_ProgrammableNonFungible, *m.TokenStandard)
		require.NotNil(t, m.RuleSet())
		require.Equal(t, "eBJLFYPxJmMGKuFwpDWkzxZeUrad92kZRC5BJLpzyT9", m.RuleSet().String())
	})
}

func TestUnmarshalTokenRecord(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString(tokenRecordB64)
	require.NoError(t, err)

	r, err := UnmarshalTokenRecord(raw)
	require.NoError(t, err)
	require.Equal(t, Key_TokenRecord, r.Key)
	require.Equal(t, uint8(255), r.Bump)
	require.Equal(t, TokenState_Unlocked, r.State)
	require.Nil(t, r.RuleSetRevision)
	require.Nil(t, r.Delegate)
	require.Nil(t, r.DelegateRole)
	require.Nil(t, r.LockedTransfer)
}
