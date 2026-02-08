package token2022

import (
	stdbinary "encoding/binary"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

type AccountState uint8

const (
	AccountStateUninitialized AccountState = iota
	AccountStateInitialized
	AccountStateFrozen
)

type AccountType uint8

const (
	AccountTypeUninitialized AccountType = iota
	AccountTypeMint
	AccountTypeAccount
)

type AuthorityType uint8

const (
	AuthorityMintTokens AuthorityType = iota
	AuthorityFreezeAccount
	AuthorityAccountOwner
	AuthorityCloseAccount
	AuthorityTransferFeeConfig
	AuthorityWithheldWithdraw
	AuthorityCloseMint
	AuthorityInterestRate
	AuthorityPermanentDelegate
	AuthorityConfidentialTransferMint
	AuthorityTransferHookProgramID
	AuthorityConfidentialTransferFeeConfig
	AuthorityMetadataPointer
	AuthorityGroupPointer
	AuthorityGroupMemberPointer
	AuthorityScaledUIAmount
	AuthorityPause
	AuthorityPermissionedBurn
)

// ExtensionType values match the on-chain Token-2022 enum (u16).
type ExtensionType uint16

const (
	ExtensionTypeUninitialized ExtensionType = iota
	ExtensionTypeTransferFeeConfig
	ExtensionTypeTransferFeeAmount
	ExtensionTypeMintCloseAuthority
	ExtensionTypeConfidentialTransferMint
	ExtensionTypeConfidentialTransferAccount
	ExtensionTypeDefaultAccountState
	ExtensionTypeImmutableOwner
	ExtensionTypeMemoTransfer
	ExtensionTypeNonTransferable
	ExtensionTypeInterestBearingConfig
	ExtensionTypeCpiGuard
	ExtensionTypePermanentDelegate
	ExtensionTypeNonTransferableAccount
	ExtensionTypeTransferHook
	ExtensionTypeTransferHookAccount
	ExtensionTypeConfidentialTransferFeeConfig
	ExtensionTypeConfidentialTransferFeeAmount
	ExtensionTypeMetadataPointer
	ExtensionTypeTokenMetadata
	ExtensionTypeGroupPointer
	ExtensionTypeTokenGroup
	ExtensionTypeGroupMemberPointer
	ExtensionTypeTokenGroupMember
	ExtensionTypeConfidentialMintBurn
	ExtensionTypeScaledUIAmount
	ExtensionTypePausable
	ExtensionTypePausableAccount
	ExtensionTypePermissionedBurn
)

// Backward-compatible aliases:
const (
	ExtensionTypeConfidentialTransferFee = ExtensionTypeConfidentialTransferFeeConfig
	ExtensionTypeScaledUiAmount          = ExtensionTypeScaledUIAmount
	ExtensionTypeScaledUIAmountConfig    = ExtensionTypeScaledUIAmount
	ExtensionTypePausableConfig          = ExtensionTypePausable
)

type TransferFee struct {
	Epoch                  uint64
	MaximumFee             uint64
	TransferFeeBasisPoints uint16
}

// ElGamal ciphertext containing an account balance (64 bytes).
type EncryptedBalance [64]byte

// Authenticated encryption containing an account balance (36 bytes).
type DecryptableBalance [36]byte

type TokenMetadataFieldType uint8

const (
	TokenMetadataFieldName TokenMetadataFieldType = iota
	TokenMetadataFieldSymbol
	TokenMetadataFieldURI
	TokenMetadataFieldKey
)

// TokenMetadataField is encoded as a Borsh enum:
// 0=Name, 1=Symbol, 2=URI, 3=Key(string).
type TokenMetadataField struct {
	Type TokenMetadataFieldType
	Key  string
}

func (f TokenMetadataField) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	if err := encoder.WriteUint8(uint8(f.Type)); err != nil {
		return err
	}
	if f.Type == TokenMetadataFieldKey {
		if err := encoder.Encode(f.Key); err != nil {
			return err
		}
	}
	return nil
}

func (f *TokenMetadataField) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	t, err := decoder.ReadUint8()
	if err != nil {
		return err
	}
	f.Type = TokenMetadataFieldType(t)
	if f.Type == TokenMetadataFieldKey {
		if err := decoder.Decode(&f.Key); err != nil {
			return err
		}
	}
	return nil
}

type MintAccount struct {
	MintAuthority   *solanago.PublicKey
	Supply          uint64
	Decimals        uint8
	IsInitialized   bool
	FreezeAuthority *solanago.PublicKey
}

func (m MintAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	if err := writeCOptionPubkey(encoder, m.MintAuthority); err != nil {
		return err
	}
	if err := encoder.WriteUint64(m.Supply, ag_binary.LE); err != nil {
		return err
	}
	if err := encoder.WriteUint8(m.Decimals); err != nil {
		return err
	}
	if err := encoder.WriteBool(m.IsInitialized); err != nil {
		return err
	}
	if err := writeCOptionPubkey(encoder, m.FreezeAuthority); err != nil {
		return err
	}
	return nil
}

func (m *MintAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	mintAuthority, err := readCOptionPubkey(decoder)
	if err != nil {
		return err
	}
	m.MintAuthority = mintAuthority

	supply, err := decoder.ReadUint64(ag_binary.LE)
	if err != nil {
		return err
	}
	m.Supply = supply

	decimals, err := decoder.ReadUint8()
	if err != nil {
		return err
	}
	m.Decimals = decimals

	initialized, err := decoder.ReadBool()
	if err != nil {
		return err
	}
	m.IsInitialized = initialized

	freezeAuthority, err := readCOptionPubkey(decoder)
	if err != nil {
		return err
	}
	m.FreezeAuthority = freezeAuthority
	return nil
}

