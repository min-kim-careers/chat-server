package chat

import (
	"log"

	"github.com/jackc/pgx/v5"
)

type RoomID string

type Room struct {
	id         RoomID
	clients    map[ClientID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	dbConn     *pgx.Conn
}

func NewRoom(id RoomID, dbConn *pgx.Conn) *Room {
	r := Room{}

	r.id = id
	r.clients = make(map[ClientID]*Client)
	r.register = make(chan *Client)
	r.unregister = make(chan *Client)
	r.broadcast = make(chan *Message)
	r.dbConn = dbConn

	return &r
}

func (r *Room) handleRegistrations(h *Hub) {
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
				log.Printf("Room <%s> is empty. Requesting removal from hub.", r.id)
				h.unregister <- r
				return
			}
		}
	}
}

func (r *Room) handleBroadcasts() {
	for msg := range r.broadcast {
		for _, c := range r.clients {
			if msg.ClientID != c.id {
				select {
				case c.channel <- msg:
				default:
					log.Printf("Client <%s> message buffer full, dropping message.", c.id)
				}
			}
		}
	}
}

func (r *Room) Run(h *Hub) {
	go r.handleRegistrations(h)
	go r.handleBroadcasts()
}
