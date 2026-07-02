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

// buildTransactionSubscribeParams assembles the JSON-RPC params array shared by
// every transactionSubscribe variant: the filter map followed by an optional
// options object. When forceDetail is non-empty it overrides
// opts.TransactionDetails, letting the detail-specific subscribe methods pin the
// notification shape they decode.
func buildTransactionSubscribeParams(
	filter TransactionSubscribeFilter,
	opts *TransactionSubscribeOpts,
	forceDetail rpc.TransactionDetailsType,
) ([]interface{}, error) {
	params := []interface{}{filter.toMap()}

	obj := make(rpc.M)
	if opts != nil {
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
		if opts.ShowRewards != nil {
			obj["showRewards"] = opts.ShowRewards
		}
		if opts.MaxSupportedTransactionVersion != nil {
			obj["maxSupportedTransactionVersion"] = *opts.MaxSupportedTransactionVersion
		}
	}

	detail := forceDetail
	if detail == "" && opts != nil {
		detail = opts.TransactionDetails
	} else if forceDetail != "" && opts != nil && opts.TransactionDetails != "" && opts.TransactionDetails != forceDetail {
		return nil, fmt.Errorf("transactionDetails %q conflicts with the pinned %q for this subscription method", opts.TransactionDetails, forceDetail)
	}
	if detail != "" {
		obj["transactionDetails"] = detail
	}

	if len(obj) > 0 {
		params = append(params, obj)
	}
	return params, nil
}

