package chat

import (
	"log"
	"sync"

	"chat-server/internal/dto/messageout"
	"chat-server/internal/service"
)

type Hub struct {
	svc               *service.Services
	rooms             map[string]*Room
	roomRegister      chan *Room
	roomUnregister    chan *Room
	clients           map[string]*Client
	clientRegister    chan *Client
	clientUnregister  chan *Client
	persistWorkerPool *HubPersistWorkerPool
	mu                sync.Mutex
}

func NewHub(svc *service.Services) *Hub {
	return &Hub{
		svc:               svc,
		rooms:             make(map[string]*Room),
		roomRegister:      make(chan *Room),
		roomUnregister:    make(chan *Room),
		clients:           make(map[string]*Client),
		clientRegister:    make(chan *Client),
		clientUnregister:  make(chan *Client),
		persistWorkerPool: NewHubPersistWorkerPool(svc),
	}
}

func (h *Hub) HandleNewClient(c *Client) {
	h.clientRegister <- c
}

func (h *Hub) getRoom(roomID string) (*Room, bool) {
	room, exists := h.rooms[roomID]
	return room, exists
}

func (h *Hub) addRoom(roomID string) {

}

func (h *Hub) setClient(c *Client) {
	h.clients[c.id] = c
}

func (h *Hub) deleteClient(c *Client) {
	delete(h.clients, c.id)
}

func (h *Hub) setRoom(r *Room) {
	h.rooms[r.id] = r
}

func (h *Hub) deleteRoom(r *Room) {
	delete(h.rooms, r.id)
}

func (h *Hub) registerRoom(r *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.setRoom(r)
	log.Printf("room <%s> registered to hub.", r.id)
}

func (h *Hub) unregisterRoom(r *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.rooms[r.id]; exists {
		h.deleteRoom(r)
		log.Printf("room <%s> unregistered from hub.", r.id)
	}
}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.setClient(c)
	log.Printf("client <%s> registered to hub.", c.id)
	p, err := messageout.ToRawMessageOut(&messageout.MessageOutEvent{
		Mode: "connected",
	})
	if err != nil {
		log.Println("error:", err)
		return
	}
	c.outbound <- p
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[c.id]; exists {
		h.deleteClient(c)
		log.Printf("client <%s> unregistered from hub.", c.id)
	}
}

func (h *Hub) Run() {
	go h.handleRoomRegistrations()
	go h.handleClientRegistrations()
}
