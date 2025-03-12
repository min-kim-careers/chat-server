package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

// ----------------------------- Client -----------------------------

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

func makeClient(conn *websocket.Conn, ID string) *Client {
	return &Client{
		ID:   ID,
		Conn: conn,
		Send: make(chan []byte),
	}
}

func (c *Client) readMessage(r *Room) {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message by client <%s>: %s\n", c.ID, err)
			c.Conn.Close()
			return
		}
		r.Broadcast <- message
	}
}

func (c *Client) sendMessage(r *Room) {
	for message := range r.Broadcast {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending message by client <%s>: %s\n", c.ID, err)
			c.Conn.Close()
			return
		}
	}
}

// ----------------------------- Room -----------------------------

type Room struct {
	ItemID    string
	Clients   map[string]*Client
	Broadcast chan []byte
}

func makeRoom(params *Params) *Room {
	return &Room{
		ItemID:    params.ItemID,
		Clients:   make(map[string]*Client),
		Broadcast: make(chan []byte),
	}
}

func (r *Room) clientExists(clientID string) bool {
	_, exists := r.Clients[clientID]
	return exists
}

func (r *Room) addClient(client *Client) {
	r.Clients[client.ID] = client
}

func (r *Room) removeClient(clientID string) {
	delete(r.Clients, clientID)
}

func (r *Room) run() {
	for {

	}
}

// ----------------------------- Hub -----------------------------

type Hub struct {
	Rooms   map[string]*Room
	Clients map[*string]bool
}

func makeHub() *Hub {
	return &Hub{
		Rooms:   make(map[string]*Room),
		Clients: make(map[*string]bool),
	}
}

func (h *Hub) ClientExists(client_id string) bool {
	_, exists := h.Clients[&client_id]
	return exists
}

// ----------------------------- Handler -----------------------------

type Params struct {
	ClientID string `schema:"client_id,required"`
}

var hub = makeHub()
var decoder = schema.NewDecoder()

// func extractRoomKey(params *Params) string {
// 	vals := []string{params.ClientID}
// 	slices.Sort(vals)
// 	return strings.Join(vals, "_")
// }

func handleWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		conn.Close()
		return
	}
	log.Printf("Successfully upgraded connection for <%v>\n", conn.RemoteAddr())
}

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	// Initializing Database

	conn, err := InitDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())
	fmt.Println("DB init successful")

	// Initializing Websocket

	flag.Parse()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocketConnection(w, r)
	})

	log.Printf("Websocket server starting on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
