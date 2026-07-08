// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

// getProgramAccountsV2 is the Helius extension of getProgramAccounts. It
// accepts the same configuration (encoding, dataSlice, filters, commitment)
// plus cursor-based pagination and incremental-sync options, and returns the
// accounts one page at a time instead of the full unbounded set.

// GetProgramAccountsV2Opts are the query options for getProgramAccountsV2.
type GetProgramAccountsV2Opts struct {
	Commitment CommitmentType `json:"commitment,omitempty"`

	Encoding solana.EncodingType `json:"encoding,omitempty"`

	// Limit the returned account data to a byte range.
	DataSlice *DataSlice `json:"dataSlice,omitempty"`

	// Filter on accounts, implicit AND between filters.
	Filters []RPCFilter `json:"filters,omitempty"`

	// Limit is the maximum number of accounts per page (1-10000).
	// 0 omits the field (server default applies).
	Limit int `json:"limit,omitempty"`

	// PaginationKey is the cursor returned by the previous page. Empty for
	// the first page.
	PaginationKey string `json:"paginationKey,omitempty"`

	// ChangedSinceSlot is a Helius extension: when set, only accounts modified
	// at or after this slot are returned. Used for incremental re-sync — the
	// response's context.slot becomes the next high-water mark. Program-scoped,
	// so an account that was closed or reassigned simply does not appear.
	// 0 omits the field (full enumeration).
	ChangedSinceSlot uint64 `json:"changedSinceSlot,omitempty"`

	// MinContextSlot requires the node to evaluate the query at or after this
	// slot, failing rather than serving a staler snapshot. 0 omits the field.
	MinContextSlot uint64 `json:"minContextSlot,omitempty"`
}

// GetProgramAccountsV2Result is the response from getProgramAccountsV2. The
// request is always sent with withContext=true, so the page is wrapped in an
// RPC context; Context.Slot is the snapshot slot of the page and serves as the
// next ChangedSinceSlot high-water mark for incremental sync.
type GetProgramAccountsV2Result struct {
	RPCContext
	Value GetProgramAccountsV2Page `json:"value"`
}

// GetProgramAccountsV2Page is a single page of program accounts.
//
// PaginationKey is null only when pagination is complete. A page may contain
// zero accounts while PaginationKey is still non-null (filters are applied
// per page), so callers must keep paginating until PaginationKey is null
// rather than stopping at the first empty page.
type GetProgramAccountsV2Page struct {
	Accounts      []*KeyedAccount `json:"accounts"`
	PaginationKey *string         `json:"paginationKey"`
}

// GetProgramAccountsV2 returns one page of accounts owned by the provided
// program publicKey, with cursor-based pagination. Pass the returned
// PaginationKey via opts to fetch the next page; a null PaginationKey marks
// the end of the set.
func (cl *Client) GetProgramAccountsV2(
	ctx context.Context,
	publicKey solana.PublicKey,
	opts *GetProgramAccountsV2Opts,
) (out *GetProgramAccountsV2Result, err error) {
	obj := M{
		"encoding": "base64",
		// Always request the context wrapper: it fixes the response to a
		// single shape and exposes the snapshot slot needed for
		// ChangedSinceSlot-based incremental sync.
		"withContext": true,
	}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if opts.Encoding != "" {
			obj["encoding"] = opts.Encoding
		}
		if opts.DataSlice != nil {
			obj["dataSlice"] = M{
				"offset": opts.DataSlice.Offset,
				"length": opts.DataSlice.Length,
			}
		}
		if len(opts.Filters) != 0 {
			obj["filters"] = opts.Filters
		}
		if opts.Limit != 0 {
			obj["limit"] = opts.Limit
		}
		if opts.PaginationKey != "" {
			obj["paginationKey"] = opts.PaginationKey
		}
		if opts.ChangedSinceSlot != 0 {
			obj["changedSinceSlot"] = opts.ChangedSinceSlot
		}
		if opts.MinContextSlot != 0 {
			obj["minContextSlot"] = opts.MinContextSlot
		}
	}

	params := []interface{}{publicKey, obj}

	err = cl.rpcClient.CallForInto(ctx, &out, "getProgramAccountsV2", params)
	return
}
