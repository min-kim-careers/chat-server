package chat

import (
	"log"

	"chat-server/internal/service"
)

func NewHub(svc *service.Services) *Hub {
	return &Hub{
		Svc:              svc,
		Rooms:            make(map[string]*Room),
		RoomRegister:     make(chan *Room),
		RoomUnregister:   make(chan *Room),
		Clients:          make(map[string]*Client),
		ClientRegister:   make(chan *Client),
		ClientUnregister: make(chan *Client),
	}
}

func (h *Hub) addRoom(roomID string) *Room {
	room, exists := h.Rooms[roomID]
	if !exists {
		newRoom := NewRoom(h, roomID)
		log.Printf("Room <%s> created.", roomID)
		h.RoomRegister <- newRoom
		newRoom.Run()
		room = newRoom
	} else {
		log.Printf("Room <%s> found in hub.", roomID)
	}
	return room
}

func (h *Hub) HandleConnection(roomID string, newClient *Client) {
	room := h.addRoom(roomID)

	h.ClientRegister <- newClient
	room.AddClient(newClient)
}

func (h *Hub) HandleRoomRegistrations() {
	for {
		select {

		case room := <-h.RoomRegister:
			h.Rooms[room.ID] = room
			log.Printf("Room <%s> registered to hub.", room.ID)

		case room := <-h.RoomUnregister:
			if _, exists := h.Rooms[room.ID]; exists {
				delete(h.Rooms, room.ID)
				log.Printf("Room <%s> unregistered from hub.", room.ID)
			}
		}
	}
}

func (h *Hub) HandleClientRegistrations() {
	for {
		select {

		case client := <-h.ClientRegister:
			h.Clients[client.id] = client
			log.Printf("Client <%s> registered to hub.", client.id)

		case client := <-h.ClientUnregister:
			if _, exists := h.Clients[client.id]; exists {
				client.cancel()
				delete(h.Clients, client.id)
				log.Printf("Client <%s> unregistered from hub.", client.id)
			}
		}
	}
}

func (h *Hub) Run() {
	go h.HandleRoomRegistrations()
	go h.HandleClientRegistrations()
}
