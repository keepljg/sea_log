package admin

import (
	"github.com/gin-gonic/gin"
)


func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("/job", AddLogJob)
	admin.DELETE("/job",  DelLogJob)
	admin.GET("/bulk/job", BulkAddLogJob)
	admin.DELETE("/bulk/job",  BulkDelLogJob)
}



func AddLogJob(ctx *gin.Context) {

}


func DelLogJob(ctx *gin.Context) {

}

func BulkAddLogJob(ctx *gin.Context) {

}


func BulkDelLogJob(ctx *gin.Context) {

}