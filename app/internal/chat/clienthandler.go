package chat

import (
	"chat-server/internal/auth"
	"chat-server/internal/constant"
	"chat-server/internal/dto"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func rehydrateChatMessage(m *dto.MessageIn, clientID string, roomID string) error {
	_roomID, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room ID format for message: %+v", m)
	}
	m.RoomID = _roomID
	m.ClientID = clientID
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return nil
}

func (c *Client) handleChatMessage(m *dto.MessageIn, roomID string) {
	err := rehydrateChatMessage(m, c.id, roomID)
	if err != nil {
		log.Printf("Error rehydrating chat message: %v", err)
		return
	}

	msgSvc := c.hub.svc.Message

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

	if msgSvc.MessageCacheIsFull(c.ctx, roomID) {
		cachedMsgs, err := msgSvc.GetMessagesFromCache(c.ctx, roomID, c.id)
		if err != nil {
			log.Printf("Error fetching messages from cache for client <%s> in room <%s>: %v", c.id, roomID, err)
			return
		}
		err = msgSvc.BulkInsertMessagesDB(c.ctx, cachedMsgs)
		if err != nil {
			log.Printf("Error bulk persisting messages for client <%s>: %v", c.id, err)
			return
		}
		msgSvc.ClearMessageCache(c.ctx, roomID)
	}
}

func (c *Client) handleRestoreMessage(m *dto.MessageIn, roomID string) {
	restoredMsgs := []*dto.MessageOutChat{}

	cachedMsgs, err := c.hub.svc.Message.GetMessageOutsFromCache(c.ctx, roomID, c.id)
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
		dbMsgs, err := c.hub.svc.Message.GetMessagesDB(
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

	_m := &dto.MessageOutRestore{
		Mode:     "restored",
		Messages: restoredMsgs,
	}

	p, err := json.Marshal(_m)
	if err != nil {
		return
	}
	c.channel <- p
}

func (c *Client) handleJoinMessage(m *dto.MessageIn) {
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
