package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
)

func CreateLease(lease clientv3.Lease, ttl int64) (clientv3.LeaseID, <-chan *clientv3.LeaseKeepAliveResponse, context.Context, context.CancelFunc, error){
	var (
		leaseGrantResp      *clientv3.LeaseGrantResponse
		leaseId             clientv3.LeaseID
		leaseKeepActiveChan <-chan *clientv3.LeaseKeepAliveResponse
		err error
	)
	ctx, cancelFunc := context.WithCancel(context.TODO())

	if leaseGrantResp, err = lease.Grant(context.Background(), ttl); err != nil {
		return leaseId, leaseKeepActiveChan, ctx, cancelFunc, err
	}

	leaseId = leaseGrantResp.ID

	// 进行续租
	if leaseKeepActiveChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		return leaseId, leaseKeepActiveChan, ctx, cancelFunc, err
	}
	return leaseId, leaseKeepActiveChan, ctx, cancelFunc, nil
}