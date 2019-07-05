package router

import (
	"github.com/gin-gonic/gin"
	"sea_log/master/v1/admin"
	"sea_log/master/v1/view"
)

func Router() *gin.Engine {
	var app = gin.New()
	app.Use(gin.Recovery())
	view.Mapping("/view", app)
	admin.Mapping("/admin", app)
	return app
}
