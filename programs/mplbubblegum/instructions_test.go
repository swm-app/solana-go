package mplbubblegum

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

type fixtureAccount struct {
	Pos      int    `json:"pos"`
	Pubkey   string `json:"pubkey"`
	Writable bool   `json:"writable"`
	Signer   bool   `json:"signer"`
}

type fixtureArgs struct {
	Root          string  `json:"root"`
	DataHash      string  `json:"data_hash"`
	CreatorHash   string  `json:"creator_hash"`
	AssetDataHash *string `json:"asset_data_hash"`
	Flags         *uint8  `json:"flags"`
	Nonce         uint64  `json:"nonce"`
	Index         uint32  `json:"index"`
}

type fixture struct {
	DataHex  string           `json:"data_hex"`
	Args     fixtureArgs      `json:"args"`
	Accounts []fixtureAccount `json:"accounts"`
}

func loadFixture(t *testing.T, path string) fixture {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var f fixture
	require.NoError(t, json.Unmarshal(raw, &f))
	return f
}

func mustHash32(t *testing.T, s string) [32]uint8 {
	t.Helper()
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	require.Len(t, b, 32)
	var out [32]uint8
	copy(out[:], b)
	return out
}

// expectedFlags is the canonical writable/signer pair this package's
// builders must produce for one account meta.
type expectedFlags struct {
	Writable bool
	Signer   bool
}

// requireExactFlags pins each canonical account meta's writable/signer pair
// against the exact values the account tables in this package's docs
// specify, independent of what any particular fixture transaction observed.
// This is what catches an under-privileging builder bug (e.g. a missing
// writable or signer flag) that requireCanonicalSubsetOfObserved cannot: the
// subset check only ever complains about over-claimed privilege.
func requireExactFlags(t *testing.T, canonical []*solanago.AccountMeta, expected []expectedFlags) {
	t.Helper()
	require.Len(t, canonical, len(expected))
	for i, c := range canonical {
		e := expected[i]
		require.Equal(t, e.Writable, c.IsWritable, "account %d writable flag", i)
		require.Equal(t, e.Signer, c.IsSigner, "account %d signer flag", i)
	}
}

// requireCanonicalSubsetOfObserved checks that the canonical account metas
// this package builds never claim less privilege than what the fixture
// transaction's message actually granted: canonical writable implies
// observed writable, canonical signer implies observed signer. A message
// may over-privilege an account beyond what the program requires (e.g. the
// V1 fixture's newLeafOwner is writable there only because the sending
// wallet happened to reuse a writable/signer account for multiple roles);
// exact equality with the observed flags is not asserted for that reason.
func requireCanonicalSubsetOfObserved(t *testing.T, canonical []*solanago.AccountMeta, observed []fixtureAccount) {
	t.Helper()
	require.Len(t, canonical, len(observed))
	for i, c := range canonical {
		o := observed[i]
		require.Equal(t, o.Pubkey, c.PublicKey.String(), "account %d pubkey", i)
		if c.IsWritable {
			require.True(t, o.Writable, "account %d expected to be writable in the observed transaction", i)
		}
		if c.IsSigner {
			require.True(t, o.Signer, "account %d expected to be a signer in the observed transaction", i)
		}
	}
}

// TestNewTransferInstruction_MatchesMainnetTransaction pins the V1 builder
// against a decoded mainnet "transfer" instruction:
// 3CdHRnDjSdpNmjFgr1KTk7DbTqpeQKPNFZFyZxd7LZXaRyGH2ei94fnRsYBsy3YZLWS1hYzuq8yB9RkgScASBinC
// (slot 434111830).
func TestNewTransferInstruction_MatchesMainnetTransaction(t *testing.T) {
	f := loadFixture(t, "testdata/fixture_transfer.json")
	require.Len(t, f.Accounts, 16)

	args := TransferArgs{
		Root:        mustHash32(t, f.Args.Root),
		DataHash:    mustHash32(t, f.Args.DataHash),
		CreatorHash: mustHash32(t, f.Args.CreatorHash),
		Nonce:       f.Args.Nonce,
		Index:       f.Args.Index,
	}

	acct := func(i int) solanago.PublicKey {
		return solanago.MustPublicKeyFromBase58(f.Accounts[i].Pubkey)
	}
	proof := make([]solanago.PublicKey, 0, len(f.Accounts)-8)
	for i := 8; i < len(f.Accounts); i++ {
		proof = append(proof, acct(i))
	}

	inst, err := NewTransferInstruction(
		args,
		acct(0), // treeConfig
		acct(1), // leafOwner
		acct(2), // leafDelegate
		acct(3), // newLeafOwner
		acct(4), // merkleTree
		acct(5), // logWrapper
		acct(6), // compressionProgram
		acct(7), // systemProgram
		proof,
	)
	require.NoError(t, err)

	data, err := inst.Data()
	require.NoError(t, err)
	require.Equal(t, f.DataHex, hex.EncodeToString(data))

	require.Equal(t, ProgramID, inst.ProgramID())
	requireCanonicalSubsetOfObserved(t, inst.Accounts(), f.Accounts)

	// Only merkleTree is writable; only leafOwner (the transfer authority in
	// this package's owner-authorized model) is a signer. All 8 proof nodes
	// are read-only, non-signer.
	requireExactFlags(t, inst.Accounts(), []expectedFlags{
		{Writable: false, Signer: false}, // treeConfig
		{Writable: false, Signer: true},  // leafOwner
		{Writable: false, Signer: false}, // leafDelegate
		{Writable: false, Signer: false}, // newLeafOwner
		{Writable: true, Signer: false},  // merkleTree
		{Writable: false, Signer: false}, // logWrapper
		{Writable: false, Signer: false}, // compressionProgram
		{Writable: false, Signer: false}, // systemProgram
		{Writable: false, Signer: false}, // proof[0]
		{Writable: false, Signer: false}, // proof[1]
		{Writable: false, Signer: false}, // proof[2]
		{Writable: false, Signer: false}, // proof[3]
		{Writable: false, Signer: false}, // proof[4]
		{Writable: false, Signer: false}, // proof[5]
		{Writable: false, Signer: false}, // proof[6]
		{Writable: false, Signer: false}, // proof[7]
	})
}

