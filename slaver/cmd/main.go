package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sea_log/logs"
	"sea_log/slaver/conf"
	"sea_log/slaver/es"
	"sea_log/slaver/etcd_ops"
	"sea_log/slaver/kafka"
	"sea_log/slaver/router"
	"sea_log/slaver/scheduler"
	"time"
)

// 初始化线程数
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 初始化配置
	if err := conf.InitConf(); err != nil {
		logs.ERROR(err)
		return
	}

	// 初始化etcd
	if err := etcd_ops.InitJobMgr(); err != nil {
		logs.ERROR(err)
		return
	}

	//初始化kafka
	if err := kafka.InitKafka(); err != nil {
		logs.ERROR(err)
		return
	}

	//初始化es
	if err := es.InitElasticClient(); err != nil {
		logs.ERROR(err)
		return
	}

	scheduler.InitScheduler()
	defer scheduler.CancelSelf()

	router := router.Router()
	srv := &http.Server{
		Addr:    ":8100",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.FATAL("listen: %s\n", err)
			return
		}
	}()

	// 等待中断信号以优雅地关闭服务器(等待5秒)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	logs.INFO("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.FATAL("Server Shutdown:", err)
		return
	}
	logs.INFO("Server exiting")
}
