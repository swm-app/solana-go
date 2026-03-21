package token2022

import (
	"bytes"
	stdbinary "encoding/binary"
	"fmt"
	"testing"

	binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

func TestResolveTransferHookAccountsFromData(t *testing.T) {
	source := pubkeyWithByte(1)
	mint := pubkeyWithByte(2)
	destination := pubkeyWithByte(3)
	authority := pubkeyWithByte(4)
	hookProgramID := pubkeyWithByte(5)
	extraMeta1 := pubkeyWithByte(6)
	extraMeta2 := pubkeyWithByte(7)
	amount := uint64(42)

	validationAccount, err := GetTransferHookExtraAccountMetaListAddress(mint, hookProgramID)
	require.NoError(t, err)

	metaConfigs := []transferHookExtraAccountMeta{
		newStaticTestExtraAccountMeta(extraMeta1, true, false),
		newStaticTestExtraAccountMeta(extraMeta2, true, false),
		newSeededTestExtraAccountMeta([]transferHookSeed{
			{kind: transferHookSeedAccountKey, index: 0},
			{kind: transferHookSeedAccountKey, index: 2},
			{kind: transferHookSeedAccountKey, index: 4},
		}, false, true),
		newSeededTestExtraAccountMeta([]transferHookSeed{
			{kind: transferHookSeedInstruction, index: 8, length: 8},
			{kind: transferHookSeedAccountKey, index: 2},
			{kind: transferHookSeedAccountKey, index: 5},
			{kind: transferHookSeedAccountKey, index: 7},
		}, false, true),
	}

	validationData := packTransferHookValidationData(metaConfigs)
	resolved, err := ResolveTransferHookAccountsFromMintStateData(
		TransferHookTransfer{
			Source:      source,
			Mint:        mint,
			Destination: destination,
			Authority:   authority,
			Amount:      amount,
		},
		hookProgramID,
		validationData,
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, resolved)
	require.Equal(t, hookProgramID, resolved.ProgramID)
	require.Equal(t, validationAccount, resolved.ValidationAccount)

	expectedPDA1, _, err := solanago.FindProgramAddress(
		[][]byte{
			source[:],
			destination[:],
			validationAccount[:],
		},
		hookProgramID,
	)
	require.NoError(t, err)

	amountBytes := make([]byte, 8)
	stdbinary.LittleEndian.PutUint64(amountBytes, amount)
	expectedPDA2, _, err := solanago.FindProgramAddress(
		[][]byte{
			amountBytes,
			destination[:],
			extraMeta1[:],
			expectedPDA1[:],
		},
		hookProgramID,
	)
	require.NoError(t, err)

	require.Len(t, resolved.ExtraAccounts, 4)
	require.Equal(t, extraMeta1, resolved.ExtraAccounts[0].PublicKey)
	require.True(t, resolved.ExtraAccounts[0].IsSigner)
	require.False(t, resolved.ExtraAccounts[0].IsWritable)
	require.Equal(t, extraMeta2, resolved.ExtraAccounts[1].PublicKey)
	require.True(t, resolved.ExtraAccounts[1].IsSigner)
	require.False(t, resolved.ExtraAccounts[1].IsWritable)
	require.Equal(t, expectedPDA1, resolved.ExtraAccounts[2].PublicKey)
	require.False(t, resolved.ExtraAccounts[2].IsSigner)
	require.True(t, resolved.ExtraAccounts[2].IsWritable)
	require.Equal(t, expectedPDA2, resolved.ExtraAccounts[3].PublicKey)
	require.False(t, resolved.ExtraAccounts[3].IsSigner)
	require.True(t, resolved.ExtraAccounts[3].IsWritable)

	remaining := resolved.RemainingAccounts()
	require.Len(t, remaining, 5)
	require.Equal(t, validationAccount, remaining[0].PublicKey)
	require.Equal(t, extraMeta1, remaining[1].PublicKey)
	require.Equal(t, extraMeta2, remaining[2].PublicKey)
	require.Equal(t, expectedPDA1, remaining[3].PublicKey)
	require.Equal(t, expectedPDA2, remaining[4].PublicKey)
}

func TestNewTransferCheckedInstructionWithTransferHook(t *testing.T) {
	source := pubkeyWithByte(11)
	mint := pubkeyWithByte(12)
	destination := pubkeyWithByte(13)
	authority := pubkeyWithByte(14)
	signer1 := pubkeyWithByte(15)
	signer2 := pubkeyWithByte(16)
	hookProgramID := pubkeyWithByte(17)
	extraMeta := pubkeyWithByte(18)

	validationAccount, err := GetTransferHookExtraAccountMetaListAddress(mint, hookProgramID)
	require.NoError(t, err)

	validationData := packTransferHookValidationData([]transferHookExtraAccountMeta{
		newStaticTestExtraAccountMeta(extraMeta, false, true),
	})
	mintData := mustEncodeMintWithTransferHook(t, hookProgramID)

	instruction, err := NewTransferCheckedInstructionWithTransferHook(
		TransferCheckedParams{
			Amount:   55,
			Decimals: 6,
		},
		TransferCheckedAccounts{
			Source:       source,
			Mint:         mint,
			Destination:  destination,
			Authority:    authority,
			MultiSigners: []solanago.PublicKey{signer1, signer2},
		},
		mintData,
		validationData,
		nil,
	)
	require.NoError(t, err)

	accounts := instruction.Accounts()
	require.Len(t, accounts, 8)
	require.Equal(t, source, accounts[0].PublicKey)
	require.Equal(t, mint, accounts[1].PublicKey)
	require.Equal(t, destination, accounts[2].PublicKey)
	require.Equal(t, authority, accounts[3].PublicKey)
	require.Equal(t, signer1, accounts[4].PublicKey)
	require.Equal(t, signer2, accounts[5].PublicKey)
	require.Equal(t, validationAccount, accounts[6].PublicKey)
	require.Equal(t, extraMeta, accounts[7].PublicKey)
	require.True(t, accounts[7].IsWritable)
}

func newStaticTestExtraAccountMeta(
	pubkey solanago.PublicKey,
	isSigner bool,
	isWritable bool,
) transferHookExtraAccountMeta {
	var addressConfig [32]byte
	copy(addressConfig[:], pubkey[:])
	return transferHookExtraAccountMeta{
		Discriminator: transferHookExtraAccountMetaStatic,
		AddressConfig: addressConfig,
		IsSigner:      isSigner,
		IsWritable:    isWritable,
	}
}

func newSeededTestExtraAccountMeta(
	seeds []transferHookSeed,
	isSigner bool,
	isWritable bool,
) transferHookExtraAccountMeta {
	return transferHookExtraAccountMeta{
		Discriminator: transferHookExtraAccountMetaSeeds,
		AddressConfig: packTestSeedConfigs(seeds),
		IsSigner:      isSigner,
		IsWritable:    isWritable,
	}
}

func packTestSeedConfigs(seeds []transferHookSeed) [32]byte {
	var out [32]byte
	offset := 0
	for _, seed := range seeds {
		switch seed.kind {
		case transferHookSeedLiteral:
			out[offset] = transferHookSeedLiteral
			out[offset+1] = uint8(len(seed.literal))
			copy(out[offset+2:], seed.literal)
			offset += 2 + len(seed.literal)
		case transferHookSeedInstruction:
			out[offset] = transferHookSeedInstruction
			out[offset+1] = seed.index
			out[offset+2] = seed.length
			offset += 3
		case transferHookSeedAccountKey:
			out[offset] = transferHookSeedAccountKey
			out[offset+1] = seed.index
			offset += 2
		case transferHookSeedAccountData:
			out[offset] = transferHookSeedAccountData
			out[offset+1] = seed.accountIndex
			out[offset+2] = seed.dataIndex
			out[offset+3] = seed.length
			offset += 4
		default:
			panic(fmt.Sprintf("unsupported test seed kind %d", seed.kind))
		}
	}
	return out
}

func packTransferHookValidationData(metas []transferHookExtraAccountMeta) []byte {
	payload := make([]byte, transferHookExtraAccountMetaCountLen+len(metas)*transferHookExtraAccountMetaLen)
	stdbinary.LittleEndian.PutUint32(payload[:transferHookExtraAccountMetaCountLen], uint32(len(metas)))
	offset := transferHookExtraAccountMetaCountLen
	for _, meta := range metas {
		payload[offset] = meta.Discriminator
		copy(payload[offset+1:offset+33], meta.AddressConfig[:])
		if meta.IsSigner {
			payload[offset+33] = 1
		}
		if meta.IsWritable {
			payload[offset+34] = 1
		}
		offset += transferHookExtraAccountMetaLen
	}
	data := make([]byte, transferHookExecuteDiscriminatorLen+transferHookTLVLengthLen+len(payload))
	copy(data[:transferHookExecuteDiscriminatorLen], transferHookExecuteDiscriminator[:])
	stdbinary.LittleEndian.PutUint32(
		data[transferHookExecuteDiscriminatorLen:transferHookExecuteDiscriminatorLen+transferHookTLVLengthLen],
		uint32(len(payload)),
	)
	copy(data[transferHookExecuteDiscriminatorLen+transferHookTLVLengthLen:], payload)
	return data
}

func mustEncodeMintWithTransferHook(t *testing.T, hookProgramID solanago.PublicKey) []byte {
	t.Helper()

	baseMint := MintAccount{
		Supply:        100,
		Decimals:      6,
		IsInitialized: true,
	}
	baseBuf := new(bytes.Buffer)
	err := binary.NewBorshEncoder(baseBuf).Encode(baseMint)
	require.NoError(t, err)
	require.Len(t, baseBuf.Bytes(), mintBaseLen)

	data := make([]byte, tokenBaseLen+1+tlvTypeLen+tlvLengthLen+64)
	copy(data[:mintBaseLen], baseBuf.Bytes())
	data[tokenBaseLen] = byte(AccountTypeMint)

	offset := tokenBaseLen + 1
	stdbinary.LittleEndian.PutUint16(data[offset:offset+tlvTypeLen], uint16(ExtensionTypeTransferHook))
	offset += tlvTypeLen
	stdbinary.LittleEndian.PutUint16(data[offset:offset+tlvLengthLen], 64)
	offset += tlvLengthLen
	copy(data[offset+32:offset+64], hookProgramID[:])
	return data
}

func pubkeyWithByte(b byte) solanago.PublicKey {
	var key solanago.PublicKey
	for i := range key {
		key[i] = b
	}
	return key
}
