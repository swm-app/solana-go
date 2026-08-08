package mplbubblegum

import (
	"bytes"
	"fmt"

	binary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

// TransferArgs are the borsh-encoded arguments of a "transfer" (V1)
// instruction.
type TransferArgs struct {
	Root        [32]uint8
	DataHash    [32]uint8
	CreatorHash [32]uint8
	Nonce       uint64
	Index       uint32
}

func (obj TransferArgs) MarshalWithEncoder(encoder *binary.Encoder) (err error) {
	if err = encoder.Encode(obj.Root); err != nil {
		return newFieldError("Root", err)
	}
	if err = encoder.Encode(obj.DataHash); err != nil {
		return newFieldError("DataHash", err)
	}
	if err = encoder.Encode(obj.CreatorHash); err != nil {
		return newFieldError("CreatorHash", err)
	}
	if err = encoder.Encode(obj.Nonce); err != nil {
		return newFieldError("Nonce", err)
	}
	if err = encoder.Encode(obj.Index); err != nil {
		return newFieldError("Index", err)
	}
	return nil
}

// TransferV2Args are the borsh-encoded arguments of a "transfer_v2"
// instruction.
type TransferV2Args struct {
	Root          [32]uint8
	DataHash      [32]uint8
	CreatorHash   [32]uint8
	AssetDataHash *[32]uint8 `bin:"optional"`
	Flags         *uint8     `bin:"optional"`
	Nonce         uint64
	Index         uint32
}

func (obj TransferV2Args) MarshalWithEncoder(encoder *binary.Encoder) (err error) {
	if err = encoder.Encode(obj.Root); err != nil {
		return newFieldError("Root", err)
	}
	if err = encoder.Encode(obj.DataHash); err != nil {
		return newFieldError("DataHash", err)
	}
	if err = encoder.Encode(obj.CreatorHash); err != nil {
		return newFieldError("CreatorHash", err)
	}
	// Serialize `AssetDataHash` (optional):
	{
		if obj.AssetDataHash == nil {
			if err = encoder.WriteOption(false); err != nil {
				return newOptionError("AssetDataHash", fmt.Errorf("error while encoding optionality: %w", err))
			}
		} else {
			if err = encoder.WriteOption(true); err != nil {
				return newOptionError("AssetDataHash", fmt.Errorf("error while encoding optionality: %w", err))
			}
			if err = encoder.Encode(*obj.AssetDataHash); err != nil {
				return newFieldError("AssetDataHash", err)
			}
		}
	}
	// Serialize `Flags` (optional):
	{
		if obj.Flags == nil {
			if err = encoder.WriteOption(false); err != nil {
				return newOptionError("Flags", fmt.Errorf("error while encoding optionality: %w", err))
			}
		} else {
			if err = encoder.WriteOption(true); err != nil {
				return newOptionError("Flags", fmt.Errorf("error while encoding optionality: %w", err))
			}
			if err = encoder.Encode(*obj.Flags); err != nil {
				return newFieldError("Flags", err)
			}
		}
	}
	if err = encoder.Encode(obj.Nonce); err != nil {
		return newFieldError("Nonce", err)
	}
	if err = encoder.Encode(obj.Index); err != nil {
		return newFieldError("Index", err)
	}
	return nil
}

// NewTransferInstruction builds a "transfer" (V1) instruction.
//
// This package only supports owner-authorized transfers: leafOwner is the
// signer. A delegate-authorized transfer would instead require the
// delegate's signature and is not exposed here.
//
// proof is appended as trailing read-only, non-signer account metas, one per
// merkle proof node.
func NewTransferInstruction(
	args TransferArgs,

	treeConfig solanago.PublicKey,
	leafOwner solanago.PublicKey,
	leafDelegate solanago.PublicKey,
	newLeafOwner solanago.PublicKey,
	merkleTree solanago.PublicKey,
	logWrapper solanago.PublicKey,
	compressionProgram solanago.PublicKey,
	systemProgram solanago.PublicKey,

	proof []solanago.PublicKey,
) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)

	if err := enc__.WriteBytes(Instruction_Transfer[:], false); err != nil {
		return nil, fmt.Errorf("failed to write instruction discriminator: %w", err)
	}
	if err := enc__.Encode(args); err != nil {
		return nil, newFieldError("args", err)
	}

	accounts__ := solanago.AccountMetaSlice{}
	accounts__.Append(solanago.NewAccountMeta(treeConfig, false, false))
	accounts__.Append(solanago.NewAccountMeta(leafOwner, false, true))
	accounts__.Append(solanago.NewAccountMeta(leafDelegate, false, false))
	accounts__.Append(solanago.NewAccountMeta(newLeafOwner, false, false))
	accounts__.Append(solanago.NewAccountMeta(merkleTree, true, false))
	accounts__.Append(solanago.NewAccountMeta(logWrapper, false, false))
	accounts__.Append(solanago.NewAccountMeta(compressionProgram, false, false))
	accounts__.Append(solanago.NewAccountMeta(systemProgram, false, false))
	for _, node := range proof {
		accounts__.Append(solanago.NewAccountMeta(node, false, false))
	}

	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}

