package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientID string

type Client struct {
	id      ClientID
	conn    *websocket.Conn
	channel chan []byte
}

func NewClient(id ClientID, conn *websocket.Conn) *Client {
	c := Client{}

	c.id = id
	c.conn = conn
	c.channel = make(chan []byte)

	return &c
}

func (c *Client) Read(r *Room) {
	for {
		_, msgJson, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %s", c.id, err)
			r.unregister <- c
			c.conn.Close()
			return
		}
		msg := Deserialize(msgJson)

		if msg.MessageType == "chat" {
			log.Printf("Message read from client <%s>: %s", c.id, msg.MessageContent)
			r.broadcast <- msgJson
		}
	}
}

func (c *Client) Send(r *Room) {
	for msgJson := range c.channel {
		err := c.conn.WriteMessage(websocket.TextMessage, msgJson)
		if err != nil {
			log.Printf("Error sending message by client <%s>: %s", c.id, err)
			r.unregister <- c
			c.conn.Close()
			return
		}
		msg := Deserialize(msgJson)
		log.Printf("Message sent to client <%s>: %s", c.id, msg.MessageContent)
	}
}

func (c *Client) Run(r *Room) {
	go c.Read(r)
	go c.Send(r)
}
