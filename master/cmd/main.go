package main

import (
	"sea_log/master/conf"
	"sea_log/master/etcd"
	"sea_log/master/router"
	"sea_log/master/schedule"
)

func main() {
	var err error
	if err = conf.InitConf(); err != nil {

	}

	if err = etcd.InitJobMgr(); err != nil {

	}
	schedule.InitSchedule()
	r := router.Router()
	r.Run(":9200")
}
