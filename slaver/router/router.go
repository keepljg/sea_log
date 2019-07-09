package router

import (
	"github.com/gin-gonic/gin"
	"sea_log/slaver/v1/log"
)

func Router() *gin.Engine {
	var app = gin.New()
	app.Use(gin.Recovery())
	log.Mapping("/log", app)
	return app
}
