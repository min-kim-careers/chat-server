package handler

import (
	"chat-server/internal/auth"
	"chat-server/internal/chat"
	"chat-server/internal/deps"
	"chat-server/internal/dto"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func authenticate(r *http.Request) *auth.VerifyTokenResponse {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("Auth token missing")
		return nil
	}

	resp, ok := auth.VerifyToken(token)
	if !ok {
		log.Printf("Authentication failed")
		return nil
	}

	return resp
}

func authorise(ctx context.Context, deps *deps.Container, userID uuid.UUID, roomID uuid.UUID) bool {
	room, err := deps.Services.Room.GetRoomById(ctx, roomID)
	log.Println(room)
	if err != nil {
		return false
	}
	if room.Client1 == userID || room.Client2 == userID {
		return true
	}
	log.Printf("Unauthorised user <%s>", userID)
	return false
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, hub *chat.Hub) {
	verResp := authenticate(r)
	if verResp == nil {
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

	_, msgJson, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error receiving message on connect: %v", err)
		conn.Close()
		return
	}

	message := dto.DeserializeMessage(msgJson)
	if message == nil {
		log.Println("Failed to deserialize message")
		conn.Close()
		return
	}

	if !authorise(r.Context(), hub.Deps, verResp.UserID, message.RoomID) {
		conn.Close()
		return
	}

	client := chat.NewClient(r.Context(), conn, verResp.UserID.String(), hub.Deps)
	hub.HandleConnection(message.RoomID, client)
}
