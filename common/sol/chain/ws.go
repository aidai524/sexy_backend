package chain

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"sexy_backend/common/log"
	"sync"
	"time"
)

func InitWs(chain *Config, commitment rpc.CommitmentType, callback func(contractName string, logResult *ws.LogResult)) {
	for {
		err := startWs(chain, commitment, callback)
		if err != nil {
			log.Error("Error starting ws: ", err)
			panic(err)
		}
	}
}

func startWs(chain *Config, commitment rpc.CommitmentType, callback func(contractName string, logResult *ws.LogResult)) (err error) {
	// 设置Solana RPC客户端
	var (
		wsClient *ws.Client
		slot     uint64
		subMap   = &sync.Map{}
	)
	defer func() {
		if wsClient != nil {
			wsClient.Close()
		}
		subMap.Range(func(key, value any) bool {
			sub := key.(*ws.LogSubscription)
			sub.Unsubscribe()
			return true
		})
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wsClient, err = ws.Connect(ctx, chain.WsEndpoint)
	if err != nil {
		log.Error("failed to connect to websocket: %v", err)
		return
	}

	chain.ContractAddressMap.Range(func(key, value any) bool {
		var sub *ws.LogSubscription
		contractAddress := key.(solana.PublicKey).String()
		contractName := value.(string)
		// 指定要监听的合约地址
		programID := solana.MustPublicKeyFromBase58(contractAddress)

		// 订阅提及该合约地址的所有交易
		sub, err = wsClient.LogsSubscribeMentions(programID, commitment)
		if err != nil {
			log.Error("failed to subscribe to logs: %v", err)
			return false
		}
		subMap.Store(sub, contractName)
		return true
	})
	if err != nil {
		log.Error("failed to subscribe to logs: %v", err)
		return
	}

	// 订阅slot更新, 确保slot是连续的
	var slotSub *ws.SlotSubscription
	slotSub, err = wsClient.SlotSubscribe()
	if err != nil {
		log.Error("failed to subscribe to logs: %v", err)
		return
	}
	subMap.Store(slotSub, "SlotSubscribe")

	errorCh := make(chan error, 1)
	subMap.Range(func(key, value any) bool {
		switch subData := key.(type) {
		case *ws.LogSubscription:
			var contractAddress string
			chain.ContractAddressMap.Range(func(k, v any) bool {
				if value.(string) == v.(string) {
					contractAddress = k.(solana.PublicKey).String()
					return false
				}
				return true
			})
			go func(sub *ws.LogSubscription, contractAddress string, errorCh chan<- error) {
				for {
					var (
						logResult *ws.LogResult
					)
					logResult, err = sub.Recv()
					if err != nil {
						errorCh <- err
						return
					}
					if logResult.Value.Err != nil {
						log.Error("failed to receive log: %v", logResult.Value.Err)
						continue
					}
					if len(logResult.Value.Logs) > 0 {
						if contractNameVolume, ok := subMap.Load(sub); ok {
							contractName := contractNameVolume.(string)
							callback(contractName, logResult)
						}
					}
				}
			}(subData, contractAddress, errorCh)
		case *ws.SlotSubscription:
			go func(sub *ws.SlotSubscription, errorCh chan<- error) {
				for {
					var (
						slotResult *ws.SlotResult
					)
					slotResult, err = sub.Recv()
					if err != nil {
						errorCh <- err
						return
					}
					if slot == 0 {
						slot = slotResult.Slot
					} else if slot+1 == slotResult.Slot {
						slot = slotResult.Slot
					} else {
						// ws 有断流, 检查网络情况
						log.Warn("failed to receive slot: %v, new slot: %v", slot, slotResult.Slot)
						errorCh <- fmt.Errorf("failed to receive slot: %v, new slot: %v", slot, slotResult.Slot)
					}
				}
			}(subData, errorCh)
		}
		return true
	})

	select {
	case err = <-errorCh:
		log.Error("error received from log subscription: %v", err)
		return err
	}
}
