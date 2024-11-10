package util

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"sexy_backend/common/log"
	"strings"
	"time"
)

func GetSlot(clients []*rpc.Client, commitment rpc.CommitmentType) (out uint64, err error) {
	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		out, err = getSlot(client, commitment)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

func getSlot(client *rpc.Client, commitment rpc.CommitmentType) (out uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel() // 确保在函数结束时取消 context，以释放资源
	out, err = client.GetSlot(ctx, commitment)
	if err != nil {
		log.Error("getSlot err: %v", err)
		return
	}
	return
}

func GetBlockWithOpts(clients []*rpc.Client, slot uint64, opts *rpc.GetBlockOpts) (block *rpc.GetBlockResult, err error) {
	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		t := time.Now().UnixMilli()
		block, err = getBlockWithOpts(client, slot, opts)
		for i, r := range clients {
			if client == r {
				log.Info("get block with opts index: %v time: %v", i, time.Now().UnixMilli()-t)
				break
			}
		}
		if err != nil {
			for i, r := range clients {
				if client == r {
					log.Error("getBlockWithOpts index: %v err: %v", i, err)
					break
				}
			}
			if strings.Contains(err.Error(), "was skipped, or missing due to ledger jump to recent snapshot") {
				err = nil
				return
			}
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

func getBlockWithOpts(client *rpc.Client, slot uint64, opts *rpc.GetBlockOpts) (block *rpc.GetBlockResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel() // 确保在函数结束时取消 context，以释放资源

	var index int
	for {
		if index > 5 && err != nil {
			return
		}
		block, err = client.GetBlockWithOpts(ctx, slot, opts)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("Block not available for slot %v", slot)) {
				time.Sleep(time.Second * time.Duration(1+index))
				index++
				continue
			}
			return
		}
		return
	}
}

func GetSignatureStatuses(clients []*rpc.Client, txHash string) (confStatus *rpc.GetSignatureStatusesResult, err error) {
	tx, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		return nil, err
	}
	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()
		confStatus, err = client.GetSignatureStatuses(
			ctx,
			true,
			tx,
		)
		if err != nil {
			return
		}
		return
	})
	return
}

func GetLatestBlockhash(clients []*rpc.Client) (recentBlockhashResp *rpc.GetLatestBlockhashResult, err error) {
	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()
		recentBlockhashResp, err = client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			return err
		}
		return
	})
	if err != nil {
		log.Error("GetRecentBlockhash err: %v", err)
	}
	return
}
