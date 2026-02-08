package token2022

import (
	stdbinary "encoding/binary"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

const (
	mintBaseLen     = 82
	tokenBaseLen    = 165
	multisigBaseLen = 355
	tlvTypeLen      = 2
	tlvLengthLen    = 2
)

type TLVExtension struct {
	Type   ExtensionType
	Length uint16
	Data   []byte
	Value  any
}

type MintStateWithExtensions struct {
	Base       MintAccount
	Extensions []TLVExtension
}

type TokenStateWithExtensions struct {
	Base       TokenAccount
	Extensions []TLVExtension
}

type MetadataPointer struct {
	Authority       solanago.PublicKey
	MetadataAddress solanago.PublicKey
}

func ParseAnyAccount(accountData []byte) (any, error) {
	if len(accountData) == multisigBaseLen {
		return ParseAccount_MultisigAccount(accountData)
	}
	if len(accountData) == mintBaseLen {
		return ParseAccount_MintAccount(accountData)
	}
	if len(accountData) == tokenBaseLen {
		return ParseAccount_TokenAccount(accountData)
	}
	if len(accountData) > tokenBaseLen {
		switch AccountType(accountData[tokenBaseLen]) {
		case AccountTypeMint:
			return ParseMintStateWithExtensions(accountData)
		case AccountTypeAccount:
			return ParseTokenStateWithExtensions(accountData)
		default:
			return nil, fmt.Errorf("unknown token2022 account type marker at offset %d: %d", tokenBaseLen, accountData[tokenBaseLen])
		}
	}
	return nil, fmt.Errorf("unable to identify token2022 account layout (len=%d)", len(accountData))
}

func ParseAccount_MintAccount(accountData []byte) (*MintAccount, error) {
	if len(accountData) < mintBaseLen {
		return nil, fmt.Errorf("mint account too small: %d", len(accountData))
	}
	dec := ag_binary.NewBorshDecoder(accountData[:mintBaseLen])
	out := new(MintAccount)
	if err := out.UnmarshalWithDecoder(dec); err != nil {
		return nil, fmt.Errorf("failed to decode mint account: %w", err)
	}
	return out, nil
}

func ParseAccount_TokenAccount(accountData []byte) (*TokenAccount, error) {
	if len(accountData) < tokenBaseLen {
		return nil, fmt.Errorf("token account too small: %d", len(accountData))
	}
	dec := ag_binary.NewBorshDecoder(accountData[:tokenBaseLen])
	out := new(TokenAccount)
	if err := out.UnmarshalWithDecoder(dec); err != nil {
		return nil, fmt.Errorf("failed to decode token account: %w", err)
	}
	return out, nil
}

func ParseAccount_MultisigAccount(accountData []byte) (*MultisigAccount, error) {
	if len(accountData) < multisigBaseLen {
		return nil, fmt.Errorf("multisig account too small: %d", len(accountData))
	}
	dec := ag_binary.NewBorshDecoder(accountData[:multisigBaseLen])
	out := new(MultisigAccount)
	if err := out.UnmarshalWithDecoder(dec); err != nil {
		return nil, fmt.Errorf("failed to decode multisig account: %w", err)
	}
	return out, nil
}

func ParseMintStateWithExtensions(accountData []byte) (*MintStateWithExtensions, error) {
	mintState, err := ParseAccount_MintAccount(accountData)
	if err != nil {
		return nil, err
	}
	res := &MintStateWithExtensions{Base: *mintState}
	if len(accountData) == mintBaseLen {
		return res, nil
	}
	// For accounts with extensions, token-2022 pads mints up to account base len (165),
	// writes account type at offset 165 and TLV starts at 166.
	if len(accountData) <= tokenBaseLen {
		return nil, fmt.Errorf("invalid extended mint size: %d", len(accountData))
	}
	if !allZero(accountData[mintBaseLen:tokenBaseLen]) {
		return nil, fmt.Errorf("invalid mint padding before account type marker")
	}
	if accountData[tokenBaseLen] != byte(AccountTypeMint) {
		return nil, fmt.Errorf("unexpected mint account type marker at offset %d: %d", tokenBaseLen, accountData[tokenBaseLen])
	}
	exts, err := ParseTLVExtensions(accountData[tokenBaseLen+1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse mint extensions: %w", err)
	}
	res.Extensions = exts
	return res, nil
}

func ParseTokenStateWithExtensions(accountData []byte) (*TokenStateWithExtensions, error) {
	tokenState, err := ParseAccount_TokenAccount(accountData)
	if err != nil {
		return nil, err
	}
	res := &TokenStateWithExtensions{Base: *tokenState}
	if len(accountData) == tokenBaseLen {
		return res, nil
	}
	if accountData[tokenBaseLen] != byte(AccountTypeAccount) {
		return nil, fmt.Errorf("unexpected token account type marker at offset %d: %d", tokenBaseLen, accountData[tokenBaseLen])
	}
	exts, err := ParseTLVExtensions(accountData[tokenBaseLen+1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse token extensions: %w", err)
	}
	res.Extensions = exts
	return res, nil
}

func ParseTLVExtensions(data []byte) ([]TLVExtension, error) {
	if len(data) == 0 {
		return nil, nil
	}
	extensions := make([]TLVExtension, 0)
	offset := 0
	for {
		if offset == len(data) {
			break
		}
		if len(data)-offset < tlvTypeLen {
			if allZero(data[offset:]) {
				break
			}
			return nil, fmt.Errorf("invalid TLV header at offset %d", offset)
		}
		typeRaw := stdbinary.LittleEndian.Uint16(data[offset : offset+tlvTypeLen])
		extType := ExtensionType(typeRaw)
		if extType == ExtensionTypeUninitialized {
			break
		}
		if len(data)-offset < tlvTypeLen+tlvLengthLen {
			return nil, fmt.Errorf("invalid TLV header length at offset %d", offset)
		}
		length := stdbinary.LittleEndian.Uint16(data[offset+tlvTypeLen : offset+tlvTypeLen+tlvLengthLen])
		offset += tlvTypeLen + tlvLengthLen

		if offset+int(length) > len(data) {
			return nil, fmt.Errorf("invalid TLV length at offset %d: %d", offset-4, length)
		}
		entryData := make([]byte, int(length))
		copy(entryData, data[offset:offset+int(length)])
		offset += int(length)

		parsedValue, _ := ParseExtensionData(extType, entryData)
		extensions = append(extensions, TLVExtension{
			Type:   extType,
			Length: length,
			Data:   entryData,
			Value:  parsedValue,
		})
	}
	return extensions, nil
}

func (s *MintStateWithExtensions) FindExtension(extType ExtensionType) *TLVExtension {
	for i := range s.Extensions {
		if s.Extensions[i].Type == extType {
			return &s.Extensions[i]
		}
	}
	return nil
}

func (s *TokenStateWithExtensions) FindExtension(extType ExtensionType) *TLVExtension {
	for i := range s.Extensions {
		if s.Extensions[i].Type == extType {
			return &s.Extensions[i]
		}
	}
	return nil
}

func (s *MintStateWithExtensions) GetMetadataPointer() (MetadataPointer, bool) {
	ext := s.FindExtension(ExtensionTypeMetadataPointer)
	if ext == nil {
		return MetadataPointer{}, false
	}
	v, err := ParseMetadataPointerExtension(ext.Data)
	if err != nil {
		return MetadataPointer{}, false
	}
	return v, true
}

func ParseMetadataPointerExtension(data []byte) (MetadataPointer, error) {
	if len(data) < 64 {
		return MetadataPointer{}, fmt.Errorf("metadata pointer extension too short: %d", len(data))
	}
	authority := decodeOptionalNonZeroPubkey(data[:32])
	metadataAddress := decodeOptionalNonZeroPubkey(data[32:64])
	var authorityValue solanago.PublicKey
	var metadataAddressValue solanago.PublicKey
	if authority != nil {
		authorityValue = *authority
	}
	if metadataAddress != nil {
		metadataAddressValue = *metadataAddress
	}
	return MetadataPointer{
		Authority:       authorityValue,
		MetadataAddress: metadataAddressValue,
	}, nil
}

func allZero(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

func decodeOptionalNonZeroPubkey(data []byte) *solanago.PublicKey {
	if len(data) != 32 || allZero(data) {
		return nil
	}
	pk := solanago.PublicKeyFromBytes(data)
	return &pk
}
