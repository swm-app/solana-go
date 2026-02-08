// Code generated from solana-program/token-2022/interface/idl.json. DO NOT EDIT.
package token2022

import (
	"bytes"
	"fmt"
	"reflect"

	binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

type InstructionType string

const (
	InstructionTypeInitializeMint                                               InstructionType = "initializeMint"
	InstructionTypeInitializeAccount                                            InstructionType = "initializeAccount"
	InstructionTypeInitializeMultisig                                           InstructionType = "initializeMultisig"
	InstructionTypeTransfer                                                     InstructionType = "transfer"
	InstructionTypeApprove                                                      InstructionType = "approve"
	InstructionTypeRevoke                                                       InstructionType = "revoke"
	InstructionTypeSetAuthority                                                 InstructionType = "setAuthority"
	InstructionTypeMintTo                                                       InstructionType = "mintTo"
	InstructionTypeBurn                                                         InstructionType = "burn"
	InstructionTypeCloseAccount                                                 InstructionType = "closeAccount"
	InstructionTypeFreezeAccount                                                InstructionType = "freezeAccount"
	InstructionTypeThawAccount                                                  InstructionType = "thawAccount"
	InstructionTypeTransferChecked                                              InstructionType = "transferChecked"
	InstructionTypeApproveChecked                                               InstructionType = "approveChecked"
	InstructionTypeMintToChecked                                                InstructionType = "mintToChecked"
	InstructionTypeBurnChecked                                                  InstructionType = "burnChecked"
	InstructionTypeInitializeAccount2                                           InstructionType = "initializeAccount2"
	InstructionTypeSyncNative                                                   InstructionType = "syncNative"
	InstructionTypeInitializeAccount3                                           InstructionType = "initializeAccount3"
	InstructionTypeInitializeMultisig2                                          InstructionType = "initializeMultisig2"
	InstructionTypeInitializeMint2                                              InstructionType = "initializeMint2"
	InstructionTypeGetAccountDataSize                                           InstructionType = "getAccountDataSize"
	InstructionTypeInitializeImmutableOwner                                     InstructionType = "initializeImmutableOwner"
	InstructionTypeAmountToUIAmount                                             InstructionType = "amountToUiAmount"
	InstructionTypeUIAmountToAmount                                             InstructionType = "uiAmountToAmount"
	InstructionTypeInitializeMintCloseAuthority                                 InstructionType = "initializeMintCloseAuthority"
	InstructionTypeInitializeTransferFeeConfig                                  InstructionType = "initializeTransferFeeConfig"
	InstructionTypeTransferCheckedWithFee                                       InstructionType = "transferCheckedWithFee"
	InstructionTypeWithdrawWithheldTokensFromMint                               InstructionType = "withdrawWithheldTokensFromMint"
	InstructionTypeWithdrawWithheldTokensFromAccounts                           InstructionType = "withdrawWithheldTokensFromAccounts"
	InstructionTypeHarvestWithheldTokensToMint                                  InstructionType = "harvestWithheldTokensToMint"
	InstructionTypeSetTransferFee                                               InstructionType = "setTransferFee"
	InstructionTypeInitializeConfidentialTransferMint                           InstructionType = "initializeConfidentialTransferMint"
	InstructionTypeUpdateConfidentialTransferMint                               InstructionType = "updateConfidentialTransferMint"
	InstructionTypeConfigureConfidentialTransferAccount                         InstructionType = "configureConfidentialTransferAccount"
	InstructionTypeApproveConfidentialTransferAccount                           InstructionType = "approveConfidentialTransferAccount"
	InstructionTypeEmptyConfidentialTransferAccount                             InstructionType = "emptyConfidentialTransferAccount"
	InstructionTypeConfidentialDeposit                                          InstructionType = "confidentialDeposit"
	InstructionTypeConfidentialWithdraw                                         InstructionType = "confidentialWithdraw"
	InstructionTypeConfidentialTransfer                                         InstructionType = "confidentialTransfer"
	InstructionTypeApplyConfidentialPendingBalance                              InstructionType = "applyConfidentialPendingBalance"
	InstructionTypeEnableConfidentialCredits                                    InstructionType = "enableConfidentialCredits"
	InstructionTypeDisableConfidentialCredits                                   InstructionType = "disableConfidentialCredits"
	InstructionTypeEnableNonConfidentialCredits                                 InstructionType = "enableNonConfidentialCredits"
	InstructionTypeDisableNonConfidentialCredits                                InstructionType = "disableNonConfidentialCredits"
	InstructionTypeConfidentialTransferWithFee                                  InstructionType = "confidentialTransferWithFee"
	InstructionTypeInitializeDefaultAccountState                                InstructionType = "initializeDefaultAccountState"
	InstructionTypeUpdateDefaultAccountState                                    InstructionType = "updateDefaultAccountState"
	InstructionTypeReallocate                                                   InstructionType = "reallocate"
	InstructionTypeEnableMemoTransfers                                          InstructionType = "enableMemoTransfers"
	InstructionTypeDisableMemoTransfers                                         InstructionType = "disableMemoTransfers"
	InstructionTypeCreateNativeMint                                             InstructionType = "createNativeMint"
	InstructionTypeInitializeNonTransferableMint                                InstructionType = "initializeNonTransferableMint"
	InstructionTypeInitializeInterestBearingMint                                InstructionType = "initializeInterestBearingMint"
	InstructionTypeUpdateRateInterestBearingMint                                InstructionType = "updateRateInterestBearingMint"
	InstructionTypeEnableCpiGuard                                               InstructionType = "enableCpiGuard"
	InstructionTypeDisableCpiGuard                                              InstructionType = "disableCpiGuard"
	InstructionTypeInitializePermanentDelegate                                  InstructionType = "initializePermanentDelegate"
	InstructionTypeInitializeTransferHook                                       InstructionType = "initializeTransferHook"
	InstructionTypeUpdateTransferHook                                           InstructionType = "updateTransferHook"
	InstructionTypeInitializeConfidentialTransferFee                            InstructionType = "initializeConfidentialTransferFee"
	InstructionTypeWithdrawWithheldTokensFromMintForConfidentialTransferFee     InstructionType = "withdrawWithheldTokensFromMintForConfidentialTransferFee"
	InstructionTypeWithdrawWithheldTokensFromAccountsForConfidentialTransferFee InstructionType = "withdrawWithheldTokensFromAccountsForConfidentialTransferFee"
	InstructionTypeHarvestWithheldTokensToMintForConfidentialTransferFee        InstructionType = "harvestWithheldTokensToMintForConfidentialTransferFee"
	InstructionTypeEnableHarvestToMint                                          InstructionType = "enableHarvestToMint"
	InstructionTypeDisableHarvestToMint                                         InstructionType = "disableHarvestToMint"
	InstructionTypeWithdrawExcessLamports                                       InstructionType = "withdrawExcessLamports"
	InstructionTypeInitializeMetadataPointer                                    InstructionType = "initializeMetadataPointer"
	InstructionTypeUpdateMetadataPointer                                        InstructionType = "updateMetadataPointer"
	InstructionTypeInitializeGroupPointer                                       InstructionType = "initializeGroupPointer"
	InstructionTypeUpdateGroupPointer                                           InstructionType = "updateGroupPointer"
	InstructionTypeInitializeGroupMemberPointer                                 InstructionType = "initializeGroupMemberPointer"
	InstructionTypeUpdateGroupMemberPointer                                     InstructionType = "updateGroupMemberPointer"
	InstructionTypeInitializeScaledUIAmountMint                                 InstructionType = "initializeScaledUiAmountMint"
	InstructionTypeUpdateMultiplierScaledUIMint                                 InstructionType = "updateMultiplierScaledUiMint"
	InstructionTypeInitializePausableConfig                                     InstructionType = "initializePausableConfig"
	InstructionTypePause                                                        InstructionType = "pause"
	InstructionTypeResume                                                       InstructionType = "resume"
	InstructionTypeInitializeTokenMetadata                                      InstructionType = "initializeTokenMetadata"
	InstructionTypeUpdateTokenMetadataField                                     InstructionType = "updateTokenMetadataField"
	InstructionTypeRemoveTokenMetadataKey                                       InstructionType = "removeTokenMetadataKey"
	InstructionTypeUpdateTokenMetadataUpdateAuthority                           InstructionType = "updateTokenMetadataUpdateAuthority"
	InstructionTypeEmitTokenMetadata                                            InstructionType = "emitTokenMetadata"
	InstructionTypeInitializeTokenGroup                                         InstructionType = "initializeTokenGroup"
	InstructionTypeUpdateTokenGroupMaxSize                                      InstructionType = "updateTokenGroupMaxSize"
	InstructionTypeUpdateTokenGroupUpdateAuthority                              InstructionType = "updateTokenGroupUpdateAuthority"
	InstructionTypeInitializeTokenGroupMember                                   InstructionType = "initializeTokenGroupMember"
	InstructionTypeUnwrapLamports                                               InstructionType = "unwrapLamports"
	InstructionTypeInitializePermissionedBurn                                   InstructionType = "initializePermissionedBurn"
	InstructionTypePermissionedBurn                                             InstructionType = "permissionedBurn"
	InstructionTypePermissionedBurnChecked                                      InstructionType = "permissionedBurnChecked"
)

type InitializeMintParams struct {
	Decimals        uint8               `json:"decimals"`
	MintAuthority   solanago.PublicKey  `json:"mintAuthority"`
	FreezeAuthority *solanago.PublicKey `json:"freezeAuthority"`
}

type InitializeMintAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
	Rent solanago.PublicKey `json:"rent"`
}

