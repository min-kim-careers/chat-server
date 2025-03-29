package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const FETCH_LIMIT = 100

type ClientID string

type Client struct {
	id     ClientID
	wsConn *websocket.Conn
}

func NewClient(id ClientID, wsConn *websocket.Conn) *Client {
	return &Client{
		id:     id,
		wsConn: wsConn,
	}
}

// Send to client
func (client *Client) Send(msgJson []byte) {
	err := client.wsConn.WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		log.Printf("Error sending message to client <%s>: %v", client.id, err)
	} else {
		log.Printf("Sent message to client <%s>: %v", client.id, string(msgJson))
	}
}

// Read from client
func (client *Client) Read(room *Room) {
	roomID := string(room.id)

	for {
		_, msg, err := client.wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %v", client.id, err)
			room.unregister <- client
			return
		}

		err = room.cache.PublishMessage(roomID, msg)
		if err != nil {
			continue
		}

		if room.cache.CacheFull(roomID) {
			room.db.AddMessages()
		}

	}
}

func (client *Client) RestoreMessages(room *Room) {
	msgs := room.cache.RestoreMessages(string(room.id), FETCH_LIMIT)
	if len(msgs) == 0 {
		return
	}

	msg := Message{
		Type:      "restore",
		RoomID:    room.id,
		ClientID:  client.id,
		Content:   msgs,
		Timestamp: Timestamp(time.Now().Format(time.RFC3339)),
	}

	msgJson := SerializeMessage(&msg)
	if msgJson == nil {
		log.Printf("Error restoring cached messages for client <%s> in room <%s>.", client.id, room.id)
		return
	}

	client.Send(msgJson)
}

func (client *Client) Run(room *Room) {
	client.RestoreMessages(room)

	go client.Read(room)
}
