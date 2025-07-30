package chat

import (
	"chat-server/internal/dto/messageout"
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Room struct {
	hub              *Hub
	id               string
	clients          map[string]*Client
	clientRegister   chan *Client
	clientUnregister chan *Client
	ctx              context.Context
	ctxCancel        context.CancelFunc
	pubsub           *redis.PubSub
	mu               sync.Mutex
}

func NewRoom(hub *Hub, id string) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	pubsub := hub.svc.Room.GetRoomChannel(ctx, id)
	if pubsub == nil {
		log.Printf("room <%s> failed to get pubsub. closing..", id)
		cancel()
		return nil
	}
	r := &Room{
		hub:              hub,
		id:               id,
		clients:          make(map[string]*Client),
		clientRegister:   make(chan *Client),
		clientUnregister: make(chan *Client),
		ctx:              ctx,
		ctxCancel:        cancel,
		pubsub:           pubsub,
	}
	r.run()
	return r
}

func (r *Room) addClient(c *Client) {
	r.clients[c.id] = c
	c.setRoom(r)
	log.Printf("client <%s> joined room <%s>.", c.id, r.id)
}

func (r *Room) deleteClient(c *Client) {
	delete(r.clients, c.id)
	c.setRoom(nil)
	log.Printf("client <%s> left room <%s>.", c.id, r.id)
}

func (r *Room) registerClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.addClient(c)

	if c.outbound != nil {
		p, err := messageout.ToRawMessageOut(&messageout.MessageOutEvent{
			Mode: "joined",
		})
		if err != nil {
			log.Println("error:", err)
			return
		}
		c.outbound <- p
	}
}

func (r *Room) unregisterClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[c.id]; exists {
		r.deleteClient(c)

		if len(r.clients) == 0 {
			log.Printf("Room <%s> is empty. Requesting removal from hub.", r.id)
			r.hub.roomUnregister <- r
			r.ctxCancel()
			return
		}

		if c.outbound != nil {
			p, err := messageout.ToRawMessageOut(&messageout.MessageOutEvent{
				Mode: "left",
			})
			if err != nil {
				log.Println("error:", err)
				return
			}
			c.outbound <- p
		}
	}

}

func (r *Room) run() {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		r.handleClientRegistrations()
	}()

	go func() {
		defer wg.Done()
		r.handleMessages()
	}()

	go func() {
		defer wg.Done()
		r.handleClose()
	}()
}
