package admin

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	"sea_log/master/balance"
	"sea_log/master/conf"
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

	runJobs := etcd.GetAllRuningJob()
	if ip, ok := runJobs[jobs.JobName]; ok { // 更新job
		err := etcd.DistributeJob(ip, jobs)
		if err != nil {
			return err
		}
	} else {
		if ip, err := balance.BlanceMapping[conf.BalanceConf.Name]().GetRightNode(); err == nil {
			err := etcd.DistributeJob(ip, jobs)
			if err != nil {
				return err
			}
		} else {
			return errors.New("get_node_error")
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
			return err
		}
	}
	ctx.JSON(common.Success())
	return nil
}

type BulkAddLogJobForm struct {
	Jobs string `form:"jobs" json:"jobs" binding:"required"`
}

func BulkAddLogJob(ctx *gin.Context) error {
	var bulkAddLogJobForm BulkAddLogJobForm
	var jobs common.Jobs
	if err := ctx.ShouldBind(&bulkAddLogJobForm); err != nil {
		return errors.New("params_error")
	}
	data_slice := strings.Split(bulkAddLogJobForm.Jobs, ";")
	runJobs := etcd.GetAllRuningJob()
	for i := range data_slice {
		if err := json.Unmarshal([]byte(data_slice[i]), &jobs); err != nil {
			return err
		}
		if ip, ok := runJobs[jobs.JobName]; ok { // 更新job
			if err := etcd.DistributeJob(ip, jobs); err != nil {
				return err
			}
		} else {
			if ip, err := balance.BlanceMapping[conf.BalanceConf.Name]().GetRightNode(); err == nil {
				if err := etcd.DistributeJob(ip, jobs); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	ctx.JSON(common.Success())
	return nil
}

type DeleteBulkStrings struct {
	JobNames string `form:"jobNames" json:"jobNames" binding:"required"`
}

func BulkDelLogJob(ctx *gin.Context) error {
	var deleteBulkStrings DeleteBulkStrings
	if err := ctx.ShouldBind(&deleteBulkStrings); err != nil {
		return errors.New("params_error")
	}

	data_slice := strings.Split(deleteBulkStrings.JobNames, ",")
	jobNodeInfo := etcd.GetAllRuningJob()
	for i := range data_slice {
		if ip, ok := jobNodeInfo[data_slice[i]]; ok {
			if err := etcd.UnDistributeJob(ip, data_slice[i]); err != nil {
				return err
			}
		}
	}
	ctx.JSON(common.Success())
	return nil
}
