package log

import (
	"github.com/gin-gonic/gin"
	"sea_log/common"
	"sea_log/common/sealog_errors"
	"sea_log/slaver/kafka"
)

func Mapping(perfix string, app *gin.Engine) {
	log := app.Group(perfix)
	log.POST("/inses", sealog_errors.MiddlewareError(LogToKafka))
}

type Data struct {
	Topic      string `form:"topic" json:"topic" binding:"required"`
	Logs       []interface{} `form:"logs" json:"logs" binding:"required"`
}

func LogToKafka(ctx *gin.Context) error {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	var (
	data Data
	err error
	kafkaLogs  []string
	)
	if err = ctx.ShouldBind(&data); err != nil{
		return err
	}
	kafkaLogs = make([]string, 0)
	for _, v := range data.Logs	{
		kafkaLogs = append(kafkaLogs, v.(string))
	}
	if err = kafka.SendToKafka(kafkaLogs, data.Topic); err != nil{
		return err
	}
	ctx.JSON(common.Success())
	return nil
}
