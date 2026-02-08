// Code generated from interface/idl.json. DO NOT EDIT.
package token2022

import (
	"fmt"

	binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

type ExtensionUninitialized struct {
}

func ParseExtensionUninitialized(data []byte) (*ExtensionUninitialized, error) {
	out := new(ExtensionUninitialized)
	_ = data
	return out, nil
}

type ExtensionTransferFeeConfig struct {
	TransferFeeConfigAuthority solanago.PublicKey `json:"transferFeeConfigAuthority"`
	WithdrawWithheldAuthority  solanago.PublicKey `json:"withdrawWithheldAuthority"`
	WithheldAmount             uint64             `json:"withheldAmount"`
	OlderTransferFee           TransferFee        `json:"olderTransferFee"`
	NewerTransferFee           TransferFee        `json:"newerTransferFee"`
}

func ParseExtensionTransferFeeConfig(data []byte) (*ExtensionTransferFeeConfig, error) {
	out := new(ExtensionTransferFeeConfig)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.TransferFeeConfigAuthority: %w", err)
	}
	out.TransferFeeConfigAuthority = solanago.PublicKeyFromBytes(tmp)
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.WithdrawWithheldAuthority: %w", err)
	}
	out.WithdrawWithheldAuthority = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.WithheldAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.WithheldAmount: %w", err)
	}
	err = decoder.Decode(&out.OlderTransferFee)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.OlderTransferFee: %w", err)
	}
	err = decoder.Decode(&out.NewerTransferFee)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.NewerTransferFee: %w", err)
	}
	return out, nil
}

type ExtensionTransferFeeAmount struct {
	WithheldAmount uint64 `json:"withheldAmount"`
}

func ParseExtensionTransferFeeAmount(data []byte) (*ExtensionTransferFeeAmount, error) {
	out := new(ExtensionTransferFeeAmount)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.WithheldAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.WithheldAmount: %w", err)
	}
	return out, nil
}

type ExtensionMintCloseAuthority struct {
	CloseAuthority solanago.PublicKey `json:"closeAuthority"`
}

func ParseExtensionMintCloseAuthority(data []byte) (*ExtensionMintCloseAuthority, error) {
	out := new(ExtensionMintCloseAuthority)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.CloseAuthority: %w", err)
	}
	out.CloseAuthority = solanago.PublicKeyFromBytes(tmp)
	return out, nil
}

type ExtensionConfidentialTransferMint struct {
	Authority              *solanago.PublicKey `json:"authority"`
	AutoApproveNewAccounts bool                `json:"autoApproveNewAccounts"`
	AuditorElgamalPubkey   *solanago.PublicKey `json:"auditorElgamalPubkey"`
}

func ParseExtensionConfidentialTransferMint(data []byte) (*ExtensionConfidentialTransferMint, error) {
	out := new(ExtensionConfidentialTransferMint)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	err = decoder.Decode(&out.AutoApproveNewAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.AutoApproveNewAccounts: %w", err)
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.AuditorElgamalPubkey: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.AuditorElgamalPubkey = &v
	}
	return out, nil
}

type ExtensionConfidentialTransferAccount struct {
	Approved                            bool               `json:"approved"`
	ElgamalPubkey                       solanago.PublicKey `json:"elgamalPubkey"`
	PendingBalanceLow                   EncryptedBalance   `json:"pendingBalanceLow"`
	PendingBalanceHigh                  EncryptedBalance   `json:"pendingBalanceHigh"`
	AvailableBalance                    EncryptedBalance   `json:"availableBalance"`
	DecryptableAvailableBalance         DecryptableBalance `json:"decryptableAvailableBalance"`
	AllowConfidentialCredits            bool               `json:"allowConfidentialCredits"`
	AllowNonConfidentialCredits         bool               `json:"allowNonConfidentialCredits"`
	PendingBalanceCreditCounter         uint64             `json:"pendingBalanceCreditCounter"`
	MaximumPendingBalanceCreditCounter  uint64             `json:"maximumPendingBalanceCreditCounter"`
	ExpectedPendingBalanceCreditCounter uint64             `json:"expectedPendingBalanceCreditCounter"`
	ActualPendingBalanceCreditCounter   uint64             `json:"actualPendingBalanceCreditCounter"`
}

