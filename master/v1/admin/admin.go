package admin

import (
	"github.com/gin-gonic/gin"
)


func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("add/job", AddLogJob)
	admin.DELETE("del/job",  DelLogJob)
	admin.GET("bulk/add/job", BulkAddLogJob)
	admin.DELETE("bulk/del/job",  BulkDelLogJob)
}



func AddLogJob(ctx *gin.Context) {

}


func DelLogJob(ctx *gin.Context) {

}

func BulkAddLogJob(ctx *gin.Context) {

}


func BulkDelLogJob(ctx *gin.Context) {

}