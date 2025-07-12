package chat

import (
	"chat-server/internal/dto"
	"chat-server/internal/helper"
	"context"
	"log"

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
}

func NewRoom(hub *Hub, id string) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	pubsub := hub.svc.Room.GetRoomChannel(ctx, id)
	if pubsub == nil {
		log.Printf("Room <%s> failed to subscribe to a channel. Disconnecting.", id)
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

func (r *Room) registerClient(c *Client) {
	r.clients[c.id] = c
	c.room = r
	log.Printf("Client <%s> joined room <%s>.", c.id, r.id)
	p, err := dto.ToRawMessageOut(&dto.MessageOut{
		Mode: "joined",
	})
	if err != nil {
		log.Printf("Error parsing joined message: %v", err)
		return
	}
	c.channel <- p
}

func (r *Room) unregisterClient(c *Client) {
	if _, exists := r.clients[c.id]; exists {
		delete(r.clients, c.id)
		c.room = nil
		log.Printf("Client <%s> left room <%s>.", c.id, r.id)
		if c.channel != nil {
			p, err := dto.ToRawMessageOut(&dto.MessageOut{
				Mode: "left",
			})
			if err != nil {
				log.Printf("Error parsing left message: %v", err)
				return
			}
			c.channel <- p
		}
	}

	if len(r.clients) == 0 {
		log.Printf("Room <%s> is empty. Requesting removal from hub.", r.id)
		r.hub.roomUnregister <- r
		r.ctxCancel()
		return
	}
}

func (r *Room) handleRegistrations() {
	for {
		select {
		case c, ok := <-r.clientRegister:
			if ok {
				r.registerClient(c)
			}
		case c, ok := <-r.clientUnregister:
			if ok {
				r.unregisterClient(c)
			}
		}

		if r.clientRegister == nil && r.clientUnregister == nil {
			return
		}
	}
}

func (r *Room) handleClients() {
	defer r.pubsub.Close()

	for data := range r.pubsub.Channel() {
		p := []byte(data.Payload)
		clientID := helper.GetSingleField(p, "clientId")
		for _, c := range r.clients {
			if c.id == string(clientID) {
				continue
			}
			c.channel <- p
		}
	}

}

func (r *Room) handleClose() {
	<-r.ctx.Done()
	r.pubsub.Close()
	r.pubsub = nil
	close(r.clientRegister)
	r.clientRegister = nil
	close(r.clientUnregister)
	r.clientUnregister = nil
}

func (r *Room) run() {
	go r.handleRegistrations()
	go r.handleClients()
	go r.handleClose()
}