func ParseExtensionConfidentialTransferAccount(data []byte) (*ExtensionConfidentialTransferAccount, error) {
	out := new(ExtensionConfidentialTransferAccount)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.Approved)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Approved: %w", err)
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.ElgamalPubkey: %w", err)
	}
	out.ElgamalPubkey = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.PendingBalanceLow)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.PendingBalanceLow: %w", err)
	}
	err = decoder.Decode(&out.PendingBalanceHigh)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.PendingBalanceHigh: %w", err)
	}
	err = decoder.Decode(&out.AvailableBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.AvailableBalance: %w", err)
	}
	err = decoder.Decode(&out.DecryptableAvailableBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.DecryptableAvailableBalance: %w", err)
	}
	err = decoder.Decode(&out.AllowConfidentialCredits)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.AllowConfidentialCredits: %w", err)
	}
	err = decoder.Decode(&out.AllowNonConfidentialCredits)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.AllowNonConfidentialCredits: %w", err)
	}
	err = decoder.Decode(&out.PendingBalanceCreditCounter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.PendingBalanceCreditCounter: %w", err)
	}
	err = decoder.Decode(&out.MaximumPendingBalanceCreditCounter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.MaximumPendingBalanceCreditCounter: %w", err)
	}
	err = decoder.Decode(&out.ExpectedPendingBalanceCreditCounter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.ExpectedPendingBalanceCreditCounter: %w", err)
	}
	err = decoder.Decode(&out.ActualPendingBalanceCreditCounter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.ActualPendingBalanceCreditCounter: %w", err)
	}
	return out, nil
}

type ExtensionDefaultAccountState struct {
	State AccountState `json:"state"`
}

func ParseExtensionDefaultAccountState(data []byte) (*ExtensionDefaultAccountState, error) {
	out := new(ExtensionDefaultAccountState)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.State)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.State: %w", err)
	}
	return out, nil
}

type ExtensionImmutableOwner struct {
}

func ParseExtensionImmutableOwner(data []byte) (*ExtensionImmutableOwner, error) {
	out := new(ExtensionImmutableOwner)
	_ = data
	return out, nil
}

type ExtensionMemoTransfer struct {
	RequireIncomingTransferMemos bool `json:"requireIncomingTransferMemos"`
}

func ParseExtensionMemoTransfer(data []byte) (*ExtensionMemoTransfer, error) {
	out := new(ExtensionMemoTransfer)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.RequireIncomingTransferMemos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.RequireIncomingTransferMemos: %w", err)
	}
	return out, nil
}

type ExtensionNonTransferable struct {
}

func ParseExtensionNonTransferable(data []byte) (*ExtensionNonTransferable, error) {
	out := new(ExtensionNonTransferable)
	_ = data
	return out, nil
}

type ExtensionInterestBearingConfig struct {
	RateAuthority           solanago.PublicKey `json:"rateAuthority"`
	InitializationTimestamp uint64             `json:"initializationTimestamp"`
	PreUpdateAverageRate    int16              `json:"preUpdateAverageRate"`
	LastUpdateTimestamp     uint64             `json:"lastUpdateTimestamp"`
	CurrentRate             int16              `json:"currentRate"`
}

func ParseExtensionInterestBearingConfig(data []byte) (*ExtensionInterestBearingConfig, error) {
	out := new(ExtensionInterestBearingConfig)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.RateAuthority: %w", err)
	}
	out.RateAuthority = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.InitializationTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.InitializationTimestamp: %w", err)
	}
	err = decoder.Decode(&out.PreUpdateAverageRate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.PreUpdateAverageRate: %w", err)
	}
	err = decoder.Decode(&out.LastUpdateTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.LastUpdateTimestamp: %w", err)
	}
	err = decoder.Decode(&out.CurrentRate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.CurrentRate: %w", err)
	}
	return out, nil
}

type ExtensionCpiGuard struct {
	LockCpi bool `json:"lockCpi"`
}

