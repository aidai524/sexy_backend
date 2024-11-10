package chain

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"sexy_backend/common/log"
	"sexy_backend/common/sol/util"
	"strings"
	"sync"
	"time"
)

type Config struct {
	RpcEndpoint          []string
	WsEndpoint           string
	RpcClient            *rpc.Client
	RpcClientList        []*rpc.Client
	NotPushKafkaMessage  bool
	BlockNumbers         []uint64
	NetworkID            string
	StartBlock           uint64
	EndBlock             uint64
	AddStartBlock        uint64
	AddEndBlock          uint64
	Concurrency          uint64
	ContractAddressList  map[string]string
	ContractAddressMap   *sync.Map
	DeltaBotAddress      string
	DeltaBotStateAddress string
	Eof                  bool
	Timeout              int
}

func (c *Config) Init() {
	for _, rpcUlr := range c.RpcEndpoint {
		client := InitRpc(rpcUlr)
		c.RpcClientList = append(c.RpcClientList, client)
		c.RpcClient = client
	}
	c.ContractAddressMap = &sync.Map{}
	for key, volume := range c.ContractAddressList {
		out, err := solana.PublicKeyFromBase58(key)
		if err != nil {
			log.Error("Failed to convert solana public key to bytes err: %v", err)
			panic(err)
		}
		c.ContractAddressMap.Store(out, volume)
		if volume == "deltabot" {
			c.DeltaBotAddress = out.String()
		}
		if volume == "deltabot_state" {
			c.DeltaBotStateAddress = out.String()
		}
	}
	if c.Timeout > 0 {
		util.TimeOut = time.Duration(c.Timeout) * time.Second
	}
}

func InitRpc(rpcUlr string) *rpc.Client {
	var client *rpc.Client
	if strings.Contains(rpcUlr, "api-v2.solscan.io") {
		client = rpc.NewWithHeaders(rpcUlr, map[string]string{
			"origin": "https://solscan.io",
		})
	} else if strings.Contains(rpcUlr, "explorer-api") {
		client = rpc.NewWithHeaders(rpcUlr, map[string]string{
			"origin":  "https://explorer.solana.com",
			"referer": "https://explorer.solana.com/",
		})
	} else if strings.Contains(rpcUlr, "jupiter") || strings.Contains(rpcUlr, "mercuria-fronten-1cd8") {
		client = rpc.NewWithHeaders(rpcUlr, map[string]string{
			"accept":             "*/*",
			"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8",
			"content-type":       "application/json",
			"origin":             "https://jup.ag",
			"priority":           "u=1, i",
			"referer":            "https://jup.ag/",
			"sec-ch-ua":          "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"127\", \"Chromium\";v=\"127\"",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": "\"macOS\"",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "cross-site",
			"solana-client":      "js/0.0.0-development",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		})
	} else if strings.Contains(rpcUlr, "greatest-old-orb.solana-mainnet.quiknode.pro") {
		client = rpc.NewWithHeaders(rpcUlr, map[string]string{
			"accept":             "*/*",
			"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8",
			"content-type":       "application/json",
			"origin":             "https://lifinity.io",
			"priority":           "u=1, i",
			"referer":            "https://lifinity.io/",
			"sec-ch-ua":          "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"127\", \"Chromium\";v=\"127\"",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": "\"macOS\"",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "cross-site",
			"solana-client":      "js/0.0.0-development",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		})
	} else if strings.Contains(rpcUlr, "ellipsis-main-98a6.mainnet.rpcpool.com") {
		client = rpc.NewWithHeaders(rpcUlr, map[string]string{
			"accept":             "*/*",
			"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8",
			"content-type":       "application/json",
			"origin":             "https://app.phoenix.trade",
			"priority":           "u=1, i",
			"referer":            "https://app.phoenix.trade/",
			"sec-ch-ua":          "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"127\", \"Chromium\";v=\"127\"",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": "\"macOS\"",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "cross-site",
			"solana-client":      "js/0.0.0-development",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		})
	} else {
		client = rpc.New(rpcUlr)
	}
	return client
}
