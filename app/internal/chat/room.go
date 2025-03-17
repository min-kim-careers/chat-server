package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type RoomID string

type Room struct {
	id         RoomID
	clients    map[ClientID]*Client
	register   chan *Client
	unregister chan *Client
}

func NewRoom(id RoomID) *Room {
	return &Room{
		id:         id,
		clients:    make(map[ClientID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (room *Room) AddClient(clientID ClientID, wsConn *websocket.Conn, cache *Cache) {
	newClient := NewClient(clientID, room.id, wsConn, cache)
	room.register <- newClient
	newClient.Run(room)
}

func (room *Room) HandleRegistrations(hub *Hub) {
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
				hub.unregister <- room
				return
			}
		}
	}
}

func (room *Room) Run(hub *Hub) {
	go room.HandleRegistrations(hub)
}
