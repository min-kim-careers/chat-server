package chat

import (
	"chat-server/internal/dto"

	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func NewClient(conn *websocket.Conn, id string, hub *Hub) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		hub:     hub,
		id:      id,
		ctx:     ctx,
		cancel:  cancel,
		conn:    conn,
		channel: make(chan []byte),
	}
	c.Run()
	return c
}

func (c *Client) hasRoom() bool {
	return c.room != nil
}

// send to client/browser
func (c *Client) send() {
	for p := range c.channel {
		m, err := dto.ToMessageOut(p, c.id)
		if err != nil {
			log.Println(err)
			continue
		}

		_p, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
			continue
		}

		err = c.conn.WriteMessage(websocket.TextMessage, _p)
		if err != nil {
			log.Printf("Error sending message to client <%s>: %v", c.id, err)
			continue
		}
		log.Printf("Sent message to client <%s>: %v", c.id, string(p))
	}

}

// read from client/browser
func (c *Client) read() {
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading payload from client <%s>: %v", c.id, err)
			return
		}

		m, err := dto.ToMessage(p)
		if err != nil {
			log.Printf("Error parsing payload from client <%s>: <%v>", c.id, err)
			continue
		}

		if c.hasRoom() {
			switch m.Mode {
			case "chat":
				c.handleChatMessage(m, c.room.id)
			case "restore":
				c.handleRestoreMessage(m, c.room.id)
			case "leave":
				c.handleLeaveMessage()
			case "disconnect":
				c.handleDisconnectMessage()
			}
			continue
		}

		switch m.Mode {
		case "join":
			c.handleJoinMessage(m)
		case "disconnect":
			c.handleDisconnectMessage()
		}

	}

}

func (c *Client) Run() {
	go c.read()
	go c.send()
}
