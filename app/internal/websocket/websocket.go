package websocket

import (
	"chat-server/internal/auth"
	"chat-server/internal/chat"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request, hub *chat.Hub) {
	clientId := auth.VerifyClient(r)
	if clientId == nil {
		log.Printf("Authentication failed")
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error:", err)
		conn.Close()
		return
	}
	log.Printf("Successfully upgraded connection for <%s>", conn.RemoteAddr())

	client := chat.NewClient(conn, clientId.String(), hub)
	hub.HandleNewClient(client)
}
