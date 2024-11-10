package service

import (
	"sexy_backend/api/conf"
	"sexy_backend/api/dao"
	"sexy_backend/common/supabase"
	"time"
)

var (
	API *Service
)

type Service struct {
	debug          bool
	Timeout        time.Duration
	Dao            *dao.Dao
	supabaseClient *supabase.Client
}

func Init(c *conf.Config) {
	d := dao.New(c)
	API = &Service{
		debug:          c.Debug,
		Dao:            d,
		Timeout:        time.Second * time.Duration(c.Timeout),
		supabaseClient: supabase.NewClient(c.Supabase.Url, c.Supabase.Key),
	}
}
