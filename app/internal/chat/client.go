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

func (client *Client) read(r *Room) {
	for {
		_, msgJson, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %s", client.id, err)
			r.unregister <- client
			client.conn.Close()
			return
		}
		msg := DeserializeMessage(msgJson)

		if msg.Type == "chat" {
			log.Printf("Message read from client <%s>: %s", client.id, msg.Content)
			r.broadcast <- msg
		}
	}
}

func (client *Client) send(r *Room) {
	for msg := range client.channel {
		msgJson := SerializeMessage(msg)
		err := client.conn.WriteMessage(websocket.TextMessage, msgJson)
		if err != nil {
			log.Printf("Error sending message by client <%s>: %s", client.id, err)
			r.unregister <- client
			client.conn.Close()
			return
		}
		log.Printf("Message sent to client <%s>: %s", client.id, msg.Content)
	}
}

func (client *Client) Run(r *Room) {
	go client.read(r)
	go client.send(r)
}
