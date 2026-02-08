package token2022

import "fmt"

const (
	AccountTypeSize = 1
	TLVTypeSize     = 2
	TLVLengthSize   = 2
)

func AddTypeAndLengthToLen(valueLen int) int {
	return valueLen + TLVTypeSize + TLVLengthSize
}

func IsVariableLengthExtensionType(extType ExtensionType) bool {
	return extType == ExtensionTypeTokenMetadata
}

func ExtensionTypeDataLen(extType ExtensionType) (int, error) {
	switch extType {
	case ExtensionTypeUninitialized:
		return 0, nil
	case ExtensionTypeTransferFeeConfig:
		return 108, nil
	case ExtensionTypeTransferFeeAmount:
		return 8, nil
	case ExtensionTypeMintCloseAuthority:
		return 32, nil
	case ExtensionTypeConfidentialTransferMint:
		return 65, nil
	case ExtensionTypeConfidentialTransferAccount:
		return 295, nil
	case ExtensionTypeDefaultAccountState:
		return 1, nil
	case ExtensionTypeImmutableOwner:
		return 0, nil
	case ExtensionTypeMemoTransfer:
		return 1, nil
	case ExtensionTypeNonTransferable:
		return 0, nil
	case ExtensionTypeInterestBearingConfig:
		return 52, nil
	case ExtensionTypeCpiGuard:
		return 1, nil
	case ExtensionTypePermanentDelegate:
		return 32, nil
	case ExtensionTypeNonTransferableAccount:
		return 0, nil
	case ExtensionTypeTransferHook:
		return 64, nil
	case ExtensionTypeTransferHookAccount:
		return 1, nil
	case ExtensionTypeConfidentialTransferFeeConfig:
		return 129, nil
	case ExtensionTypeConfidentialTransferFeeAmount:
		return 64, nil
	case ExtensionTypeMetadataPointer:
		return 64, nil
	case ExtensionTypeTokenMetadata:
		return 0, fmt.Errorf("token metadata is variable-length")
	case ExtensionTypeGroupPointer:
		return 64, nil
	case ExtensionTypeTokenGroup:
		return 80, nil
	case ExtensionTypeGroupMemberPointer:
		return 64, nil
	case ExtensionTypeTokenGroupMember:
		return 72, nil
	case ExtensionTypeConfidentialMintBurn:
		return 196, nil
	case ExtensionTypeScaledUIAmount:
		return 56, nil
	case ExtensionTypePausable:
		return 33, nil
	case ExtensionTypePausableAccount:
		return 0, nil
	case ExtensionTypePermissionedBurn:
		return 32, nil
	default:
		return 0, fmt.Errorf("unknown extension type: %d", extType)
	}
}

func IsMintExtensionType(extType ExtensionType) bool {
	switch extType {
	case ExtensionTypeTransferFeeConfig,
		ExtensionTypeMintCloseAuthority,
		ExtensionTypeConfidentialTransferMint,
		ExtensionTypeDefaultAccountState,
		ExtensionTypeNonTransferable,
		ExtensionTypeInterestBearingConfig,
		ExtensionTypePermanentDelegate,
		ExtensionTypeTransferHook,
		ExtensionTypeConfidentialTransferFeeConfig,
		ExtensionTypeMetadataPointer,
		ExtensionTypeTokenMetadata,
		ExtensionTypeGroupPointer,
		ExtensionTypeTokenGroup,
		ExtensionTypeGroupMemberPointer,
		ExtensionTypeTokenGroupMember,
		ExtensionTypeConfidentialMintBurn,
		ExtensionTypeScaledUIAmount,
		ExtensionTypePausable,
		ExtensionTypePermissionedBurn:
		return true
	default:
		return false
	}
}

func IsAccountExtensionType(extType ExtensionType) bool {
	switch extType {
	case ExtensionTypeTransferFeeAmount,
		ExtensionTypeConfidentialTransferAccount,
		ExtensionTypeImmutableOwner,
		ExtensionTypeMemoTransfer,
		ExtensionTypeCpiGuard,
		ExtensionTypeNonTransferableAccount,
		ExtensionTypeTransferHookAccount,
		ExtensionTypeConfidentialTransferFeeAmount,
		ExtensionTypePausableAccount:
		return true
	default:
		return false
	}
}

