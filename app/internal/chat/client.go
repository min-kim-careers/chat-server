package chat

import (
	"chat-server/internal/cache"
	"chat-server/internal/db"
	"chat-server/internal/models"
	"chat-server/internal/utils"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	wsConn *websocket.Conn
	lock   sync.Mutex
}

func NewClient(id string, wsConn *websocket.Conn) *Client {
	return &Client{
		id:     id,
		wsConn: wsConn,
	}
}

func (client *Client) initializeMessages(room *Room) {
	client.lock.Lock()
	defer client.lock.Unlock()

	msgs := utils.ReverseOrder(room.cache.Restore(room.id, utils.RESTORE_LIMIT))

	if len(msgs) < utils.RESTORE_LIMIT {
		var lastTimestamp string
		if len(msgs) == 0 {
			lastTimestamp = time.Now().Format(utils.TIMESTAMP_FORMAT)
		} else {
			lastTimestamp = msgs[len(msgs)-1].Timestamp
		}

		delta := utils.RESTORE_LIMIT - len(msgs)
		dbMsgs := room.db.RestoreMessages(room.id, lastTimestamp, delta)
		if dbMsgs != nil {
			msgs = append(msgs, dbMsgs...)
		}
	}

	msg := &models.Message{
		Type:     "restore",
		RoomID:   room.id,
		ClientID: client.id,
		Content:  msgs,
	}

	msgJson := models.SerializeMessage(msg)
	if msgJson == nil {
		log.Printf("Error restoring cached messages for client <%s> in room <%s>.", client.id, room.id)
		return
	}

	client.Send(msgJson)
}

func (client *Client) handleChatMessage(roomID string, msgJson []byte, cache *cache.Cache, db *db.DB) {
	client.lock.Lock()
	defer client.lock.Unlock()

	if !cache.Publish(roomID, msgJson) {
		return
	}

	if !cache.Add(roomID, msgJson) {
		return
	}

	if !cache.IsFull(roomID, utils.CACHE_LIMIT) {
		return
	}

	cachedMsgs := cache.Restore(roomID, utils.RESTORE_LIMIT)
	if len(cachedMsgs) == 0 {
		return
	}

	if !db.BulkInsertMessages(cachedMsgs) {
		return
	}

	cache.Clear(roomID)
}

func (client *Client) handleRestoreMessage(roomID, timestamp string, db *db.DB) {
	client.lock.Lock()
	defer client.lock.Unlock()

	dbMsgs := db.RestoreMessages(roomID, timestamp, utils.RESTORE_LIMIT)

	newMsg := &models.Message{
		Type:     "restore",
		RoomID:   roomID,
		ClientID: client.id,
		Content:  dbMsgs,
	}

	msgJson := models.SerializeMessage(newMsg)
	if msgJson == nil {
		log.Printf("Error restoring cached messages for client <%s> in room <%s>.", client.id, roomID)
		return
	}

	client.Send(msgJson)
}

// Send to client/browser
func (client *Client) Send(msgJson []byte) {
	err := client.wsConn.WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		log.Printf("Error sending message to client <%s>: %v", client.id, err)
	} else {
		log.Printf("Sent message to client <%s>: %v", client.id, string(msgJson))
	}
}

// Read from client/browser
func (client *Client) Read(room *Room) {
	roomID := room.id
	db := room.db
	cache := room.cache

	for {
		_, msgJson, err := client.wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %v", client.id, err)
			room.unregister <- client
			return
		}

		msg := models.DeserializeMessage(msgJson)
		if msg == nil {
			continue
		}

		if msg.Type == "chat" {
			client.handleChatMessage(roomID, msgJson, cache, db)
		} else if msg.Type == "restore" {
			client.handleRestoreMessage(roomID, msg.Timestamp, db)
		}

	}
}

func (client *Client) Run(room *Room) {
	client.initializeMessages(room)

	go client.Read(room)
}
