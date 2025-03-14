package chat

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type Hub struct {
	dbConn     *pgx.Conn
	rooms      map[RoomID]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub(dbConn *pgx.Conn) *Hub {
	h := Hub{}

	h.dbConn = dbConn
	h.rooms = make(map[RoomID]*Room)
	h.register = make(chan *Room)
	h.unregister = make(chan *Room)

	return &h
}

func (h *Hub) HandleWsConnection(wsConnMsg *Message, wsConn *websocket.Conn) {
	room, roomExists := h.rooms[RoomID(wsConnMsg.RoomID)]
	if roomExists {
		log.Printf("Room <%s> found in hub.", wsConnMsg.RoomID)
	} else {
		newRoom := NewRoom(wsConnMsg.RoomID, h.dbConn)
		log.Printf("Room <%s> created.", wsConnMsg.RoomID)
		h.register <- newRoom
		room = newRoom

		newRoom.Run(h)
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

func (h *Hub) handleRegistrations() {
	for {
		select {

		case room := <-h.register:
			h.rooms[room.id] = room
			log.Printf("Room <%s> registered to hub.", room.id)

		case room := <-h.unregister:
			if _, exists := h.rooms[room.id]; exists {
				delete(h.rooms, room.id)
				log.Printf("Room <%s> unregistered from hub.", room.id)
			}

		}
	}
}

func (h *Hub) Run() {
	go h.handleRegistrations()
}
