package chat

import (
	"chat-server/internal/constant"
	"chat-server/internal/dto"
	"chat-server/internal/helper"
	"chat-server/internal/service"
	"context"
	"encoding/json"
	"log"

	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(conn *websocket.Conn, id string, svc *service.Services) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		id:     id,
		svc:    svc,
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
	}
}

func (c *Client) handleChatMessage(roomID string, m *dto.Message) {
	m.ClientID = c.id
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	p, err := json.Marshal(m)
	if err != nil {
		log.Printf("Failed to marshal chat message from <%s>: %v", c.id, err)
		return
	}

	if !c.svc.Message.PublishMessage(c.ctx, roomID, p) {
		return
	}

	if !c.svc.Message.CacheMessage(c.ctx, roomID, p) {
		return
	}

	if !c.svc.Message.MessageCacheIsFull(c.ctx, roomID) {
		return
	}

	cachedMsgs, err := c.svc.Message.GetMessagesFromCache(c.ctx, roomID, c.id)
	if err != nil {
		log.Printf("Error fetching messages from cache for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}
	if len(cachedMsgs) == 0 {
		return
	}

	// if !c.services().BulkAddMessages(cachedMsgs) {
	// 	return
	// }

	c.svc.Message.ClearMessageCache(c.ctx, roomID)
}

func (c *Client) handleRestoreMessage(roomID string, createdAt time.Time) {
	restoredMsgs := []*dto.MessageOut{}

	cachedMsgs, err := c.svc.Message.GetMessagesFromCache(c.ctx, roomID, c.id)
	if err != nil {
		log.Printf("Error fetching messages from cache for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}
	restoredMsgs = append(restoredMsgs, cachedMsgs...)

	if len(cachedMsgs) != constant.RESTORE_LIMIT {
		_roomID, err := uuid.Parse(roomID)
		if err != nil {
			log.Println("Error parsing room ID:", err)
			return
		}
		dbMsgs, err := c.svc.Message.GetMessagesFromDB(
			c.ctx,
			_roomID,
			createdAt,
			constant.RESTORE_LIMIT-len(cachedMsgs),
			c.id,
		)
		if err != nil {
			log.Printf("Error fetching messages from db for client <%s> in room <%s>: %v", c.id, roomID, err)
			return
		}
		restoredMsgs = append(restoredMsgs, dbMsgs...)
	}

	data, err := helper.ToRawMessages(restoredMsgs)
	if err != nil {
		log.Printf("Error encoding restored messages for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}

	p, err := json.Marshal(&dto.MessageOut{
		Mode: "restore",
		Data: data,
	})
	if err != nil {
		log.Printf("Error marshalling restored messages for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}

	err = c.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		log.Printf("Error sending restored messages to client <%s>: %v", c.id, err)
	}
	log.Printf("Sent restored messages to client <%s>: %v", c.id, string(p))
}

// Send to client/browser
func (c *Client) Send(p []byte) {
	m, err := dto.ToMessageOut(p, c.id)
	if err != nil {
		return
	}

	_p, err := json.Marshal(m)
	if err != nil {
		return
	}

	err = c.conn.WriteMessage(websocket.TextMessage, _p)
	if err != nil {
		log.Printf("Error sending message to client <%s>: %v", c.id, err)
	}
	log.Printf("Sent message to client <%s>: %v", c.id, string(p))
}

// Read from client/browser
func (c *Client) Read(r *Room) {

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading payload from client <%s>: %v", c.id, err)
			r.Unregister <- c
			return
		}

		m, err := dto.ToMessage(p)
		if err != nil {
			log.Println("Error parsing payload:", err)
			continue
		}

		switch m.Mode {
		case "chat":
			c.handleChatMessage(r.ID, m)
		case "restore":
			c.handleRestoreMessage(r.ID, m.CreatedAt)
		}
	}
}

func (c *Client) Run(room *Room) {
	go c.Read(room)
}
