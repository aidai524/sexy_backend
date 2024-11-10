package util

import (
	"errors"
	bin "github.com/gagliardetto/binary"
	tokenmetadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"sexy_backend/common/log"
	"sexy_backend/common/sexy/model"
	"strings"
)

// GetTokenInfo 获取Token的基本信息
func GetTokenInfo(clients []*rpc.Client, tokenAddress string) (tokenInfo *model.TokenInfo, err error) {
	tokenMintAddress, err := solana.PublicKeyFromBase58(tokenAddress)
	if err != nil {
		return
	}

	var mint *token.Mint
	mintAccountInfo, err := GetAccountInfo(clients, tokenMintAddress)
	if err != nil {
		log.Error("")
		return
	}

	err = bin.NewBinDecoder(mintAccountInfo.Value.Data.GetBinary()).Decode(&mint)
	if err != nil {
		return
	}
	tokenInfo = &model.TokenInfo{
		Address:  tokenAddress,
		Decimals: int32(mint.Decimals),
	}

	var metadataAddress solana.PublicKey
	metadataAddress, _, err = solana.FindTokenMetadataAddress(tokenMintAddress)
	if err != nil {
		log.Error("GetTokenInfo - FindTokenMetadataAddress: %v err: %v", tokenMintAddress, err)
		return
	}

	var metadataAddressResult *rpc.GetAccountInfoResult
	metadataAddressResult, err = GetAccountInfo(clients, metadataAddress)
	if err != nil {
		log.Error("Failed to execute thread err: %v", err)
		return
	}
	d := bin.NewBorshDecoder(metadataAddressResult.Value.Data.GetBinary())

	var metadata *tokenmetadata.Metadata
	err = d.Decode(&metadata)
	if err != nil {
		log.Error("Failed to execute thread err: %v", err)
		return
	}
	tokenInfo.Symbol = strings.TrimRight(metadata.Data.Symbol, "\u0000")
	tokenInfo.Name = strings.TrimRight(metadata.Data.Name, "\u0000")
	if tokenInfo.Symbol == "" || tokenInfo.Name == "" {
		err = errors.New("token info not found")
		return
	}
	return
}

func GetTokenMintInfo(clients []*rpc.Client, tokenAddress string) (mint *token.Mint, err error) {
	tokenMintAddress, err := solana.PublicKeyFromBase58(tokenAddress)
	if err != nil {
		return
	}

	mintAccountInfo, err := GetAccountInfo(clients, tokenMintAddress)
	if err != nil {
		log.Error("")
		return
	}

	err = bin.NewBinDecoder(mintAccountInfo.Value.Data.GetBinary()).Decode(&mint)
	if err != nil {
		return
	}
	return
}
