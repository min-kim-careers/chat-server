package chat

import (
	"log"
	"maps"
	"sync"
	"time"

	"chat-server/internal/constant"
	"chat-server/internal/dto/messageout"
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
	mu               sync.Mutex
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

func (h *Hub) getRoom(roomID string) (*Room, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[roomID]
	return room, exists
}

func (h *Hub) registerRoom(r *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.rooms[r.id] = r
	log.Printf("Room <%s> registered to hub.", r.id)
}

func (h *Hub) unregisterRoom(r *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.rooms[r.id]; exists {
		delete(h.rooms, r.id)
		log.Printf("Room <%s> unregistered from hub.", r.id)
	}
}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c.id] = c
	log.Printf("Client <%s> registered to hub.", c.id)
	p, err := messageout.ToRawMessageOut(&messageout.MessageOutEvent{
		Mode: "connected",
	})
	if err != nil {
		log.Printf("error parsing connected message: %v", err)
		return
	}
	c.channel <- p
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[c.id]; exists {
		delete(h.clients, c.id)
		log.Printf("Client <%s> unregistered from hub.", c.id)
	}
}

func (h *Hub) handleRoomRegistrations() {
	for {
		select {
		case r := <-h.roomRegister:
			h.registerRoom(r)
		case r := <-h.roomUnregister:
			h.unregisterRoom(r)
		}
	}
}

func (h *Hub) handleClientRegistrations() {
	for {
		select {
		case c := <-h.clientRegister:
			h.registerClient(c)
		case c := <-h.clientUnregister:
			h.unregisterClient(c)
		}
	}
}

func (h *Hub) getRoomsSnapshot() map[string]*Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	snap := make(map[string]*Room, len(h.rooms))
	maps.Copy(snap, h.rooms)
	return snap
}

func (h *Hub) handleFlush() {
	interval := 5 * time.Second
	delay := 5 * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for _, r := range h.getRoomsSnapshot() {
			h.mu.Lock()
			age := now.Sub(r.lastActivity)
			h.mu.Unlock()

			if age < delay {
				continue
			}

			if r.getCacheSize() < constant.CACHE_LIMIT {
				continue
			}

			r.flushRoom()
		}
	}
}

func (h *Hub) Run() {
	go h.handleRoomRegistrations()
	go h.handleClientRegistrations()
	go h.handleFlush()
}
