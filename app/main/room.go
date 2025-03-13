package main

import (
	"log"
)

type RoomID string

type Room struct {
	id         RoomID
	clients    map[ClientID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewRoom(id RoomID) *Room {
	r := Room{}

	r.id = id
	r.clients = make(map[ClientID]*Client)
	r.register = make(chan *Client)
	r.unregister = make(chan *Client)
	r.broadcast = make(chan []byte)

	return &r
}

func (r *Room) Run(h *Hub) {
	for {
		select {
		case c := <-r.register:
			r.clients[c.id] = c
			log.Printf("Client <%s> registered to room <%s>.", c.id, r.id)

		case c := <-r.unregister:
			if _, exists := r.clients[c.id]; exists {
				delete(r.clients, c.id)
				log.Printf("Client <%s> unregistered from room <%s>.", c.id, r.id)
			}

			if len(r.clients) == 0 {
				log.Printf("Room <%s> is empty.", r.id)
				h.unregister <- r
			}

		case msg := <-r.broadcast:
			for _, c := range r.clients {
				c.channel <- msg
			}
		}
	}
}