func NewInitializeMintInstruction(params InitializeMintParams, accounts InitializeMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	if err := enc__.Encode(params.MintAuthority); err != nil {
		return nil, fmt.Errorf("failed to encode mintAuthority: %w", err)
	}
	if params.FreezeAuthority == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode freezeAuthority option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode freezeAuthority option prefix: %w", err)
		}
		if err := enc__.Encode(*params.FreezeAuthority); err != nil {
			return nil, fmt.Errorf("failed to encode freezeAuthority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	rentPK := accounts.Rent
	if rentPK.IsZero() {
		rentPK = solanago.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(rentPK, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeAccountParams struct {
}

type InitializeAccountAccounts struct {
	Account solanago.PublicKey `json:"account"`
	Mint    solanago.PublicKey `json:"mint"`
	Owner   solanago.PublicKey `json:"owner"`
	Rent    solanago.PublicKey `json:"rent"`
}

func NewInitializeAccountInstruction(params InitializeAccountParams, accounts InitializeAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, false))
	rentPK := accounts.Rent
	if rentPK.IsZero() {
		rentPK = solanago.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(rentPK, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeMultisigParams struct {
	M uint8 `json:"m"`
}

type InitializeMultisigAccounts struct {
	Multisig solanago.PublicKey   `json:"multisig"`
	Rent     solanago.PublicKey   `json:"rent"`
	Signers  []solanago.PublicKey `json:"signers"`
}

func NewInitializeMultisigInstruction(params InitializeMultisigParams, accounts InitializeMultisigAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeMultisig); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.M); err != nil {
		return nil, fmt.Errorf("failed to encode m: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Multisig.IsZero() {
		return nil, fmt.Errorf("account multisig is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Multisig, true, false))
	rentPK := accounts.Rent
	if rentPK.IsZero() {
		rentPK = solanago.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(rentPK, false, false))
	for _, pk := range accounts.Signers {
		accounts__.Append(solanago.NewAccountMeta(pk, false, false))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type TransferParams struct {
	Amount uint64 `json:"amount"`
}

type TransferAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Destination  solanago.PublicKey   `json:"destination"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewTransferInstruction(params TransferParams, accounts TransferAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeTransfer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ApproveParams struct {
	Amount uint64 `json:"amount"`
}

type ApproveAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Delegate     solanago.PublicKey   `json:"delegate"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewApproveInstruction(params ApproveParams, accounts ApproveAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeApprove); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Delegate.IsZero() {
		return nil, fmt.Errorf("account delegate is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Delegate, false, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type RevokeParams struct {
}

type RevokeAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewRevokeInstruction(params RevokeParams, accounts RevokeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeRevoke); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type SetAuthorityParams struct {
	AuthorityType AuthorityType       `json:"authorityType"`
	NewAuthority  *solanago.PublicKey `json:"newAuthority"`
}

type SetAuthorityAccounts struct {
	Owned        solanago.PublicKey   `json:"owned"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewSetAuthorityInstruction(params SetAuthorityParams, accounts SetAuthorityAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeSetAuthority); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.AuthorityType); err != nil {
		return nil, fmt.Errorf("failed to encode authorityType: %w", err)
	}
	if params.NewAuthority == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode newAuthority option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode newAuthority option prefix: %w", err)
		}
		if err := enc__.Encode(*params.NewAuthority); err != nil {
			return nil, fmt.Errorf("failed to encode newAuthority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Owned.IsZero() {
		return nil, fmt.Errorf("account owned is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owned, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type MintToParams struct {
	Amount uint64 `json:"amount"`
}

type MintToAccounts struct {
	Mint          solanago.PublicKey   `json:"mint"`
	Token         solanago.PublicKey   `json:"token"`
	MintAuthority solanago.PublicKey   `json:"mintAuthority"`
	MultiSigners  []solanago.PublicKey `json:"multiSigners"`
}

func NewMintToInstruction(params MintToParams, accounts MintToAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeMintTo); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.MintAuthority.IsZero() {
		return nil, fmt.Errorf("account mintAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MintAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type BurnParams struct {
	Amount uint64 `json:"amount"`
}

type BurnAccounts struct {
	Account      solanago.PublicKey   `json:"account"`
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewBurnInstruction(params BurnParams, accounts BurnAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeBurn); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type CloseAccountParams struct {
}

type CloseAccountAccounts struct {
	Account      solanago.PublicKey   `json:"account"`
	Destination  solanago.PublicKey   `json:"destination"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewCloseAccountInstruction(params CloseAccountParams, accounts CloseAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeCloseAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type FreezeAccountParams struct {
}

type FreezeAccountAccounts struct {
	Account      solanago.PublicKey   `json:"account"`
	Mint         solanago.PublicKey   `json:"mint"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewFreezeAccountInstruction(params FreezeAccountParams, accounts FreezeAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeFreezeAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ThawAccountParams struct {
}

type ThawAccountAccounts struct {
	Account      solanago.PublicKey   `json:"account"`
	Mint         solanago.PublicKey   `json:"mint"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewThawAccountInstruction(params ThawAccountParams, accounts ThawAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeThawAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type TransferCheckedParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
}

type TransferCheckedAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Mint         solanago.PublicKey   `json:"mint"`
	Destination  solanago.PublicKey   `json:"destination"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewTransferCheckedInstruction(params TransferCheckedParams, accounts TransferCheckedAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeTransferChecked); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ApproveCheckedParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
}

type ApproveCheckedAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Mint         solanago.PublicKey   `json:"mint"`
	Delegate     solanago.PublicKey   `json:"delegate"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewApproveCheckedInstruction(params ApproveCheckedParams, accounts ApproveCheckedAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeApproveChecked); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Delegate.IsZero() {
		return nil, fmt.Errorf("account delegate is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Delegate, false, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type MintToCheckedParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
}

type MintToCheckedAccounts struct {
	Mint          solanago.PublicKey   `json:"mint"`
	Token         solanago.PublicKey   `json:"token"`
	MintAuthority solanago.PublicKey   `json:"mintAuthority"`
	MultiSigners  []solanago.PublicKey `json:"multiSigners"`
}

func NewMintToCheckedInstruction(params MintToCheckedParams, accounts MintToCheckedAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeMintToChecked); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.MintAuthority.IsZero() {
		return nil, fmt.Errorf("account mintAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MintAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type BurnCheckedParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
}

type BurnCheckedAccounts struct {
	Account      solanago.PublicKey   `json:"account"`
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewBurnCheckedInstruction(params BurnCheckedParams, accounts BurnCheckedAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeBurnChecked); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeAccount2Params struct {
	Owner solanago.PublicKey `json:"owner"`
}

type InitializeAccount2Accounts struct {
	Account solanago.PublicKey `json:"account"`
	Mint    solanago.PublicKey `json:"mint"`
	Rent    solanago.PublicKey `json:"rent"`
}

func NewInitializeAccount2Instruction(params InitializeAccount2Params, accounts InitializeAccount2Accounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeAccount2); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Owner); err != nil {
		return nil, fmt.Errorf("failed to encode owner: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	rentPK := accounts.Rent
	if rentPK.IsZero() {
		rentPK = solanago.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(rentPK, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type SyncNativeParams struct {
}

type SyncNativeAccounts struct {
	Account solanago.PublicKey `json:"account"`
}

func NewSyncNativeInstruction(params SyncNativeParams, accounts SyncNativeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeSyncNative); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeAccount3Params struct {
	Owner solanago.PublicKey `json:"owner"`
}

type InitializeAccount3Accounts struct {
	Account solanago.PublicKey `json:"account"`
	Mint    solanago.PublicKey `json:"mint"`
}

func NewInitializeAccount3Instruction(params InitializeAccount3Params, accounts InitializeAccount3Accounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeAccount3); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Owner); err != nil {
		return nil, fmt.Errorf("failed to encode owner: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeMultisig2Params struct {
	M uint8 `json:"m"`
}

type InitializeMultisig2Accounts struct {
	Multisig solanago.PublicKey   `json:"multisig"`
	Signers  []solanago.PublicKey `json:"signers"`
}

func NewInitializeMultisig2Instruction(params InitializeMultisig2Params, accounts InitializeMultisig2Accounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeMultisig2); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.M); err != nil {
		return nil, fmt.Errorf("failed to encode m: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Multisig.IsZero() {
		return nil, fmt.Errorf("account multisig is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Multisig, true, false))
	for _, pk := range accounts.Signers {
		accounts__.Append(solanago.NewAccountMeta(pk, false, false))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeMint2Params struct {
	Decimals        uint8               `json:"decimals"`
	MintAuthority   solanago.PublicKey  `json:"mintAuthority"`
	FreezeAuthority *solanago.PublicKey `json:"freezeAuthority"`
}

type InitializeMint2Accounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeMint2Instruction(params InitializeMint2Params, accounts InitializeMint2Accounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeMint2); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	if err := enc__.Encode(params.MintAuthority); err != nil {
		return nil, fmt.Errorf("failed to encode mintAuthority: %w", err)
	}
	if params.FreezeAuthority == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode freezeAuthority option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode freezeAuthority option prefix: %w", err)
		}
		if err := enc__.Encode(*params.FreezeAuthority); err != nil {
			return nil, fmt.Errorf("failed to encode freezeAuthority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type GetAccountDataSizeParams struct {
}

type GetAccountDataSizeAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewGetAccountDataSizeInstruction(params GetAccountDataSizeParams, accounts GetAccountDataSizeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeGetAccountDataSize); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeImmutableOwnerParams struct {
}

type InitializeImmutableOwnerAccounts struct {
	Account solanago.PublicKey `json:"account"`
}

func NewInitializeImmutableOwnerInstruction(params InitializeImmutableOwnerParams, accounts InitializeImmutableOwnerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeImmutableOwner); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type AmountToUIAmountParams struct {
	Amount uint64 `json:"amount"`
}

type AmountToUIAmountAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewAmountToUIAmountInstruction(params AmountToUIAmountParams, accounts AmountToUIAmountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeAmountToUIAmount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UIAmountToAmountParams struct {
	UIAmount string `json:"uiAmount"`
}

type UIAmountToAmountAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewUIAmountToAmountInstruction(params UIAmountToAmountParams, accounts UIAmountToAmountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUIAmountToAmount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.UIAmount); err != nil {
		return nil, fmt.Errorf("failed to encode uiAmount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeMintCloseAuthorityParams struct {
	CloseAuthority *solanago.PublicKey `json:"closeAuthority"`
}

type InitializeMintCloseAuthorityAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeMintCloseAuthorityInstruction(params InitializeMintCloseAuthorityParams, accounts InitializeMintCloseAuthorityAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeMintCloseAuthority); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.CloseAuthority == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode closeAuthority option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode closeAuthority option prefix: %w", err)
		}
		if err := enc__.Encode(*params.CloseAuthority); err != nil {
			return nil, fmt.Errorf("failed to encode closeAuthority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeTransferFeeConfigParams struct {
	TransferFeeConfigAuthority *solanago.PublicKey `json:"transferFeeConfigAuthority"`
	WithdrawWithheldAuthority  *solanago.PublicKey `json:"withdrawWithheldAuthority"`
	TransferFeeBasisPoints     uint16              `json:"transferFeeBasisPoints"`
	MaximumFee                 uint64              `json:"maximumFee"`
}

type InitializeTransferFeeConfigAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeTransferFeeConfigInstruction(params InitializeTransferFeeConfigParams, accounts InitializeTransferFeeConfigAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeTransferFeeConfig); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.TransferFeeConfigAuthority == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode transferFeeConfigAuthority option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode transferFeeConfigAuthority option prefix: %w", err)
		}
		if err := enc__.Encode(*params.TransferFeeConfigAuthority); err != nil {
			return nil, fmt.Errorf("failed to encode transferFeeConfigAuthority: %w", err)
		}
	}
	if params.WithdrawWithheldAuthority == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode withdrawWithheldAuthority option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode withdrawWithheldAuthority option prefix: %w", err)
		}
		if err := enc__.Encode(*params.WithdrawWithheldAuthority); err != nil {
			return nil, fmt.Errorf("failed to encode withdrawWithheldAuthority: %w", err)
		}
	}
	if err := enc__.Encode(params.TransferFeeBasisPoints); err != nil {
		return nil, fmt.Errorf("failed to encode transferFeeBasisPoints: %w", err)
	}
	if err := enc__.Encode(params.MaximumFee); err != nil {
		return nil, fmt.Errorf("failed to encode maximumFee: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type TransferCheckedWithFeeParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
	Fee      uint64 `json:"fee"`
}

type TransferCheckedWithFeeAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Mint         solanago.PublicKey   `json:"mint"`
	Destination  solanago.PublicKey   `json:"destination"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewTransferCheckedWithFeeInstruction(params TransferCheckedWithFeeParams, accounts TransferCheckedWithFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeTransferCheckedWithFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	if err := enc__.Encode(params.Fee); err != nil {
		return nil, fmt.Errorf("failed to encode fee: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type WithdrawWithheldTokensFromMintParams struct {
}

type WithdrawWithheldTokensFromMintAccounts struct {
	Mint                      solanago.PublicKey   `json:"mint"`
	FeeReceiver               solanago.PublicKey   `json:"feeReceiver"`
	WithdrawWithheldAuthority solanago.PublicKey   `json:"withdrawWithheldAuthority"`
	MultiSigners              []solanago.PublicKey `json:"multiSigners"`
}

func NewWithdrawWithheldTokensFromMintInstruction(params WithdrawWithheldTokensFromMintParams, accounts WithdrawWithheldTokensFromMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeWithdrawWithheldTokensFromMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.FeeReceiver.IsZero() {
		return nil, fmt.Errorf("account feeReceiver is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.FeeReceiver, true, false))
	if accounts.WithdrawWithheldAuthority.IsZero() {
		return nil, fmt.Errorf("account withdrawWithheldAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.WithdrawWithheldAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type WithdrawWithheldTokensFromAccountsParams struct {
	NumTokenAccounts uint8 `json:"numTokenAccounts"`
}

type WithdrawWithheldTokensFromAccountsAccounts struct {
	Mint                      solanago.PublicKey   `json:"mint"`
	FeeReceiver               solanago.PublicKey   `json:"feeReceiver"`
	WithdrawWithheldAuthority solanago.PublicKey   `json:"withdrawWithheldAuthority"`
	MultiSigners              []solanago.PublicKey `json:"multiSigners"`
	Sources                   []solanago.PublicKey `json:"sources"`
}

func NewWithdrawWithheldTokensFromAccountsInstruction(params WithdrawWithheldTokensFromAccountsParams, accounts WithdrawWithheldTokensFromAccountsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeWithdrawWithheldTokensFromAccounts); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.NumTokenAccounts); err != nil {
		return nil, fmt.Errorf("failed to encode numTokenAccounts: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.FeeReceiver.IsZero() {
		return nil, fmt.Errorf("account feeReceiver is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.FeeReceiver, true, false))
	if accounts.WithdrawWithheldAuthority.IsZero() {
		return nil, fmt.Errorf("account withdrawWithheldAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.WithdrawWithheldAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	for _, pk := range accounts.Sources {
		accounts__.Append(solanago.NewAccountMeta(pk, false, false))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type HarvestWithheldTokensToMintParams struct {
}

type HarvestWithheldTokensToMintAccounts struct {
	Mint    solanago.PublicKey   `json:"mint"`
	Sources []solanago.PublicKey `json:"sources"`
}

func NewHarvestWithheldTokensToMintInstruction(params HarvestWithheldTokensToMintParams, accounts HarvestWithheldTokensToMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeHarvestWithheldTokensToMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	for _, pk := range accounts.Sources {
		accounts__.Append(solanago.NewAccountMeta(pk, false, false))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type SetTransferFeeParams struct {
	TransferFeeBasisPoints uint16 `json:"transferFeeBasisPoints"`
	MaximumFee             uint64 `json:"maximumFee"`
}

type SetTransferFeeAccounts struct {
	Mint                       solanago.PublicKey   `json:"mint"`
	TransferFeeConfigAuthority solanago.PublicKey   `json:"transferFeeConfigAuthority"`
	MultiSigners               []solanago.PublicKey `json:"multiSigners"`
}

func NewSetTransferFeeInstruction(params SetTransferFeeParams, accounts SetTransferFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeSetTransferFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.TransferFeeBasisPoints); err != nil {
		return nil, fmt.Errorf("failed to encode transferFeeBasisPoints: %w", err)
	}
	if err := enc__.Encode(params.MaximumFee); err != nil {
		return nil, fmt.Errorf("failed to encode maximumFee: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.TransferFeeConfigAuthority.IsZero() {
		return nil, fmt.Errorf("account transferFeeConfigAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.TransferFeeConfigAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeConfidentialTransferMintParams struct {
	Authority              *solanago.PublicKey `json:"authority"`
	AutoApproveNewAccounts bool                `json:"autoApproveNewAccounts"`
	AuditorElgamalPubkey   *solanago.PublicKey `json:"auditorElgamalPubkey"`
}

type InitializeConfidentialTransferMintAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeConfidentialTransferMintInstruction(params InitializeConfidentialTransferMintParams, accounts InitializeConfidentialTransferMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeConfidentialTransferMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if err := enc__.Encode(params.AutoApproveNewAccounts); err != nil {
		return nil, fmt.Errorf("failed to encode autoApproveNewAccounts: %w", err)
	}
	if params.AuditorElgamalPubkey == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode auditorElgamalPubkey: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.AuditorElgamalPubkey[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode auditorElgamalPubkey: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateConfidentialTransferMintParams struct {
	AutoApproveNewAccounts bool                `json:"autoApproveNewAccounts"`
	AuditorElgamalPubkey   *solanago.PublicKey `json:"auditorElgamalPubkey"`
}

type UpdateConfidentialTransferMintAccounts struct {
	Mint      solanago.PublicKey `json:"mint"`
	Authority solanago.PublicKey `json:"authority"`
}

func NewUpdateConfidentialTransferMintInstruction(params UpdateConfidentialTransferMintParams, accounts UpdateConfidentialTransferMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateConfidentialTransferMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.AutoApproveNewAccounts); err != nil {
		return nil, fmt.Errorf("failed to encode autoApproveNewAccounts: %w", err)
	}
	if params.AuditorElgamalPubkey == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode auditorElgamalPubkey: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.AuditorElgamalPubkey[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode auditorElgamalPubkey: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ConfigureConfidentialTransferAccountParams struct {
	DecryptableZeroBalance             DecryptableBalance `json:"decryptableZeroBalance"`
	MaximumPendingBalanceCreditCounter uint64             `json:"maximumPendingBalanceCreditCounter"`
	ProofInstructionOffset             int8               `json:"proofInstructionOffset"`
}

type ConfigureConfidentialTransferAccountAccounts struct {
	Token                            solanago.PublicKey   `json:"token"`
	Mint                             solanago.PublicKey   `json:"mint"`
	InstructionsSysvarOrContextState solanago.PublicKey   `json:"instructionsSysvarOrContextState"`
	Record                           *solanago.PublicKey  `json:"record"`
	Authority                        solanago.PublicKey   `json:"authority"`
	MultiSigners                     []solanago.PublicKey `json:"multiSigners"`
}

func NewConfigureConfidentialTransferAccountInstruction(params ConfigureConfidentialTransferAccountParams, accounts ConfigureConfidentialTransferAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeConfigureConfidentialTransferAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.DecryptableZeroBalance); err != nil {
		return nil, fmt.Errorf("failed to encode decryptableZeroBalance: %w", err)
	}
	if err := enc__.Encode(params.MaximumPendingBalanceCreditCounter); err != nil {
		return nil, fmt.Errorf("failed to encode maximumPendingBalanceCreditCounter: %w", err)
	}
	if err := enc__.Encode(params.ProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode proofInstructionOffset: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	instructionsSysvarOrContextStatePK := accounts.InstructionsSysvarOrContextState
	if instructionsSysvarOrContextStatePK.IsZero() {
		instructionsSysvarOrContextStatePK = solanago.MustPublicKeyFromBase58("Sysvar1nstructions1111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(instructionsSysvarOrContextStatePK, false, false))
	if accounts.Record != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.Record, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ApproveConfidentialTransferAccountParams struct {
}

type ApproveConfidentialTransferAccountAccounts struct {
	Token     solanago.PublicKey `json:"token"`
	Mint      solanago.PublicKey `json:"mint"`
	Authority solanago.PublicKey `json:"authority"`
}

func NewApproveConfidentialTransferAccountInstruction(params ApproveConfidentialTransferAccountParams, accounts ApproveConfidentialTransferAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeApproveConfidentialTransferAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EmptyConfidentialTransferAccountParams struct {
	ProofInstructionOffset int8 `json:"proofInstructionOffset"`
}

type EmptyConfidentialTransferAccountAccounts struct {
	Token                            solanago.PublicKey   `json:"token"`
	InstructionsSysvarOrContextState solanago.PublicKey   `json:"instructionsSysvarOrContextState"`
	Record                           *solanago.PublicKey  `json:"record"`
	Authority                        solanago.PublicKey   `json:"authority"`
	MultiSigners                     []solanago.PublicKey `json:"multiSigners"`
}

func NewEmptyConfidentialTransferAccountInstruction(params EmptyConfidentialTransferAccountParams, accounts EmptyConfidentialTransferAccountAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEmptyConfidentialTransferAccount); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.ProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode proofInstructionOffset: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	instructionsSysvarOrContextStatePK := accounts.InstructionsSysvarOrContextState
	if instructionsSysvarOrContextStatePK.IsZero() {
		instructionsSysvarOrContextStatePK = solanago.MustPublicKeyFromBase58("Sysvar1nstructions1111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(instructionsSysvarOrContextStatePK, false, false))
	if accounts.Record != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.Record, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ConfidentialDepositParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
}

type ConfidentialDepositAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewConfidentialDepositInstruction(params ConfidentialDepositParams, accounts ConfidentialDepositAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeConfidentialDeposit); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ConfidentialWithdrawParams struct {
	Amount                         uint64             `json:"amount"`
	Decimals                       uint8              `json:"decimals"`
	NewDecryptableAvailableBalance DecryptableBalance `json:"newDecryptableAvailableBalance"`
	EqualityProofInstructionOffset int8               `json:"equalityProofInstructionOffset"`
	RangeProofInstructionOffset    int8               `json:"rangeProofInstructionOffset"`
}

type ConfidentialWithdrawAccounts struct {
	Token              solanago.PublicKey   `json:"token"`
	Mint               solanago.PublicKey   `json:"mint"`
	InstructionsSysvar *solanago.PublicKey  `json:"instructionsSysvar"`
	EqualityRecord     *solanago.PublicKey  `json:"equalityRecord"`
	RangeRecord        *solanago.PublicKey  `json:"rangeRecord"`
	Authority          solanago.PublicKey   `json:"authority"`
	MultiSigners       []solanago.PublicKey `json:"multiSigners"`
}

func NewConfidentialWithdrawInstruction(params ConfidentialWithdrawParams, accounts ConfidentialWithdrawAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeConfidentialWithdraw); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	if err := enc__.Encode(params.NewDecryptableAvailableBalance); err != nil {
		return nil, fmt.Errorf("failed to encode newDecryptableAvailableBalance: %w", err)
	}
	if err := enc__.Encode(params.EqualityProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode equalityProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.RangeProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode rangeProofInstructionOffset: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.InstructionsSysvar != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.InstructionsSysvar, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.EqualityRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.EqualityRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.RangeRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.RangeRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ConfidentialTransferParams struct {
	NewSourceDecryptableAvailableBalance     DecryptableBalance `json:"newSourceDecryptableAvailableBalance"`
	EqualityProofInstructionOffset           int8               `json:"equalityProofInstructionOffset"`
	CiphertextValidityProofInstructionOffset int8               `json:"ciphertextValidityProofInstructionOffset"`
	RangeProofInstructionOffset              int8               `json:"rangeProofInstructionOffset"`
}

type ConfidentialTransferAccounts struct {
	SourceToken              solanago.PublicKey   `json:"sourceToken"`
	Mint                     solanago.PublicKey   `json:"mint"`
	DestinationToken         solanago.PublicKey   `json:"destinationToken"`
	InstructionsSysvar       *solanago.PublicKey  `json:"instructionsSysvar"`
	EqualityRecord           *solanago.PublicKey  `json:"equalityRecord"`
	CiphertextValidityRecord *solanago.PublicKey  `json:"ciphertextValidityRecord"`
	RangeRecord              *solanago.PublicKey  `json:"rangeRecord"`
	Authority                solanago.PublicKey   `json:"authority"`
	MultiSigners             []solanago.PublicKey `json:"multiSigners"`
}

func NewConfidentialTransferInstruction(params ConfidentialTransferParams, accounts ConfidentialTransferAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeConfidentialTransfer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.NewSourceDecryptableAvailableBalance); err != nil {
		return nil, fmt.Errorf("failed to encode newSourceDecryptableAvailableBalance: %w", err)
	}
	if err := enc__.Encode(params.EqualityProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode equalityProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.CiphertextValidityProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode ciphertextValidityProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.RangeProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode rangeProofInstructionOffset: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.SourceToken.IsZero() {
		return nil, fmt.Errorf("account sourceToken is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.SourceToken, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.DestinationToken.IsZero() {
		return nil, fmt.Errorf("account destinationToken is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.DestinationToken, true, false))
	if accounts.InstructionsSysvar != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.InstructionsSysvar, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.EqualityRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.EqualityRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.CiphertextValidityRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.CiphertextValidityRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.RangeRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.RangeRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ApplyConfidentialPendingBalanceParams struct {
	ExpectedPendingBalanceCreditCounter uint64             `json:"expectedPendingBalanceCreditCounter"`
	NewDecryptableAvailableBalance      DecryptableBalance `json:"newDecryptableAvailableBalance"`
}

type ApplyConfidentialPendingBalanceAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewApplyConfidentialPendingBalanceInstruction(params ApplyConfidentialPendingBalanceParams, accounts ApplyConfidentialPendingBalanceAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeApplyConfidentialPendingBalance); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.ExpectedPendingBalanceCreditCounter); err != nil {
		return nil, fmt.Errorf("failed to encode expectedPendingBalanceCreditCounter: %w", err)
	}
	if err := enc__.Encode(params.NewDecryptableAvailableBalance); err != nil {
		return nil, fmt.Errorf("failed to encode newDecryptableAvailableBalance: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EnableConfidentialCreditsParams struct {
}

type EnableConfidentialCreditsAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewEnableConfidentialCreditsInstruction(params EnableConfidentialCreditsParams, accounts EnableConfidentialCreditsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEnableConfidentialCredits); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type DisableConfidentialCreditsParams struct {
}

type DisableConfidentialCreditsAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewDisableConfidentialCreditsInstruction(params DisableConfidentialCreditsParams, accounts DisableConfidentialCreditsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeDisableConfidentialCredits); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EnableNonConfidentialCreditsParams struct {
}

type EnableNonConfidentialCreditsAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewEnableNonConfidentialCreditsInstruction(params EnableNonConfidentialCreditsParams, accounts EnableNonConfidentialCreditsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEnableNonConfidentialCredits); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type DisableNonConfidentialCreditsParams struct {
}

type DisableNonConfidentialCreditsAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewDisableNonConfidentialCreditsInstruction(params DisableNonConfidentialCreditsParams, accounts DisableNonConfidentialCreditsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeDisableNonConfidentialCredits); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ConfidentialTransferWithFeeParams struct {
	NewSourceDecryptableAvailableBalance                   DecryptableBalance `json:"newSourceDecryptableAvailableBalance"`
	EqualityProofInstructionOffset                         int8               `json:"equalityProofInstructionOffset"`
	TransferAmountCiphertextValidityProofInstructionOffset int8               `json:"transferAmountCiphertextValidityProofInstructionOffset"`
	FeeSigmaProofInstructionOffset                         int8               `json:"feeSigmaProofInstructionOffset"`
	FeeCiphertextValidityProofInstructionOffset            int8               `json:"feeCiphertextValidityProofInstructionOffset"`
	RangeProofInstructionOffset                            int8               `json:"rangeProofInstructionOffset"`
}

type ConfidentialTransferWithFeeAccounts struct {
	SourceToken                            solanago.PublicKey   `json:"sourceToken"`
	Mint                                   solanago.PublicKey   `json:"mint"`
	DestinationToken                       solanago.PublicKey   `json:"destinationToken"`
	InstructionsSysvar                     *solanago.PublicKey  `json:"instructionsSysvar"`
	EqualityRecord                         *solanago.PublicKey  `json:"equalityRecord"`
	TransferAmountCiphertextValidityRecord *solanago.PublicKey  `json:"transferAmountCiphertextValidityRecord"`
	FeeSigmaRecord                         *solanago.PublicKey  `json:"feeSigmaRecord"`
	FeeCiphertextValidityRecord            *solanago.PublicKey  `json:"feeCiphertextValidityRecord"`
	RangeRecord                            *solanago.PublicKey  `json:"rangeRecord"`
	Authority                              solanago.PublicKey   `json:"authority"`
	MultiSigners                           []solanago.PublicKey `json:"multiSigners"`
}

func NewConfidentialTransferWithFeeInstruction(params ConfidentialTransferWithFeeParams, accounts ConfidentialTransferWithFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeConfidentialTransferWithFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.NewSourceDecryptableAvailableBalance); err != nil {
		return nil, fmt.Errorf("failed to encode newSourceDecryptableAvailableBalance: %w", err)
	}
	if err := enc__.Encode(params.EqualityProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode equalityProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.TransferAmountCiphertextValidityProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode transferAmountCiphertextValidityProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.FeeSigmaProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode feeSigmaProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.FeeCiphertextValidityProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode feeCiphertextValidityProofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.RangeProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode rangeProofInstructionOffset: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.SourceToken.IsZero() {
		return nil, fmt.Errorf("account sourceToken is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.SourceToken, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.DestinationToken.IsZero() {
		return nil, fmt.Errorf("account destinationToken is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.DestinationToken, true, false))
	if accounts.InstructionsSysvar != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.InstructionsSysvar, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.EqualityRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.EqualityRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.TransferAmountCiphertextValidityRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.TransferAmountCiphertextValidityRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.FeeSigmaRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.FeeSigmaRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.FeeCiphertextValidityRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.FeeCiphertextValidityRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.RangeRecord != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.RangeRecord, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeDefaultAccountStateParams struct {
	State AccountState `json:"state"`
}

type InitializeDefaultAccountStateAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeDefaultAccountStateInstruction(params InitializeDefaultAccountStateParams, accounts InitializeDefaultAccountStateAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeDefaultAccountState); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.State); err != nil {
		return nil, fmt.Errorf("failed to encode state: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateDefaultAccountStateParams struct {
	State AccountState `json:"state"`
}

type UpdateDefaultAccountStateAccounts struct {
	Mint            solanago.PublicKey   `json:"mint"`
	FreezeAuthority solanago.PublicKey   `json:"freezeAuthority"`
	MultiSigners    []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateDefaultAccountStateInstruction(params UpdateDefaultAccountStateParams, accounts UpdateDefaultAccountStateAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateDefaultAccountState); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.State); err != nil {
		return nil, fmt.Errorf("failed to encode state: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.FreezeAuthority.IsZero() {
		return nil, fmt.Errorf("account freezeAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.FreezeAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ReallocateParams struct {
	NewExtensionTypes []ExtensionType `json:"newExtensionTypes"`
}

type ReallocateAccounts struct {
	Token         solanago.PublicKey   `json:"token"`
	Payer         solanago.PublicKey   `json:"payer"`
	SystemProgram solanago.PublicKey   `json:"systemProgram"`
	Owner         solanago.PublicKey   `json:"owner"`
	MultiSigners  []solanago.PublicKey `json:"multiSigners"`
}

func NewReallocateInstruction(params ReallocateParams, accounts ReallocateAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeReallocate); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	for _, item := range params.NewExtensionTypes {
		if err := enc__.Encode(item); err != nil {
			return nil, fmt.Errorf("failed to encode newExtensionTypes item: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Payer.IsZero() {
		return nil, fmt.Errorf("account payer is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Payer, true, true))
	systemProgramPK := accounts.SystemProgram
	if systemProgramPK.IsZero() {
		systemProgramPK = solanago.MustPublicKeyFromBase58("11111111111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(systemProgramPK, false, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EnableMemoTransfersParams struct {
}

type EnableMemoTransfersAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewEnableMemoTransfersInstruction(params EnableMemoTransfersParams, accounts EnableMemoTransfersAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEnableMemoTransfers); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type DisableMemoTransfersParams struct {
}

type DisableMemoTransfersAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewDisableMemoTransfersInstruction(params DisableMemoTransfersParams, accounts DisableMemoTransfersAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeDisableMemoTransfers); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type CreateNativeMintParams struct {
}

type CreateNativeMintAccounts struct {
	Payer         solanago.PublicKey `json:"payer"`
	NativeMint    solanago.PublicKey `json:"nativeMint"`
	SystemProgram solanago.PublicKey `json:"systemProgram"`
}

func NewCreateNativeMintInstruction(params CreateNativeMintParams, accounts CreateNativeMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeCreateNativeMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Payer.IsZero() {
		return nil, fmt.Errorf("account payer is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Payer, true, true))
	if accounts.NativeMint.IsZero() {
		return nil, fmt.Errorf("account nativeMint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.NativeMint, true, false))
	systemProgramPK := accounts.SystemProgram
	if systemProgramPK.IsZero() {
		systemProgramPK = solanago.MustPublicKeyFromBase58("11111111111111111111111111111111")
	}
	accounts__.Append(solanago.NewAccountMeta(systemProgramPK, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeNonTransferableMintParams struct {
}

type InitializeNonTransferableMintAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeNonTransferableMintInstruction(params InitializeNonTransferableMintParams, accounts InitializeNonTransferableMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeNonTransferableMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeInterestBearingMintParams struct {
	RateAuthority *solanago.PublicKey `json:"rateAuthority"`
	Rate          int16               `json:"rate"`
}

type InitializeInterestBearingMintAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeInterestBearingMintInstruction(params InitializeInterestBearingMintParams, accounts InitializeInterestBearingMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeInterestBearingMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.RateAuthority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode rateAuthority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.RateAuthority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode rateAuthority: %w", err)
		}
	}
	if err := enc__.Encode(params.Rate); err != nil {
		return nil, fmt.Errorf("failed to encode rate: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateRateInterestBearingMintParams struct {
	Rate int16 `json:"rate"`
}

type UpdateRateInterestBearingMintAccounts struct {
	Mint          solanago.PublicKey   `json:"mint"`
	RateAuthority solanago.PublicKey   `json:"rateAuthority"`
	MultiSigners  []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateRateInterestBearingMintInstruction(params UpdateRateInterestBearingMintParams, accounts UpdateRateInterestBearingMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateRateInterestBearingMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Rate); err != nil {
		return nil, fmt.Errorf("failed to encode rate: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.RateAuthority.IsZero() {
		return nil, fmt.Errorf("account rateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.RateAuthority, true, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EnableCpiGuardParams struct {
}

type EnableCpiGuardAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewEnableCpiGuardInstruction(params EnableCpiGuardParams, accounts EnableCpiGuardAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEnableCpiGuard); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type DisableCpiGuardParams struct {
}

type DisableCpiGuardAccounts struct {
	Token        solanago.PublicKey   `json:"token"`
	Owner        solanago.PublicKey   `json:"owner"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewDisableCpiGuardInstruction(params DisableCpiGuardParams, accounts DisableCpiGuardAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeDisableCpiGuard); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Token.IsZero() {
		return nil, fmt.Errorf("account token is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Token, true, false))
	if accounts.Owner.IsZero() {
		return nil, fmt.Errorf("account owner is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Owner, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializePermanentDelegateParams struct {
	Delegate solanago.PublicKey `json:"delegate"`
}

type InitializePermanentDelegateAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializePermanentDelegateInstruction(params InitializePermanentDelegateParams, accounts InitializePermanentDelegateAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializePermanentDelegate); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Delegate); err != nil {
		return nil, fmt.Errorf("failed to encode delegate: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeTransferHookParams struct {
	Authority *solanago.PublicKey `json:"authority"`
	ProgramID *solanago.PublicKey `json:"programId"`
}

type InitializeTransferHookAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeTransferHookInstruction(params InitializeTransferHookParams, accounts InitializeTransferHookAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeTransferHook); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if params.ProgramID == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode programId: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.ProgramID[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode programId: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateTransferHookParams struct {
	ProgramID *solanago.PublicKey `json:"programId"`
}

type UpdateTransferHookAccounts struct {
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateTransferHookInstruction(params UpdateTransferHookParams, accounts UpdateTransferHookAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateTransferHook); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.ProgramID == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode programId: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.ProgramID[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode programId: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeConfidentialTransferFeeParams struct {
	Authority                              *solanago.PublicKey `json:"authority"`
	WithdrawWithheldAuthorityElGamalPubkey *solanago.PublicKey `json:"withdrawWithheldAuthorityElGamalPubkey"`
}

type InitializeConfidentialTransferFeeAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeConfidentialTransferFeeInstruction(params InitializeConfidentialTransferFeeParams, accounts InitializeConfidentialTransferFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeConfidentialTransferFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if params.WithdrawWithheldAuthorityElGamalPubkey == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode withdrawWithheldAuthorityElGamalPubkey: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.WithdrawWithheldAuthorityElGamalPubkey[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode withdrawWithheldAuthorityElGamalPubkey: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type WithdrawWithheldTokensFromMintForConfidentialTransferFeeParams struct {
	ProofInstructionOffset         int8               `json:"proofInstructionOffset"`
	NewDecryptableAvailableBalance DecryptableBalance `json:"newDecryptableAvailableBalance"`
}

type WithdrawWithheldTokensFromMintForConfidentialTransferFeeAccounts struct {
	Mint                             solanago.PublicKey   `json:"mint"`
	Destination                      solanago.PublicKey   `json:"destination"`
	InstructionsSysvarOrContextState solanago.PublicKey   `json:"instructionsSysvarOrContextState"`
	Record                           *solanago.PublicKey  `json:"record"`
	Authority                        solanago.PublicKey   `json:"authority"`
	MultiSigners                     []solanago.PublicKey `json:"multiSigners"`
}

func NewWithdrawWithheldTokensFromMintForConfidentialTransferFeeInstruction(params WithdrawWithheldTokensFromMintForConfidentialTransferFeeParams, accounts WithdrawWithheldTokensFromMintForConfidentialTransferFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeWithdrawWithheldTokensFromMintForConfidentialTransferFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.ProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode proofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.NewDecryptableAvailableBalance); err != nil {
		return nil, fmt.Errorf("failed to encode newDecryptableAvailableBalance: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.InstructionsSysvarOrContextState.IsZero() {
		return nil, fmt.Errorf("account instructionsSysvarOrContextState is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.InstructionsSysvarOrContextState, false, false))
	if accounts.Record != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.Record, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type WithdrawWithheldTokensFromAccountsForConfidentialTransferFeeParams struct {
	NumTokenAccounts               uint8              `json:"numTokenAccounts"`
	ProofInstructionOffset         int8               `json:"proofInstructionOffset"`
	NewDecryptableAvailableBalance DecryptableBalance `json:"newDecryptableAvailableBalance"`
}

type WithdrawWithheldTokensFromAccountsForConfidentialTransferFeeAccounts struct {
	Mint                             solanago.PublicKey   `json:"mint"`
	Destination                      solanago.PublicKey   `json:"destination"`
	InstructionsSysvarOrContextState solanago.PublicKey   `json:"instructionsSysvarOrContextState"`
	Record                           *solanago.PublicKey  `json:"record"`
	Authority                        solanago.PublicKey   `json:"authority"`
	MultiSigners                     []solanago.PublicKey `json:"multiSigners"`
}

func NewWithdrawWithheldTokensFromAccountsForConfidentialTransferFeeInstruction(params WithdrawWithheldTokensFromAccountsForConfidentialTransferFeeParams, accounts WithdrawWithheldTokensFromAccountsForConfidentialTransferFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeWithdrawWithheldTokensFromAccountsForConfidentialTransferFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.NumTokenAccounts); err != nil {
		return nil, fmt.Errorf("failed to encode numTokenAccounts: %w", err)
	}
	if err := enc__.Encode(params.ProofInstructionOffset); err != nil {
		return nil, fmt.Errorf("failed to encode proofInstructionOffset: %w", err)
	}
	if err := enc__.Encode(params.NewDecryptableAvailableBalance); err != nil {
		return nil, fmt.Errorf("failed to encode newDecryptableAvailableBalance: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.InstructionsSysvarOrContextState.IsZero() {
		return nil, fmt.Errorf("account instructionsSysvarOrContextState is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.InstructionsSysvarOrContextState, false, false))
	if accounts.Record != nil {
		accounts__.Append(solanago.NewAccountMeta(*accounts.Record, false, false))
	} else {
		accounts__.Append(solanago.NewAccountMeta(ProgramID, false, false))
	}
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type HarvestWithheldTokensToMintForConfidentialTransferFeeParams struct {
}

type HarvestWithheldTokensToMintForConfidentialTransferFeeAccounts struct {
	Mint    solanago.PublicKey   `json:"mint"`
	Sources []solanago.PublicKey `json:"sources"`
}

func NewHarvestWithheldTokensToMintForConfidentialTransferFeeInstruction(params HarvestWithheldTokensToMintForConfidentialTransferFeeParams, accounts HarvestWithheldTokensToMintForConfidentialTransferFeeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeHarvestWithheldTokensToMintForConfidentialTransferFee); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	for _, pk := range accounts.Sources {
		accounts__.Append(solanago.NewAccountMeta(pk, false, false))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EnableHarvestToMintParams struct {
}

type EnableHarvestToMintAccounts struct {
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewEnableHarvestToMintInstruction(params EnableHarvestToMintParams, accounts EnableHarvestToMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEnableHarvestToMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type DisableHarvestToMintParams struct {
}

type DisableHarvestToMintAccounts struct {
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewDisableHarvestToMintInstruction(params DisableHarvestToMintParams, accounts DisableHarvestToMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeDisableHarvestToMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type WithdrawExcessLamportsParams struct {
}

type WithdrawExcessLamportsAccounts struct {
	SourceAccount      solanago.PublicKey   `json:"sourceAccount"`
	DestinationAccount solanago.PublicKey   `json:"destinationAccount"`
	Authority          solanago.PublicKey   `json:"authority"`
	MultiSigners       []solanago.PublicKey `json:"multiSigners"`
}

func NewWithdrawExcessLamportsInstruction(params WithdrawExcessLamportsParams, accounts WithdrawExcessLamportsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeWithdrawExcessLamports); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.SourceAccount.IsZero() {
		return nil, fmt.Errorf("account sourceAccount is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.SourceAccount, true, false))
	if accounts.DestinationAccount.IsZero() {
		return nil, fmt.Errorf("account destinationAccount is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.DestinationAccount, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeMetadataPointerParams struct {
	Authority       *solanago.PublicKey `json:"authority"`
	MetadataAddress *solanago.PublicKey `json:"metadataAddress"`
}

type InitializeMetadataPointerAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeMetadataPointerInstruction(params InitializeMetadataPointerParams, accounts InitializeMetadataPointerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeMetadataPointer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if params.MetadataAddress == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode metadataAddress: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.MetadataAddress[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode metadataAddress: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateMetadataPointerParams struct {
	MetadataAddress *solanago.PublicKey `json:"metadataAddress"`
}

type UpdateMetadataPointerAccounts struct {
	Mint                     solanago.PublicKey   `json:"mint"`
	MetadataPointerAuthority solanago.PublicKey   `json:"metadataPointerAuthority"`
	MultiSigners             []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateMetadataPointerInstruction(params UpdateMetadataPointerParams, accounts UpdateMetadataPointerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateMetadataPointer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.MetadataAddress == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode metadataAddress: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.MetadataAddress[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode metadataAddress: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.MetadataPointerAuthority.IsZero() {
		return nil, fmt.Errorf("account metadataPointerAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MetadataPointerAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeGroupPointerParams struct {
	Authority    *solanago.PublicKey `json:"authority"`
	GroupAddress *solanago.PublicKey `json:"groupAddress"`
}

type InitializeGroupPointerAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeGroupPointerInstruction(params InitializeGroupPointerParams, accounts InitializeGroupPointerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeGroupPointer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if params.GroupAddress == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode groupAddress: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.GroupAddress[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode groupAddress: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateGroupPointerParams struct {
	GroupAddress *solanago.PublicKey `json:"groupAddress"`
}

type UpdateGroupPointerAccounts struct {
	Mint                  solanago.PublicKey   `json:"mint"`
	GroupPointerAuthority solanago.PublicKey   `json:"groupPointerAuthority"`
	MultiSigners          []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateGroupPointerInstruction(params UpdateGroupPointerParams, accounts UpdateGroupPointerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateGroupPointer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.GroupAddress == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode groupAddress: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.GroupAddress[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode groupAddress: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.GroupPointerAuthority.IsZero() {
		return nil, fmt.Errorf("account groupPointerAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.GroupPointerAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeGroupMemberPointerParams struct {
	Authority     *solanago.PublicKey `json:"authority"`
	MemberAddress *solanago.PublicKey `json:"memberAddress"`
}

type InitializeGroupMemberPointerAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeGroupMemberPointerInstruction(params InitializeGroupMemberPointerParams, accounts InitializeGroupMemberPointerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeGroupMemberPointer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if params.MemberAddress == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode memberAddress: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.MemberAddress[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode memberAddress: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateGroupMemberPointerParams struct {
	MemberAddress *solanago.PublicKey `json:"memberAddress"`
}

type UpdateGroupMemberPointerAccounts struct {
	Mint                        solanago.PublicKey   `json:"mint"`
	GroupMemberPointerAuthority solanago.PublicKey   `json:"groupMemberPointerAuthority"`
	MultiSigners                []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateGroupMemberPointerInstruction(params UpdateGroupMemberPointerParams, accounts UpdateGroupMemberPointerAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateGroupMemberPointer); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.MemberAddress == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode memberAddress: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.MemberAddress[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode memberAddress: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.GroupMemberPointerAuthority.IsZero() {
		return nil, fmt.Errorf("account groupMemberPointerAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.GroupMemberPointerAuthority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeScaledUIAmountMintParams struct {
	Authority  *solanago.PublicKey `json:"authority"`
	Multiplier float64             `json:"multiplier"`
}

type InitializeScaledUIAmountMintAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializeScaledUIAmountMintInstruction(params InitializeScaledUIAmountMintParams, accounts InitializeScaledUIAmountMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeScaledUIAmountMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	if err := enc__.Encode(params.Multiplier); err != nil {
		return nil, fmt.Errorf("failed to encode multiplier: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateMultiplierScaledUIMintParams struct {
	Multiplier         float64 `json:"multiplier"`
	EffectiveTimestamp int64   `json:"effectiveTimestamp"`
}

type UpdateMultiplierScaledUIMintAccounts struct {
	Mint         solanago.PublicKey   `json:"mint"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewUpdateMultiplierScaledUIMintInstruction(params UpdateMultiplierScaledUIMintParams, accounts UpdateMultiplierScaledUIMintAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateMultiplierScaledUIMint); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Multiplier); err != nil {
		return nil, fmt.Errorf("failed to encode multiplier: %w", err)
	}
	if err := enc__.Encode(params.EffectiveTimestamp); err != nil {
		return nil, fmt.Errorf("failed to encode effectiveTimestamp: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, true, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializePausableConfigParams struct {
	Authority *solanago.PublicKey `json:"authority"`
}

type InitializePausableConfigAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializePausableConfigInstruction(params InitializePausableConfigParams, accounts InitializePausableConfigAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializePausableConfig); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Authority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.Authority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode authority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type PauseParams struct {
}

type PauseAccounts struct {
	Mint      solanago.PublicKey `json:"mint"`
	Authority solanago.PublicKey `json:"authority"`
}

func NewPauseInstruction(params PauseParams, accounts PauseAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypePause); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type ResumeParams struct {
}

type ResumeAccounts struct {
	Mint      solanago.PublicKey `json:"mint"`
	Authority solanago.PublicKey `json:"authority"`
}

func NewResumeInstruction(params ResumeParams, accounts ResumeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeResume); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeTokenMetadataParams struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Uri    string `json:"uri"`
}

type InitializeTokenMetadataAccounts struct {
	Metadata        solanago.PublicKey `json:"metadata"`
	UpdateAuthority solanago.PublicKey `json:"updateAuthority"`
	Mint            solanago.PublicKey `json:"mint"`
	MintAuthority   solanago.PublicKey `json:"mintAuthority"`
}

func NewInitializeTokenMetadataInstruction(params InitializeTokenMetadataParams, accounts InitializeTokenMetadataAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeTokenMetadata); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Name); err != nil {
		return nil, fmt.Errorf("failed to encode name: %w", err)
	}
	if err := enc__.Encode(params.Symbol); err != nil {
		return nil, fmt.Errorf("failed to encode symbol: %w", err)
	}
	if err := enc__.Encode(params.Uri); err != nil {
		return nil, fmt.Errorf("failed to encode uri: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Metadata.IsZero() {
		return nil, fmt.Errorf("account metadata is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Metadata, true, false))
	if accounts.UpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account updateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.UpdateAuthority, false, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.MintAuthority.IsZero() {
		return nil, fmt.Errorf("account mintAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MintAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateTokenMetadataFieldParams struct {
	Field TokenMetadataField `json:"field"`
	Value string             `json:"value"`
}

type UpdateTokenMetadataFieldAccounts struct {
	Metadata        solanago.PublicKey `json:"metadata"`
	UpdateAuthority solanago.PublicKey `json:"updateAuthority"`
}

func NewUpdateTokenMetadataFieldInstruction(params UpdateTokenMetadataFieldParams, accounts UpdateTokenMetadataFieldAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateTokenMetadataField); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Field); err != nil {
		return nil, fmt.Errorf("failed to encode field: %w", err)
	}
	if err := enc__.Encode(params.Value); err != nil {
		return nil, fmt.Errorf("failed to encode value: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Metadata.IsZero() {
		return nil, fmt.Errorf("account metadata is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Metadata, true, false))
	if accounts.UpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account updateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.UpdateAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type RemoveTokenMetadataKeyParams struct {
	IDempotent bool   `json:"idempotent"`
	Key        string `json:"key"`
}

type RemoveTokenMetadataKeyAccounts struct {
	Metadata        solanago.PublicKey `json:"metadata"`
	UpdateAuthority solanago.PublicKey `json:"updateAuthority"`
}

func NewRemoveTokenMetadataKeyInstruction(params RemoveTokenMetadataKeyParams, accounts RemoveTokenMetadataKeyAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeRemoveTokenMetadataKey); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.IDempotent); err != nil {
		return nil, fmt.Errorf("failed to encode idempotent: %w", err)
	}
	if err := enc__.Encode(params.Key); err != nil {
		return nil, fmt.Errorf("failed to encode key: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Metadata.IsZero() {
		return nil, fmt.Errorf("account metadata is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Metadata, true, false))
	if accounts.UpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account updateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.UpdateAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateTokenMetadataUpdateAuthorityParams struct {
	NewUpdateAuthority *solanago.PublicKey `json:"newUpdateAuthority"`
}

type UpdateTokenMetadataUpdateAuthorityAccounts struct {
	Metadata        solanago.PublicKey `json:"metadata"`
	UpdateAuthority solanago.PublicKey `json:"updateAuthority"`
}

func NewUpdateTokenMetadataUpdateAuthorityInstruction(params UpdateTokenMetadataUpdateAuthorityParams, accounts UpdateTokenMetadataUpdateAuthorityAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateTokenMetadataUpdateAuthority); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.NewUpdateAuthority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode newUpdateAuthority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.NewUpdateAuthority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode newUpdateAuthority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Metadata.IsZero() {
		return nil, fmt.Errorf("account metadata is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Metadata, true, false))
	if accounts.UpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account updateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.UpdateAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type EmitTokenMetadataParams struct {
	Start *uint64 `json:"start"`
	End   *uint64 `json:"end"`
}

type EmitTokenMetadataAccounts struct {
	Metadata solanago.PublicKey `json:"metadata"`
}

func NewEmitTokenMetadataInstruction(params EmitTokenMetadataParams, accounts EmitTokenMetadataAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeEmitTokenMetadata); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Start == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode start option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode start option prefix: %w", err)
		}
		if err := enc__.Encode(*params.Start); err != nil {
			return nil, fmt.Errorf("failed to encode start: %w", err)
		}
	}
	if params.End == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode end option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode end option prefix: %w", err)
		}
		if err := enc__.Encode(*params.End); err != nil {
			return nil, fmt.Errorf("failed to encode end: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Metadata.IsZero() {
		return nil, fmt.Errorf("account metadata is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Metadata, false, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeTokenGroupParams struct {
	UpdateAuthority *solanago.PublicKey `json:"updateAuthority"`
	MaxSize         uint64              `json:"maxSize"`
}

type InitializeTokenGroupAccounts struct {
	Group         solanago.PublicKey `json:"group"`
	Mint          solanago.PublicKey `json:"mint"`
	MintAuthority solanago.PublicKey `json:"mintAuthority"`
}

func NewInitializeTokenGroupInstruction(params InitializeTokenGroupParams, accounts InitializeTokenGroupAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeTokenGroup); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.UpdateAuthority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode updateAuthority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.UpdateAuthority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode updateAuthority: %w", err)
		}
	}
	if err := enc__.Encode(params.MaxSize); err != nil {
		return nil, fmt.Errorf("failed to encode maxSize: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Group.IsZero() {
		return nil, fmt.Errorf("account group is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Group, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, false, false))
	if accounts.MintAuthority.IsZero() {
		return nil, fmt.Errorf("account mintAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MintAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateTokenGroupMaxSizeParams struct {
	MaxSize uint64 `json:"maxSize"`
}

type UpdateTokenGroupMaxSizeAccounts struct {
	Group           solanago.PublicKey `json:"group"`
	UpdateAuthority solanago.PublicKey `json:"updateAuthority"`
}

func NewUpdateTokenGroupMaxSizeInstruction(params UpdateTokenGroupMaxSizeParams, accounts UpdateTokenGroupMaxSizeAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateTokenGroupMaxSize); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.MaxSize); err != nil {
		return nil, fmt.Errorf("failed to encode maxSize: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Group.IsZero() {
		return nil, fmt.Errorf("account group is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Group, true, false))
	if accounts.UpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account updateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.UpdateAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UpdateTokenGroupUpdateAuthorityParams struct {
	NewUpdateAuthority *solanago.PublicKey `json:"newUpdateAuthority"`
}

type UpdateTokenGroupUpdateAuthorityAccounts struct {
	Group           solanago.PublicKey `json:"group"`
	UpdateAuthority solanago.PublicKey `json:"updateAuthority"`
}

func NewUpdateTokenGroupUpdateAuthorityInstruction(params UpdateTokenGroupUpdateAuthorityParams, accounts UpdateTokenGroupUpdateAuthorityAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUpdateTokenGroupUpdateAuthority); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.NewUpdateAuthority == nil {
		var zero solanago.PublicKey
		if err := enc__.WriteBytes(zero[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode newUpdateAuthority: %w", err)
		}
	} else {
		if err := enc__.WriteBytes(params.NewUpdateAuthority[:], false); err != nil {
			return nil, fmt.Errorf("failed to encode newUpdateAuthority: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Group.IsZero() {
		return nil, fmt.Errorf("account group is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Group, true, false))
	if accounts.UpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account updateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.UpdateAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializeTokenGroupMemberParams struct {
}

type InitializeTokenGroupMemberAccounts struct {
	Member               solanago.PublicKey `json:"member"`
	MemberMint           solanago.PublicKey `json:"memberMint"`
	MemberMintAuthority  solanago.PublicKey `json:"memberMintAuthority"`
	Group                solanago.PublicKey `json:"group"`
	GroupUpdateAuthority solanago.PublicKey `json:"groupUpdateAuthority"`
}

func NewInitializeTokenGroupMemberInstruction(params InitializeTokenGroupMemberParams, accounts InitializeTokenGroupMemberAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializeTokenGroupMember); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Member.IsZero() {
		return nil, fmt.Errorf("account member is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Member, true, false))
	if accounts.MemberMint.IsZero() {
		return nil, fmt.Errorf("account memberMint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MemberMint, false, false))
	if accounts.MemberMintAuthority.IsZero() {
		return nil, fmt.Errorf("account memberMintAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.MemberMintAuthority, false, true))
	if accounts.Group.IsZero() {
		return nil, fmt.Errorf("account group is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Group, true, false))
	if accounts.GroupUpdateAuthority.IsZero() {
		return nil, fmt.Errorf("account groupUpdateAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.GroupUpdateAuthority, false, true))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type UnwrapLamportsParams struct {
	Amount *uint64 `json:"amount"`
}

type UnwrapLamportsAccounts struct {
	Source       solanago.PublicKey   `json:"source"`
	Destination  solanago.PublicKey   `json:"destination"`
	Authority    solanago.PublicKey   `json:"authority"`
	MultiSigners []solanago.PublicKey `json:"multiSigners"`
}

func NewUnwrapLamportsInstruction(params UnwrapLamportsParams, accounts UnwrapLamportsAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeUnwrapLamports); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if params.Amount == nil {
		if err := enc__.WriteUint8(0); err != nil {
			return nil, fmt.Errorf("failed to encode amount option prefix: %w", err)
		}
	} else {
		if err := enc__.WriteUint8(1); err != nil {
			return nil, fmt.Errorf("failed to encode amount option prefix: %w", err)
		}
		if err := enc__.Encode(*params.Amount); err != nil {
			return nil, fmt.Errorf("failed to encode amount: %w", err)
		}
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Source.IsZero() {
		return nil, fmt.Errorf("account source is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Source, true, false))
	if accounts.Destination.IsZero() {
		return nil, fmt.Errorf("account destination is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Destination, true, false))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type InitializePermissionedBurnParams struct {
	Authority solanago.PublicKey `json:"authority"`
}

type InitializePermissionedBurnAccounts struct {
	Mint solanago.PublicKey `json:"mint"`
}

func NewInitializePermissionedBurnInstruction(params InitializePermissionedBurnParams, accounts InitializePermissionedBurnAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypeInitializePermissionedBurn); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Authority); err != nil {
		return nil, fmt.Errorf("failed to encode authority: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type PermissionedBurnParams struct {
	Amount uint64 `json:"amount"`
}

type PermissionedBurnAccounts struct {
	Account                   solanago.PublicKey   `json:"account"`
	Mint                      solanago.PublicKey   `json:"mint"`
	PermissionedBurnAuthority solanago.PublicKey   `json:"permissionedBurnAuthority"`
	Authority                 solanago.PublicKey   `json:"authority"`
	MultiSigners              []solanago.PublicKey `json:"multiSigners"`
}

func NewPermissionedBurnInstruction(params PermissionedBurnParams, accounts PermissionedBurnAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypePermissionedBurn); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.PermissionedBurnAuthority.IsZero() {
		return nil, fmt.Errorf("account permissionedBurnAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.PermissionedBurnAuthority, false, true))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

type PermissionedBurnCheckedParams struct {
	Amount   uint64 `json:"amount"`
	Decimals uint8  `json:"decimals"`
}

type PermissionedBurnCheckedAccounts struct {
	Account                   solanago.PublicKey   `json:"account"`
	Mint                      solanago.PublicKey   `json:"mint"`
	PermissionedBurnAuthority solanago.PublicKey   `json:"permissionedBurnAuthority"`
	Authority                 solanago.PublicKey   `json:"authority"`
	MultiSigners              []solanago.PublicKey `json:"multiSigners"`
}

func NewPermissionedBurnCheckedInstruction(params PermissionedBurnCheckedParams, accounts PermissionedBurnCheckedAccounts) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)
	if err := encodeInstructionTypeDiscriminator(enc__, InstructionTypePermissionedBurnChecked); err != nil {
		return nil, fmt.Errorf("failed to encode discriminator: %w", err)
	}
	if err := enc__.Encode(params.Amount); err != nil {
		return nil, fmt.Errorf("failed to encode amount: %w", err)
	}
	if err := enc__.Encode(params.Decimals); err != nil {
		return nil, fmt.Errorf("failed to encode decimals: %w", err)
	}
	accounts__ := solanago.AccountMetaSlice{}
	if accounts.Account.IsZero() {
		return nil, fmt.Errorf("account account is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Account, true, false))
	if accounts.Mint.IsZero() {
		return nil, fmt.Errorf("account mint is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Mint, true, false))
	if accounts.PermissionedBurnAuthority.IsZero() {
		return nil, fmt.Errorf("account permissionedBurnAuthority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.PermissionedBurnAuthority, false, true))
	if accounts.Authority.IsZero() {
		return nil, fmt.Errorf("account authority is required")
	}
	accounts__.Append(solanago.NewAccountMeta(accounts.Authority, false, len(accounts.MultiSigners) == 0))
	for _, pk := range accounts.MultiSigners {
		accounts__.Append(solanago.NewAccountMeta(pk, false, true))
	}
	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

func IdentifyInstructionType(data []byte) (InstructionType, error) {
	return identifyInstructionTypeByDiscriminator(data)
}

type DecodedInstruction struct {
	Type           InstructionType
	Accounts       []*solanago.AccountMeta
	Data           []byte
	Params         any
	ParsedAccounts any
}

func DecodeInstruction(accounts []*solanago.AccountMeta, data []byte) (*DecodedInstruction, error) {
	t, err := IdentifyInstructionType(data)
	if err != nil {
		return nil, err
	}

	out := &DecodedInstruction{
		Type:     t,
		Accounts: accounts,
		Data:     append([]byte(nil), data...),
	}

	paramsProto := newParamsByInstructionType(t)
	if paramsProto != nil {
		decodedParams, err := decodeInstructionParams(t, paramsProto, data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %s params: %w", t, err)
		}
		out.Params = decodedParams
	}

	accountsProto := newAccountsByInstructionType(t)
	if accountsProto != nil {
		decodedAccounts, err := decodeInstructionAccounts(t, accountsProto, accounts, out.Params)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %s accounts: %w", t, err)
		}
		out.ParsedAccounts = decodedAccounts
	}
	return out, nil
}

func registryDecodeInstruction(accounts []*solanago.AccountMeta, data []byte) (interface{}, error) {
	return DecodeInstruction(accounts, data)
}

func decodeInstructionParams(instructionType InstructionType, paramsProto any, data []byte) (any, error) {
	paramsType := reflect.TypeOf(paramsProto)
	if paramsType == nil || paramsType.Kind() != reflect.Ptr || paramsType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("params prototype must be pointer to struct")
	}

	prefixLen, err := instructionDiscriminatorLenByType(instructionType)
	if err != nil {
		return nil, err
	}
	if len(data) < prefixLen {
		return nil, fmt.Errorf("instruction data too short for discriminator: expected at least %d bytes, got %d", prefixLen, len(data))
	}
	inst := reflect.New(paramsType.Elem()).Interface()
	dec := binary.NewBorshDecoder(data[prefixLen:])
	if err := dec.Decode(inst); err != nil {
		return nil, err
	}
	if dec.Remaining() != 0 {
		return nil, fmt.Errorf("unable to decode params for %s: %d trailing bytes", instructionType, dec.Remaining())
	}
	return reflect.ValueOf(inst).Elem().Interface(), nil
}

func decodeInstructionAccounts(
	instructionType InstructionType,
	accountsProto any,
	accountMetas []*solanago.AccountMeta,
	decodedParams any,
) (any, error) {
	// Special-case instructions with multiple variadic account arrays.
	if instructionType == InstructionTypeWithdrawWithheldTokensFromAccounts {
		return decodeWithdrawWithheldTokensFromAccountsAccounts(accountMetas, decodedParams)
	}

	accountsType := reflect.TypeOf(accountsProto)
	if accountsType == nil || accountsType.Kind() != reflect.Ptr || accountsType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("accounts prototype must be pointer to struct")
	}

	out := reflect.New(accountsType.Elem()).Elem()
	publicKeyType := reflect.TypeOf(solanago.PublicKey{})
	publicKeyPtrType := reflect.TypeOf(&solanago.PublicKey{})
	publicKeySliceType := reflect.TypeOf([]solanago.PublicKey{})

	index := 0
	for i := 0; i < out.NumField(); i++ {
		field := out.Field(i)
		fieldType := field.Type()
		fieldName := out.Type().Field(i).Name

		switch fieldType {
		case publicKeyType:
			if index >= len(accountMetas) {
				return nil, fmt.Errorf("missing account for field %s", fieldName)
			}
			field.Set(reflect.ValueOf(accountMetas[index].PublicKey))
			index++
		case publicKeyPtrType:
			if index >= len(accountMetas) {
				return nil, fmt.Errorf("missing optional account for field %s", fieldName)
			}
			pk := accountMetas[index].PublicKey
			index++
			// Optional accounts are encoded as ProgramID placeholder.
			if pk.Equals(ProgramID) {
				field.Set(reflect.Zero(fieldType))
			} else {
				copied := pk
				field.Set(reflect.ValueOf(&copied))
			}
		case publicKeySliceType:
			remaining := len(accountMetas) - index
			slice := reflect.MakeSlice(fieldType, remaining, remaining)
			for j := 0; j < remaining; j++ {
				slice.Index(j).Set(reflect.ValueOf(accountMetas[index+j].PublicKey))
			}
			field.Set(slice)
			index = len(accountMetas)
		default:
			return nil, fmt.Errorf("unsupported accounts field type for %s: %s", fieldName, fieldType)
		}
	}

	if index != len(accountMetas) {
		return nil, fmt.Errorf("unexpected trailing accounts: consumed=%d total=%d", index, len(accountMetas))
	}

	return out.Interface(), nil
}

func decodeWithdrawWithheldTokensFromAccountsAccounts(
	accountMetas []*solanago.AccountMeta,
	decodedParams any,
) (any, error) {
	params, ok := decodedParams.(WithdrawWithheldTokensFromAccountsParams)
	if !ok {
		return nil, fmt.Errorf("missing or invalid params for withdrawWithheldTokensFromAccounts")
	}
	if len(accountMetas) < 3 {
		return nil, fmt.Errorf("not enough accounts: expected at least 3, got %d", len(accountMetas))
	}

	numSources := int(params.NumTokenAccounts)
	numRemaining := len(accountMetas) - 3
	if numSources > numRemaining {
		return nil, fmt.Errorf("numTokenAccounts=%d exceeds remaining accounts=%d", numSources, numRemaining)
	}
	numMultiSigners := numRemaining - numSources

	out := WithdrawWithheldTokensFromAccountsAccounts{
		Mint:                      accountMetas[0].PublicKey,
		FeeReceiver:               accountMetas[1].PublicKey,
		WithdrawWithheldAuthority: accountMetas[2].PublicKey,
		MultiSigners:              make([]solanago.PublicKey, 0, numMultiSigners),
		Sources:                   make([]solanago.PublicKey, 0, numSources),
	}

	start := 3
	for i := 0; i < numMultiSigners; i++ {
		out.MultiSigners = append(out.MultiSigners, accountMetas[start+i].PublicKey)
	}
	start += numMultiSigners
	for i := 0; i < numSources; i++ {
		out.Sources = append(out.Sources, accountMetas[start+i].PublicKey)
	}
	return out, nil
}
