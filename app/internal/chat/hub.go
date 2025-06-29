package chat

import (
	"context"
	"log"

	"chat-server/internal/deps"
	"chat-server/internal/dto"

	"github.com/gorilla/websocket"
)

type Hub struct {
	deps       *deps.Container
	rooms      map[string]*Room
	register   chan *Room
	unregister chan *Room
}

func NewHub(deps *deps.Container) *Hub {
	return &Hub{
		deps:       deps,
		rooms:      make(map[string]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}

func (h *Hub) HandleConnection(ctx context.Context, conn *websocket.Conn, message *dto.Message) {
	roomID := message.RoomID
	room, roomExists := h.rooms[roomID]
	if !roomExists {
		newRoom := NewRoom(roomID, h)
		log.Printf("Room <%s> created.", roomID)
		h.register <- newRoom
		room = newRoom
		room.Run(ctx, h)
	} else {
		log.Printf("Room <%s> found in hub.", roomID)
	}

	room.AddClient(ctx, conn, message.ClientID)
}

func (h *Hub) HandleRegistrations() {
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
	go h.HandleRegistrations()
}
