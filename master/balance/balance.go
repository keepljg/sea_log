package balance

import (
	"errors"
	"sea_log/master/etcd"
	"sea_log/master/utils"
)

type Balancer interface {
	GetRightNode() (string, error)
}

var BlanceMapping = map[string]func()Balancer {
	"random": Random,
	"lessConn": LessConn,
	"lessPree": LessPree,
}

// 随机策略
type random struct {
}

func (this *random) GetRightNode() (string, error) {
	ips, err := etcd.GetAllNode()
	if err != nil {
		return "", err
	}
	index := utils.RangeInt(0, len(ips) -1)
	return ips[index], nil
}

func Random() Balancer {
	return &random{}
}

// 最少连接策略
type lessConn struct {
}

func (this *lessConn)  GetRightNode() (string, error) {
	nodeinfos :=  etcd.GetNodeInfo()
	if len(nodeinfos) == 0 {
		return "", errors.New("have no node")
	}
	var (
		currIp string
		minJobCount int
	)
	for ip, nodeInfo := range nodeinfos {
		if minJobCount <= nodeInfo.JobCount {
			currIp = ip
			minJobCount = nodeInfo.JobCount
		}
	}
	return currIp, nil
}

func LessConn() Balancer {
	return &lessConn{}
}

type lessPree struct {
}

func (this *lessPree) GetRightNode() (string, error) {
	nodeinfos :=  etcd.GetNodeInfo()
	if len(nodeinfos) == 0 {
		return "", errors.New("have no node")
	}
	var (
		currIp string
		minLogCount int
	)
	for ip, nodeInfo := range nodeinfos {
		if int64(minLogCount) <= nodeInfo.LogCount {
			currIp = ip
			minLogCount = nodeInfo.JobCount
		}
	}
	return currIp, nil
}

func LessPree() Balancer {
	return &lessPree{}
}


