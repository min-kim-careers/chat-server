package chat

import (
	"chat-server/internal/dto"
	"context"
	"log"

	"github.com/google/uuid"
)

func NewRoom(hub *Hub, id uuid.UUID) *Room {
	return &Room{
		Hub:        hub,
		ID:         id,
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Room) AddClient(client *Client) {
	r.Register <- client
	client.Run(r)
}

func (r *Room) HandleRegistrations() {
	for {
		select {
		case c := <-r.Register:
			r.Clients[c.id] = c
			log.Printf("Client <%s> joined room <%s>.", c.id, r.ID)

		case c := <-r.Unregister:
			if _, exists := r.Clients[c.id]; exists {
				delete(r.Clients, c.id)
				log.Printf("Client <%s> left room <%s>.", c.id, r.ID)
			}

			if len(r.Clients) == 0 {
				log.Printf("Room <%s> is empty. Requesting removal from hub.", r.ID)
				r.Hub.RoomUnregister <- r
				return
			}
		}
	}
}

func (r *Room) HandleClients() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pubsub := r.Hub.Deps.Cache.PubSub(ctx, r.ID.String())
	if pubsub == nil {
		log.Printf("Room <%s> failed to subscribe to a channel. Disconnecting.", r.ID)
		r.Hub.RoomUnregister <- r
		return
	}
	defer pubsub.Close()

	channel := pubsub.Channel()

	for data := range channel {
		msgJson := []byte(data.Payload)

		msg := dto.DeserializeMessage(msgJson)
		if msg == nil {
			log.Printf("Error deserializing message in room <%s>. Message: %s", r.ID, msgJson)
			continue
		}

		for clientID, client := range r.Clients {
			if clientID == msg.ClientID {
				continue
			}

			client.Send(msgJson)
		}
	}

}

func (r *Room) Run() {
	go r.HandleRegistrations()
	go r.HandleClients()
}
