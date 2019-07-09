package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

func TestInitJobMgr(t *testing.T) {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	config = clientv3.Config{
		Endpoints:            []string{"192.168.183.103:2379"},
		DialKeepAliveTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		panic(err)
	}
	//job := common.Jobs{
	//	JobName:   "test",
	//	Topic:     "test",
	//	IndexType: "doc",
	//	Pipeline:  "",
	//}
	//jobByte, _ := common.PackJob(job)
	//client.KV.Put(context.Background(), "/master/jobs/test", string(jobByte))

	resp, err :=client.KV.Delete(context.Background(), "/testlog", clientv3.WithPrefix(), clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range resp.PrevKvs {
		fmt.Println(string(v.Key))
	}
	//resp := clientv3.OpGet(utils.JOB_LOCK_DIR + "tutuapp_test",clientv3.WithPrefix())
	//fmt.Println(string(resp.KeyBytes()))
	//fmt.Println(string(resp.ValueBytes()))
}
