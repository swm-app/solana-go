// Code generated from interface/src/error.rs. DO NOT EDIT.
package token2022

import "fmt"

type ErrorCode uint32

const (
	ErrorNotRentExempt ErrorCode = iota
	ErrorInsufficientFunds
	ErrorInvalidMint
	ErrorMintMismatch
	ErrorOwnerMismatch
	ErrorFixedSupply
	ErrorAlreadyInUse
	ErrorInvalidNumberOfProvidedSigners
	ErrorInvalidNumberOfRequiredSigners
	ErrorUninitializedState
	ErrorNativeNotSupported
	ErrorNonNativeHasBalance
	ErrorInvalidInstruction
	ErrorInvalidState
	ErrorOverflow
	ErrorAuthorityTypeNotSupported
	ErrorMintCannotFreeze
	ErrorAccountFrozen
	ErrorMintDecimalsMismatch
	ErrorNonNativeNotSupported
	ErrorExtensionTypeMismatch
	ErrorExtensionBaseMismatch
	ErrorExtensionAlreadyInitialized
	ErrorConfidentialTransferAccountHasBalance
	ErrorConfidentialTransferAccountNotApproved
	ErrorConfidentialTransferDepositsAndTransfersDisabled
	ErrorConfidentialTransferElGamalPubkeyMismatch
	ErrorConfidentialTransferBalanceMismatch
	ErrorMintHasSupply
	ErrorNoAuthorityExists
	ErrorTransferFeeExceedsMaximum
	ErrorMintRequiredForTransfer
	ErrorFeeMismatch
	ErrorFeeParametersMismatch
	ErrorImmutableOwner
	ErrorAccountHasWithheldTransferFees
	ErrorNoMemo
	ErrorNonTransferable
	ErrorNonTransferableNeedsImmutableOwnership
	ErrorMaximumPendingBalanceCreditCounterExceeded
	ErrorMaximumDepositAmountExceeded
	ErrorCpiGuardSettingsLocked
	ErrorCpiGuardTransferBlocked
	ErrorCpiGuardBurnBlocked
	ErrorCpiGuardCloseAccountBlocked
	ErrorCpiGuardApproveBlocked
	ErrorCpiGuardSetAuthorityBlocked
	ErrorCpiGuardOwnerChangeBlocked
	ErrorExtensionNotFound
	ErrorNonConfidentialTransfersDisabled
	ErrorConfidentialTransferFeeAccountHasWithheldFee
	ErrorInvalidExtensionCombination
	ErrorInvalidLengthForAlloc
	ErrorAccountDecryption
	ErrorProofGeneration
	ErrorInvalidProofInstructionOffset
	ErrorHarvestToMintDisabled
	ErrorSplitProofContextStateAccountsNotSupported
	ErrorNotEnoughProofContextStateAccounts
	ErrorMalformedCiphertext
	ErrorCiphertextArithmeticFailed
	ErrorPedersenCommitmentMismatch
	ErrorRangeProofLengthMismatch
	ErrorIllegalBitLength
	ErrorFeeCalculation
	ErrorIllegalMintBurnConversion
	ErrorInvalidScale
	ErrorMintPaused
	ErrorPendingBalanceNonZero
)