// TestNewTransferV2Instruction_MatchesMainnetTransaction pins the V2 builder
// against a decoded mainnet "transfer_v2" instruction:
// 1PNeVUqGhzG4KtN294uRnzY2AfU3ynFYPcAM6FFoXXiTAkqw1RqNf9qqyNQ2yDRm9i2TBEkc6uXemED9sFZKFzF
// (slot 434112017).
func TestNewTransferV2Instruction_MatchesMainnetTransaction(t *testing.T) {
	f := loadFixture(t, "testdata/fixture_transfer_v2.json")
	require.Len(t, f.Accounts, 18)
	require.NotNil(t, f.Args.AssetDataHash)
	require.Nil(t, f.Args.Flags)

	assetDataHash := mustHash32(t, *f.Args.AssetDataHash)
	args := TransferV2Args{
		Root:          mustHash32(t, f.Args.Root),
		DataHash:      mustHash32(t, f.Args.DataHash),
		CreatorHash:   mustHash32(t, f.Args.CreatorHash),
		AssetDataHash: &assetDataHash,
		Flags:         nil,
		Nonce:         f.Args.Nonce,
		Index:         f.Args.Index,
	}

	acct := func(i int) solanago.PublicKey {
		return solanago.MustPublicKeyFromBase58(f.Accounts[i].Pubkey)
	}
	proof := make([]solanago.PublicKey, 0, len(f.Accounts)-11)
	for i := 11; i < len(f.Accounts); i++ {
		proof = append(proof, acct(i))
	}

	// The fixture's authority account (pos 2) is the Bubblegum program ID
	// itself: the on-chain transaction omitted authority, so leafOwner
	// signed instead.
	require.Equal(t, ProgramID.String(), f.Accounts[2].Pubkey)

	inst, err := NewTransferV2Instruction(
		args,
		acct(0),  // treeConfig
		acct(1),  // payer
		acct(2),  // authority (omitted: equals ProgramID)
		acct(3),  // leafOwner
		acct(4),  // leafDelegate
		acct(5),  // newLeafOwner
		acct(6),  // merkleTree
		acct(7),  // coreCollection
		acct(8),  // logWrapper
		acct(9),  // compressionProgram
		acct(10), // systemProgram
		proof,
	)
	require.NoError(t, err)

	data, err := inst.Data()
	require.NoError(t, err)
	require.Equal(t, f.DataHex, hex.EncodeToString(data))

	require.Equal(t, ProgramID, inst.ProgramID())
	requireCanonicalSubsetOfObserved(t, inst.Accounts(), f.Accounts)

	// treeConfig/payer/merkleTree are writable. Only payer and leafOwner are
	// signers in this fixture: authority is omitted (equals ProgramID), so
	// leafOwner signs in its place instead of authority. All 7 proof nodes
	// are read-only, non-signer.
	requireExactFlags(t, inst.Accounts(), []expectedFlags{
		{Writable: true, Signer: false},  // treeConfig
		{Writable: true, Signer: true},   // payer
		{Writable: false, Signer: false}, // authority (omitted)
		{Writable: false, Signer: true},  // leafOwner (signs in authority's place)
		{Writable: false, Signer: false}, // leafDelegate
		{Writable: false, Signer: false}, // newLeafOwner
		{Writable: true, Signer: false},  // merkleTree
		{Writable: false, Signer: false}, // coreCollection
		{Writable: false, Signer: false}, // logWrapper
		{Writable: false, Signer: false}, // compressionProgram
		{Writable: false, Signer: false}, // systemProgram
		{Writable: false, Signer: false}, // proof[0]
		{Writable: false, Signer: false}, // proof[1]
		{Writable: false, Signer: false}, // proof[2]
		{Writable: false, Signer: false}, // proof[3]
		{Writable: false, Signer: false}, // proof[4]
		{Writable: false, Signer: false}, // proof[5]
		{Writable: false, Signer: false}, // proof[6]
	})
}

// TestFindTreeConfigAddress pins tree config PDA derivation against the two
// tree accounts backing the fixtures above.
func TestFindTreeConfigAddress(t *testing.T) {
	tests := []struct {
		tree     string
		expected string
	}{
		{"GYGG71sQjYSZjxiXXfDjTLwHj4cNi4vUJvLsreANYAep", "3XnzMMXQwE3BqSK1RStLbxUCtZU1J6dWZ5axnVi3kcmu"},
		{"Dw4aZw7qkgbwHyEPA55p9uN818S6inD3L1h8E4Z8s45H", "7mWBfPKHipiXha3qsY8u9jJ2n9WQCmkXa861Ncejg2AU"},
	}
	for _, tt := range tests {
		t.Run(tt.tree, func(t *testing.T) {
			addr, err := FindTreeConfigAddress(solanago.MustPublicKeyFromBase58(tt.tree))
			require.NoError(t, err)
			require.Equal(t, tt.expected, addr.String())
		})
	}
}
