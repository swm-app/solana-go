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

package ws

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// TransactionResult is a transactionSubscribe notification payload.
//
// NOTE: transactionSubscribe is a Helius extension to the standard Solana
// websocket API and is only available on endpoints that implement it.
type TransactionResult struct {
	Transaction struct {
		// Transaction holds the encoded transaction. For the base58 and base64
		// encodings it carries the decoded binary; for jsonParsed it carries a
		// parsed *solana.Transaction.
		//
		// NOTE: the jsonParsed encoding returns instructions in a parsed shape
		// that does not map cleanly onto solana.Transaction. Prefer base58 or
		// base64 when you need the full, faithful transaction.
		Transaction *rpc.TransactionResultEnvelope `json:"transaction"`
		Meta        *rpc.TransactionMeta           `json:"meta,omitempty"`
	} `json:"transaction"`
	// The transaction signature.
	Signature solana.Signature `json:"signature"`
	// The slot this transaction was processed in.
	Slot uint64 `json:"slot"`
	// The index of the transaction within the block.
	TransactionIndex uint64 `json:"transactionIndex"`
}

// TransactionSubscribeTokenAccounts controls whether the notification is
// expanded with token account data.
type TransactionSubscribeTokenAccounts string

const (
	// Include only token accounts whose balance changed.
	TransactionSubscribeTokenAccountsBalanceChanged TransactionSubscribeTokenAccounts = "balanceChanged"
	// Include all token accounts.
	TransactionSubscribeTokenAccountsAll TransactionSubscribeTokenAccounts = "all"
	// Do not include token accounts.
	TransactionSubscribeTokenAccountsNone TransactionSubscribeTokenAccounts = "none"
)

// TransactionSubscribeFilter selects which transactions to receive updates for.
// All fields are optional; leaving a field unset omits it from the request.
type TransactionSubscribeFilter struct {
	// Include (true) or exclude (false) vote-related transactions.
	Vote *bool
	// Include (true) or exclude (false) transactions that failed.
	Failed *bool
	// Filter updates to a single transaction by its signature.
	Signature string
	// Receive updates for transactions mentioning any of these accounts (OR).
	// Up to 50,000 addresses.
	AccountInclude []solana.PublicKey
	// Exclude transactions mentioning any of these accounts.
	// Up to 50,000 addresses.
	AccountExclude []solana.PublicKey
	// Only receive transactions that mention all of these accounts (AND).
	// Up to 50,000 addresses.
	AccountRequired []solana.PublicKey
	// Expansion mode for token account data in the notification.
	TokenAccounts TransactionSubscribeTokenAccounts
}

func (f TransactionSubscribeFilter) toMap() rpc.M {
	obj := make(rpc.M)
	if f.Vote != nil {
		obj["vote"] = *f.Vote
	}
	if f.Failed != nil {
		obj["failed"] = *f.Failed
	}
	if f.Signature != "" {
		obj["signature"] = f.Signature
	}
	if len(f.AccountInclude) > 0 {
		obj["accountInclude"] = publicKeysToStrings(f.AccountInclude)
	}
	if len(f.AccountExclude) > 0 {
		obj["accountExclude"] = publicKeysToStrings(f.AccountExclude)
	}
	if len(f.AccountRequired) > 0 {
		obj["accountRequired"] = publicKeysToStrings(f.AccountRequired)
	}
	if f.TokenAccounts != "" {
		obj["tokenAccounts"] = f.TokenAccounts
	}
	return obj
}

func publicKeysToStrings(keys []solana.PublicKey) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = k.String()
	}
	return out
}

// TransactionSubscribeOpts holds the optional configuration of the subscription.
type TransactionSubscribeOpts struct {
	Commitment rpc.CommitmentType

	// Encoding for the returned transaction.
	// Supported: solana.EncodingBase58, solana.EncodingBase64, solana.EncodingJSONParsed.
	Encoding solana.EncodingType

	// Level of transaction detail to return.
	TransactionDetails rpc.TransactionDetailsType

	// Whether to include reward data in the updates.
	ShowRewards *bool

	// Max transaction version to return in responses.
	// Required when TransactionDetails is "full" or "accounts".
	MaxSupportedTransactionVersion *uint64
}

// TransactionSubscribe subscribes to transaction updates matching the provided
// filter.
//
// NOTE: this is a Helius extension to the standard Solana websocket API; the
// endpoint must support the transactionSubscribe method.
func (cl *Client) TransactionSubscribe(
	filter TransactionSubscribeFilter,
	opts *TransactionSubscribeOpts,
) (*Subscription[TransactionResult], error) {
	params := []interface{}{filter.toMap()}

	if opts != nil {
		obj := make(rpc.M)
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Encoding != "" {
			if !solana.IsAnyOfEncodingType(
				opts.Encoding,
				solana.EncodingBase58,
				solana.EncodingBase64,
				solana.EncodingJSONParsed,
			) {
				return nil, fmt.Errorf("provided encoding is not supported: %s", opts.Encoding)
			}
			obj["encoding"] = opts.Encoding
		}
		if opts.TransactionDetails != "" {
			obj["transactionDetails"] = opts.TransactionDetails
		}
		if opts.ShowRewards != nil {
			obj["showRewards"] = opts.ShowRewards
		}
		if opts.MaxSupportedTransactionVersion != nil {
			obj["maxSupportedTransactionVersion"] = *opts.MaxSupportedTransactionVersion
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}

	genSub, err := cl.subscribe(
		params,
		nil,
		"transactionSubscribe",
		"transactionUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res TransactionResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &Subscription[TransactionResult]{
		sub:       genSub,
		closeFunc: func() { genSub.closeFunc(nil) },
	}, nil
}
