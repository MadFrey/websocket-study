package service

import (
	"github.com/gorilla/websocket"
	"websocket/util"
)

func Chat(hub *util.Hub,conn *websocket.Conn) {
	client := &util.Client{Message: make(chan []byte)}
	client.Conn = conn
	hub.Register <- client
	go client.Read(hub)
	go client.Write()
}