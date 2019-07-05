package view

import "github.com/gin-gonic/gin"

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.GET("list/all", ListAllLogJob)
	admin.GET("list/nodes",  ListNodes)
	admin.GET("list/nodes/pree", ListNodePree)

}


func ListAllLogJob(ctx *gin.Context) {

}


func ListNodes(ctx *gin.Context) {

}


func ListNodePree(ctx *gin.Context) {

}