// NewTransferV2Instruction builds a "transfer_v2" instruction.
//
// authority couples to leafOwner's signer flag: when authority is omitted
// (passed as ProgramID), the program falls back to treating leafOwner as the
// authority, so leafOwner must sign instead of authority. Passing a non-zero
// authority makes authority the signer and leafOwner a plain read-only
// account.
//
// proof is appended as trailing read-only, non-signer account metas, one per
// merkle proof node.
func NewTransferV2Instruction(
	args TransferV2Args,

	treeConfig solanago.PublicKey,
	payer solanago.PublicKey,
	authority solanago.PublicKey,
	leafOwner solanago.PublicKey,
	leafDelegate solanago.PublicKey,
	newLeafOwner solanago.PublicKey,
	merkleTree solanago.PublicKey,
	coreCollection solanago.PublicKey,
	logWrapper solanago.PublicKey,
	compressionProgram solanago.PublicKey,
	systemProgram solanago.PublicKey,

	proof []solanago.PublicKey,
) (solanago.Instruction, error) {
	buf__ := new(bytes.Buffer)
	enc__ := binary.NewBorshEncoder(buf__)

	if err := enc__.WriteBytes(Instruction_TransferV2[:], false); err != nil {
		return nil, fmt.Errorf("failed to write instruction discriminator: %w", err)
	}
	if err := enc__.Encode(args); err != nil {
		return nil, newFieldError("args", err)
	}

	authorityProvided := !authority.Equals(ProgramID)

	accounts__ := solanago.AccountMetaSlice{}
	accounts__.Append(solanago.NewAccountMeta(treeConfig, true, false))
	accounts__.Append(solanago.NewAccountMeta(payer, true, true))
	accounts__.Append(solanago.NewAccountMeta(authority, false, authorityProvided))
	accounts__.Append(solanago.NewAccountMeta(leafOwner, false, !authorityProvided))
	accounts__.Append(solanago.NewAccountMeta(leafDelegate, false, false))
	accounts__.Append(solanago.NewAccountMeta(newLeafOwner, false, false))
	accounts__.Append(solanago.NewAccountMeta(merkleTree, true, false))
	accounts__.Append(solanago.NewAccountMeta(coreCollection, false, false))
	accounts__.Append(solanago.NewAccountMeta(logWrapper, false, false))
	accounts__.Append(solanago.NewAccountMeta(compressionProgram, false, false))
	accounts__.Append(solanago.NewAccountMeta(systemProgram, false, false))
	for _, node := range proof {
		accounts__.Append(solanago.NewAccountMeta(node, false, false))
	}

	return solanago.NewInstruction(ProgramID, accounts__, buf__.Bytes()), nil
}
