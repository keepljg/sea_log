package main

import (
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"runtime"
	"sea_log/logs"
	"sea_log/slaver/conf"
	"sea_log/slaver/es"
	"sea_log/slaver/etcd"
	"sea_log/slaver/httpServer"
	"sea_log/slaver/kafka"
	"sea_log/slaver/scheduler"
)

// 初始化线程数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err  error
		quit chan os.Signal
		s    *fasthttp.Server
	)
	// 初始化配置
	if err = conf.InitConf(); err != nil {
		goto ERR
	}

	initEnv()

	// 初始化etcd
	if err = etcd.InitJobMgr(); err != nil {
		goto ERR
	}

	//
	if err = kafka.InitKafka(); err != nil {
		goto ERR
	}

	if err = es.InitElasticClient(); err != nil {
		goto ERR
	}

	scheduler.InitScheduler()
	panic()
	defer scheduler.CancelSelf()

	s = httpServer.InitHttpServer()
	quit = make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logs.INFO("Shutdown Server ...")

	if err := s.Shutdown(); err != nil {
		logs.ERROR(err)
	}

	logs.INFO("Server exiting")

ERR:
	logs.ERROR(err)
	return
}
