package chat

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RestoreCursor struct {
	LastCacheID string
	LastDBID    time.Time
	NoMessages  bool
}

func DefaultRestoreCursor() RestoreCursor {
	return RestoreCursor{
		LastCacheID: "0",
		LastDBID:    time.Now(),
		NoMessages:  false,
	}
}

type Client struct {
	hub          *Hub
	room         *Room
	id           string
	ctx          context.Context
	ctxCancel    context.CancelFunc
	wg           sync.WaitGroup
	conn         *websocket.Conn
	outbound     chan []byte
	inbound      chan []byte
	cursor       RestoreCursor
	lastActivity time.Time
	mu           sync.Mutex
}

func NewClient(conn *websocket.Conn, id string, hub *Hub) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		hub:          hub,
		id:           id,
		ctx:          ctx,
		ctxCancel:    cancel,
		wg:           sync.WaitGroup{},
		conn:         conn,
		cursor:       DefaultRestoreCursor(),
		lastActivity: time.Now(),
		outbound:     make(chan []byte),
		inbound:      make(chan []byte),
	}
	c.run()
	return c
}

func (c *Client) touch() {
	c.lastActivity = time.Now()
}

func (c *Client) resetCursor() {
	c.cursor = DefaultRestoreCursor()
}

func (c *Client) hasNoMessages() bool {
	return c.cursor.NoMessages
}

func (c *Client) hasRoom() bool {
	return c.room != nil
}

func (c *Client) setRoom(r *Room) {
	c.room = r
}

func (c *Client) run() {
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()
		c.handleInbound()
	}()

	go func() {
		defer wg.Done()
		c.handleOutbound()
	}()

	go func() {
		defer wg.Done()
		c.handleMessages()
	}()

	go func() {
		defer wg.Done()
		c.handleIdleTimeout()
	}()

	go func() {
		defer wg.Done()
		c.handleClose()
	}()

}
