package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
	"sexy_backend/common/log"
	"sort"
	"strings"
	"sync"
	"time"
)

type TxParams struct {
	FromAccountOwner string          `json:"solSourceTokenAccountOwner"`
	ToAccountOwner   string          `json:"solDestTokenAccountOwner"`
	ToAccount        string          `json:"toAddress"`
	Amount           decimal.Decimal `json:"amount"`
	Fee              decimal.Decimal `json:"fee"`
	Memo             string          `json:"memo"`
	RecentBlockHash  string          `json:"solRecentBlockHash"`
	ContractAddress  string          `json:"commonContractAddress"`
	FromAccounts     []TokenAccount  `json:"fromAccounts"`
	TokenDecimal     int32           `json:"commonTokenDecimal"`
}

type TokenAccount struct {
	TokenAccount string          `json:"tokenAccount"`
	Amount       decimal.Decimal `json:"amount"`
	MintAddress  string          `json:"mintAddress"`
}

func CreateTransaction(clients []*rpc.Client, params TxParams) (*solana.Transaction, error) {
	extraParams := params
	sourceOwnerKey, err := solana.PublicKeyFromBase58(extraParams.FromAccountOwner)
	if err != nil {
		return nil, err
	}
	builder1 := solana.NewTransactionBuilder()
	if extraParams.ToAccount == "" && extraParams.ToAccountOwner == "" {
		return nil, fmt.Errorf("must provide either to account or toAccountOwner to create transaction")
	}
	if extraParams.ToAccount == "" {
		ownerKey, err := solana.PublicKeyFromBase58(extraParams.ToAccountOwner)
		if err != nil {
			return nil, err
		}
		mintAccount, err := solana.PublicKeyFromBase58(extraParams.ContractAddress)
		if err != nil {
			return nil, err
		}
		addr, _, err := solana.FindAssociatedTokenAddress(ownerKey, mintAccount)
		if err != nil {
			return nil, err
		}
		// 判断用户是否创建了ata地址
		ataAccountInfo, err := GetAccountInfo(clients, addr)
		if err != nil {
			if err.Error() == "not found" {
				err = nil
			} else {
				return nil, err
			}
		}
		if ataAccountInfo != nil {
			extraParams.ToAccount = addr.String()
		} else {
			extraParams.ToAccount = addr.String()
			builder := associatedtokenaccount.NewCreateInstructionBuilder()
			builder.SetMint(mintAccount)
			builder.SetPayer(sourceOwnerKey)
			meta := solana.NewAccountMeta(addr, true, false)
			builder.Append(meta)
			builder.SetWallet(ownerKey)
			builder1.AddInstruction(builder.Build())
		}
	}
	// Create a timestamp memo to distinguish different transactions.
	meta := solana.NewAccountMeta(sourceOwnerKey, true, true)
	instruction := solana.NewInstruction(solana.MemoProgramID, solana.AccountMetaSlice{meta}, []byte(fmt.Sprint(time.Now().UnixMilli())))
	builder1.AddInstruction(instruction)
	//Create on-chain memo
	if extraParams.Memo != "" {
		builder1.AddInstruction(solana.NewInstruction(solana.MemoProgramID, solana.AccountMetaSlice{meta}, []byte(extraParams.Memo)))
	}
	if len(extraParams.FromAccounts) == 0 {
		return nil, fmt.Errorf("must provide either from account or fromAccountOwner to create transaction")
	}
	sort.SliceStable(extraParams.FromAccounts, func(i, j int) bool {
		return extraParams.FromAccounts[i].Amount.LessThan(extraParams.FromAccounts[j].Amount)
	})
	//Calculate balances for accounts with the same token type, and provide transfers if the balance is not zero.
	var transferAccounts []TokenAccount
	var amountSum = decimal.NewFromInt(0)
	var amountDrops = decimal.NewFromFloat(float64(10)).
		Pow(decimal.NewFromFloat(float64(extraParams.TokenDecimal))).
		Mul(extraParams.Amount).RoundFloor(0)
	for _, account := range extraParams.FromAccounts {
		if decimal.NewFromInt(0).LessThan(account.Amount) {
			amountSum = amountSum.Add(account.Amount)
			if amountSum.GreaterThanOrEqual(extraParams.Amount) {
				sub := amountSum.Sub(amountDrops)
				d := account.Amount.Sub(sub)
				account.Amount = d
				transferAccounts = append(transferAccounts, account)
				break
			} else {
				transferAccounts = append(transferAccounts, account)
			}
		}
	}
	//transaction
	if amountSum.LessThan(extraParams.Amount) {
		return nil, fmt.Errorf("must provide either from account or fromAccountOwner to create transaction")
	}
	for _, account := range transferAccounts {
		builder := token.NewTransferInstructionBuilder()
		builder.SetAmount(uint64(account.Amount.IntPart()))
		builder.SetOwnerAccount(sourceOwnerKey, sourceOwnerKey)
		sourceAccount, err := solana.PublicKeyFromBase58(account.TokenAccount)
		if err != nil {
			return nil, err
		}
		builder.SetSourceAccount(sourceAccount)
		toAccount, err := solana.PublicKeyFromBase58(extraParams.ToAccount)
		if err != nil {
			return nil, err
		}
		builder.SetDestinationAccount(toAccount)
		build, err := builder.ValidateAndBuild()
		builder1.AddInstruction(build)
	}

	builder1.SetFeePayer(sourceOwnerKey)
	fromBase58, err := solana.HashFromBase58(extraParams.RecentBlockHash)
	if err != nil {
		return nil, err
	}
	builder1.SetRecentBlockHash(fromBase58)
	transaction, err := builder1.Build()
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func CreateTransactionNativeSol(clients []*rpc.Client, amount decimal.Decimal, from solana.PublicKey, to solana.PublicKey) (*solana.Transaction, error) {
	// 指定转账金额（以 lamports 为单位，1 SOL = 1,000,000,000 lamports）
	amount = amount.Mul(decimal.NewFromInt(1000000000)).RoundFloor(0)

	if amount.IsZero() {
		return nil, fmt.Errorf("must provide amount")
	}

	if from.IsZero() {
		return nil, fmt.Errorf("must provide fromAccountOwner")
	}

	// 创建转账指令
	transferInstruction := system.NewTransferInstruction(
		uint64(amount.IntPart()),
		from,
		to,
	).Build()

	recentBlockHashResp, err := GetLatestBlockhash(clients)
	if err != nil {
		log.Error("CreateTransactionNativeSol - GetLatestBlockhash err: %v", err)
		return nil, err
	}

	// 创建交易
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			transferInstruction,
		},
		recentBlockHashResp.Value.Blockhash,
		// 设置发件人的公钥作为交易的 fee payer
		solana.TransactionPayer(from),
	)
	if err != nil {
		log.Error("CreateTransactionNativeSol - err: %v", err)
		return nil, err
	}
	return tx, nil
}

