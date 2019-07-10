package log

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	"sea_log/slaver/kafka"
)

func Mapping(perfix string, app *gin.Engine) {
	log := app.Group(perfix)
	log.POST("/inses", sealog_errors.MiddlewareError(LogToKafka))
}

//将日志写入kafka
func LogToKafka(ctx *gin.Context) error {
	topic := ctx.PostForm("topic")
	logs := ctx.PostFormArray("logs")
	if topic == "" || len(logs) == 0 {
		return errors.New("params_error")
	}
	kafkaLogs := make([]string, 0)
	for _, log := range logs {
		kafkaLogs = append(kafkaLogs, log)
	}
	if err := kafka.SendToKafka(kafkaLogs, topic); err != nil {
		return err
	}
	ctx.JSON(common.Success())
	return nil
}
