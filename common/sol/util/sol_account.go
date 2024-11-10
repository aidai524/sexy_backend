package util

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetAccountInfo(clients []*rpc.Client, account solana.PublicKey) (out *rpc.GetAccountInfoResult, err error) {
	err = WithAllClients(clients, func(client *rpc.Client) error {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()
		out, err = client.GetAccountInfo(ctx, account)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}
