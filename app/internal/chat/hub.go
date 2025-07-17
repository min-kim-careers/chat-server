package chat

import (
	"context"
	"fmt"
	"log"
	"maps"
	"sync"

	"chat-server/internal/cache"
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

func (h *Hub) handlePersist(parentCtx context.Context, workerID string) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	for _, r := range h.getRoomsSnapshot() {
		msgs, err := h.svc.Message.GetChatMessagePersistStream(ctx, workerID, r.id)
		if err != nil {
			log.Printf("error getting persist stream for room <%s>: %v", r.id, err)
			continue
		}

		cachedMsgs := make([]cache.CacheMessage, len(msgs))
		for i, m := range msgs[0].Messages {
			c, err := cache.StreamToCacheMessage(m.Values)
			if err != nil {
				log.Printf("error parsing stream to cache <%s>: %v", m.Values, err)
				break
			}
			cachedMsgs[i] = *c
		}
		if len(cachedMsgs) == len(msgs) {
			err = h.svc.Message.FlushCachedMessagesToDB(ctx, cachedMsgs)
			if err != nil {
				log.Printf("error persisting room <%s>", r.id)
				continue
			}
		}
	}
}

func (h *Hub) startPersistWorkers(numWorkers int) {
	parentCtx := context.Background()
	for i := range numWorkers {
		workerID := fmt.Sprintf("worker-%d", i)
		go h.handlePersist(parentCtx, workerID)
	}
}

func (h *Hub) Run() {
	go h.handleRoomRegistrations()
	go h.handleClientRegistrations()
	go h.startPersistWorkers(1)
}
