package server

import (
	"chat-go/internal/chat"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request, hub *chat.Hub) {
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

	wsConnMsg, err := chat.DeserializeMessage(msgJson)
	if err != nil {
		wsConn.Close()
		return
	}

	hub.HandleWsConnection(wsConnMsg, wsConn)

}