func ParseExtensionCpiGuard(data []byte) (*ExtensionCpiGuard, error) {
	out := new(ExtensionCpiGuard)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.LockCpi)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.LockCpi: %w", err)
	}
	return out, nil
}

type ExtensionPermanentDelegate struct {
	Delegate solanago.PublicKey `json:"delegate"`
}

func ParseExtensionPermanentDelegate(data []byte) (*ExtensionPermanentDelegate, error) {
	out := new(ExtensionPermanentDelegate)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Delegate: %w", err)
	}
	out.Delegate = solanago.PublicKeyFromBytes(tmp)
	return out, nil
}

type ExtensionNonTransferableAccount struct {
}

func ParseExtensionNonTransferableAccount(data []byte) (*ExtensionNonTransferableAccount, error) {
	out := new(ExtensionNonTransferableAccount)
	_ = data
	return out, nil
}

type ExtensionTransferHook struct {
	Authority solanago.PublicKey `json:"authority"`
	ProgramID solanago.PublicKey `json:"programId"`
}

func ParseExtensionTransferHook(data []byte) (*ExtensionTransferHook, error) {
	out := new(ExtensionTransferHook)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	out.Authority = solanago.PublicKeyFromBytes(tmp)
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.ProgramID: %w", err)
	}
	out.ProgramID = solanago.PublicKeyFromBytes(tmp)
	return out, nil
}

type ExtensionTransferHookAccount struct {
	Transferring bool `json:"transferring"`
}

func ParseExtensionTransferHookAccount(data []byte) (*ExtensionTransferHookAccount, error) {
	out := new(ExtensionTransferHookAccount)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.Transferring)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Transferring: %w", err)
	}
	return out, nil
}

type ExtensionConfidentialTransferFee struct {
	Authority            *solanago.PublicKey `json:"authority"`
	ElgamalPubkey        solanago.PublicKey  `json:"elgamalPubkey"`
	HarvestToMintEnabled bool                `json:"harvestToMintEnabled"`
	WithheldAmount       EncryptedBalance    `json:"withheldAmount"`
}

func ParseExtensionConfidentialTransferFee(data []byte) (*ExtensionConfidentialTransferFee, error) {
	out := new(ExtensionConfidentialTransferFee)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.ElgamalPubkey: %w", err)
	}
	out.ElgamalPubkey = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.HarvestToMintEnabled)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.HarvestToMintEnabled: %w", err)
	}
	err = decoder.Decode(&out.WithheldAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.WithheldAmount: %w", err)
	}
	return out, nil
}

type ExtensionConfidentialTransferFeeAmount struct {
	WithheldAmount EncryptedBalance `json:"withheldAmount"`
}

func ParseExtensionConfidentialTransferFeeAmount(data []byte) (*ExtensionConfidentialTransferFeeAmount, error) {
	out := new(ExtensionConfidentialTransferFeeAmount)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	err = decoder.Decode(&out.WithheldAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.WithheldAmount: %w", err)
	}
	return out, nil
}

type ExtensionMetadataPointer struct {
	Authority       *solanago.PublicKey `json:"authority"`
	MetadataAddress *solanago.PublicKey `json:"metadataAddress"`
}

func ParseExtensionMetadataPointer(data []byte) (*ExtensionMetadataPointer, error) {
	out := new(ExtensionMetadataPointer)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.MetadataAddress: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.MetadataAddress = &v
	}
	return out, nil
}

type ExtensionTokenMetadata struct {
	UpdateAuthority    *solanago.PublicKey `json:"updateAuthority"`
	Mint               solanago.PublicKey  `json:"mint"`
	Name               string              `json:"name"`
	Symbol             string              `json:"symbol"`
	Uri                string              `json:"uri"`
	AdditionalMetadata map[string]string   `json:"additionalMetadata"`
}

