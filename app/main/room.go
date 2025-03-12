package main

import "log"

type Room struct {
	itemID     string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewRoom(itemID string, clients map[*Client]bool, broadcast chan []byte) *Room {
	r := Room{}
	r.itemID = itemID
	r.clients = clients
	r.register = make(chan *Client)
	r.unregister = make(chan *Client)
	r.broadcast = broadcast
	return &r
}

func (r *Room) ClientExists(c *Client) bool {
	_, exists := r.clients[c]
	return exists
}

func (r *Room) AddClient(c *Client) {
	r.clients[c] = true
}

func (r *Room) RemoveClient(c *Client) {
	delete(r.clients, c)
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			log.Printf("Client %s registered to room", client.id)
		case client := <-r.unregister:
			if r.ClientExists(client) {
				delete(r.clients, client)
				log.Printf("Client %s unregisted to room", client.id)
			}
		}
	}
}
