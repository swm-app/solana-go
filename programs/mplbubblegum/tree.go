package mplbubblegum

import (
	"encoding/binary"
	"fmt"
	"math/bits"

	solanago "github.com/gagliardetto/solana-go"
)

// TreeHeaderSize is the fixed-size prefix of a ConcurrentMerkleTree account,
// shared by both the SPL and MPL Account Compression programs:
//
//	byte    0    account type (must be 1: ConcurrentMerkleTree)
//	byte    1    header version
//	bytes   2:6  max_buffer_size (u32 LE)
//	bytes   6:10 max_depth (u32 LE)
//	bytes  10:42 authority pubkey
//	bytes  42:50 creation slot (u64 LE)
//	byte   50    is_batch_initialized
//	bytes 51:56  padding
const TreeHeaderSize = 56

// concurrentMerkleTreeAccountType is the only valid value of the account
// type byte for a ConcurrentMerkleTree account.
const concurrentMerkleTreeAccountType = 1

// TreeAccount is the parsed header of a ConcurrentMerkleTree account.
type TreeAccount struct {
	MaxDepth      uint32
	MaxBufferSize uint32
	CanopyDepth   uint32
	Authority     solanago.PublicKey
	CreationSlot  uint64
}

// ParseTreeAccount parses a ConcurrentMerkleTree account from its fixed-size
// header and the account's total data length.
//
// The tree body (change log ring buffer + rightmost proof) is sized
// entirely by max_depth and max_buffer_size, so any bytes past the body are
// the canopy — cached proof nodes for the top canopyDepth levels of the
// tree, stored so clients don't need to supply those nodes in every
// instruction's proof. The canopy has no length field of its own; its size
// is derived from what's left in the account after the header and the body.
//
// Taking the length separately from the bytes lets a caller read a tree
// through a dataSlice covering only the header: a canopy runs to tens of
// kilobytes and none of it is needed here, while an RPC response reports the
// account's full data length regardless of any slice applied to its data.
// header may carry more than TreeHeaderSize bytes (the whole account, say);
// everything past the header is ignored.
func ParseTreeAccount(header []byte, accountLen uint64) (*TreeAccount, error) {
	if len(header) < TreeHeaderSize {
		return nil, fmt.Errorf("mplbubblegum: tree account data too short: got %d bytes, need at least %d", len(header), TreeHeaderSize)
	}
	if accountLen < uint64(len(header)) {
		return nil, fmt.Errorf("mplbubblegum: tree account length %d is shorter than the %d bytes of data supplied for it", accountLen, len(header))
	}

	accountType := header[0]
	if accountType != concurrentMerkleTreeAccountType {
		return nil, fmt.Errorf("mplbubblegum: unexpected tree account type %d, want %d", accountType, concurrentMerkleTreeAccountType)
	}

	maxBufferSize := binary.LittleEndian.Uint32(header[2:6])
	maxDepth := binary.LittleEndian.Uint32(header[6:10])
	var authority solanago.PublicKey
	copy(authority[:], header[10:42])
	creationSlot := binary.LittleEndian.Uint64(header[42:50])

	bodySize := int64(24) + int64(maxBufferSize)*(40+32*int64(maxDepth)) + (32*int64(maxDepth) + 40)
	canopyBytes := int64(accountLen) - TreeHeaderSize - bodySize
	if canopyBytes < 0 {
		return nil, fmt.Errorf("mplbubblegum: tree account data too short for max_depth=%d max_buffer_size=%d: need at least %d body bytes past the header", maxDepth, maxBufferSize, bodySize)
	}
	if canopyBytes%32 != 0 {
		return nil, fmt.Errorf("mplbubblegum: canopy size %d is not a multiple of 32", canopyBytes)
	}

	canopyNodes := canopyBytes/32 + 2
	if canopyNodes&(canopyNodes-1) != 0 {
		return nil, fmt.Errorf("mplbubblegum: canopy size %d does not correspond to a complete canopy (node count %d is not a power of two)", canopyBytes, canopyNodes)
	}
	canopyDepth := uint32(bits.Len64(uint64(canopyNodes)) - 2)

	return &TreeAccount{
		MaxDepth:      maxDepth,
		MaxBufferSize: maxBufferSize,
		CanopyDepth:   canopyDepth,
		Authority:     authority,
		CreationSlot:  creationSlot,
	}, nil
}

// TransferVersion identifies which Bubblegum transfer instruction a tree
// routes to, determined by which Account Compression program owns the tree.
type TransferVersion int

const (
	// TransferVersionV1 is the "transfer" instruction, used by trees owned
	// by the SPL Account Compression program.
	TransferVersionV1 TransferVersion = iota + 1
	// TransferVersionV2 is the "transfer_v2" instruction, used by trees
	// owned by the MPL Account Compression program.
	TransferVersionV2
)

// TransferVersionForTreeOwner maps a tree account's owner program to the
// Bubblegum instruction version that operates on it.
func TransferVersionForTreeOwner(owner solanago.PublicKey) (TransferVersion, error) {
	switch {
	case owner.Equals(SPLAccountCompressionProgramID):
		return TransferVersionV1, nil
	case owner.Equals(MPLAccountCompressionProgramID):
		return TransferVersionV2, nil
	default:
		return 0, fmt.Errorf("mplbubblegum: %s is not a known Account Compression program", owner)
	}
}
