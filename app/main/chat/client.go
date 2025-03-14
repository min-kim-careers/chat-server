package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientID string

type Client struct {
	id      ClientID
	conn    *websocket.Conn
	channel chan *Message
}

func NewClient(id ClientID, wsConn *websocket.Conn) *Client {
	c := Client{}

	c.id = id
	c.conn = wsConn
	c.channel = make(chan *Message)

	return &c
}

func (c *Client) read(r *Room) {
	for {
		_, msgJson, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %s", c.id, err)
			r.unregister <- c
			c.conn.Close()
			return
		}
		msg := DeserializeMessage(msgJson)

		if msg.Type == "chat" {
			log.Printf("Message read from client <%s>: %s", c.id, msg.Content)
			r.broadcast <- msg
		}
	}
}

func (c *Client) send(r *Room) {
	for msg := range c.channel {
		msgJson := SerializeMessage(msg)
		err := c.conn.WriteMessage(websocket.TextMessage, msgJson)
		if err != nil {
			log.Printf("Error sending message by client <%s>: %s", c.id, err)
			r.unregister <- c
			c.conn.Close()
			return
		}
		log.Printf("Message sent to client <%s>: %s", c.id, msg.Content)
	}
}

func (c *Client) Run(r *Room) {
	go c.read(r)
	go c.send(r)
}
