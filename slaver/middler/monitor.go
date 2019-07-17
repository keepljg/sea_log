package middler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sea_log/myLoadGen"
	"sea_log/myLoadGen/lib"
	"time"
)

var gen lib.Generator

func InitGen() {
	pset := myLoadGen.ParamSet{
		Caller: nil,
		DurationNS: 10 * time.Second,
	}
	var err error
	gen, err = myLoadGen.NewGenerator(pset)
	if err != nil {
		fmt.Println(err)
	}
	gen.Start()
}

func Mintor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now().UnixNano()
		ctx.Next()
		end := time.Now().UnixNano()
		elapsedTime := time.Duration(end - start)
		gen.SendResult(&lib.CallResult{
			FuncName: ctx.Request.URL.Path,
			//Code:     0,
			//Msg:      "",
			Elapse: elapsedTime,
		})
	}
}
