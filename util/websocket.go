package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"websocket/model"
)

type Client struct {
	Id       int
	UserName string
	Message  chan []byte
	Conn     *websocket.Conn
}

type Hub struct {
	BroadcastChan chan []byte
	Register      chan *Client
	Unregister    chan *Client
	Clients       map[*Client]bool
}

func CreateHub() *Hub {
	return &Hub{
		BroadcastChan: make(chan []byte),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Clients:       make(map[*Client]bool),
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(55 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Println(err)
					return
				}
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println(err)
				return
			}
			//心跳
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Read(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for true {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		bytes1, err := json.Marshal(model.Message{Id: c.Id, Content: string(message)})
		if err != nil {
			log.Println(err)
			return
		}
		bytes1 = bytes.TrimSpace(bytes.Replace(bytes1, []byte("\n"), []byte(" "), -1))
		hub.Send(bytes1, c)
	}
}

func (hub *Hub) Send(message []byte, c *Client) {
	fmt.Println(hub.Clients)
	for client := range hub.Clients {
		if client != c {
			client.Message <- message
		}
	}
}

func (hub *Hub) Start() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true

		case client := <-hub.Unregister:
			if _, ok := hub.Clients[client]; ok {
				delete(hub.Clients, client)
				close(client.Message)
			}

		case message := <-hub.BroadcastChan:
			for client := range hub.Clients {
				select {
				case client.Message <- message:

				default:
					close(client.Message)
					delete(hub.Clients, client)
				}
			}
		}
	}
}

var NewHub =CreateHub()