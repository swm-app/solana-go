package token2022

import (
	"crypto/sha256"
	stdbinary "encoding/binary"
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
)

const (
	transferHookExecuteDiscriminatorLen    = 8
	transferHookTLVLengthLen               = 4
	transferHookExtraAccountMetaLen        = 35
	transferHookExtraAccountMetaCountLen   = 4
	transferHookExtraAccountMetaSeed       = "extra-account-metas"
	transferHookInterfaceDiscriminatorSeed = "spl-transfer-hook-interface"
)

const (
	transferHookExtraAccountMetaStatic          = 0
	transferHookExtraAccountMetaSeeds           = 1
	transferHookExtraAccountMetaPubkeyData      = 2
	transferHookExtraAccountMetaExternalPDABase = 128
)

const (
	transferHookSeedUninitialized = 0
	transferHookSeedLiteral       = 1
	transferHookSeedInstruction   = 2
	transferHookSeedAccountKey    = 3
	transferHookSeedAccountData   = 4
)

var transferHookExecuteDiscriminator = deriveTransferHookInstructionDiscriminator("execute")

type TransferHookTransfer struct {
	Source      solanago.PublicKey
	Mint        solanago.PublicKey
	Destination solanago.PublicKey
	Authority   solanago.PublicKey
	Amount      uint64
}

type TransferHookResolvedAccounts struct {
	ProgramID         solanago.PublicKey
	ValidationAccount solanago.PublicKey
	ExtraAccounts     []*solanago.AccountMeta
}

func (r *TransferHookResolvedAccounts) RemainingAccounts() []*solanago.AccountMeta {
	if r == nil {
		return nil
	}
	out := make([]*solanago.AccountMeta, 0, 1+len(r.ExtraAccounts))
	out = append(out, solanago.NewAccountMeta(r.ValidationAccount, false, false))
	for _, meta := range r.ExtraAccounts {
		out = append(out, solanago.NewAccountMeta(meta.PublicKey, meta.IsWritable, meta.IsSigner))
	}
	return out
}

func GetTransferHookProgramID(mintState *MintStateWithExtensions) (solanago.PublicKey, bool) {
	if mintState == nil {
		return solanago.PublicKey{}, false
	}
	ext := mintState.FindExtension(ExtensionTypeTransferHook)
	if ext == nil {
		return solanago.PublicKey{}, false
	}
	parsed, err := ParseExtensionTransferHook(ext.Data)
	if err != nil || parsed == nil || parsed.ProgramID.IsZero() {
		return solanago.PublicKey{}, false
	}
	return parsed.ProgramID, true
}

func ParseTransferHookProgramID(mintData []byte) (solanago.PublicKey, bool, error) {
	mintState, err := ParseMintStateWithExtensions(mintData)
	if err != nil {
		return solanago.PublicKey{}, false, fmt.Errorf("parse mint state with extensions: %w", err)
	}
	programID, ok := GetTransferHookProgramID(mintState)
	return programID, ok, nil
}

func GetTransferHookExtraAccountMetaListAddress(
	mint solanago.PublicKey,
	hookProgramID solanago.PublicKey,
) (solanago.PublicKey, error) {
	address, _, err := solanago.FindProgramAddress(
		[][]byte{
			[]byte(transferHookExtraAccountMetaSeed),
			mint[:],
		},
		hookProgramID,
	)
	if err != nil {
		return solanago.PublicKey{}, err
	}
	return address, nil
}

