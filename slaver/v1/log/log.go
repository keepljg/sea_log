package log

import (
	"github.com/gin-gonic/gin"
	"sea_log/common/sealog_errors"
)

func Mapping(perfix string, app *gin.Engine) {
	log := app.Group(perfix)
	log.POST("/inses", sealog_errors.MiddlewareError(LogToKafka))
}

func LogToKafka(ctx *gin.Context) error {
	return nil
}