type TokenAccount struct {
	Mint            solanago.PublicKey
	Owner           solanago.PublicKey
	Amount          uint64
	Delegate        *solanago.PublicKey
	State           AccountState
	IsNative        *uint64
	DelegatedAmount uint64
	CloseAuthority  *solanago.PublicKey
}

func (a TokenAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	if err := encoder.WriteBytes(a.Mint[:], false); err != nil {
		return err
	}
	if err := encoder.WriteBytes(a.Owner[:], false); err != nil {
		return err
	}
	if err := encoder.WriteUint64(a.Amount, ag_binary.LE); err != nil {
		return err
	}
	if err := writeCOptionPubkey(encoder, a.Delegate); err != nil {
		return err
	}
	if err := encoder.WriteUint8(uint8(a.State)); err != nil {
		return err
	}
	if err := writeCOptionU64(encoder, a.IsNative); err != nil {
		return err
	}
	if err := encoder.WriteUint64(a.DelegatedAmount, ag_binary.LE); err != nil {
		return err
	}
	if err := writeCOptionPubkey(encoder, a.CloseAuthority); err != nil {
		return err
	}
	return nil
}

func (a *TokenAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	mintBytes, err := decoder.ReadNBytes(32)
	if err != nil {
		return err
	}
	a.Mint = solanago.PublicKeyFromBytes(mintBytes)

	ownerBytes, err := decoder.ReadNBytes(32)
	if err != nil {
		return err
	}
	a.Owner = solanago.PublicKeyFromBytes(ownerBytes)

	amount, err := decoder.ReadUint64(ag_binary.LE)
	if err != nil {
		return err
	}
	a.Amount = amount

	delegate, err := readCOptionPubkey(decoder)
	if err != nil {
		return err
	}
	a.Delegate = delegate

	state, err := decoder.ReadUint8()
	if err != nil {
		return err
	}
	a.State = AccountState(state)

	isNative, err := readCOptionU64(decoder)
	if err != nil {
		return err
	}
	a.IsNative = isNative

	delegatedAmount, err := decoder.ReadUint64(ag_binary.LE)
	if err != nil {
		return err
	}
	a.DelegatedAmount = delegatedAmount

	closeAuthority, err := readCOptionPubkey(decoder)
	if err != nil {
		return err
	}
	a.CloseAuthority = closeAuthority
	return nil
}

type MultisigAccount struct {
	M             uint8
	N             uint8
	IsInitialized bool
	Signers       [MAX_SIGNERS]solanago.PublicKey
}

func (m MultisigAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	if err := encoder.WriteUint8(m.M); err != nil {
		return err
	}
	if err := encoder.WriteUint8(m.N); err != nil {
		return err
	}
	if err := encoder.WriteBool(m.IsInitialized); err != nil {
		return err
	}
	for _, signer := range m.Signers {
		if err := encoder.WriteBytes(signer[:], false); err != nil {
			return err
		}
	}
	return nil
}

func (m *MultisigAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	v, err := decoder.ReadUint8()
	if err != nil {
		return err
	}
	m.M = v

	v, err = decoder.ReadUint8()
	if err != nil {
		return err
	}
	m.N = v

	initialized, err := decoder.ReadBool()
	if err != nil {
		return err
	}
	m.IsInitialized = initialized

	for i := 0; i < MAX_SIGNERS; i++ {
		pk, err := decoder.ReadNBytes(32)
		if err != nil {
			return err
		}
		m.Signers[i] = solanago.PublicKeyFromBytes(pk)
	}
	return nil
}

func writeCOptionPubkey(encoder *ag_binary.Encoder, key *solanago.PublicKey) error {
	if key == nil {
		if err := encoder.WriteUint32(0, stdbinary.LittleEndian); err != nil {
			return err
		}
		var zero solanago.PublicKey
		return encoder.WriteBytes(zero[:], false)
	}
	if err := encoder.WriteUint32(1, stdbinary.LittleEndian); err != nil {
		return err
	}
	return encoder.WriteBytes(key[:], false)
}

func readCOptionPubkey(decoder *ag_binary.Decoder) (*solanago.PublicKey, error) {
	present, err := decoder.ReadUint32(stdbinary.LittleEndian)
	if err != nil {
		return nil, err
	}
	pk, err := decoder.ReadNBytes(32)
	if err != nil {
		return nil, err
	}
	if present == 1 {
		v := solanago.PublicKeyFromBytes(pk)
		return &v, nil
	}
	if present != 0 {
		return nil, fmt.Errorf("invalid coption discriminant for pubkey: %d", present)
	}
	return nil, nil
}

func writeCOptionU64(encoder *ag_binary.Encoder, value *uint64) error {
	if value == nil {
		if err := encoder.WriteUint32(0, stdbinary.LittleEndian); err != nil {
			return err
		}
		return encoder.WriteUint64(0, ag_binary.LE)
	}
	if err := encoder.WriteUint32(1, stdbinary.LittleEndian); err != nil {
		return err
	}
	return encoder.WriteUint64(*value, ag_binary.LE)
}

func readCOptionU64(decoder *ag_binary.Decoder) (*uint64, error) {
	present, err := decoder.ReadUint32(stdbinary.LittleEndian)
	if err != nil {
		return nil, err
	}
	value, err := decoder.ReadUint64(ag_binary.LE)
	if err != nil {
		return nil, err
	}
	if present == 1 {
		v := value
		return &v, nil
	}
	if present != 0 {
		return nil, fmt.Errorf("invalid coption discriminant for u64: %d", present)
	}
	return nil, nil
}
