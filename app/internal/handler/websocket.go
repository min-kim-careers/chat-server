package handler

import (
	"chat-server/internal/chat"

	"chat-server/internal/dto"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request, hub *chat.Hub) {
	clientId := mockVerifyClient()
	if clientId == nil {
		log.Printf("Authentication failed")
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		conn.Close()
		return
	}
	log.Printf("Successfully upgraded connection for <%s>", conn.RemoteAddr())

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error receiving message on connect: %v", err)
		conn.Close()
		return
	}

	m, err := dto.ToMessageDTO(p)
	if err != nil {
		log.Println("Error parsing connection payload")
		conn.Close()
		return
	}

	if !isAuthorised(r.Context(), hub.Svc.Room, *clientId, m.RoomID) {
		conn.Close()
		return
	}

	client := chat.NewClient(conn, clientId.String(), hub.Svc)
	hub.HandleConnection(m.RoomID.String(), client)
}
