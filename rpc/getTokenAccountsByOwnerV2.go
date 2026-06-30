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
	"encoding/base64"
	stdjson "encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// getTokenAccountsByOwnerV2 is the Helius extension of getTokenAccountsByOwner.
// It takes the same [ownerAddress, filter, options] parameter layout and
// returns the token accounts wrapped in a context object together with a
// cursor for pagination.
//
// The types below model the jsonParsed response shape: every account carries
// its fully parsed SPL Token / Token-2022 state, including account extensions.
// Token-2022 accounts whose data cannot be expressed as jsonParsed are returned
// as a raw base64 array and exposed through RawData for caller-side decoding.

// TokenAccountV2Filter selects which token accounts to return. Exactly one of
// ProgramID or Mint is set.
type TokenAccountV2Filter struct {
	ProgramID string `json:"programId,omitempty"`
	Mint      string `json:"mint,omitempty"`
}

// TokenAccountV2Options are the query options for getTokenAccountsByOwnerV2.
type TokenAccountV2Options struct {
	Encoding      string `json:"encoding,omitempty"`      // "jsonParsed"
	Limit         int    `json:"limit,omitempty"`         // 1-10000
	PaginationKey string `json:"paginationKey,omitempty"` // cursor from previous response
	Commitment    string `json:"commitment,omitempty"`    // "processed", "confirmed", or "finalized"

	// ChangedSinceSlot is a Helius extension: when set, only accounts modified
	// at or after this slot are returned. Used for incremental re-sync after a
	// stream reconnect — the response's context.slot becomes the next high-water
	// mark. Owner-scoped, so a closed account simply does not appear (there is no
	// "closed" status for owner enumeration). 0 omits the field (full enumeration).
	ChangedSinceSlot uint64 `json:"changedSinceSlot,omitempty"`

	// MinContextSlot requires the node to evaluate the query at or after this
	// slot, failing rather than serving a staler snapshot. 0 omits the field.
	MinContextSlot uint64 `json:"minContextSlot,omitempty"`
}

// TokenAccountV2Response is the response from getTokenAccountsByOwnerV2.
type TokenAccountV2Response struct {
	Context TokenAccountV2Context   `json:"context"`
	Value   TokenAccountV2ValueWrap `json:"value"`
}

// TokenAccountV2Context contains the RPC context info.
type TokenAccountV2Context struct {
	APIVersion string `json:"apiVersion"`
	Slot       uint64 `json:"slot"`
}

// TokenAccountV2ValueWrap wraps the accounts array with pagination info.
type TokenAccountV2ValueWrap struct {
	Accounts      []TokenAccountV2Value `json:"accounts"`
	PaginationKey *string               `json:"paginationKey"` // null when no more pages
	Count         int                   `json:"count"`
}

// TokenAccountV2Value represents a single token account from the V2 API.
type TokenAccountV2Value struct {
	Pubkey  string             `json:"pubkey"`  // token account address
	Account TokenAccountV2Data `json:"account"` // account data
}

// TokenAccountV2Data contains the account envelope.
type TokenAccountV2Data struct {
	Data       TokenAccountV2DataParsed `json:"data"`
	Executable bool                     `json:"executable"`
	Lamports   uint64                   `json:"lamports"`
	Owner      string                   `json:"owner"` // token program ID
	RentEpoch  uint64                   `json:"rentEpoch"`
	Space      uint64                   `json:"space"`
}

// TokenAccountV2DataParsed contains the parsed field with program info.
// It implements custom UnmarshalJSON because the RPC may return a raw array
// ["<base64>", "base64"] instead of a parsed object for Token-2022 accounts.
type TokenAccountV2DataParsed struct {
	Parsed  TokenAccountV2Parsed `json:"parsed"`
	Program string               `json:"program"` // "spl-token" or "spl-token-2022"
	Space   uint64               `json:"space"`
	RawData []byte               `json:"-"` // populated when data comes as base64 array
}

func (d *TokenAccountV2DataParsed) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '[' {
		// Raw encoding: ["<base64_data>", "base64"]
		var arr []string
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		if len(arr) >= 1 {
			decoded, err := base64.StdEncoding.DecodeString(arr[0])
			if err != nil {
				return fmt.Errorf("decode base64 account data: %w", err)
			}
			d.RawData = decoded
		}
		return nil
	}
	type alias TokenAccountV2DataParsed
	return json.Unmarshal(data, (*alias)(d))
}

// TokenAccountV2Parsed contains the info field.
type TokenAccountV2Parsed struct {
	Info TokenAccountV2Info `json:"info"`
	Type string             `json:"type"` // "account"
}

// TokenAccountV2Info contains mint, owner, state, delegate and extension info.
type TokenAccountV2Info struct {
	IsNative        bool                      `json:"isNative"`
	Mint            string                    `json:"mint"`
	Owner           string                    `json:"owner"`
	State           string                    `json:"state"` // "initialized" or "frozen"
	TokenAmount     TokenAccountV2Amount      `json:"tokenAmount"`
	Delegate        string                    `json:"delegate,omitempty"`
	DelegatedAmount *TokenAccountV2Amount     `json:"delegatedAmount,omitempty"`
	Extensions      []TokenAccountV2Extension `json:"extensions,omitempty"` // Token-2022 account extensions
}

// TokenAccountV2Extension represents a single Token-2022 extension from a
// jsonParsed response.
type TokenAccountV2Extension struct {
	Extension string             `json:"extension"` // e.g. "transferFeeAmount", "immutableOwner"
	State     stdjson.RawMessage `json:"state"`     // extension-specific data, shape varies
}

// TokenAccountV2Amount contains balance information.
type TokenAccountV2Amount struct {
	Amount         string  `json:"amount"` // raw balance as string
	Decimals       int     `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}

// GetTokenAccountsByOwnerV2WithOpts returns the token accounts owned by owner
// for a single program (or mint), with cursor-based pagination. The accounts
// are returned fully parsed; pass options.PaginationKey to fetch the next page.
func (cl *Client) GetTokenAccountsByOwnerV2WithOpts(
	ctx context.Context,
	owner solana.PublicKey,
	filter *TokenAccountV2Filter,
	options *TokenAccountV2Options,
) (out *TokenAccountV2Response, err error) {
	params := []interface{}{owner.String(), filter, options}
	err = cl.rpcClient.CallForInto(ctx, &out, "getTokenAccountsByOwnerV2", params)
	return
}
