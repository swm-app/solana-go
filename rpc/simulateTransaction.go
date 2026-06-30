// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package rpc

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type SimulateTransactionResponse struct {
	RPCContext
	Value *SimulateTransactionResult `json:"value"`
}

type SimulateTransactionResult struct {
	// Error if transaction failed, null if transaction succeeded.
	Err interface{} `json:"err,omitempty"`

	// Array of log messages the transaction instructions output during execution,
	// null if simulation failed before the transaction was able to execute
	// (for example due to an invalid blockhash or signature verification failure)
	Logs []string `json:"logs,omitempty"`

	// Array of accounts with the same length as the accounts.addresses array in the request.
	Accounts []*Account `json:"accounts"`

	// The number of compute budget units consumed during the processing of this transaction.
	UnitsConsumed *uint64 `json:"unitsConsumed,omitempty"`

	// Fee this transaction would be charged (base fee plus any priority fee
	// implied by ComputeBudget instructions in the simulated transaction).
	Fee *uint64 `json:"fee,omitempty"`

	// Lamport balances of the transaction's accounts before and after the
	// simulated transaction, indexed by the transaction's account list.
	PreBalances  []uint64 `json:"preBalances,omitempty"`
	PostBalances []uint64 `json:"postBalances,omitempty"`

	// SPL / Token-2022 balances before and after the simulated transaction,
	// for the token accounts the transaction touched.
	PreTokenBalances  []TokenBalance `json:"preTokenBalances,omitempty"`
	PostTokenBalances []TokenBalance `json:"postTokenBalances,omitempty"`

	// Inner (CPI) instructions, populated only when InnerInstructions was requested.
	InnerInstructions []InnerInstruction `json:"innerInstructions,omitempty"`

	// Return data from the most recent instruction that produced any.
	ReturnData *ReturnData `json:"returnData,omitempty"`

	// Blockhash the node substituted for simulation when ReplaceRecentBlockhash
	// was requested.
	ReplacementBlockhash *ReplacementBlockhash `json:"replacementBlockhash,omitempty"`

	// Addresses resolved from Address Lookup Tables during simulation.
	LoadedAddresses *LoadedAddresses `json:"loadedAddresses,omitempty"`

	// Total size in bytes of all accounts loaded by the transaction.
	LoadedAccountsDataSize *uint64 `json:"loadedAccountsDataSize,omitempty"`
}

// ReplacementBlockhash is returned when SimulateTransactionOpts.ReplaceRecentBlockhash
// was set: the recent blockhash the node substituted before simulation.
type ReplacementBlockhash struct {
	Blockhash            solana.Hash `json:"blockhash"`
	LastValidBlockHeight uint64      `json:"lastValidBlockHeight"`
}

// SimulateTransaction simulates sending a transaction.
func (cl *Client) SimulateTransaction(
	ctx context.Context,
	transaction *solana.Transaction,
) (out *SimulateTransactionResponse, err error) {
	return cl.SimulateTransactionWithOpts(
		ctx,
		transaction,
		nil,
	)
}

type SimulateTransactionOpts struct {
	// If true the transaction signatures will be verified
	// (default: false, conflicts with ReplaceRecentBlockhash)
	SigVerify bool

	// Commitment level to simulate the transaction at.
	// (default: "finalized").
	Commitment CommitmentType

	// If true the transaction recent blockhash will be replaced with the most recent blockhash.
	// (default: false, conflicts with SigVerify)
	ReplaceRecentBlockhash bool

	Accounts *SimulateTransactionAccountsOpts

	// MinContextSlot is the minimum slot the request can be evaluated at.
	// Zero means unset (no minimum).
	MinContextSlot uint64

	// If true, inner (CPI) instructions are included in the result.
	// (default: false)
	InnerInstructions bool
}

type SimulateTransactionAccountsOpts struct {
	// (optional) Encoding for returned Account data,
	// either "base64" (default), "base64+zstd" or "jsonParsed".
	// - "jsonParsed" encoding attempts to use program-specific state parsers
	//   to return more human-readable and explicit account state data.
	//   If "jsonParsed" is requested but a parser cannot be found,
	//   the field falls back to binary encoding, detectable when
	//   the data field is type <string>.
	Encoding solana.EncodingType

	// An array of accounts to return.
	Addresses []solana.PublicKey
}

// SimulateTransaction simulates sending a transaction.
func (cl *Client) SimulateTransactionWithOpts(
	ctx context.Context,
	transaction *solana.Transaction,
	opts *SimulateTransactionOpts,
) (out *SimulateTransactionResponse, err error) {
	txData, err := transaction.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("send transaction: encode transaction: %w", err)
	}
	return cl.SimulateRawTransactionWithOpts(ctx, txData, opts)
}

func (cl *Client) SimulateRawTransactionWithOpts(
	ctx context.Context,
	txData []byte,
	opts *SimulateTransactionOpts,
) (out *SimulateTransactionResponse, err error) {
	obj := M{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.SigVerify {
			obj["sigVerify"] = opts.SigVerify
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.ReplaceRecentBlockhash {
			obj["replaceRecentBlockhash"] = opts.ReplaceRecentBlockhash
		}
		if opts.Accounts != nil {
			obj["accounts"] = M{
				"encoding":  opts.Accounts.Encoding,
				"addresses": opts.Accounts.Addresses,
			}
		}
		if opts.MinContextSlot > 0 {
			obj["minContextSlot"] = opts.MinContextSlot
		}
		if opts.InnerInstructions {
			obj["innerInstructions"] = opts.InnerInstructions
		}
	}

	b64Data := base64.StdEncoding.EncodeToString(txData)
	params := []interface{}{
		b64Data,
		obj,
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "simulateTransaction", params)
	return
}
