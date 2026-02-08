package token2022

import (
	"fmt"

	binary "github.com/gagliardetto/binary"
)

var instructionDiscriminatorByType = map[InstructionType][]byte{
	InstructionTypeAmountToUIAmount:                                             []byte{23},
	InstructionTypeApplyConfidentialPendingBalance:                              []byte{27, 8},
	InstructionTypeApprove:                                                      []byte{4},
	InstructionTypeApproveChecked:                                               []byte{13},
	InstructionTypeApproveConfidentialTransferAccount:                           []byte{27, 3},
	InstructionTypeBurn:                                                         []byte{8},
	InstructionTypeBurnChecked:                                                  []byte{15},
	InstructionTypeCloseAccount:                                                 []byte{9},
	InstructionTypeConfidentialDeposit:                                          []byte{27, 5},
	InstructionTypeConfidentialTransfer:                                         []byte{27, 7},
	InstructionTypeConfidentialTransferWithFee:                                  []byte{27, 13},
	InstructionTypeConfidentialWithdraw:                                         []byte{27, 6},
	InstructionTypeConfigureConfidentialTransferAccount:                         []byte{27, 2},
	InstructionTypeCreateNativeMint:                                             []byte{31},
	InstructionTypeDisableConfidentialCredits:                                   []byte{27, 10},
	InstructionTypeDisableCpiGuard:                                              []byte{34, 1},
	InstructionTypeDisableHarvestToMint:                                         []byte{37, 5},
	InstructionTypeDisableMemoTransfers:                                         []byte{30, 1},
	InstructionTypeDisableNonConfidentialCredits:                                []byte{27, 12},
	InstructionTypeEmitTokenMetadata:                                            []byte{250, 166, 180, 250, 13, 12, 184, 70},
	InstructionTypeEmptyConfidentialTransferAccount:                             []byte{27, 4},
	InstructionTypeEnableConfidentialCredits:                                    []byte{27, 9},
	InstructionTypeEnableCpiGuard:                                               []byte{34, 0},
	InstructionTypeEnableHarvestToMint:                                          []byte{37, 4},
	InstructionTypeEnableMemoTransfers:                                          []byte{30, 0},
	InstructionTypeEnableNonConfidentialCredits:                                 []byte{27, 11},
	InstructionTypeFreezeAccount:                                                []byte{10},
	InstructionTypeGetAccountDataSize:                                           []byte{21},
	InstructionTypeHarvestWithheldTokensToMint:                                  []byte{26, 4},
	InstructionTypeHarvestWithheldTokensToMintForConfidentialTransferFee:        []byte{37, 3},
	InstructionTypeInitializeAccount:                                            []byte{1},
	InstructionTypeInitializeAccount2:                                           []byte{16},
	InstructionTypeInitializeAccount3:                                           []byte{18},
	InstructionTypeInitializeConfidentialTransferFee:                            []byte{37, 0},
	InstructionTypeInitializeConfidentialTransferMint:                           []byte{27, 0},
	InstructionTypeInitializeDefaultAccountState:                                []byte{28, 0},
	InstructionTypeInitializeGroupMemberPointer:                                 []byte{41, 0},
	InstructionTypeInitializeGroupPointer:                                       []byte{40, 0},
	InstructionTypeInitializeImmutableOwner:                                     []byte{22},
	InstructionTypeInitializeInterestBearingMint:                                []byte{33, 0},
	InstructionTypeInitializeMetadataPointer:                                    []byte{39, 0},
	InstructionTypeInitializeMint:                                               []byte{0},
	InstructionTypeInitializeMint2:                                              []byte{20},
	InstructionTypeInitializeMintCloseAuthority:                                 []byte{25},
	InstructionTypeInitializeMultisig:                                           []byte{2},
	InstructionTypeInitializeMultisig2:                                          []byte{19},
	InstructionTypeInitializeNonTransferableMint:                                []byte{32},
	InstructionTypeInitializePausableConfig:                                     []byte{44, 0},
	InstructionTypeInitializePermanentDelegate:                                  []byte{35},
	InstructionTypeInitializePermissionedBurn:                                   []byte{46, 0},
	InstructionTypeInitializeScaledUIAmountMint:                                 []byte{43, 0},
	InstructionTypeInitializeTokenGroup:                                         []byte{121, 113, 108, 39, 54, 51, 0, 4},
	InstructionTypeInitializeTokenGroupMember:                                   []byte{152, 32, 222, 176, 223, 237, 116, 134},
	InstructionTypeInitializeTokenMetadata:                                      []byte{210, 225, 30, 162, 88, 184, 77, 141},
	InstructionTypeInitializeTransferFeeConfig:                                  []byte{26, 0},
	InstructionTypeInitializeTransferHook:                                       []byte{36, 0},
	InstructionTypeMintTo:                                                       []byte{7},
	InstructionTypeMintToChecked:                                                []byte{14},
	InstructionTypePause:                                                        []byte{44, 1},
	InstructionTypePermissionedBurn:                                             []byte{46, 1},
	InstructionTypePermissionedBurnChecked:                                      []byte{46, 2},
	InstructionTypeReallocate:                                                   []byte{29},
	InstructionTypeRemoveTokenMetadataKey:                                       []byte{234, 18, 32, 56, 89, 141, 37, 181},
	InstructionTypeResume:                                                       []byte{44, 2},
	InstructionTypeRevoke:                                                       []byte{5},
	InstructionTypeSetAuthority:                                                 []byte{6},
	InstructionTypeSetTransferFee:                                               []byte{26, 5},
	InstructionTypeSyncNative:                                                   []byte{17},
	InstructionTypeThawAccount:                                                  []byte{11},
	InstructionTypeTransfer:                                                     []byte{3},
	InstructionTypeTransferChecked:                                              []byte{12},
	InstructionTypeTransferCheckedWithFee:                                       []byte{26, 1},
	InstructionTypeUIAmountToAmount:                                             []byte{24},
	InstructionTypeUnwrapLamports:                                               []byte{45},
	InstructionTypeUpdateConfidentialTransferMint:                               []byte{27, 1},
	InstructionTypeUpdateDefaultAccountState:                                    []byte{28, 1},
	InstructionTypeUpdateGroupMemberPointer:                                     []byte{41, 1},
	InstructionTypeUpdateGroupPointer:                                           []byte{40, 1},
	InstructionTypeUpdateMetadataPointer:                                        []byte{39, 1},
	InstructionTypeUpdateMultiplierScaledUIMint:                                 []byte{43, 1},
	InstructionTypeUpdateRateInterestBearingMint:                                []byte{33, 1},
	InstructionTypeUpdateTokenGroupMaxSize:                                      []byte{108, 37, 171, 143, 248, 30, 18, 110},
	InstructionTypeUpdateTokenGroupUpdateAuthority:                              []byte{161, 105, 88, 1, 237, 221, 216, 203},
	InstructionTypeUpdateTokenMetadataField:                                     []byte{221, 233, 49, 45, 181, 202, 220, 200},
	InstructionTypeUpdateTokenMetadataUpdateAuthority:                           []byte{215, 228, 166, 228, 84, 100, 86, 123},
	InstructionTypeUpdateTransferHook:                                           []byte{36, 1},
	InstructionTypeWithdrawExcessLamports:                                       []byte{38},
	InstructionTypeWithdrawWithheldTokensFromAccounts:                           []byte{26, 3},
	InstructionTypeWithdrawWithheldTokensFromAccountsForConfidentialTransferFee: []byte{37, 2},
	InstructionTypeWithdrawWithheldTokensFromMint:                               []byte{26, 2},
	InstructionTypeWithdrawWithheldTokensFromMintForConfidentialTransferFee:     []byte{37, 1},
}