// TransactionSubscribe subscribes to transaction updates matching the provided
// filter, decoding the full-detail notification shape.
//
// Each transactionDetails level below "full" changes the notification shape and
// has its own typed entry point: TransactionSubscribeAccounts,
// TransactionSubscribeSignatures, TransactionSubscribeNone. Requesting any of
// those levels here is rejected rather than silently decoded into an empty
// TransactionResult.
//
// NOTE: this is a Helius extension to the standard Solana websocket API; the
// endpoint must support the transactionSubscribe method.
func (cl *Client) TransactionSubscribe(
	filter TransactionSubscribeFilter,
	opts *TransactionSubscribeOpts,
) (*Subscription[TransactionResult], error) {
	if opts != nil {
		switch opts.TransactionDetails {
		case rpc.TransactionDetailsAccounts:
			return nil, fmt.Errorf("transactionDetails %q requires TransactionSubscribeAccounts", opts.TransactionDetails)
		case rpc.TransactionDetailsSignatures:
			return nil, fmt.Errorf("transactionDetails %q requires TransactionSubscribeSignatures", opts.TransactionDetails)
		case rpc.TransactionDetailsNone:
			return nil, fmt.Errorf("transactionDetails %q requires TransactionSubscribeNone", opts.TransactionDetails)
		}
	}

	params, err := buildTransactionSubscribeParams(filter, opts, "")
	if err != nil {
		return nil, err
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

// ParsedAccountKey is one entry of the flattened account list returned at
// transactionDetails "accounts". The list already includes addresses loaded
// from address lookup tables (Source "lookupTable"), so an accountIndex from a
// token balance entry indexes directly into it.
type ParsedAccountKey struct {
	Pubkey   solana.PublicKey `json:"pubkey"`
	Writable bool             `json:"writable"`
	Signer   bool             `json:"signer"`
	// Source is "transaction" for keys carried by the message and "lookupTable"
	// for keys resolved from an address lookup table.
	Source string `json:"source"`
}

// TransactionAccountsResult is a transactionSubscribe notification decoded at
// transactionDetails "accounts": the flattened account keys plus transaction
// metadata (pre/post balances and token balances), without the full instruction
// set.
type TransactionAccountsResult struct {
	Signature        solana.Signature `json:"signature"`
	Slot             uint64           `json:"slot"`
	TransactionIndex uint64           `json:"transactionIndex"`
	Transaction      struct {
		Transaction struct {
			Signatures  []solana.Signature `json:"signatures"`
			AccountKeys []ParsedAccountKey `json:"accountKeys"`
		} `json:"transaction"`
		Meta    *rpc.TransactionMeta   `json:"meta,omitempty"`
		Version rpc.TransactionVersion `json:"version"`
	} `json:"transaction"`
}

// TransactionSubscribeAccounts subscribes to transaction updates matching the
// filter, pinning transactionDetails to "accounts". The notification carries the
// flattened account keys and transaction metadata but not the full instruction
// set.
//
// opts.MaxSupportedTransactionVersion is mandatory: without it the endpoint
// cannot deliver versioned (v0) transactions, and account keys loaded from
// address lookup tables would be missing from the flattened list. It is
// required here rather than defaulted so the omission fails loudly.
//
// NOTE: this is a Helius extension to the standard Solana websocket API; the
// endpoint must support the transactionSubscribe method.
func (cl *Client) TransactionSubscribeAccounts(
	filter TransactionSubscribeFilter,
	opts *TransactionSubscribeOpts,
) (*Subscription[TransactionAccountsResult], error) {
	if opts == nil || opts.MaxSupportedTransactionVersion == nil {
		return nil, fmt.Errorf("TransactionSubscribeAccounts requires opts.MaxSupportedTransactionVersion")
	}

	params, err := buildTransactionSubscribeParams(filter, opts, rpc.TransactionDetailsAccounts)
	if err != nil {
		return nil, err
	}

	genSub, err := cl.subscribe(
		params,
		nil,
		"transactionSubscribe",
		"transactionUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res TransactionAccountsResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &Subscription[TransactionAccountsResult]{
		sub:       genSub,
		closeFunc: func() { genSub.closeFunc(nil) },
	}, nil
}

// TransactionSignaturesResult is a transactionSubscribe notification decoded at
// transactionDetails "signatures": only the transaction's primary signature and
// its position, without any transaction body or metadata.
type TransactionSignaturesResult struct {
	Signature        solana.Signature `json:"signature"`
	Slot             uint64           `json:"slot"`
	TransactionIndex uint64           `json:"transactionIndex"`
}

// TransactionSubscribeSignatures subscribes to transaction updates matching the
// filter, pinning transactionDetails to "signatures". The notification carries
// only the transaction signature and its slot/index.
//
// NOTE: this is a Helius extension to the standard Solana websocket API; the
// endpoint must support the transactionSubscribe method.
func (cl *Client) TransactionSubscribeSignatures(
	filter TransactionSubscribeFilter,
	opts *TransactionSubscribeOpts,
) (*Subscription[TransactionSignaturesResult], error) {
	params, err := buildTransactionSubscribeParams(filter, opts, rpc.TransactionDetailsSignatures)
	if err != nil {
		return nil, err
	}

	genSub, err := cl.subscribe(
		params,
		nil,
		"transactionSubscribe",
		"transactionUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res TransactionSignaturesResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &Subscription[TransactionSignaturesResult]{
		sub:       genSub,
		closeFunc: func() { genSub.closeFunc(nil) },
	}, nil
}

// TransactionNoneResult is a transactionSubscribe notification decoded at
// transactionDetails "none": only the slot and the transaction's position
// within the block. No signature or transaction body is delivered at this level.
type TransactionNoneResult struct {
	Slot             uint64 `json:"slot"`
	TransactionIndex uint64 `json:"transactionIndex"`
}

// TransactionSubscribeNone subscribes to transaction updates matching the
// filter, pinning transactionDetails to "none". The notification carries only
// the slot and the transaction's index within the block.
//
// NOTE: this is a Helius extension to the standard Solana websocket API; the
// endpoint must support the transactionSubscribe method.
func (cl *Client) TransactionSubscribeNone(
	filter TransactionSubscribeFilter,
	opts *TransactionSubscribeOpts,
) (*Subscription[TransactionNoneResult], error) {
	params, err := buildTransactionSubscribeParams(filter, opts, rpc.TransactionDetailsNone)
	if err != nil {
		return nil, err
	}

	genSub, err := cl.subscribe(
		params,
		nil,
		"transactionSubscribe",
		"transactionUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res TransactionNoneResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &Subscription[TransactionNoneResult]{
		sub:       genSub,
		closeFunc: func() { genSub.closeFunc(nil) },
	}, nil
}
