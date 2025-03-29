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
	cache      *Cache
	db         *DB
}

func NewRoom(id RoomID, cache *Cache, db *DB) *Room {
	return &Room{
		id:         id,
		clients:    make(map[ClientID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		cache:      cache,
		db:         db,
	}
}

func (room *Room) AddClient(wsConn *websocket.Conn, clientID ClientID) {
	newClient := NewClient(clientID, wsConn)
	room.register <- newClient
	newClient.Run(room)
}

func (room *Room) HandleRegistrations(hub *Hub) {
	for {
		select {
		case c := <-room.register:
			room.clients[c.id] = c
			log.Printf("Client <%s> connected to room <%s>.", c.id, room.id)

		case c := <-room.unregister:
			if _, exists := room.clients[c.id]; exists {
				delete(room.clients, c.id)
				log.Printf("Client <%s> disconnected from room <%s>.", c.id, room.id)
			}

			if len(room.clients) == 0 {
				log.Printf("Room <%s> is empty. Requesting removal from hub.", room.id)
				hub.unregister <- room
				return
			}
		}
	}
}

func (room *Room) HandleClients(hub *Hub) {
	pubsub := hub.cache.PubSub(string(room.id))
	if pubsub == nil {
		log.Printf("Room <%s> failed to subscribe to a channel. Disconnecting.", room.id)
		hub.unregister <- room
		return
	}
	defer pubsub.Close()

	channel := pubsub.Channel()

	for data := range channel {
		msgJson := []byte(data.Payload)

		msg := DeserializeMessage(msgJson)
		if msg == nil {
			log.Printf("Error deserializing message in room <%s>. Message: %s", room.id, msgJson)
			continue
		}

		for clientID, client := range room.clients {
			if clientID == msg.ClientID {
				continue
			}

			client.Send(msgJson)
		}
	}
}

func (room *Room) Run(hub *Hub) {
	go room.HandleRegistrations(hub)
	go room.HandleClients(hub)
}
