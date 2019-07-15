package sealog_errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var errorMsgs = map[string]errorMsg{
	"params_error": errorMsg{
		Code:    1002,
		Message: "参数解析有误",
	},

	"distributeJob_error": errorMsg{
		Code:    1003,
		Message: "节点注册job失败",
	},

	"get_node_error": errorMsg{
		Code:    1004,
		Message: "获取节点失败",
	},

	"undistributeJob_error": errorMsg{
		Code:    1005,
		Message: "节点删除job失败",
	},
}

func GetError(err error) errorMsg {
	if _, ok := errorMsgs[err.Error()]; !ok {
		return errorMsg{
			Code:    1001,
			Message: err.Error(),
		}
	}
	return errorMsgs[err.Error()]
}

type errorHandlefunc func(*gin.Context) error

func MiddlewareError(handlefunc errorHandlefunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err1 := handlefunc(ctx)
		if err1 != nil {
			err2 := GetError(err1)
			ctx.AbortWithStatusJSON(http.StatusOK, err2)
		}
	}
}
