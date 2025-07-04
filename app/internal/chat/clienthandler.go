package chat

import (
	"chat-server/internal/auth"
	"chat-server/internal/constant"
	"chat-server/internal/dto"
	"chat-server/internal/helper"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

func (c *Client) handleChatMessage(m *dto.Message, roomID string) {
	msgSvc := c.hub.svc.Message

	m.ClientID = c.id
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	p, err := json.Marshal(m)
	if err != nil {
		log.Printf("Failed to marshal chat message from <%s>: %v", c.id, err)
		return
	}

	if !msgSvc.PublishMessage(c.ctx, roomID, p) {
		return
	}

	if !msgSvc.CacheMessage(c.ctx, roomID, p) {
		return
	}

	if !msgSvc.MessageCacheIsFull(c.ctx, roomID) {
		return
	}

	cachedMsgs, err := msgSvc.GetMessagesFromCache(c.ctx, roomID, c.id)
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

	msgSvc.ClearMessageCache(c.ctx, roomID)
}

func (c *Client) handleRestoreMessage(m *dto.Message, roomID string) {
	restoredMsgs := []*dto.MessageOut{}

	cachedMsgs, err := c.hub.svc.Message.GetMessagesFromCache(c.ctx, roomID, c.id)
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
		dbMsgs, err := c.hub.svc.Message.GetMessagesFromDB(
			c.ctx,
			_roomID,
			m.CreatedAt,
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

	p, err := dto.NewMessagePayload(&dto.MessageOut{
		Mode: "restored",
		Data: data,
	})
	if err != nil {
		return
	}
	c.channel <- p
}

func (c *Client) handleJoinMessage(m *dto.Message) {
	err := auth.IsAuthorised(c.ctx, c.hub.svc.Room, c.id, m.RoomID)
	if err != nil {
		log.Printf("Unauthorised user <%s>: %v", c.id, err)
		return
	}

	room, exists := c.hub.rooms[m.RoomID.String()]
	if exists {
		room.clientRegister <- c
		c.room = room
		return
	}

	newRoom := NewRoom(c.hub, m.RoomID.String())
	c.hub.roomRegister <- newRoom
	newRoom.clientRegister <- c
	c.room = newRoom
}

func (c *Client) handleLeaveMessage() {
	if c.hasRoom() {
		c.room.clientUnregister <- c
	}
}

func (c *Client) handleDisconnectMessage() {
	if c.hasRoom() {
		c.room.clientUnregister <- c
	}
	c.hub.clientUnregister <- c
	c.conn.Close()
}
