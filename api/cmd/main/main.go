package main

import (
	"flag"
	"os"
	"os/signal"
	"sexy_backend/api/conf"
	"sexy_backend/api/http"
	"sexy_backend/api/service"
	"sexy_backend/common/log"
	"sexy_backend/common/shutdown"
	"syscall"
)

// @title deltaTrade
// @version 1.0
// @description deltaTrade API 文档
// @BasePath /api
func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Init(conf.Conf.Log, conf.Conf.Debug)
	log.Info("API service start")

	service.Init(conf.Conf)
	log.Info("API service init")
	//service.InitConsumer(conf.Conf.Kafka)
	log.Info("API service init consumer")
	http.Init()

	log.Info("API http init")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	s := <-c
	shutdown.StopAndWaitAll()
	log.Info("API exit for signal %v", s)
}
