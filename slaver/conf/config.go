package conf

import (
	"fmt"
	"github.com/go-ini/ini"
	"sea_log/helper"
	"sea_log/logs"
	"strings"
)

type Etcd struct {
	Addr        string
	EtcdTimeOut int
}

type Kafka struct {
	Addr string
}

type Job struct {
	JobSave string
	JobLock string
	JobInfo string
}

type Es struct {
	Addr []string
}

var (
	KafkaConf Kafka
	EsConf    Es
	EtcdConf  Etcd
	JobConf   Job
)

func InitConf() (err error) {
	//confPath := GetRootPath() + "/slaver/conf/conf.ini"
	confPath := strings.Replace(helper.GetRootPath(), "cmd", "conf/conf.ini", -1)
	cfg, err := ini.Load(confPath)
	logs.INFO(confPath)
	err = cfg.Section("kafka").MapTo(&KafkaConf)
	if err != nil {
		logs.ERROR("cfg.MapTo kafka settings err: ", err)
	}
	err = cfg.Section("etcd").MapTo(&EtcdConf)
	if err != nil {
		logs.ERROR("cfg.MapTo Search settings err: ", err)
	}
	err = cfg.Section("job").MapTo(&JobConf)
	if err != nil {
		logs.ERROR("cfg.MapTo Job settings err: ", err)
	}
	esAddr, err := cfg.Section("es").GetKey("Addr")
	if err != nil {
		logs.ERROR("cfg.MapTo Es settings err: ", err)
	}
	EsConf.Addr = strings.Split(esAddr.String(), ",")
	localIp, err := helper.LocalIPv4s()
	if err != nil {
		logs.ERROR("Get LocalIPv4s err:", err)
	}
	JobConf.JobLock = fmt.Sprintf("%s%s/", JobConf.JobLock, localIp)
	JobConf.JobInfo = fmt.Sprintf("%s%s/", JobConf.JobInfo, localIp)
	JobConf.JobSave = fmt.Sprintf("%s%s/", JobConf.JobSave, localIp)
	return
}