func ResolveTransferHookAccounts(
	transfer TransferHookTransfer,
	mintData []byte,
	validationAccountData []byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (*TransferHookResolvedAccounts, error) {
	if err := validateTransferHookTransfer(transfer); err != nil {
		return nil, err
	}
	hookProgramID, ok, err := ParseTransferHookProgramID(mintData)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return ResolveTransferHookAccountsFromMintStateData(
		transfer,
		hookProgramID,
		validationAccountData,
		additionalAccountData,
	)
}

func ResolveTransferHookAccountsFromMintStateData(
	transfer TransferHookTransfer,
	hookProgramID solanago.PublicKey,
	validationAccountData []byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (*TransferHookResolvedAccounts, error) {
	if hookProgramID.IsZero() {
		return nil, fmt.Errorf("transfer hook program id is required")
	}
	if err := validateTransferHookTransfer(transfer); err != nil {
		return nil, err
	}
	validationAccount, err := GetTransferHookExtraAccountMetaListAddress(transfer.Mint, hookProgramID)
	if err != nil {
		return nil, fmt.Errorf("derive transfer hook validation account: %w", err)
	}
	extraMetas, err := parseTransferHookValidationAccount(validationAccountData)
	if err != nil {
		return nil, fmt.Errorf("parse transfer hook validation account data: %w", err)
	}

	accountDataCache := map[solanago.PublicKey][]byte{}
	baseAccounts := []*solanago.AccountMeta{
		solanago.NewAccountMeta(transfer.Source, false, false),
		solanago.NewAccountMeta(transfer.Mint, false, false),
		solanago.NewAccountMeta(transfer.Destination, false, false),
		solanago.NewAccountMeta(transfer.Authority, false, false),
		solanago.NewAccountMeta(validationAccount, false, false),
	}
	resolvedExtraAccounts := make([]*solanago.AccountMeta, 0, len(extraMetas))
	currentAccounts := append([]*solanago.AccountMeta{}, baseAccounts...)
	instructionData := encodeTransferHookExecuteInstructionData(transfer.Amount)
	for i, extraMeta := range extraMetas {
		resolvedMeta, err := extraMeta.resolve(
			instructionData,
			hookProgramID,
			currentAccounts,
			accountDataCache,
			additionalAccountData,
		)
		if err != nil {
			return nil, fmt.Errorf("resolve extra account meta %d: %w", i, err)
		}
		resolvedExtraAccounts = append(resolvedExtraAccounts, resolvedMeta)
		currentAccounts = append(currentAccounts, resolvedMeta)
	}

	return &TransferHookResolvedAccounts{
		ProgramID:         hookProgramID,
		ValidationAccount: validationAccount,
		ExtraAccounts:     resolvedExtraAccounts,
	}, nil
}

func AddTransferHookAccountsToInstruction(
	instruction solanago.Instruction,
	transfer TransferHookTransfer,
	mintData []byte,
	validationAccountData []byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (solanago.Instruction, error) {
	resolved, err := ResolveTransferHookAccounts(
		transfer,
		mintData,
		validationAccountData,
		additionalAccountData,
	)
	if err != nil {
		return nil, err
	}
	if resolved == nil {
		return instruction, nil
	}
	return appendInstructionAccounts(instruction, resolved.RemainingAccounts())
}

func NewTransferCheckedInstructionWithTransferHook(
	params TransferCheckedParams,
	accounts TransferCheckedAccounts,
	mintData []byte,
	validationAccountData []byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (solanago.Instruction, error) {
	instruction, err := NewTransferCheckedInstruction(params, accounts)
	if err != nil {
		return nil, err
	}
	return AddTransferHookAccountsToInstruction(
		instruction,
		TransferHookTransfer{
			Source:      accounts.Source,
			Mint:        accounts.Mint,
			Destination: accounts.Destination,
			Authority:   accounts.Authority,
			Amount:      params.Amount,
		},
		mintData,
		validationAccountData,
		additionalAccountData,
	)
}

func NewTransferCheckedWithFeeInstructionWithTransferHook(
	params TransferCheckedWithFeeParams,
	accounts TransferCheckedWithFeeAccounts,
	mintData []byte,
	validationAccountData []byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (solanago.Instruction, error) {
	instruction, err := NewTransferCheckedWithFeeInstruction(params, accounts)
	if err != nil {
		return nil, err
	}
	return AddTransferHookAccountsToInstruction(
		instruction,
		TransferHookTransfer{
			Source:      accounts.Source,
			Mint:        accounts.Mint,
			Destination: accounts.Destination,
			Authority:   accounts.Authority,
			Amount:      params.Amount,
		},
		mintData,
		validationAccountData,
		additionalAccountData,
	)
}

type transferHookExtraAccountMeta struct {
	Discriminator uint8
	AddressConfig [32]byte
	IsSigner      bool
	IsWritable    bool
}

type transferHookSeed struct {
	kind         uint8
	literal      []byte
	index        uint8
	accountIndex uint8
	dataIndex    uint8
	length       uint8
}

func (m transferHookExtraAccountMeta) resolve(
	instructionData []byte,
	programID solanago.PublicKey,
	currentAccounts []*solanago.AccountMeta,
	accountDataCache map[solanago.PublicKey][]byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (*solanago.AccountMeta, error) {
	switch {
	case m.Discriminator == transferHookExtraAccountMetaStatic:
		return solanago.NewAccountMeta(
			solanago.PublicKeyFromBytes(m.AddressConfig[:]),
			m.IsWritable,
			m.IsSigner,
		), nil
	case m.Discriminator == transferHookExtraAccountMetaSeeds:
		address, err := resolveTransferHookPDA(
			instructionData,
			programID,
			currentAccounts,
			m.AddressConfig,
			accountDataCache,
			additionalAccountData,
		)
		if err != nil {
			return nil, err
		}
		return solanago.NewAccountMeta(address, m.IsWritable, m.IsSigner), nil
	case m.Discriminator >= transferHookExtraAccountMetaExternalPDABase:
		programIndex := int(m.Discriminator - transferHookExtraAccountMetaExternalPDABase)
		if programIndex >= len(currentAccounts) {
			return nil, fmt.Errorf("external PDA program index %d out of range", programIndex)
		}
		address, err := resolveTransferHookPDA(
			instructionData,
			currentAccounts[programIndex].PublicKey,
			currentAccounts,
			m.AddressConfig,
			accountDataCache,
			additionalAccountData,
		)
		if err != nil {
			return nil, err
		}
		return solanago.NewAccountMeta(address, m.IsWritable, m.IsSigner), nil
	case m.Discriminator == transferHookExtraAccountMetaPubkeyData:
		return nil, fmt.Errorf("pubkey-data transfer hook account metas are not supported yet")
	default:
		return nil, fmt.Errorf("unsupported transfer hook extra account meta discriminator: %d", m.Discriminator)
	}
}

func parseTransferHookValidationAccount(data []byte) ([]transferHookExtraAccountMeta, error) {
	offset := 0
	for offset < len(data) {
		if allZero(data[offset:]) {
			break
		}
		if len(data)-offset < transferHookExecuteDiscriminatorLen+transferHookTLVLengthLen {
			return nil, fmt.Errorf("invalid TLV header at offset %d", offset)
		}
		discriminator := data[offset : offset+transferHookExecuteDiscriminatorLen]
		offset += transferHookExecuteDiscriminatorLen
		entryLen := int(stdbinary.LittleEndian.Uint32(data[offset : offset+transferHookTLVLengthLen]))
		offset += transferHookTLVLengthLen
		if entryLen < 0 || offset+entryLen > len(data) {
			return nil, fmt.Errorf("invalid TLV entry length %d", entryLen)
		}
		entryData := data[offset : offset+entryLen]
		offset += entryLen
		if !equalBytes(discriminator, transferHookExecuteDiscriminator[:]) {
			continue
		}
		return parseTransferHookExtraAccountMetas(entryData)
	}
	return nil, fmt.Errorf("execute extra account meta list not found")
}

func parseTransferHookExtraAccountMetas(data []byte) ([]transferHookExtraAccountMeta, error) {
	if len(data) < transferHookExtraAccountMetaCountLen {
		return nil, fmt.Errorf("extra account meta list is too short: %d", len(data))
	}
	count := int(stdbinary.LittleEndian.Uint32(data[:transferHookExtraAccountMetaCountLen]))
	expectedLen := transferHookExtraAccountMetaCountLen + count*transferHookExtraAccountMetaLen
	if len(data) < expectedLen {
		return nil, fmt.Errorf("extra account meta list has invalid size: want at least %d, got %d", expectedLen, len(data))
	}
	out := make([]transferHookExtraAccountMeta, 0, count)
	offset := transferHookExtraAccountMetaCountLen
	for i := 0; i < count; i++ {
		chunk := data[offset : offset+transferHookExtraAccountMetaLen]
		var addressConfig [32]byte
		copy(addressConfig[:], chunk[1:33])
		out = append(out, transferHookExtraAccountMeta{
			Discriminator: chunk[0],
			AddressConfig: addressConfig,
			IsSigner:      chunk[33] != 0,
			IsWritable:    chunk[34] != 0,
		})
		offset += transferHookExtraAccountMetaLen
	}
	return out, nil
}

func resolveTransferHookPDA(
	instructionData []byte,
	programID solanago.PublicKey,
	currentAccounts []*solanago.AccountMeta,
	addressConfig [32]byte,
	accountDataCache map[solanago.PublicKey][]byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) (solanago.PublicKey, error) {
	seeds, err := unpackTransferHookSeedConfigs(addressConfig)
	if err != nil {
		return solanago.PublicKey{}, err
	}
	seedBytes := make([][]byte, 0, len(seeds))
	for _, seed := range seeds {
		value, err := resolveTransferHookSeed(
			instructionData,
			currentAccounts,
			seed,
			accountDataCache,
			additionalAccountData,
		)
		if err != nil {
			return solanago.PublicKey{}, err
		}
		seedBytes = append(seedBytes, value)
	}
	address, _, err := solanago.FindProgramAddress(seedBytes, programID)
	if err != nil {
		return solanago.PublicKey{}, err
	}
	return address, nil
}

func unpackTransferHookSeedConfigs(addressConfig [32]byte) ([]transferHookSeed, error) {
	out := make([]transferHookSeed, 0)
	offset := 0
	for offset < len(addressConfig) {
		discriminator := addressConfig[offset]
		if discriminator == transferHookSeedUninitialized {
			break
		}
		switch discriminator {
		case transferHookSeedLiteral:
			if offset+2 > len(addressConfig) {
				return nil, fmt.Errorf("invalid literal seed config at offset %d", offset)
			}
			literalLen := int(addressConfig[offset+1])
			if offset+2+literalLen > len(addressConfig) {
				return nil, fmt.Errorf("invalid literal seed size at offset %d", offset)
			}
			literal := make([]byte, literalLen)
			copy(literal, addressConfig[offset+2:offset+2+literalLen])
			out = append(out, transferHookSeed{
				kind:    transferHookSeedLiteral,
				literal: literal,
			})
			offset += 2 + literalLen
		case transferHookSeedInstruction:
			if offset+3 > len(addressConfig) {
				return nil, fmt.Errorf("invalid instruction-data seed config at offset %d", offset)
			}
			out = append(out, transferHookSeed{
				kind:   transferHookSeedInstruction,
				index:  addressConfig[offset+1],
				length: addressConfig[offset+2],
			})
			offset += 3
		case transferHookSeedAccountKey:
			if offset+2 > len(addressConfig) {
				return nil, fmt.Errorf("invalid account-key seed config at offset %d", offset)
			}
			out = append(out, transferHookSeed{
				kind:  transferHookSeedAccountKey,
				index: addressConfig[offset+1],
			})
			offset += 2
		case transferHookSeedAccountData:
			if offset+4 > len(addressConfig) {
				return nil, fmt.Errorf("invalid account-data seed config at offset %d", offset)
			}
			out = append(out, transferHookSeed{
				kind:         transferHookSeedAccountData,
				accountIndex: addressConfig[offset+1],
				dataIndex:    addressConfig[offset+2],
				length:       addressConfig[offset+3],
			})
			offset += 4
		default:
			return nil, fmt.Errorf("unsupported seed discriminator %d at offset %d", discriminator, offset)
		}
	}
	return out, nil
}

func resolveTransferHookSeed(
	instructionData []byte,
	currentAccounts []*solanago.AccountMeta,
	seed transferHookSeed,
	accountDataCache map[solanago.PublicKey][]byte,
	additionalAccountData map[solanago.PublicKey][]byte,
) ([]byte, error) {
	switch seed.kind {
	case transferHookSeedLiteral:
		return seed.literal, nil
	case transferHookSeedInstruction:
		start := int(seed.index)
		end := start + int(seed.length)
		if end > len(instructionData) {
			return nil, fmt.Errorf("instruction-data seed [%d:%d] out of range for data len=%d", start, end, len(instructionData))
		}
		out := make([]byte, int(seed.length))
		copy(out, instructionData[start:end])
		return out, nil
	case transferHookSeedAccountKey:
		index := int(seed.index)
		if index >= len(currentAccounts) {
			return nil, fmt.Errorf("account-key seed index %d out of range", index)
		}
		key := currentAccounts[index].PublicKey
		out := make([]byte, solanago.PublicKeyLength)
		copy(out, key[:])
		return out, nil
	case transferHookSeedAccountData:
		index := int(seed.accountIndex)
		if index >= len(currentAccounts) {
			return nil, fmt.Errorf("account-data seed account index %d out of range", index)
		}
		accountKey := currentAccounts[index].PublicKey
		data, ok := accountDataCache[accountKey]
		if !ok {
			if additionalAccountData == nil {
				return nil, fmt.Errorf("missing account data for %s required by account-data seed", accountKey)
			}
			provided, ok := additionalAccountData[accountKey]
			if !ok {
				return nil, fmt.Errorf("missing account data for %s required by account-data seed", accountKey)
			}
			accountDataCache[accountKey] = provided
			data = provided
		}
		start := int(seed.dataIndex)
		end := start + int(seed.length)
		if end > len(data) {
			return nil, fmt.Errorf("account-data seed [%d:%d] out of range for account %s with data len=%d", start, end, accountKey, len(data))
		}
		out := make([]byte, int(seed.length))
		copy(out, data[start:end])
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported seed kind %d", seed.kind)
	}
}

func appendInstructionAccounts(
	instruction solanago.Instruction,
	extraAccounts []*solanago.AccountMeta,
) (solanago.Instruction, error) {
	data, err := instruction.Data()
	if err != nil {
		return nil, err
	}
	accounts := make(solanago.AccountMetaSlice, 0, len(instruction.Accounts())+len(extraAccounts))
	accounts = append(accounts, instruction.Accounts()...)
	accounts = append(accounts, extraAccounts...)
	return solanago.NewInstruction(instruction.ProgramID(), accounts, data), nil
}

func encodeTransferHookExecuteInstructionData(amount uint64) []byte {
	out := make([]byte, transferHookExecuteDiscriminatorLen+8)
	copy(out[:transferHookExecuteDiscriminatorLen], transferHookExecuteDiscriminator[:])
	stdbinary.LittleEndian.PutUint64(out[transferHookExecuteDiscriminatorLen:], amount)
	return out
}

func deriveTransferHookInstructionDiscriminator(name string) [transferHookExecuteDiscriminatorLen]byte {
	sum := sha256.Sum256([]byte(transferHookInterfaceDiscriminatorSeed + ":" + name))
	var out [transferHookExecuteDiscriminatorLen]byte
	copy(out[:], sum[:transferHookExecuteDiscriminatorLen])
	return out
}

func validateTransferHookTransfer(transfer TransferHookTransfer) error {
	switch {
	case transfer.Source.IsZero():
		return fmt.Errorf("transfer hook source is required")
	case transfer.Mint.IsZero():
		return fmt.Errorf("transfer hook mint is required")
	case transfer.Destination.IsZero():
		return fmt.Errorf("transfer hook destination is required")
	case transfer.Authority.IsZero():
		return fmt.Errorf("transfer hook authority is required")
	default:
		return nil
	}
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
