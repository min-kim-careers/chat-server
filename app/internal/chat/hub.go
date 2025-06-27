package chat

import (
	"log"

	"chat-server/internal/cache"
	"chat-server/internal/db"
	"chat-server/internal/models"

	"github.com/gorilla/websocket"
)

type Hub struct {
	db         *db.DB
	cache      *cache.Cache
	rooms      map[string]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub(db *db.DB, cache *cache.Cache) *Hub {
	return &Hub{
		db:         db,
		cache:      cache,
		rooms:      make(map[string]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}

func (hub *Hub) HandleWsConnection(wsConn *websocket.Conn, connMsg *models.Message) {
	roomID := connMsg.RoomID
	room, roomExists := hub.rooms[roomID]
	if !roomExists {
		newRoom := NewRoom(roomID, hub.cache, hub.db)
		log.Printf("Room <%s> created.", roomID)
		hub.register <- newRoom
		room = newRoom
		room.Run(hub)
	} else {
		log.Printf("Room <%s> found in hub.", roomID)
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
