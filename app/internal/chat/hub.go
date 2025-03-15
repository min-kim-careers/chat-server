package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	db         *DB
	rooms      map[RoomID]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub(db *DB) *Hub {
	h := Hub{}
	h.db = db
	h.rooms = make(map[RoomID]*Room)
	h.register = make(chan *Room)
	h.unregister = make(chan *Room)
	return &h
}

func (hub *Hub) HandleWsConnection(wsConnMsg *Message, wsConn *websocket.Conn) {
	room, roomExists := hub.rooms[RoomID(wsConnMsg.RoomID)]
	if roomExists {
		log.Printf("Room <%s> found in hub.", wsConnMsg.RoomID)
	} else {
		newRoom := NewRoom(wsConnMsg.RoomID, hub)
		log.Printf("Room <%s> created.", wsConnMsg.RoomID)
		hub.register <- newRoom
		room = newRoom

		newRoom.Run(hub)
	}

	_, clientExists := room.clients[wsConnMsg.ClientID]
	if clientExists {
		log.Printf("Client <%s> already in room. Disconnecting.", wsConnMsg.ClientID)
		wsConn.Close()
		return
	}
	newClient := NewClient(wsConnMsg.ClientID, wsConn)
	log.Printf("Client <%s> created.", wsConnMsg.ClientID)
	room.register <- newClient

	newClient.Run(room)
}

func (hub *Hub) handleRegistrations() {
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
	go hub.handleRegistrations()
}
