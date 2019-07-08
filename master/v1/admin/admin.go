package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sea_log/common/sealog_errors"
	"sea_log/master/utils"
	"sea_log/master/v1/forms"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("/job", adminError(AddLogJob))
	admin.DELETE("/job", adminError(DelLogJob))
	admin.GET("bulk/job", adminError(BulkAddLogJob))
	admin.DELETE("bulk/job", adminError(BulkDelLogJob))
}

type adminHandlefunc func(*gin.Context) error

func adminError(handlefunc adminHandlefunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err1 := handlefunc(ctx)
		if err1 == nil {
			return
		}
		err2 := sealog_errors.GetError(err1)
		ctx.AbortWithStatusJSON(http.StatusOK, err2)
	}
}

func AddLogJob(ctx *gin.Context) error {
	var addLogJobForm forms.AddLogJobForm
	if ctx.ShouldBind(&addLogJobForm) != nil {
		return errors.New("params_error")
	}
	//ip, err := balance.BlanceMapping[conf.BalanceConf.Name]().GetRightNode()
	//runJobs := etcd.GetAllRuningJob()
	//if ip, ok := runJobs[addLogJobForm.JobName]; ok { // 更新job
	//	etcd.DistributeJob(ip, job)
	//} else {
	//	if ip, err := this.b.GetRightNode(); err == nil {
	//		etcd.DistributeJob(ip, job)
	//	}
	//}
	//etcd.DistributeJob( )
	ctx.JSON(utils.Success())
	return nil
}

func DelLogJob(ctx *gin.Context) error {
	return nil
}

func BulkAddLogJob(ctx *gin.Context) error {
	return nil
}

func BulkDelLogJob(ctx *gin.Context) error {
	return nil
}
