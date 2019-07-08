package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	"sea_log/master/balance"
	"sea_log/master/conf"
	"sea_log/master/etcd"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("/job", sealog_errors.MiddlewareError(AddLogJob))
	admin.DELETE("/job", sealog_errors.MiddlewareError(DelLogJob))
	admin.POST("bulk/job", sealog_errors.MiddlewareError(BulkAddLogJob))
	admin.DELETE("bulk/job", sealog_errors.MiddlewareError(BulkDelLogJob))
}

//添加任务
func AddLogJob(ctx *gin.Context) error {
	var jobs common.Jobs
	if ctx.ShouldBind(&jobs) != nil {
		return errors.New("params_error")
	}

	runJobs := etcd.GetAllRuningJob()
	if ip, ok := runJobs[jobs.JobName]; ok { // 更新job
		err := etcd.DistributeJob(ip, jobs)
		if err != nil {
			return errors.New("distributeJob_error")
		}
	} else {
		if ip, err := balance.BlanceMapping[conf.BalanceConf.Name]().GetRightNode(); err == nil {
			err := etcd.DistributeJob(ip, jobs)
			if err != nil {
				return errors.New("distributeJob_error")
			}
		}
	}
	ctx.JSON(common.Success())
	return nil
}

type DelLogJobForm struct {
	JobName string `form:"jobName" binding:"required"`
}

//删除任务
func DelLogJob(ctx *gin.Context) error {
	var delLogJobForm DelLogJobForm
	if ctx.ShouldBind(&delLogJobForm) != nil {
		return errors.New("params_error")
	}

	jobNodeInfo := etcd.GetAllRuningJob()
	if ip, ok := jobNodeInfo[delLogJobForm.JobName]; ok {
		err := etcd.UnDistributeJob(ip, delLogJobForm.JobName)
		if err != nil {
			return errors.New("distributeJob_error")
		}
	}
	ctx.JSON(common.Success())
	return nil
}

func BulkAddLogJob(ctx *gin.Context) error {

	return nil
}

func BulkDelLogJob(ctx *gin.Context) error {

	return nil
}
