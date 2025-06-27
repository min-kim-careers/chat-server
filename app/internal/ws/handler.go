package ws

import (
	"chat-server/internal/chat"
	"chat-server/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request, hub *chat.Hub) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		wsConn.Close()
		return
	}
	log.Printf("Successfully upgraded connection for <%s>", wsConn.RemoteAddr())

	_, msgJson, err := wsConn.ReadMessage()
	if err != nil {
		log.Printf("Error receiving message on connect: %v", err)
		wsConn.Close()
		return
	}

	connMsg := models.DeserializeMessage(msgJson)
	if connMsg == nil {
		wsConn.Close()
		return
	}

	hub.HandleWsConnection(wsConn, connMsg)

}
