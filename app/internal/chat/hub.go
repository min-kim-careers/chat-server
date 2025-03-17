package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	db         *DB
	cache      *Cache
	rooms      map[RoomID]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub(db *DB, cache *Cache) *Hub {
	return &Hub{
		db:         db,
		cache:      cache,
		rooms:      make(map[RoomID]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}

func (hub *Hub) HandleWsConnection(wsConnMsg *Message, wsConn *websocket.Conn) {
	room, roomExists := hub.rooms[RoomID(wsConnMsg.RoomID)]
	if roomExists {
		log.Printf("Room <%s> found in hub.", wsConnMsg.RoomID)
	} else {
		newRoom := NewRoom(wsConnMsg.RoomID)
		log.Printf("Room <%s> created.", wsConnMsg.RoomID)
		hub.register <- newRoom
		room = newRoom

		room.Run(hub)
	}

	room.AddClient(wsConnMsg.ClientID, wsConn, hub.cache)
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
