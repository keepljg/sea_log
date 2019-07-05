package common

type NodeInfo struct {
	LogCount int64 `json:"log_count"`
	JobCount int `json:"job_count"`
}



// log 任务
type Jobs struct {
	JobName   string `json:"job_name"`
	Topic     string `json:"topic"`
	IndexType string `json:"index_type"`
	Pipeline  string `json:"pipeline"`
}


type NodeLogInfos []NodeInfo
type NodeJobInfos []NodeInfo


// 安装node的任务数排序
func (this NodeLogInfos) Len() int {
	return len(this)
}

func (this NodeLogInfos) Less(i, j int) bool {
	if this[i].LogCount < this[j].LogCount {
		return true
	}
	return false
}

func (this NodeLogInfos) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}


// 安装node 单位时间log数排序
func (this NodeJobInfos) Len() int {
	return len(this)
}

func (this NodeJobInfos) Less(i, j int) bool {
	if this[i].JobCount < this[j].JobCount {
		return true
	}
	return false
}

func (this NodeJobInfos) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
