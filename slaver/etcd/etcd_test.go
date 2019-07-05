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
	//fmt.Println(client.KV.Put(context.Background(), "/test/1", ""))
	clientv3.WithKeysOnly()
	fmt.Println(client.KV.Put(context.Background(), "/test/2/2", ""))
	fmt.Println(client.KV.Put(context.Background(), "/test/1/3", ""))
	//client.KV.Delete(context.Background(), "/test/1/", clientv3.WithPrefix())
	//fmt.Println(client.KV.Delete(context.Background(), "test/1"))
	resp, err :=client.KV.Get(context.Background(), "/test/", clientv3.WithKeysOnly(), clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range resp.Kvs {
		fmt.Println(string(v.Key))
	}
	//resp := clientv3.OpGet(utils.JOB_LOCK_DIR + "tutuapp_test",clientv3.WithPrefix())
	//fmt.Println(string(resp.KeyBytes()))
	//fmt.Println(string(resp.ValueBytes()))
}
