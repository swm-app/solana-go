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
	"testing"

	stdjson "encoding/json"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetTokenAccountsByOwnerV2WithOpts_Parsed(t *testing.T) {
	responseBody := `{"context":{"apiVersion":"2.0.15","slot":1114},"value":{` +
		`"accounts":[{"pubkey":"CnPoSPKXu7wJqxe59Fs72tkBeALovhsCxYeFwPCQH9TD","account":{"data":{` +
		`"program":"spl-token","space":165,"parsed":{"type":"account","info":{` +
		`"isNative":false,"mint":"3wyAj7Rt1TWVPZVteFJPLa26JmLvdb1CAKEFZm3NY75E",` +
		`"owner":"4Qkev8aNZcqFNSRhQzwyLMFSsi94jHqE8WNVTJzTP99F","state":"initialized",` +
		`"tokenAmount":{"amount":"600","decimals":6,"uiAmount":0.0006,"uiAmountString":"0.0006"}}}},` +
		`"executable":false,"lamports":1726080,"owner":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",` +
		`"rentEpoch":4,"space":165}}],"paginationKey":"NextPageCursor111","count":1}}`
	server, closer := mockJSONRPC(t, stdjson.RawMessage(wrapIntoRPC(responseBody)))
	defer closer()
	client := New(server.URL)

	ownerString := "7xLk17EQQ5KLDLDe44wCmupJKJjTGd8hs3eSVVhCx932"
	owner := solana.MustPublicKeyFromBase58(ownerString)
	programIDString := "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"

	out, err := client.GetTokenAccountsByOwnerV2WithOpts(
		context.Background(),
		owner,
		&TokenAccountV2Filter{ProgramID: programIDString},
		&TokenAccountV2Options{
			Encoding:      "jsonParsed",
			Limit:         100,
			Commitment:    "confirmed",
			PaginationKey: "PrevPageCursor000",
		},
	)
	require.NoError(t, err)

	reqBody := server.RequestBody(t)
	assert.NotNil(t, reqBody["id"])
	reqBody["id"] = any(nil)

	assert.Equal(t,
		map[string]interface{}{
			"id":      any(nil),
			"jsonrpc": "2.0",
			"method":  "getTokenAccountsByOwnerV2",
			"params": []interface{}{
				ownerString,
				map[string]interface{}{
					"programId": programIDString,
				},
				map[string]interface{}{
					"encoding":      "jsonParsed",
					"limit":         float64(100),
					"commitment":    "confirmed",
					"paginationKey": "PrevPageCursor000",
				},
			},
		},
		reqBody,
	)

	require.Len(t, out.Value.Accounts, 1)
	acc := out.Value.Accounts[0]
	assert.Equal(t, "CnPoSPKXu7wJqxe59Fs72tkBeALovhsCxYeFwPCQH9TD", acc.Pubkey)
	assert.Equal(t, "spl-token", acc.Account.Data.Program)
	assert.Equal(t, "3wyAj7Rt1TWVPZVteFJPLa26JmLvdb1CAKEFZm3NY75E", acc.Account.Data.Parsed.Info.Mint)
	assert.Equal(t, "initialized", acc.Account.Data.Parsed.Info.State)
	assert.Equal(t, "600", acc.Account.Data.Parsed.Info.TokenAmount.Amount)
	assert.Equal(t, 6, acc.Account.Data.Parsed.Info.TokenAmount.Decimals)

	require.NotNil(t, out.Value.PaginationKey)
	assert.Equal(t, "NextPageCursor111", *out.Value.PaginationKey)
	assert.Equal(t, uint64(1114), out.Context.Slot)
}

func TestClient_GetTokenAccountsByOwnerV2WithOpts_RawToken2022(t *testing.T) {
	// Token-2022 accounts that cannot be expressed as jsonParsed come back as a
	// ["<base64>", "base64"] array; the decoded bytes must land in RawData.
	responseBody := `{"context":{"slot":10},"value":{"accounts":[{` +
		`"pubkey":"CnPoSPKXu7wJqxe59Fs72tkBeALovhsCxYeFwPCQH9TD","account":{` +
		`"data":["AQID","base64"],"executable":false,"lamports":2039280,` +
		`"owner":"TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb","rentEpoch":0,"space":182}}],` +
		`"paginationKey":null,"count":1}}`
	server, closer := mockJSONRPC(t, stdjson.RawMessage(wrapIntoRPC(responseBody)))
	defer closer()
	client := New(server.URL)

	out, err := client.GetTokenAccountsByOwnerV2WithOpts(
		context.Background(),
		solana.MustPublicKeyFromBase58("7xLk17EQQ5KLDLDe44wCmupJKJjTGd8hs3eSVVhCx932"),
		&TokenAccountV2Filter{ProgramID: "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"},
		&TokenAccountV2Options{Encoding: "jsonParsed", Limit: 100},
	)
	require.NoError(t, err)

	require.Len(t, out.Value.Accounts, 1)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, out.Value.Accounts[0].Account.Data.RawData)
	assert.Nil(t, out.Value.PaginationKey)
}
