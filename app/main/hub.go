package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms      map[RoomID]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub() *Hub {
	h := Hub{}

	h.rooms = make(map[RoomID]*Room)
	h.register = make(chan *Room)
	h.unregister = make(chan *Room)

	return &h
}

func (h *Hub) HandleConnection(connMsg *Message, conn *websocket.Conn) {
	room, roomExists := h.rooms[RoomID(connMsg.RoomID)]
	if roomExists {
		log.Printf("Room <%s> found in hub.", connMsg.RoomID)
	} else {
		newRoom := NewRoom(connMsg.RoomID)
		log.Printf("Room <%s> created.", connMsg.RoomID)
		h.register <- newRoom
		room = newRoom

		go newRoom.Run(h)
	}

	_, clientExists := room.clients[connMsg.ClientID]
	if clientExists {
		log.Printf("Client <%s> already in room. Disconnecting.", connMsg.ClientID)
		conn.Close()
		return
	}
	newClient := NewClient(connMsg.ClientID, conn)
	log.Printf("Client <%s> created.", connMsg.ClientID)
	room.register <- newClient

	go newClient.Run(room)

}

func (h *Hub) Run() {
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
