package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	wsConn *websocket.Conn
}

func NewClient(id string, wsConn *websocket.Conn) *Client {
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
	roomID := room.id
	db := room.db
	cache := room.cache

	for {
		_, msg, err := client.wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %v", client.id, err)
			room.unregister <- client
			return
		}

		err = cache.Publish(roomID, msg)
		if err != nil {
			continue
		}

		err = cache.Add(roomID, msg)
		if err != nil {
			continue
		}

		if !cache.IsFull(roomID) {
			continue
		}

		cachedMsgs := cache.Restore(roomID, RESTORE_LIMIT)
		if len(cachedMsgs) == 0 {
			continue
		}

		err = db.AddBulk(cachedMsgs)
		if err != nil {
			log.Printf("Error persisting messages for room <%s>: %v", roomID, err)
		} else {
			cache.Clear(roomID)
		}

	}
}

func (client *Client) RestoreMessages(room *Room) {
	msgs := room.cache.Restore(room.id, RESTORE_LIMIT)

	if len(msgs) < RESTORE_LIMIT {
		delta := RESTORE_LIMIT - len(msgs)
		dbMsgs := room.db.Restore(room.id, delta)
		if dbMsgs != nil {
			msgs = append(msgs, dbMsgs...)
		}
	}

	if len(msgs) == 0 {
		return
	}

	msg := Message{
		Type:      "restore",
		RoomID:    room.id,
		ClientID:  client.id,
		Content:   msgs,
		Timestamp: time.Now().Format(time.RFC3339),
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
