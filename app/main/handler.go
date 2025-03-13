package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request, hub *Hub) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		log.Println("Disconnecting client.")
		conn.Close()
		return
	}
	log.Printf("Successfully upgraded connection for <%s>", conn.RemoteAddr())

	_, msgJson, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error receiving message on connect:", err)
		log.Println("Disconnecting client.")
		conn.Close()
		return
	}

	connMsg := Deserialize(msgJson)
	if connMsg == nil {
		log.Println("Disconnecting client.")
		conn.Close()
		return
	}

	hub.HandleConnection(connMsg, conn)

}
