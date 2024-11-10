package conf

import "sexy_backend/common/xtime"

type Mysql struct {
	DSN         string
	Opts        string
	Active      int
	Idle        int
	IdleTimeout int
}

type Pgsql struct {
	DSN            string
	Opts           string
	Active         int
	Idle           int
	IdleTimeout    xtime.Duration
	NotAutoMigrate bool
}

type Redis struct {
	Proto       string
	Addr        string
	Auth        string
	Active      int // pool
	Idle        int // pool
	IdleTimeout int // connect max life time.
	DB          int
	Password    string
}

type Kafka struct {
	Addr                    string
	Group                   string
	AutoCommit              bool // true/false
	IgnoreOffsets           []int64
	IgnoreRequestOffsets    []int64
	FetchMinBytes           int
	FetchMaxBytes           int
	FetchWaitMaxMs          int
	ReadMessageTimeout      int
	MaxFetchWaitMaxMs       int
	BatchSize               int
	BatchBytes              int64
	BatchTimeout            int64
	ProducerTimeout         int64
	JudgeIndexRequestOffset uint64
}

type KafkaBusiness struct {
	Topic                string
	SkipOffsetArray      []uint64
	SkipOffset           map[uint64]bool
	MsgTypeReceivedArray []string
	MsgTypeReceivedMap   map[string]bool
	StartOffset          uint64
}