var instructionTypeByDiscriminator = func() map[string]InstructionType {
	out := make(map[string]InstructionType, len(instructionDiscriminatorByType))
	for instructionType, discriminator := range instructionDiscriminatorByType {
		out[string(discriminator)] = instructionType
	}
	return out
}()

var instructionDiscriminatorPrefixLengths = []int{8, 2, 1}

func InstructionDiscriminator(instructionType InstructionType) ([]byte, error) {
	discriminator, ok := instructionDiscriminatorByType[instructionType]
	if !ok {
		return nil, fmt.Errorf("unknown token2022 instruction type: %s", instructionType)
	}
	return append([]byte(nil), discriminator...), nil
}

func instructionDiscriminatorLenByType(instructionType InstructionType) (int, error) {
	discriminator, ok := instructionDiscriminatorByType[instructionType]
	if !ok {
		return 0, fmt.Errorf("unknown token2022 instruction type: %s", instructionType)
	}
	return len(discriminator), nil
}

func encodeInstructionTypeDiscriminator(enc *binary.Encoder, instructionType InstructionType) error {
	discriminator, ok := instructionDiscriminatorByType[instructionType]
	if !ok {
		return fmt.Errorf("unknown token2022 instruction type: %s", instructionType)
	}
	return enc.WriteBytes(discriminator, false)
}

func identifyInstructionTypeByDiscriminator(data []byte) (InstructionType, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty instruction data")
	}
	for _, prefixLen := range instructionDiscriminatorPrefixLengths {
		if len(data) < prefixLen {
			continue
		}
		if instructionType, ok := instructionTypeByDiscriminator[string(data[:prefixLen])]; ok {
			return instructionType, nil
		}
	}
	return "", fmt.Errorf("unknown token2022 instruction discriminator")
}
