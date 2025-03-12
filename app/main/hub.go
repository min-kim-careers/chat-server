package main

import "log"

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	rooms      map[string]*Room
}

func NewHub(rooms map[string]*Room, clients map[string]*Client) *Hub {
	h := Hub{}
	h.clients = clients
	h.register = make(chan *Client)
	h.unregister = make(chan *Client)
	h.rooms = rooms
	return &h
}

func (h *Hub) AddClient(c *Client) {
	h.register <- c
}

func (h *Hub) RemoveClient(c *Client) {
	h.unregister <- c
}

func (h *Hub) ClientExists(id string) bool {
	_, exists := h.clients[id]
	return exists
}

func (h *Hub) Run() {
	log.Println("Hub started.")
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
			log.Printf("Client <%s> registered to hub", client.id)
		case client := <-h.unregister:
			if h.ClientExists(client.id) {
				delete(h.clients, client.id)
				log.Printf("Client <%s> unregisted to hub", client.id)
			}
		}
	}
}
