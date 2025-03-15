package chat

import (
	"log"
)

type RoomID string

type Room struct {
	id         RoomID
	clients    map[ClientID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	hub        *Hub
}

func NewRoom(id RoomID, hub *Hub) *Room {
	r := Room{}
	r.id = id
	r.clients = make(map[ClientID]*Client)
	r.register = make(chan *Client)
	r.unregister = make(chan *Client)
	r.broadcast = make(chan *Message)
	r.hub = hub
	return &r
}

func (room *Room) handleRegistrations() {
	for {
		select {
		case c := <-room.register:
			room.clients[c.id] = c
			log.Printf("Client <%s> registered to room <%s>.", c.id, room.id)

		case c := <-room.unregister:
			if _, exists := room.clients[c.id]; exists {
				delete(room.clients, c.id)
				log.Printf("Client <%s> unregistered from room <%s>.", c.id, room.id)
			}

			if len(room.clients) == 0 {
				log.Printf("Room <%s> is empty. Requesting removal from hub.", room.id)
				room.hub.unregister <- room
				return
			}
		}
	}
}

func (room *Room) handleBroadcasts() {
	for msg := range room.broadcast {
		room.hub.db.AddMessage(msg)
		for _, c := range room.clients {
			if msg.ClientID != c.id {
				select {
				case c.channel <- msg:
				default:
					log.Printf("Client <%s> message channel full, dropping message.", c.id)
				}
			}
		}
	}
}

func (room *Room) Run(h *Hub) {
	go room.handleRegistrations()
	go room.handleBroadcasts()
}
