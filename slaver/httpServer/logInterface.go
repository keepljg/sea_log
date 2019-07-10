package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"sea_log/slaver/kafka"
)

func NewRouter() *fasthttprouter.Router {
	var router *fasthttprouter.Router
	router = fasthttprouter.New()
	router.POST("/log/inses", LogToKafka)
	return router
}

func LogToKafka(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			DoJSONWrite(ctx, 400, GenerateResp("", -2, "failed"))
			return
		}
	}()
	var (
		postRes    []byte
		mapResults map[string]interface{}
		err        error
		topic      string
		logs       []interface{}
		kafkaLogs  []string
		resp       map[string]map[string]interface{}
	)
	postRes = ctx.PostBody()
	if err = json.Unmarshal(postRes, &mapResults); err != nil {
		goto ERR
	}
	topic = mapResults["topic"].(string)
	logs = mapResults["logs"].([]interface{})
	kafkaLogs = make([]string, 0)
	for _, v := range logs {
		kafkaLogs = append(kafkaLogs, v.(string))
	}
	if err = kafka.SendToKafka(kafkaLogs, topic); err != nil {
		goto ERR
	}
	resp = GenerateResp("插入成功", 0, "success")
	DoJSONWrite(ctx, 200, resp)
	return
ERR:
	resp = GenerateResp(err.Error(), -1, "failed")
	DoJSONWrite(ctx, 400, resp)
	return
}

func InitHttpServer() *fasthttp.Server {
	var (
		router *fasthttprouter.Router
		err    error
	)
	router = NewRouter()
	s := &fasthttp.Server{
		Handler: router.Handler,
	}
	go func() {
		if err = s.ListenAndServe(":9100"); err != nil {
			panic(fmt.Sprintf("start fasthttp fail: %v", err))
		}
	}()

	return s
}