func GetRequiredInitAccountExtensions(mintExtensionTypes []ExtensionType) []ExtensionType {
	out := make([]ExtensionType, 0)
	for _, extType := range mintExtensionTypes {
		switch extType {
		case ExtensionTypeTransferFeeConfig:
			out = appendIfMissingExtensionType(out, ExtensionTypeTransferFeeAmount)
		case ExtensionTypeNonTransferable:
			out = appendIfMissingExtensionType(out, ExtensionTypeNonTransferableAccount)
			out = appendIfMissingExtensionType(out, ExtensionTypeImmutableOwner)
		case ExtensionTypeTransferHook:
			out = appendIfMissingExtensionType(out, ExtensionTypeTransferHookAccount)
		case ExtensionTypePausable:
			out = appendIfMissingExtensionType(out, ExtensionTypePausableAccount)
		}
	}
	return out
}

func CalculateMintAccountLen(
	extensionTypes []ExtensionType,
	variableLengthExtensions map[ExtensionType]int,
) (int, error) {
	return calculateAccountLen(mintBaseLen, extensionTypes, variableLengthExtensions, IsMintExtensionType)
}

func CalculateTokenAccountLen(extensionTypes []ExtensionType) (int, error) {
	return calculateAccountLen(tokenBaseLen, extensionTypes, nil, IsAccountExtensionType)
}

func GetNewAccountLenForExtensionLen(
	currentAccountLen int,
	extensions []TLVExtension,
	extensionType ExtensionType,
	newExtensionLen int,
) int {
	currentEntryLen := 0
	for _, ext := range extensions {
		if ext.Type == extensionType {
			currentEntryLen = AddTypeAndLengthToLen(int(ext.Length))
			break
		}
	}
	newEntryLen := AddTypeAndLengthToLen(newExtensionLen)
	newLen := currentAccountLen - currentEntryLen + newEntryLen
	if newLen == multisigBaseLen {
		return newLen + TLVTypeSize
	}
	return newLen
}

func calculateAccountLen(
	baseSize int,
	extensionTypes []ExtensionType,
	variableLengthExtensions map[ExtensionType]int,
	validateType func(ExtensionType) bool,
) (int, error) {
	if len(extensionTypes) == 0 && len(variableLengthExtensions) == 0 {
		return baseSize, nil
	}

	totalTLVLen := 0
	dedup := make([]ExtensionType, 0, len(extensionTypes))
	for _, extType := range extensionTypes {
		if !validateType(extType) {
			return 0, fmt.Errorf("extension type %d is invalid for this account kind", extType)
		}
		if containsExtensionType(dedup, extType) {
			continue
		}
		dedup = append(dedup, extType)

		valueLen, err := ExtensionTypeDataLen(extType)
		if err != nil {
			return 0, err
		}
		totalTLVLen += AddTypeAndLengthToLen(valueLen)
	}

	for extType, valueLen := range variableLengthExtensions {
		if !validateType(extType) {
			return 0, fmt.Errorf("variable extension type %d is invalid for this account kind", extType)
		}
		if !IsVariableLengthExtensionType(extType) {
			return 0, fmt.Errorf("extension type %d is not variable-length", extType)
		}
		totalTLVLen += AddTypeAndLengthToLen(valueLen)
	}

	total := tokenBaseLen + AccountTypeSize + totalTLVLen
	if total == multisigBaseLen {
		total += TLVTypeSize
	}
	return total, nil
}

func appendIfMissingExtensionType(list []ExtensionType, extType ExtensionType) []ExtensionType {
	if containsExtensionType(list, extType) {
		return list
	}
	return append(list, extType)
}

func containsExtensionType(list []ExtensionType, extType ExtensionType) bool {
	for _, v := range list {
		if v == extType {
			return true
		}
	}
	return false
}
