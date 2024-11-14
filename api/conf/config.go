package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"sexy_backend/common/conf"
	"sexy_backend/common/log"
	"sexy_backend/common/sol/chain"
	"sexy_backend/common/supabase"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	Debug        bool
	Timeout      int64
	Port         int
	Log          *log.Config
	Pgsql        *Pgsql
	DB           *gorm.DB
	Redis        *conf.Redis
	Red          *redis.Pool
	Cors         bool
	AllowOrigins []string
	Supabase     *supabase.Supabase
	SolChain     *chain.Config
}

type Pgsql struct {
	DB *conf.Pgsql
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if err != nil {
		log.Error("error decoding [%v]:%v", confPath, err)
		return
	}

	Conf.SolChain.Init()
	return
}