var errorCodeNames = map[ErrorCode]string{
	ErrorNotRentExempt:                                    "NotRentExempt",
	ErrorInsufficientFunds:                                "InsufficientFunds",
	ErrorInvalidMint:                                      "InvalidMint",
	ErrorMintMismatch:                                     "MintMismatch",
	ErrorOwnerMismatch:                                    "OwnerMismatch",
	ErrorFixedSupply:                                      "FixedSupply",
	ErrorAlreadyInUse:                                     "AlreadyInUse",
	ErrorInvalidNumberOfProvidedSigners:                   "InvalidNumberOfProvidedSigners",
	ErrorInvalidNumberOfRequiredSigners:                   "InvalidNumberOfRequiredSigners",
	ErrorUninitializedState:                               "UninitializedState",
	ErrorNativeNotSupported:                               "NativeNotSupported",
	ErrorNonNativeHasBalance:                              "NonNativeHasBalance",
	ErrorInvalidInstruction:                               "InvalidInstruction",
	ErrorInvalidState:                                     "InvalidState",
	ErrorOverflow:                                         "Overflow",
	ErrorAuthorityTypeNotSupported:                        "AuthorityTypeNotSupported",
	ErrorMintCannotFreeze:                                 "MintCannotFreeze",
	ErrorAccountFrozen:                                    "AccountFrozen",
	ErrorMintDecimalsMismatch:                             "MintDecimalsMismatch",
	ErrorNonNativeNotSupported:                            "NonNativeNotSupported",
	ErrorExtensionTypeMismatch:                            "ExtensionTypeMismatch",
	ErrorExtensionBaseMismatch:                            "ExtensionBaseMismatch",
	ErrorExtensionAlreadyInitialized:                      "ExtensionAlreadyInitialized",
	ErrorConfidentialTransferAccountHasBalance:            "ConfidentialTransferAccountHasBalance",
	ErrorConfidentialTransferAccountNotApproved:           "ConfidentialTransferAccountNotApproved",
	ErrorConfidentialTransferDepositsAndTransfersDisabled: "ConfidentialTransferDepositsAndTransfersDisabled",
	ErrorConfidentialTransferElGamalPubkeyMismatch:        "ConfidentialTransferElGamalPubkeyMismatch",
	ErrorConfidentialTransferBalanceMismatch:              "ConfidentialTransferBalanceMismatch",
	ErrorMintHasSupply:                                    "MintHasSupply",
	ErrorNoAuthorityExists:                                "NoAuthorityExists",
	ErrorTransferFeeExceedsMaximum:                        "TransferFeeExceedsMaximum",
	ErrorMintRequiredForTransfer:                          "MintRequiredForTransfer",
	ErrorFeeMismatch:                                      "FeeMismatch",
	ErrorFeeParametersMismatch:                            "FeeParametersMismatch",
	ErrorImmutableOwner:                                   "ImmutableOwner",
	ErrorAccountHasWithheldTransferFees:                   "AccountHasWithheldTransferFees",
	ErrorNoMemo:                                           "NoMemo",
	ErrorNonTransferable:                                  "NonTransferable",
	ErrorNonTransferableNeedsImmutableOwnership:           "NonTransferableNeedsImmutableOwnership",
	ErrorMaximumPendingBalanceCreditCounterExceeded:       "MaximumPendingBalanceCreditCounterExceeded",
	ErrorMaximumDepositAmountExceeded:                     "MaximumDepositAmountExceeded",
	ErrorCpiGuardSettingsLocked:                           "CpiGuardSettingsLocked",
	ErrorCpiGuardTransferBlocked:                          "CpiGuardTransferBlocked",
	ErrorCpiGuardBurnBlocked:                              "CpiGuardBurnBlocked",
	ErrorCpiGuardCloseAccountBlocked:                      "CpiGuardCloseAccountBlocked",
	ErrorCpiGuardApproveBlocked:                           "CpiGuardApproveBlocked",
	ErrorCpiGuardSetAuthorityBlocked:                      "CpiGuardSetAuthorityBlocked",
	ErrorCpiGuardOwnerChangeBlocked:                       "CpiGuardOwnerChangeBlocked",
	ErrorExtensionNotFound:                                "ExtensionNotFound",
	ErrorNonConfidentialTransfersDisabled:                 "NonConfidentialTransfersDisabled",
	ErrorConfidentialTransferFeeAccountHasWithheldFee:     "ConfidentialTransferFeeAccountHasWithheldFee",
	ErrorInvalidExtensionCombination:                      "InvalidExtensionCombination",
	ErrorInvalidLengthForAlloc:                            "InvalidLengthForAlloc",
	ErrorAccountDecryption:                                "AccountDecryption",
	ErrorProofGeneration:                                  "ProofGeneration",
	ErrorInvalidProofInstructionOffset:                    "InvalidProofInstructionOffset",
	ErrorHarvestToMintDisabled:                            "HarvestToMintDisabled",
	ErrorSplitProofContextStateAccountsNotSupported:       "SplitProofContextStateAccountsNotSupported",
	ErrorNotEnoughProofContextStateAccounts:               "NotEnoughProofContextStateAccounts",
	ErrorMalformedCiphertext:                              "MalformedCiphertext",
	ErrorCiphertextArithmeticFailed:                       "CiphertextArithmeticFailed",
	ErrorPedersenCommitmentMismatch:                       "PedersenCommitmentMismatch",
	ErrorRangeProofLengthMismatch:                         "RangeProofLengthMismatch",
	ErrorIllegalBitLength:                                 "IllegalBitLength",
	ErrorFeeCalculation:                                   "FeeCalculation",
	ErrorIllegalMintBurnConversion:                        "IllegalMintBurnConversion",
	ErrorInvalidScale:                                     "InvalidScale",
	ErrorMintPaused:                                       "MintPaused",
	ErrorPendingBalanceNonZero:                            "PendingBalanceNonZero",
}

func (c ErrorCode) Name() string {
	if n, ok := errorCodeNames[c]; ok {
		return n
	}
	return "Unknown"
}

func (c ErrorCode) Message() string {
	return c.Name()
}

func (c ErrorCode) Error() string {
	return fmt.Sprintf("token2022 error %d: %s", uint32(c), c.Message())
}

func ParseErrorCode(code uint32) (ErrorCode, bool) {
	c := ErrorCode(code)
	_, ok := errorCodeNames[c]
	return c, ok
}
