package router

import (
	"github.com/gin-gonic/gin"
	"websocket/controller"
)

func InitRouter() *gin.Engine {
	r:=gin.Default()
	r.GET("/ws",controller.Chat)
	return r
}
