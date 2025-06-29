package chat

import (
	"chat-server/internal/dto"
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type Room struct {
	id         string
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	hub        *Hub
}

func NewRoom(id string, hub *Hub) *Room {
	return &Room{
		id:         id,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		hub:        hub,
	}
}

func (r *Room) AddClient(ctx context.Context, conn *websocket.Conn, clientID string) {
	newClient := NewClient(ctx, conn, clientID, r)
	r.register <- newClient
	newClient.Run(r)
}

func (r *Room) HandleRegistrations(hub *Hub) {
	for {
		select {
		case c := <-r.register:
			r.clients[c.id] = c
			log.Printf("Client <%s> connected to room <%s>.", c.id, r.id)

		case c := <-r.unregister:
			if _, exists := r.clients[c.id]; exists {
				delete(r.clients, c.id)
				log.Printf("Client <%s> disconnected from room <%s>.", c.id, r.id)
			}

			if len(r.clients) == 0 {
				log.Printf("Room <%s> is empty. Requesting removal from hub.", r.id)
				hub.unregister <- r
				return
			}
		}
	}
}

func (r *Room) HandleClients(ctx context.Context, hub *Hub) {
	pubsub := hub.deps.Cache.PubSub(ctx, r.id)
	if pubsub == nil {
		log.Printf("Room <%s> failed to subscribe to a channel. Disconnecting.", r.id)
		hub.unregister <- r
		return
	}
	defer pubsub.Close()

	channel := pubsub.Channel()

	for data := range channel {
		msgJson := []byte(data.Payload)

		msg := dto.DeserializeMessage(msgJson)
		if msg == nil {
			log.Printf("Error deserializing message in room <%s>. Message: %s", r.id, msgJson)
			continue
		}

		for clientID, client := range r.clients {
			if clientID == msg.ClientID {
				continue
			}

			client.Send(msgJson)
		}
	}
}

func (r *Room) Run(ctx context.Context, hub *Hub) {
	go r.HandleRegistrations(hub)
	go r.HandleClients(ctx, hub)
}
