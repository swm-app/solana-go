package token2022

import (
	"context"
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func FetchAnyAccount(ctx context.Context, rpcCli *rpc.Client, address solanago.PublicKey) (any, error) {
	data, _, err := fetchToken2022AccountData(ctx, rpcCli, address)
	if err != nil {
		return nil, err
	}
	return ParseAnyAccount(data)
}

func FetchMintAccount(ctx context.Context, rpcCli *rpc.Client, address solanago.PublicKey) (*MintAccount, error) {
	data, _, err := fetchToken2022AccountData(ctx, rpcCli, address)
	if err != nil {
		return nil, err
	}
	return ParseAccount_MintAccount(data)
}

func FetchTokenAccount(ctx context.Context, rpcCli *rpc.Client, address solanago.PublicKey) (*TokenAccount, error) {
	data, _, err := fetchToken2022AccountData(ctx, rpcCli, address)
	if err != nil {
		return nil, err
	}
	return ParseAccount_TokenAccount(data)
}

func FetchMultisigAccount(ctx context.Context, rpcCli *rpc.Client, address solanago.PublicKey) (*MultisigAccount, error) {
	data, _, err := fetchToken2022AccountData(ctx, rpcCli, address)
	if err != nil {
		return nil, err
	}
	return ParseAccount_MultisigAccount(data)
}

func FetchMintStateWithExtensions(ctx context.Context, rpcCli *rpc.Client, address solanago.PublicKey) (*MintStateWithExtensions, error) {
	data, _, err := fetchToken2022AccountData(ctx, rpcCli, address)
	if err != nil {
		return nil, err
	}
	return ParseMintStateWithExtensions(data)
}

func FetchTokenStateWithExtensions(ctx context.Context, rpcCli *rpc.Client, address solanago.PublicKey) (*TokenStateWithExtensions, error) {
	data, _, err := fetchToken2022AccountData(ctx, rpcCli, address)
	if err != nil {
		return nil, err
	}
	return ParseTokenStateWithExtensions(data)
}

func FetchAnyAccounts(ctx context.Context, rpcCli *rpc.Client, addresses []solanago.PublicKey) ([]any, error) {
	if len(addresses) == 0 {
		return nil, nil
	}
	resp, err := rpcCli.GetMultipleAccounts(ctx, addresses...)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}
	out := make([]any, len(resp.Value))
	for i, acct := range resp.Value {
		if acct == nil {
			return nil, fmt.Errorf("account not found at index %d", i)
		}
		if !acct.Owner.Equals(ProgramID) {
			return nil, fmt.Errorf("account at index %d has unexpected owner %s (expected %s)", i, acct.Owner, ProgramID)
		}
		decoded, err := ParseAnyAccount(acct.Data.GetBinary())
		if err != nil {
			return nil, fmt.Errorf("failed to parse account at index %d: %w", i, err)
		}
		out[i] = decoded
	}
	return out, nil
}

func fetchToken2022AccountData(
	ctx context.Context,
	rpcCli *rpc.Client,
	address solanago.PublicKey,
) ([]byte, *rpc.Account, error) {
	resp, err := rpcCli.GetAccountInfo(ctx, address)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil || resp.Value == nil {
		return nil, nil, fmt.Errorf("account not found: %s", address)
	}
	if !resp.Value.Owner.Equals(ProgramID) {
		return nil, nil, fmt.Errorf("account %s has unexpected owner %s (expected %s)", address, resp.Value.Owner, ProgramID)
	}
	if resp.Value.Data == nil {
		return nil, nil, fmt.Errorf("account %s has empty data", address)
	}
	return resp.Value.Data.GetBinary(), resp.Value, nil
}
