package shutdown

import (
	"os"
	"time"
)

const (
	lockExpiration  = 10 * time.Second
	refreshInterval = 5 * time.Second
)

// Single instance keys
const (
	PersistService              = "singleInstance:persistService"
	CacheService                = "singleInstance:cacheService"
	OrderBookWsJob              = "singleInstance:orderBookWsJob"
	OrderBookApiJob             = "singleInstance:orderBookApiJob"
	OrderBookQuotesJob          = "singleInstance:orderBookQuotesJob"
	OrderBookMining             = "singleInstance:orderBookMining"
	DataPanel                   = "singleInstance:dataPanel"
	Operator                    = "singleInstance:operator"   // not in this repo
	OperatorGo                  = "singleInstance:operatorGo" // not in this repo
	StateCheck                  = "singleInstance:stateCheck" // not in this repo
	OrderBookConsumer           = "singleInstance:orderBookConsumer"
	OrderBookReward             = "singleInstance:orderBookReward"
	OrderBookArbSubmitBlocks    = "singleInstance:orderBookArbSubmitBlocks"
	OrderBookArbSubmitBlocksJob = "singleInstance:orderBookArbSubmitBlocksJob"
	OrderBookMiningConsumer     = "singleInstance:orderBookMiningConsumer"
	OrderBookGasJob             = "singleInstance:orderBookGasJob"
	EthScanner                  = "singleInstance:ethScanner"
	OrderBookBoss               = "singleInstance:orderBookBoss" // TODO: Update after merging newest BOSS branch into testnet
)

func HandleReleaseLockSignal(releaseLockCh chan os.Signal, sig os.Signal) {
	releaseLockCh <- sig

	// wait for the single instance key to be deleted
	time.Sleep(time.Second * 1)
}
