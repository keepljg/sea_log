package conf

import (
	"github.com/go-ini/ini"
	"sea_log/helper"
	"sea_log/logs"
	"strings"
)

type Etcd struct {
	Addr        string
	EtcdTimeOut int
}

type Job struct {
	JobSave string
	JobLock string
	JobInfo string
}

type Master struct {
	Jobs string
}

type Balance struct {
	Name string
}

var (
	EtcdConf    Etcd
	JobConf     Job
	MasterConf  Master
	BalanceConf Balance
)

func InitConf() (err error) {
	//confPath := GetRootPath() + "/slaver/conf/conf.ini"
	confPath := strings.Replace(helper.GetRootPath(), "cmd", "conf/conf.ini", 1)
	cfg, err := ini.Load(confPath)
	logs.INFO(confPath)
	err = cfg.Section("etcd").MapTo(&EtcdConf)
	if err != nil {
		logs.ERROR("cfg.MapTo Search settings err: ", err)
	}

	err = cfg.Section("job").MapTo(&JobConf)
	if err != nil {
		logs.ERROR("cfg.MapTo Job settings err: ", err)
	}

	err = cfg.Section("master").MapTo(&MasterConf)
	if err != nil {
		logs.ERROR("cfg.MapTo master settings err: ", err)
	}

	err = cfg.Section("balance").MapTo(&BalanceConf)
	if err != nil {
		logs.ERROR("cfg.MapTo balance settings err: ", err)
	}
	return
}
