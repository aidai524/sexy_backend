package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"net/http"
	"sexy_backend/common/ecode"
	"sexy_backend/common/log"
	"sexy_backend/common/sexy/dao"
	"sexy_backend/common/sexyerror"
	"time"
)

func (s *Service) GetAccountToken(address string, date int64, signature string) (token string, err error) {
	if len(signature) == 0 || len(address) == 0 || date == 0 {
		err = &sexyerror.Error{Code: http.StatusUnauthorized, Message: "signature and address and time required"}
		return
	}

	nowTime := time.Now().UnixMilli()
	diffTime := nowTime - date
	if diffTime > 60000 || diffTime < (0-60000) {
		log.Warn("checkAuth Time exceeds the setting error, now time: %v, user time: %v, auth mix time: %v", nowTime, date, 60000)
		err = &sexyerror.Error{Code: http.StatusUnauthorized, Message: "Time exceeds the setting"}
		return
	}

	message := fmt.Sprintf("login sexy,time:%v", date)
	messageBytes := []byte(message)

	if messageBytes == nil || len(messageBytes) == 0 {
		err = &sexyerror.Error{Code: http.StatusUnauthorized, Message: "signature message error"}
		return
	}

	publicKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		err = &sexyerror.Error{Code: http.StatusUnauthorized, Message: "address error"}
		return
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		err = &sexyerror.Error{Code: http.StatusUnauthorized, Message: "signature base64 error"}
		return
	}

	signatureBytes := solana.SignatureFromBytes(decodedBytes)

	// 校验签名
	isVerified := publicKey.Verify(messageBytes, signatureBytes)
	if !isVerified {
		err = &sexyerror.Error{Code: http.StatusUnauthorized, Message: "Signature Verification error"}
		return
	}

	token, err = GenerateToken(64)
	if err != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: err.Error()}
		return
	}
	err = dao.SetAuth(s.Dao.GetRedis(), token, address)
	if err != nil {
		err = &sexyerror.Error{Code: ecode.UnknownError, Message: err.Error()}
		return
	}
	return
}

// GenerateToken generates a random token of the specified length in bytes.
func GenerateToken(length int) (string, error) {
	// Create a byte slice with the specified length
	token := make([]byte, length)

	// Fill the byte slice with random data
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a hexadecimal string
	return hex.EncodeToString(token), nil
}
