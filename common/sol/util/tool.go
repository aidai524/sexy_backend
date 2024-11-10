package util

import (
	"crypto/sha256"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"math/rand"
	"regexp"
	"sexy_backend/common/log"
	"strings"
	"time"
)

var (
	DiscriminatorLength = 8
	TimeOut             = time.Duration(3) * time.Second
)

// SolLimitsExceeded 请求被限流判断
var SolLimitsExceeded = func(err error) bool {
	return strings.Contains(err.Error(), "Too many requests for a specific RPC call") ||
		strings.Contains(err.Error(), "Connection rate limits exceeded") ||
		strings.Contains(err.Error(), "Too many requests from your IP")
}

func GetDiscriminator(prefix string) []byte {
	hash := sha256.Sum256([]byte(prefix))
	return hash[:DiscriminatorLength]
}

var solCache *cache.Cache

func WithAllClients(clients []*rpc.Client, action func(client *rpc.Client) (err error)) (err error) {
	if solCache == nil {
		solCache = cache.New(11*time.Second, 1*time.Second)
	}

	var newClients []*rpc.Client
	for {
		for _, client := range clients {
			is := false
			for _, v := range solCache.Items() {
				cc := v.Object.(*rpc.Client)
				if cc == client {
					is = true
					break
				}
			}
			if !is {
				newClients = append(newClients, client)
			}
		}

		if len(newClients) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	var index = time.Duration(0)
	sorted := append([]*rpc.Client{}, newClients...)
	shuffleClients(sorted)
	for _, client := range sorted {
		err = action(client)
		if err == nil {
			return
		} else {
			if SolLimitsExceeded(err) {
				solCache.SetDefault(uuid.NewString(), client)
			}
		}
	}
	if SolLimitsExceeded(err) {
		if index > 10 {
			log.Error("严重限流, 检查ip是否被禁")
			return err
		}
		time.Sleep((time.Duration(index) + 11) * time.Second)
		index++
		err = nil
		return WithAllClients(clients, action)
	}
	return
}

func shuffleClients(slice []*rpc.Client) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 创建独立的随机数生成器
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// 日志转换相关
var programInvocation = regexp.MustCompile(`^Program\s([a-zA-Z0-9]+)?\sinvoke\s\[\d\]$`)
var programFinished = regexp.MustCompile(`^Program\s([a-zA-Z0-9]+)?\s(?:success|error)$`)
var programLogEvent = regexp.MustCompile(`^Program\s(?:log|data):\s([+/0-9A-Za-z]+={0,2})?$`)

func ExtractEvents(logs []string, programIDBase58 string) []string {
	var (
		output          []string
		invocationStack []string
	)
	for _, logStr := range logs {
		if matches := programInvocation.FindStringSubmatch(logStr); matches != nil {
			invokedProgramID := matches[1]
			invocationStack = append(invocationStack, invokedProgramID)
			continue
		}
		if matches := programLogEvent.FindStringSubmatch(logStr); matches != nil {
			currentProgramID := invocationStack[len(invocationStack)-1]
			if programIDBase58 == currentProgramID {
				output = append(output, matches[1])
			}
			continue
		}
		if matches := programFinished.FindStringSubmatch(logStr); matches != nil {
			if len(invocationStack) == 0 {
				break // incorrect execution trace.
			}
			finishedProgramID := matches[1]
			if invocationStack[len(invocationStack)-1] == finishedProgramID {
				invocationStack = invocationStack[:len(invocationStack)-1]
			}
		}
	}
	return output
}
