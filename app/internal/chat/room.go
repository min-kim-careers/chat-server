package chat

import (
	"chat-server/internal/cache"
	"chat-server/internal/db"
	"chat-server/internal/models"
	"log"

	"github.com/gorilla/websocket"
)

type Room struct {
	id         string
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	cache      *cache.Cache
	db         *db.DB
}

func NewRoom(id string, cache *cache.Cache, db *db.DB) *Room {
	return &Room{
		id:         id,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		cache:      cache,
		db:         db,
	}
}

func (room *Room) AddClient(wsConn *websocket.Conn, clientID string) {
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
	pubsub := hub.cache.PubSub(room.id)
	if pubsub == nil {
		log.Printf("Room <%s> failed to subscribe to a channel. Disconnecting.", room.id)
		hub.unregister <- room
		return
	}
	defer pubsub.Close()

	channel := pubsub.Channel()

	for data := range channel {
		msgJson := []byte(data.Payload)

		msg := models.DeserializeMessage(msgJson)
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
