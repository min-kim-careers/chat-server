package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	db         *DB
	cache      *Cache
	rooms      map[string]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub(db *DB, cache *Cache) *Hub {
	return &Hub{
		db:         db,
		cache:      cache,
		rooms:      make(map[string]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}

func (hub *Hub) HandleWsConnection(wsConn *websocket.Conn, connMsg *Message) {
	room, roomExists := hub.rooms[connMsg.RoomID]
	if !roomExists {
		newRoom := NewRoom(connMsg.RoomID, hub.cache, hub.db)
		log.Printf("Room <%s> created.", connMsg.RoomID)
		hub.register <- newRoom
		room = newRoom
		room.Run(hub)
	} else {
		log.Printf("Room <%s> found in hub.", connMsg.RoomID)
	}

	room.AddClient(wsConn, connMsg.ClientID)
}

func (hub *Hub) HandleRegistrations() {
	for {
		select {

		case room := <-hub.register:
			hub.rooms[room.id] = room
			log.Printf("Room <%s> registered to hub.", room.id)

		case room := <-hub.unregister:
			if _, exists := hub.rooms[room.id]; exists {
				delete(hub.rooms, room.id)
				log.Printf("Room <%s> unregistered from hub.", room.id)
			}
		}
	}
}

func (hub *Hub) Run() {
	go hub.HandleRegistrations()
}
