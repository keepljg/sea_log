package log

import (
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	logs1 "sea_log/logs"
	"sea_log/slaver/kafka"
)

func Mapping(perfix string, app *gin.Engine) {
	log := app.Group(perfix)
	log.POST("/inses", sealog_errors.MiddlewareError(LogToKafka))
}

type LogToKafkaForm struct {
	Topic string `form:"topic" binding:"required"`
	Logs  string `form:"logs" binding:"required"`
}

//将日志写入kafka
func LogToKafka(ctx *gin.Context) error {
	var logToKafkaForm LogToKafkaForm
	logs := make([]map[string]interface{}, 0)
	kafkaLogs := make([]string, 0)

	if err := ctx.ShouldBind(&logToKafkaForm); err != nil {
		return err
	}

	if err := jsoniter.UnmarshalFromString(logToKafkaForm.Logs, &logs); err != nil {
		return err
	}
	for _, log := range logs {
		logstr, err := jsoniter.MarshalToString(log)
		if err != nil {
			return err
		}
		kafkaLogs = append(kafkaLogs, logstr)
	}
	if err := kafka.SendToKafka(kafkaLogs, logToKafkaForm.Topic); err != nil {
		return err
	}
	logs1.INFO("log to kafka!!!")
	ctx.JSON(common.Success())
	return nil
}
