package handler

import (
	"chat-server/internal/auth"
	"chat-server/internal/chat"
	"chat-server/internal/dto"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func authenticate(r *http.Request) bool {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("Auth token missing")
		return false
	}

	_, ok := auth.VerifyToken(token)
	if !ok {
		log.Printf("Authentication failed")
		return false
	}

	return true
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, hub *chat.Hub) {
	// if !authenticate(r) {
	// 	return
	// }

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

	_, msgJson, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error receiving message on connect: %v", err)
		conn.Close()
		return
	}

	message := dto.DeserializeMessage(msgJson)
	if message == nil {
		conn.Close()
		return
	}

	ctx := context.Background()
	hub.HandleConnection(ctx, conn, message)
}
