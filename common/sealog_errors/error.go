package sealog_errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type errorHandlefunc func(*gin.Context) error

type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var errorMsgs = map[string]errorMsg{
	"default_error": errorMsg{
		Code:    1001,
		Message: "default_error",
	},

	"params_error": errorMsg{
		Code:    1002,
		Message: "参数解析有误",
	},

	"distributeJob_error": errorMsg{
		Code:    1003,
		Message: "节点注册job失败",
	},

	"undistributeJob_error": errorMsg{
		Code:    1004,
		Message: "节点删除job失败",
	},
}

func GetError(err error) errorMsg {
	if err == nil {
		return errorMsgs["default_error"]
	}
	return errorMsgs[err.Error()]
}

func MiddlewareError(handlefunc errorHandlefunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err1 := handlefunc(ctx)
		if err1 == nil {
			return
		}
		err2 := GetError(err1)
		ctx.AbortWithStatusJSON(http.StatusOK, err2)
	}
}
