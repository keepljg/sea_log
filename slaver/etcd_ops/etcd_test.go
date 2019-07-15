package etcd_ops

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"math/rand"
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

	resp, err :=client.KV.Delete(context.Background(), "/master", clientv3.WithPrefix(), clientv3.WithPrevKV())
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

func RangeInt(start, end int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(end-start+1) + start
}


func BenchmarkRandom(b *testing.B)  {
	//arr := make([]int, 0, 1000)
	//for i :=0; i < 1000; i ++ {
	//	arr = append(arr, i)
	//}
	for i := 0; i < 40; i ++ {
		arr := make([]int, 0, 40)
		arr = append(arr, RangeInt(0, 9999))
	}
}

func TestCreateLease(t *testing.T) {
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
	lease := clientv3.NewLease(client)
	leaseGrantResp, err := lease.Grant(context.Background(), 60)
	if err != nil {
		panic(err)
	}

	leaseId := leaseGrantResp.ID
	_, err =client.KV.Put(context.Background(), "/test/123", "v", clientv3.WithLease(leaseId))
	if err != nil {
		panic(err)
	}
}