func ParseExtensionTokenMetadata(data []byte) (*ExtensionTokenMetadata, error) {
	out := new(ExtensionTokenMetadata)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.UpdateAuthority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.UpdateAuthority = &v
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Mint: %w", err)
	}
	out.Mint = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Name: %w", err)
	}
	err = decoder.Decode(&out.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Symbol: %w", err)
	}
	err = decoder.Decode(&out.Uri)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Uri: %w", err)
	}
	var count uint32
	count, err = decoder.ReadUint32(binary.LE)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.AdditionalMetadata: %w", err)
	}
	out.AdditionalMetadata = make(map[string]string, count)
	for i := uint32(0); i < count; i++ {
		var k string
		var v string
		err = decoder.Decode(&k)
		if err != nil {
			return nil, fmt.Errorf("failed to decode map key: %w", err)
		}
		err = decoder.Decode(&v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode map value: %w", err)
		}
		out.AdditionalMetadata[k] = v
	}
	return out, nil
}

type ExtensionGroupPointer struct {
	Authority    *solanago.PublicKey `json:"authority"`
	GroupAddress *solanago.PublicKey `json:"groupAddress"`
}

func ParseExtensionGroupPointer(data []byte) (*ExtensionGroupPointer, error) {
	out := new(ExtensionGroupPointer)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.GroupAddress: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.GroupAddress = &v
	}
	return out, nil
}

type ExtensionTokenGroup struct {
	UpdateAuthority *solanago.PublicKey `json:"updateAuthority"`
	Mint            solanago.PublicKey  `json:"mint"`
	Size            uint64              `json:"size"`
	MaxSize         uint64              `json:"maxSize"`
}

func ParseExtensionTokenGroup(data []byte) (*ExtensionTokenGroup, error) {
	out := new(ExtensionTokenGroup)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.UpdateAuthority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.UpdateAuthority = &v
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Mint: %w", err)
	}
	out.Mint = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Size: %w", err)
	}
	err = decoder.Decode(&out.MaxSize)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.MaxSize: %w", err)
	}
	return out, nil
}

type ExtensionGroupMemberPointer struct {
	Authority     *solanago.PublicKey `json:"authority"`
	MemberAddress *solanago.PublicKey `json:"memberAddress"`
}

func ParseExtensionGroupMemberPointer(data []byte) (*ExtensionGroupMemberPointer, error) {
	out := new(ExtensionGroupMemberPointer)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.MemberAddress: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.MemberAddress = &v
	}
	return out, nil
}

type ExtensionTokenGroupMember struct {
	Mint         solanago.PublicKey `json:"mint"`
	Group        solanago.PublicKey `json:"group"`
	MemberNumber uint64             `json:"memberNumber"`
}

func ParseExtensionTokenGroupMember(data []byte) (*ExtensionTokenGroupMember, error) {
	out := new(ExtensionTokenGroupMember)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Mint: %w", err)
	}
	out.Mint = solanago.PublicKeyFromBytes(tmp)
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Group: %w", err)
	}
	out.Group = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.MemberNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.MemberNumber: %w", err)
	}
	return out, nil
}

type ExtensionConfidentialMintBurn struct {
}

func ParseExtensionConfidentialMintBurn(data []byte) (*ExtensionConfidentialMintBurn, error) {
	out := new(ExtensionConfidentialMintBurn)
	_ = data
	return out, nil
}

type ExtensionScaledUIAmountConfig struct {
	Authority                       solanago.PublicKey `json:"authority"`
	Multiplier                      float64            `json:"multiplier"`
	NewMultiplierEffectiveTimestamp uint64             `json:"newMultiplierEffectiveTimestamp"`
	NewMultiplier                   float64            `json:"newMultiplier"`
}

func ParseExtensionScaledUIAmountConfig(data []byte) (*ExtensionScaledUIAmountConfig, error) {
	out := new(ExtensionScaledUIAmountConfig)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	out.Authority = solanago.PublicKeyFromBytes(tmp)
	err = decoder.Decode(&out.Multiplier)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Multiplier: %w", err)
	}
	err = decoder.Decode(&out.NewMultiplierEffectiveTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.NewMultiplierEffectiveTimestamp: %w", err)
	}
	err = decoder.Decode(&out.NewMultiplier)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.NewMultiplier: %w", err)
	}
	return out, nil
}

type ExtensionPausableConfig struct {
	Authority *solanago.PublicKey `json:"authority"`
	Paused    bool                `json:"paused"`
}