// SendTransactionWithOpts 发送交易, 交易使用并发的方式发送
func SendTransactionWithOpts(clients []*rpc.Client, transaction *solana.Transaction) (txHash solana.Signature, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstHash solana.Signature
	var firstErr error
	var success bool // 标志是否有成功的交易

	// 为每个客户端并发发送交易
	for _, client := range clients {
		wg.Add(1)
		go func(client *rpc.Client) {
			defer wg.Done()
			var maxRetries = uint(2)

			t := time.Now().UnixMilli()
			hash, sendErr := client.SendTransactionWithOpts(ctx, transaction, rpc.TransactionOpts{
				SkipPreflight:       true,
				PreflightCommitment: rpc.CommitmentFinalized,
				MaxRetries:          &maxRetries,
				MinContextSlot:      nil,
			})
			if sendErr != nil {
				log.Error("Failed to send transaction: %v", sendErr)
			}
			for i, c := range clients {
				if c == client {
					log.Info("并发交易发送: %v, hash: %v, time: %v, err: %v", i, hash, time.Now().UnixMilli()-t, sendErr)
					break
				}
			}
			mu.Lock()
			defer mu.Unlock()

			if sendErr == nil && !success {
				// 如果有成功的交易，将第一个成功的hash存储起来
				firstHash = hash
				success = true
			} else if sendErr != nil && firstErr == nil {
				// 如果出现错误且还没有记录第一个错误，记录下来
				firstErr = sendErr
			}
		}(client)
	}

	// 等待所有 goroutines 完成
	wg.Wait()

	// 检查是否有成功的交易
	if success {
		// 返回第一个成功的交易hash
		return firstHash, nil
	}

	// 如果没有成功的交易，返回第一个捕获到的错误
	if firstErr != nil {
		return solana.Signature{}, firstErr
	}

	// 理论上不应该到达这里，但防止出现未知的情况
	return solana.Signature{}, fmt.Errorf("unknown error: no transaction success and no error found")
}

// CheckTransactionStats 查询交易状态
func CheckTransactionStats(clients []*rpc.Client, txHash string) (err error) {
	return CheckTransactionStateRetry(clients, txHash, 50)
}

func CheckTransactionStateRetry(clients []*rpc.Client, txHash string, retryNumber int) (err error) {
	var index int
	for {
		if index > retryNumber {
			// 超过50s 大概率是没有成功发出去的交易
			log.Warn("CheckTransactionStats time out: %v", txHash)
			return fmt.Errorf("failed to execute thread err time out")
		}
		index++
		// 查询交易确认状态
		var confStatus *rpc.GetSignatureStatusesResult
		confStatus, err = GetSignatureStatuses(clients, txHash)
		if err != nil {
			log.Error("Failed to get transaction status: %v", err)
			return
		}

		// 检查返回结果
		if confStatus.Value[0] == nil {
			log.Warn("Transaction status not available yet for hash: %s", txHash)
			time.Sleep(1000 * time.Millisecond)
			continue
		} else {
			status := confStatus.Value[0]
			if status.ConfirmationStatus == rpc.ConfirmationStatusConfirmed || status.ConfirmationStatus == rpc.ConfirmationStatusFinalized {
				if status.Err != nil {
					log.Error("Failed to execute thread err: %v", status.Err)
					err = errors.New("failed to execute thread err")
					return
				}
				err = nil
				break
			} else {
				log.Warn("Transaction %s is still pending confirmation", txHash)
				break
			}
		}
	}
	return
}

func GetTransactionSync(clients []*rpc.Client, txHash string) (transaction *rpc.GetTransactionResult, err error) {
	var tx solana.Signature
	tx, err = solana.SignatureFromBase58(txHash)
	if err != nil {
		log.Error("Failed to get transaction signature: %v, %v", txHash, err)
		return
	}
	if err = CheckTransactionStats(clients, txHash); err != nil {
		if strings.Contains(err.Error(), "failed to execute thread err") {
			// 失败的交易
			return
		}
		return
	}

	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()
		transaction, err = client.GetTransaction(
			ctx,
			tx,
			&rpc.GetTransactionOpts{
				Commitment: rpc.CommitmentConfirmed,
			},
		)
		if err != nil {
			log.Error("Failed to get transaction: %v, err: %v", txHash, err)
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}
