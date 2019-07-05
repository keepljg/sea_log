package utils

import (
	"context"
	"sea_log/common"
)


// log 任务事件
type JobEvent struct {
	EventType int
	Job       *common.Jobs
}

type JobWorkInfo struct {
	Job        *common.Jobs
	ConText    context.Context
	CancelFunc context.CancelFunc
}
