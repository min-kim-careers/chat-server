package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	conn *websocket.Conn
}

func NewClient(id string, conn *websocket.Conn) *Client {
	c := Client{}
	c.id = id
	c.conn = conn
	return &c
}

func (c *Client) Read(r *Room) {
	// defer func() {
	// 	log.Printf("Client <%s> disconnected.\n", c.id)
	// 	r.RemoveClient(c.id)
	// 	c.conn.Close()
	// }()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %s\n", c.id, err)
		}
		log.Printf("Message read from client <%s>: %s\n", c.id, msg)
		r.broadcast <- msg
	}
}

func (c *Client) Send(r *Room) {
	for msg := range r.broadcast {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error sending message by client <%s>: %s\n", c.id, err)
		}
	}
}
