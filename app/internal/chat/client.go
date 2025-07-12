package chat

import (
	"chat-server/internal/dto"

	"context"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub       *Hub
	room      *Room
	id        string
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *websocket.Conn
	channel   chan []byte
}

func NewClient(conn *websocket.Conn, id string, hub *Hub) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		hub:       hub,
		id:        id,
		ctx:       ctx,
		ctxCancel: cancel,
		conn:      conn,
		channel:   make(chan []byte),
	}
	c.run()
	return c
}

func (c *Client) hasRoom() bool {
	return c.room != nil
}

// send to client/browser
func (c *Client) send() {
	for p := range c.channel {
		err := c.conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			log.Printf("Error sending message to client <%s>: %v", c.id, err)
			c.ctxCancel()
			return
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
			c.ctxCancel()
			return
		}

		m, err := dto.ToMessageIn(p)
		if err != nil {
			log.Printf("Error parsing payload from client <%s>: <%v>", c.id, err)
			continue
		}

		switch m.Mode {
		case "chat":
			c.handleChat(m)
		case "restore":
			c.handleRestore(m)
		case "leave":
			c.handleLeave()
		case "join":
			c.handleJoin(m)
		case "typing":
			c.handleTyping(m)
		case "not_typing":
			c.handleNotTyping(m)
		}
	}
}

func (c *Client) handleClose() {
	<-c.ctx.Done()
	c.handleDisconnect()
	close(c.channel)
	c.channel = nil
}

func (c *Client) run() {
	go c.read()
	go c.send()
	go c.handleClose()
}
