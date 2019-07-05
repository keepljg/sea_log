package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"sea_log/common"
	"sea_log/logs"
	"sea_log/master/conf"
	"sea_log/master/utils"
)

// 获取当前全部节点
func GetAllNode() ([]string, error) {
	var ips []string
	getResp, err := EtcdClient.KV.Get(context.Background(), conf.JobConf.JobInfo, clientv3.WithPrefix())
	if err != nil {
		return ips, err
	}
	for _, v := range getResp.Kvs {
		ips = append(ips, utils.GetNodeIp(string(v.Key)))
	}
	return ips, err
}

// 获取所有节点压力情况
func GetNodeInfo() map[string]common.NodeInfo {
	nodeinfos := make(map[string]common.NodeInfo)
	getResp, err := EtcdClient.KV.Get(context.Background(), conf.JobConf.JobInfo, clientv3.WithPrefix())
	if err != nil {
		return nodeinfos
	}
	for _, v := range getResp.Kvs {
		if nodeInfo, err := common.UnmarshalJobInfo(v.Value); err != nil {
			nodeinfos[utils.GetNodeIp(string(v.Key))] = nodeInfo
		}
	}
	return nodeinfos
}

// 获取所有注册的任务
func GetAllJob() map[string]*common.Jobs {
	var jobsNmae = make(map[string]*common.Jobs)
	getResp, err := EtcdClient.KV.Get(context.Background(), conf.MasterConf.Jobs, clientv3.WithPrefix())
	if err != nil {
		return jobsNmae
	}

	for _, v := range getResp.Kvs {
		if job, err := common.UnPackJob(v.Value); err == nil {
			jobsNmae[string(v.Key)] = job
		}

	}

	return jobsNmae
}

// 获取所有在运行的job
func GetAllRuningJob() map[string]string {
	res := make(map[string]string)
	getResp, err := EtcdClient.KV.Get(context.Background(), conf.JobConf.JobLock, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		logs.ERROR(err)
		return res
	}
	for _, v := range getResp.Kvs {
		if ipName, err := utils.ExtractRuningJob(string(v.Key)); err == nil {
			res[ipName[1]] = ipName[0]
		}
	}
	return res
}

// 向某个节点注册job
func DistributeJob(ip string, jobs common.Jobs) error {
	jobBytes, err := common.PackJob(jobs)
	if err != nil {
		logs.ERROR(err)
		return err
	}
	if _, err = EtcdClient.KV.Put(context.Background(), utils.ExtractJobSave(ip, jobs.JobName), string(jobBytes)); err != nil {
		logs.ERROR(err)
		return err
	}
	return nil
}

// 向某个节点注销job
func UnDistributeJob(ip string, jobName string) error {
	if _, err := EtcdClient.KV.Delete(context.Background(), utils.ExtractJobSave(ip, jobName)); err != nil {
		logs.ERROR(err)
		return err
	}
	return nil
}
