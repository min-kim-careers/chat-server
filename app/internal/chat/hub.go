package chat

import (
	"log"

	"chat-server/internal/dto"
	"chat-server/internal/service"
)

type Hub struct {
	svc              *service.Services
	rooms            map[string]*Room
	roomRegister     chan *Room
	roomUnregister   chan *Room
	clients          map[string]*Client
	clientRegister   chan *Client
	clientUnregister chan *Client
}

func NewHub(svc *service.Services) *Hub {
	return &Hub{
		svc:              svc,
		rooms:            make(map[string]*Room),
		roomRegister:     make(chan *Room),
		roomUnregister:   make(chan *Room),
		clients:          make(map[string]*Client),
		clientRegister:   make(chan *Client),
		clientUnregister: make(chan *Client),
	}
}

func (h *Hub) HandleNewClient(c *Client) {
	h.clientRegister <- c
}

func (h *Hub) registerRoom(r *Room) {
	h.rooms[r.id] = r
	log.Printf("Room <%s> registered to hub.", r.id)
}

func (h *Hub) unregisterRoom(r *Room) {
	if _, exists := h.rooms[r.id]; exists {
		delete(h.rooms, r.id)
		log.Printf("Room <%s> unregistered from hub.", r.id)
	}
}

func (h *Hub) registerClient(c *Client) {
	h.clients[c.id] = c
	log.Printf("Client <%s> registered to hub.", c.id)
	p, err := dto.ToRawMessageOut(&dto.MessageOut{
		Mode: "connected",
	})
	if err != nil {
		return
	}
	c.channel <- p
}

func (h *Hub) unregisterClient(c *Client) {
	if _, exists := h.clients[c.id]; exists {
		c.cancel()
		delete(h.clients, c.id)
		log.Printf("Client <%s> unregistered from hub.", c.id)
	}
}

func (h *Hub) HandleRoomRegistrations() {
	for {
		select {
		case r := <-h.roomRegister:
			h.registerRoom(r)
		case r := <-h.roomUnregister:
			h.unregisterRoom(r)
		}
	}
}

func (h *Hub) HandleClientRegistrations() {
	for {
		select {
		case c := <-h.clientRegister:
			h.registerClient(c)
		case c := <-h.clientUnregister:
			h.unregisterClient(c)
		}
	}
}

func (h *Hub) Run() {
	go h.HandleRoomRegistrations()
	go h.HandleClientRegistrations()
}
