package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"sea_log/master/conf"
	"time"
)

var EtcdClient *clientv3.Client
var EtcdWatch clientv3.Watcher

func InitJobMgr() error {
	var (
		config clientv3.Config
		err    error
	)
	config = clientv3.Config{
		Endpoints:            []string{conf.EtcdConf.Addr},
		DialKeepAliveTimeout: time.Duration(conf.EtcdConf.EtcdTimeOut) * time.Second,
	}
	if EtcdClient, err = clientv3.New(config); err != nil {
		return err
	}
	EtcdWatch = clientv3.NewWatcher(EtcdClient)
	return nil
}
