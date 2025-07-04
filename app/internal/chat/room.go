package chat

import (
	"chat-server/internal/dto"
	"context"
	"log"
)

func NewRoom(hub *Hub, id string) *Room {
	r := &Room{
		hub:              hub,
		id:               id,
		clients:          make(map[string]*Client),
		clientRegister:   make(chan *Client),
		clientUnregister: make(chan *Client),
	}
	r.Run()
	return r
}

func (r *Room) registerClient(c *Client) {
	r.clients[c.id] = c
	c.room = r
	log.Printf("Client <%s> joined room <%s>.", c.id, r.id)
	p, err := dto.NewMessagePayload(&dto.MessageOut{
		Mode: "joined",
	})
	if err != nil {
		return
	}
	c.channel <- p
}

func (r *Room) unregisterClient(c *Client) {
	if _, exists := r.clients[c.id]; exists {
		delete(r.clients, c.id)
		c.room = nil
		log.Printf("Client <%s> left room <%s>.", c.id, r.id)
		p, err := dto.NewMessagePayload(&dto.MessageOut{
			Mode: "left",
		})
		if err != nil {
			return
		}
		c.channel <- p
	}

	if len(r.clients) == 0 {
		log.Printf("Room <%s> is empty. Requesting removal from hub.", r.id)
		r.hub.roomUnregister <- r
		return
	}
}

func (r *Room) HandleRegistrations() {
	for {
		select {
		case c := <-r.clientRegister:
			r.registerClient(c)
		case c := <-r.clientUnregister:
			r.unregisterClient(c)
		}
	}
}

func (r *Room) HandleClients() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pubsub := r.hub.svc.Room.GetRoomChannel(ctx, r.id)
	if pubsub == nil {
		log.Printf("Room <%s> failed to subscribe to a channel. Disconnecting.", r.id)
		r.hub.roomUnregister <- r
		return
	}
	defer pubsub.Close()

	for data := range pubsub.Channel() {
		p := []byte(data.Payload)
		for _, c := range r.clients {
			c.channel <- p
		}
	}

}

func (r *Room) Run() {
	go r.HandleRegistrations()
	go r.HandleClients()
}
