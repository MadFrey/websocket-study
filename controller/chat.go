package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"websocket/service"
	"websocket/util"
)

var UP=websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(c *gin.Context) {
	conn, err := UP.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	service.Chat(util.NewHub,conn)
}

