package util

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
	"sexy_backend/common/log"
	"strconv"
	"strings"
)

func FindATA(wallet solana.PrivateKey, mint solana.PublicKey) (ataAccount solana.PublicKey, err error) {
	ataAccount, _, err = solana.FindAssociatedTokenAddress(wallet.PublicKey(), mint)
	return
}

func FindOrCreateATA(client *rpc.Client, wallet solana.PrivateKey, mint solana.PublicKey) (solana.PublicKey, error) {
	// cSol turbo sol AVxnqyCameKsKTCGVKeyJMA7vjHnxJit6afC8AM9MdMj
	// cUSDC turbo usdc HKijBKC2zKcV2BXA9CuNemmWUhTuFkPLLgvQBP7zrQjL
	// mint = solana.MustPublicKeyFromBase58("AVxnqyCameKsKTCGVKeyJMA7vjHnxJit6afC8AM9MdMj")
	ataAccount, _, _ := solana.FindAssociatedTokenAddress(wallet.PublicKey(), mint)
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel() // 确保在函数结束时取消 context，以释放资源
	tokens, err := client.GetTokenAccountsByOwner(ctx, wallet.PublicKey(),
		&rpc.GetTokenAccountsConfig{
			Mint: &mint,
			// ProgramId: solana.TokenProgramID.ToPointer(),
		},

		&rpc.GetTokenAccountsOpts{
			Commitment: "",
			Encoding:   solana.EncodingBase64,
			DataSlice:  nil,
		})
	if err != nil {
		log.Error("got error during searching", err)
	}
	for _, tok := range tokens.Value {
		log.Error("found ata", tok.Pubkey)
		return tok.Pubkey, nil
	}
	log.Error("not found ata, creating ata")
	i := associatedtokenaccount.NewCreateInstruction(
		wallet.PublicKey(),
		wallet.PublicKey(),
		mint,
	).Build()
	for _, tok := range tokens.Value {
		log.Info("token: %v", tok.Pubkey)
	}
	ctx1, cancel1 := context.WithTimeout(context.Background(), TimeOut)
	defer cancel1() // 确保在函数结束时取消 context，以释放资源
	recent, err := client.GetRecentBlockhash(ctx1, rpc.CommitmentFinalized)
	if err != nil {
		return solana.PublicKey{}, err
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			i,
		},
		recent.Value.Blockhash, //NONCE
		solana.TransactionPayer(wallet.PublicKey()),
	)
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if wallet.PublicKey().Equals(key) {
				return &wallet
			}
			return nil
		},
	)
	if err != nil {
		log.Error("unable to sign transaction: %w", err)
		return solana.PublicKey{}, err
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), TimeOut)
	defer cancel2() // 确保在函数结束时取消 context，以释放资源
	sig, err := client.SendTransactionWithOpts(ctx2, tx,
		rpc.TransactionOpts{
			Encoding:            "",
			SkipPreflight:       false,
			PreflightCommitment: "",
			MaxRetries:          nil,
			MinContextSlot:      nil,
		},
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	log.Error("tx for creating ata:", sig)
	// return sig
	return ataAccount, nil
}

// GetBalance
// 获取Solana链的SPL与SOL余额，两者均支持
// calculateDecimal: 是否计算小数点，如果为true，会除以小数位数，得到可视的余额数量
func GetBalance(clients []*rpc.Client, tokenAddr, userAddr string, calculateDecimal bool) (amountDec decimal.Decimal, err error) {
	user, err := solana.PublicKeyFromBase58(userAddr)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	if strings.ToLower(tokenAddr) == "sol" {
		return getSolBalance(clients, user, calculateDecimal)
	} else {
		return getSPLBalance(clients, user, tokenAddr, calculateDecimal)
	}
}

func getSolBalance(clients []*rpc.Client, user solana.PublicKey, calculateDecimal bool) (amountDec decimal.Decimal, err error) {
	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()
		solBalance, err := client.GetBalance(ctx, user, rpc.CommitmentFinalized)
		if err != nil {
			return
		}
		amountDec, err = decimal.NewFromString(strconv.FormatUint(solBalance.Value, 10))
		if err != nil {
			return
		}
		if calculateDecimal {
			amountDec = amountDec.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(9)))
		}
		return
	})
	return
}

func getSPLBalance(clients []*rpc.Client, user solana.PublicKey, tokenAddr string, calculateDecimal bool) (amountDec decimal.Decimal, err error) {
	tokenMint, err := solana.PublicKeyFromBase58(tokenAddr)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	_, _, amount, err := getTokenBalance(clients, tokenMint, user)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	amountDec, err = decimal.NewFromString(strconv.FormatUint(amount, 10))
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	if calculateDecimal {
		tokenInfo, err := GetTokenInfo(clients, tokenAddr)
		if err != nil {
			return decimal.NewFromInt(0), err
		}
		amountDec = amountDec.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt32(tokenInfo.Decimals)))
	}
	return amountDec, nil
}

func getTokenBalance(clients []*rpc.Client, tokenPublicMint, owner solana.PublicKey) (tokenAdd solana.PublicKey, tokenMint solana.PublicKey, amount uint64, err error) {
	err = WithAllClients(clients, func(client *rpc.Client) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()
		tokens, err := client.GetTokenAccountsByOwner(ctx, owner,
			&rpc.GetTokenAccountsConfig{
				ProgramId: solana.TokenProgramID.ToPointer(),
			},
			&rpc.GetTokenAccountsOpts{
				Encoding: solana.EncodingBase64,
			})
		if err != nil {
			return err
		}
		for _, tk := range tokens.Value {
			var ta token.Account
			borshDec := bin.NewBorshDecoder(tk.Account.Data.GetBinary())
			err = borshDec.Decode(&ta)
			if err != nil {
				return
			}
			if ta.Mint == tokenPublicMint {
				tokenAdd = tk.Pubkey
				tokenMint = ta.Mint
				amount = ta.Amount
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	return
}
