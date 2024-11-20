package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	Socket         *websocket.Conn
	Receive        chan *domain.Message
	ChatController *ChatHandler
}

func (c *Client) Read(userID string) {
	defer c.Socket.Close()
	for {
		msg := &domain.Message{}
		_, jsonMessage, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		err = json.Unmarshal(jsonMessage, msg)
		if err != nil {
			return
		}
		msg.SenderID = userID
		c.ChatController.Messages <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Receive {
		msg.CreatedAt = time.Now()
		jsonForSend, err := json.Marshal(msg)
		if err != nil {
			return
		}
		err = c.Socket.WriteMessage(websocket.TextMessage, jsonForSend)
		if err != nil {
			return
		}
	}
}
