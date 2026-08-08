package mplbubblegum

import solanago "github.com/gagliardetto/solana-go"

// ProgramID is the Metaplex Bubblegum program (compressed NFTs).
//
// The address is often mistyped with a lowercase "m" in "BGUmAp9..." — that
// account does not exist on-chain. The correct address has a capital M in
// "BGUM".
var ProgramID = solanago.MustPublicKeyFromBase58("BGUMAp9Gq7iTEuizy4pqaxsTyUCBK68MDfK752saRPUY")

// SPLNoopProgramID is the SPL Noop program, used by V1 (transfer) to log
// leaf data for indexers via CPI self-invocation.
var SPLNoopProgramID = solanago.MustPublicKeyFromBase58("noopb9bkMVfRPU8AsbpTUg8AQkHtKwMYZiFUjNRtMmV")

// SPLAccountCompressionProgramID is the SPL Account Compression program,
// used by V1 (transfer) to verify and mutate the merkle tree.
var SPLAccountCompressionProgramID = solanago.MustPublicKeyFromBase58("cmtDvXumGCrqC1Age74AVPhSRVXJMd8PJS91L8KbNCK")

// MPLNoopProgramID is the MPL Account Compression fork's noop program, used
// by V2 (transfer_v2) to log leaf data for indexers via CPI self-invocation.
var MPLNoopProgramID = solanago.MustPublicKeyFromBase58("mnoopTCrg4p8ry25e4bcWA9XZjbNjMTfgYVGGEdRsf3")

// MPLAccountCompressionProgramID is the MPL Account Compression fork, used
// by V2 (transfer_v2) to verify and mutate the merkle tree. V2 requires this
// fork rather than the SPL program because it added MPL Core collection
// verification to the compression layer.
var MPLAccountCompressionProgramID = solanago.MustPublicKeyFromBase58("mcmt6YrQEMKw8Mw43FmpRLmf7BqRnFMKmAcbxE3xkAW")

// FindTreeConfigAddress derives the tree config PDA for a given merkle tree:
// seeds [merkleTree].
func FindTreeConfigAddress(merkleTree solanago.PublicKey) (solanago.PublicKey, error) {
	addr, _, err := solanago.FindProgramAddress(
		[][]byte{merkleTree.Bytes()},
		ProgramID,
	)
	return addr, err
}
