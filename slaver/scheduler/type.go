package scheduler

import (
	"context"
	"sea_log/slaver/utils"
)

type Scheduler struct {
	logCount     chan int
	JobEventChan chan *utils.JobEvent
	JobWorkTable map[string]*utils.JobWorkInfo
	HeartCancelFunc context.CancelFunc
}

var Gscheduler *Scheduler
