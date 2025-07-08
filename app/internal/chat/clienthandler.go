package chat

import (
	"chat-server/internal/auth"
	"chat-server/internal/constant"
	"chat-server/internal/dto"
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
	m.RoomID = &_roomID
	m.ClientID = clientID
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return nil
}

func (c *Client) handleChatMessage(m *dto.MessageIn, roomID string) {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	err := rehydrateChatMessage(m, c.id, roomID)
	if err != nil {
		log.Printf("Error rehydrating chat message: %v", err)
		return
	}

	msgSvc := c.hub.svc.Message

	if !msgSvc.PublishChatMessage(c.ctx, roomID, m) {
		return
	}

	if !msgSvc.CacheChatMessage(c.ctx, roomID, m) {
		return
	}

	cacheSize := msgSvc.GetCacheSize(c.ctx, roomID)
	if cacheSize >= constant.CACHE_LIMIT {
		err = msgSvc.FlushCachedMessagesToDB(c.ctx, roomID, c.id, cacheSize)
		if err != nil {
			log.Printf("Error bulk persisting messages for client <%s>: %v", c.id, err)
			return
		}
		msgSvc.ClearChatMessageCache(c.ctx, roomID)
	}
}

func (c *Client) handleRestoreMessage(m *dto.MessageIn, roomID string) {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	restoredMsgs := []*dto.MessageOutChat{}

	cachedMsgs, err := c.hub.svc.Message.GetCachedChatMessages(c.ctx, roomID, c.id)
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
		dbMsgs, err := c.hub.svc.Message.GetDBMessages(
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

	p, err := dto.ToRawMessageOut(&dto.MessageOutRestore{
		Mode:     "restored",
		Messages: restoredMsgs,
	})
	if err != nil {
		log.Printf("Error parsing restore message: %v", err)
		return
	}
	c.channel <- p
}

func (c *Client) handleJoinMessage(m *dto.MessageIn) {
	if c.hasRoom() {
		log.Printf("Client is already in a room")
		return
	}

	err := auth.IsAuthorised(c.ctx, c.hub.svc.Room, c.id, *m.RoomID)
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
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

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
