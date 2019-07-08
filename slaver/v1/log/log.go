package log

import (
	"github.com/gin-gonic/gin"
	"sea_log/common/sealog_errors"
)

func Mapping(perfix string, app *gin.Engine) {
	log := app.Group(perfix)
	log.POST("/inses", sealog_errors.MiddlewareError(LogToKafka))
	log.GET("/list", sealog_errors.MiddlewareError(LogJobList))
	log.POST("/putjob", sealog_errors.MiddlewareError(LogPutJob))
	log.GET("/deljob", sealog_errors.MiddlewareError(LogDelJob))
	log.GET("/delall", sealog_errors.MiddlewareError(LogAllDelJob))
	log.POST("/delbulk", sealog_errors.MiddlewareError(LogBulkDelJob))
	log.GET("/runwork", sealog_errors.MiddlewareError(GetRuningTopic))
}

func LogToKafka(ctx *gin.Context) error {
	return nil
}

func LogJobList(ctx *gin.Context) error {
	return nil
}

func LogPutJob(ctx *gin.Context) error {
	return nil
}

func LogDelJob(ctx *gin.Context) error {
	return nil
}

func LogAllDelJob(ctx *gin.Context) error {
	return nil
}

func LogBulkDelJob(ctx *gin.Context) error {
	return nil
}

func GetRuningTopic(ctx *gin.Context) error {
	return nil
}
