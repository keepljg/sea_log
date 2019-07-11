package view

import (
	"github.com/gin-gonic/gin"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	"sea_log/master/etcd"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.GET("jobs", ListAllLogJob)
	admin.GET("nodes", sealog_errors.MiddlewareError(ListNodes))
	admin.GET("nodes/pree", ListNodePree)
	admin.GET("regjobs", GetAllDistributejob)
}

//获取所有已经注册的任务
func GetAllDistributejob(ctx *gin.Context) {
	data := etcd.GetAllDistributeJob()
	ctx.JSON(common.SuccessWithDate(data))
	return
}

func ListAllLogJob(ctx *gin.Context) {
	data := etcd.GetAllJob()
	ctx.JSON(common.SuccessWithDate(data))
	return
}

func ListNodes(ctx *gin.Context) error {
	data, err := etcd.GetAllNode()
	if err != nil {
		return err
	}
	ctx.JSON(common.SuccessWithDate(data))
	return nil
}

func ListNodePree(ctx *gin.Context) {
	data := etcd.GetNodeInfo()
	ctx.JSON(common.SuccessWithDate(data))
	return
}