func ParseExtensionPausableConfig(data []byte) (*ExtensionPausableConfig, error) {
	out := new(ExtensionPausableConfig)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	err = decoder.Decode(&out.Paused)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Paused: %w", err)
	}
	return out, nil
}

type ExtensionPausableAccount struct {
}

func ParseExtensionPausableAccount(data []byte) (*ExtensionPausableAccount, error) {
	out := new(ExtensionPausableAccount)
	_ = data
	return out, nil
}

type ExtensionPermissionedBurn struct {
	Authority *solanago.PublicKey `json:"authority"`
}

func ParseExtensionPermissionedBurn(data []byte) (*ExtensionPermissionedBurn, error) {
	out := new(ExtensionPermissionedBurn)
	decoder := binary.NewBorshDecoder(data)
	var err error
	var tmp []byte
	_ = tmp
	tmp, err = decoder.ReadNBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode out.Authority: %w", err)
	}
	if !allZero(tmp) {
		v := solanago.PublicKeyFromBytes(tmp)
		out.Authority = &v
	}
	return out, nil
}

func ParseExtensionData(extType ExtensionType, data []byte) (any, error) {
	switch extType {
	case ExtensionTypeUninitialized:
		return ParseExtensionUninitialized(data)
	case ExtensionTypeTransferFeeConfig:
		return ParseExtensionTransferFeeConfig(data)
	case ExtensionTypeTransferFeeAmount:
		return ParseExtensionTransferFeeAmount(data)
	case ExtensionTypeMintCloseAuthority:
		return ParseExtensionMintCloseAuthority(data)
	case ExtensionTypeConfidentialTransferMint:
		return ParseExtensionConfidentialTransferMint(data)
	case ExtensionTypeConfidentialTransferAccount:
		return ParseExtensionConfidentialTransferAccount(data)
	case ExtensionTypeDefaultAccountState:
		return ParseExtensionDefaultAccountState(data)
	case ExtensionTypeImmutableOwner:
		return ParseExtensionImmutableOwner(data)
	case ExtensionTypeMemoTransfer:
		return ParseExtensionMemoTransfer(data)
	case ExtensionTypeNonTransferable:
		return ParseExtensionNonTransferable(data)
	case ExtensionTypeInterestBearingConfig:
		return ParseExtensionInterestBearingConfig(data)
	case ExtensionTypeCpiGuard:
		return ParseExtensionCpiGuard(data)
	case ExtensionTypePermanentDelegate:
		return ParseExtensionPermanentDelegate(data)
	case ExtensionTypeNonTransferableAccount:
		return ParseExtensionNonTransferableAccount(data)
	case ExtensionTypeTransferHook:
		return ParseExtensionTransferHook(data)
	case ExtensionTypeTransferHookAccount:
		return ParseExtensionTransferHookAccount(data)
	case ExtensionTypeConfidentialTransferFeeConfig:
		return ParseExtensionConfidentialTransferFee(data)
	case ExtensionTypeConfidentialTransferFeeAmount:
		return ParseExtensionConfidentialTransferFeeAmount(data)
	case ExtensionTypeMetadataPointer:
		return ParseExtensionMetadataPointer(data)
	case ExtensionTypeTokenMetadata:
		return ParseExtensionTokenMetadata(data)
	case ExtensionTypeGroupPointer:
		return ParseExtensionGroupPointer(data)
	case ExtensionTypeTokenGroup:
		return ParseExtensionTokenGroup(data)
	case ExtensionTypeGroupMemberPointer:
		return ParseExtensionGroupMemberPointer(data)
	case ExtensionTypeTokenGroupMember:
		return ParseExtensionTokenGroupMember(data)
	case ExtensionTypeConfidentialMintBurn:
		return ParseExtensionConfidentialMintBurn(data)
	case ExtensionTypeScaledUIAmount:
		return ParseExtensionScaledUIAmountConfig(data)
	case ExtensionTypePausable:
		return ParseExtensionPausableConfig(data)
	case ExtensionTypePausableAccount:
		return ParseExtensionPausableAccount(data)
	case ExtensionTypePermissionedBurn:
		return ParseExtensionPermissionedBurn(data)
	default:
		return nil, fmt.Errorf("unsupported extension type: %d", extType)
	}
}
