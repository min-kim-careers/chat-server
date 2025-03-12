package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request, hub *Hub) {
	query := r.URL.Query()
	id := query.Get("client")
	if id == "" {
		log.Println("Client ID not provided. Disconnecting.")
		return
	}

	if hub.ClientExists(id) {
		log.Printf("Client <%s> already exists in hub. Disconnecting.", id)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		conn.Close()
		return
	}
	log.Printf("Successfully upgraded connection for <%s>\n", id)

	client := NewClient(id, conn)
	hub.AddClient(client)
}
