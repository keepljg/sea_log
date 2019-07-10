package admin

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	"sea_log/master/etcd"
	"strings"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("/job", sealog_errors.MiddlewareError(AddLogJob))
	admin.DELETE("/job", sealog_errors.MiddlewareError(DelLogJob))
	admin.POST("/jobs", sealog_errors.MiddlewareError(BulkAddLogJob))
	admin.DELETE("/jobs", sealog_errors.MiddlewareError(BulkDelLogJob))
}

//添加任务
func AddLogJob(ctx *gin.Context) error {
	var jobs common.Jobs
	if ctx.ShouldBind(&jobs) != nil {
		return errors.New("params_error")
	}

	err := etcd.DistributeMasterJob(jobs)
	if err != nil {
		return err
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
			return err
		}
	}
	ctx.JSON(common.Success())
	return nil
}

type BulkAddLogJobForm struct {
	Jobs string `form:"jobs" json:"jobs" binding:"required"`
}

//批量添加job(body传参)
func BulkAddLogJob(ctx *gin.Context) error {
	var bulkAddLogJobForm BulkAddLogJobForm
	var jobs common.Jobs
	if err := ctx.ShouldBind(&bulkAddLogJobForm); err != nil {
		return errors.New("params_error")
	}
	Jobss := strings.Split(bulkAddLogJobForm.Jobs, ";")
	for i := range Jobss {
		if err := json.Unmarshal([]byte(Jobss[i]), &jobs); err != nil {
			return err
		}

		if err := etcd.DistributeMasterJob(jobs); err != nil {
			return err
		}
	}
	ctx.JSON(common.Success())
	return nil
}

type DeleteBulkStrings struct {
	JobNames string `form:"jobNames" json:"jobNames" binding:"required"`
}

//批量删除job(body传参)
func BulkDelLogJob(ctx *gin.Context) error {
	var deleteBulkStrings DeleteBulkStrings
	if err := ctx.ShouldBind(&deleteBulkStrings); err != nil {
		return errors.New("params_error")
	}

	JobNamess := strings.Split(deleteBulkStrings.JobNames, ",")
	jobNodeInfo := etcd.GetAllRuningJob()
	for i := range JobNamess {
		if ip, ok := jobNodeInfo[JobNamess[i]]; ok {
			if err := etcd.UnDistributeJob(ip, JobNamess[i]); err != nil {
				return err
			}
		}
	}
	ctx.JSON(common.Success())
	return nil
}
