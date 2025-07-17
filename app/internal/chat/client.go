package chat

import (
	"chat-server/internal/dto/messagein"
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RestoreCursor struct {
	LastCacheID string
	LastDBID    time.Time
}

func DefaultRestoreCursor() RestoreCursor {
	return RestoreCursor{
		LastCacheID: "0",
		LastDBID:    time.Time{},
	}
}

type Client struct {
	hub       *Hub
	room      *Room
	id        string
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *websocket.Conn
	channel   chan []byte
	cursor    RestoreCursor
	mu        sync.Mutex
}

func NewClient(conn *websocket.Conn, id string, hub *Hub) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		hub:       hub,
		id:        id,
		ctx:       ctx,
		ctxCancel: cancel,
		conn:      conn,
		cursor:    DefaultRestoreCursor(),
		channel:   make(chan []byte),
	}
	c.run()
	return c
}

func (c *Client) hasRoom() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.room == nil {
		log.Printf("room not found")
		return false
	}
	return true
}

func (c *Client) setRoom(r *Room) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.room = r
}

// send to client/browser
func (c *Client) send() {
	for p := range c.channel {
		err := c.conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			log.Printf("error sending message to client <%s>: %v", c.id, err)
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
			log.Printf("error reading message: %v", err)
			c.ctxCancel()
			return
		}

		m, err := messagein.ToMessageIn(p)
		if err != nil {
			log.Printf("error parsing message in: %v", err)
			continue
		}

		switch v := m.(type) {
		case *messagein.MessageInChat:
			c.handleChat(v)
		case *messagein.MessageInJoin:
			c.handleJoin(v)

		case *messagein.MessageInEvent:
			switch v.Mode {
			case "restore":
				c.handleRestore()
			case "leave":
				c.handleLeave()
			case "typing":
				c.handleTyping()
			case "not_typing":
				c.handleNotTyping()
			}
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
