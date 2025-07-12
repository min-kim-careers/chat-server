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

func hydrateChatMessage(m *dto.MessageIn, c *Client) error {
	_roomID, err := uuid.Parse(c.room.id)
	if err != nil {
		return fmt.Errorf("invalid room ID format for message: %+v", m)
	}
	m.RoomID = &_roomID
	m.ClientID = c.id
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return nil
}

func (c *Client) handleChat(m *dto.MessageIn) {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	err := hydrateChatMessage(m, c)
	if err != nil {
		log.Printf("Error hydrating chat message: %v", err)
		return
	}

	msgSvc := c.hub.svc.Message

	if !msgSvc.PublishChatMessage(c.ctx, c.room.id, m) {
		return
	}

	if !msgSvc.CacheChatMessage(c.ctx, c.room.id, m) {
		return
	}

	cacheSize := msgSvc.GetCacheSize(c.ctx, c.room.id)
	if cacheSize >= constant.CACHE_LIMIT {
		err = msgSvc.FlushCachedMessagesToDB(c.ctx, c.room.id, c.id, cacheSize)
		if err != nil {
			log.Printf("Error bulk persisting messages for client <%s>: %v", c.id, err)
			return
		}
		msgSvc.ClearChatMessageCache(c.ctx, c.room.id)
	}
}

func (c *Client) handleRestore(m *dto.MessageIn) {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	restoredMsgs := []*dto.MessageOutChat{}

	cachedMsgs, err := c.hub.svc.Message.GetCachedChatMessages(c.ctx, c.room.id, c.id)
	if err != nil {
		log.Printf("Error fetching messages from cache for client <%s> in room <%s>: %v", c.id, c.room.id, err)
		return
	}
	restoredMsgs = append(restoredMsgs, cachedMsgs...)

	if len(cachedMsgs) != constant.RESTORE_LIMIT {
		_roomID, err := uuid.Parse(c.room.id)
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
			log.Printf("Error fetching messages from db for client <%s> in room <%s>: %v", c.id, c.room.id, err)
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

func (c *Client) handleJoin(m *dto.MessageIn) {
	if c.hasRoom() {
		c.handleLeave()
		c.handleJoin(m)
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
	if newRoom == nil {

	}
	c.hub.roomRegister <- newRoom
	newRoom.clientRegister <- c
	c.room = newRoom
}

func (c *Client) handleLeave() {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	if c.hasRoom() {
		c.room.clientUnregister <- c
		c.room = nil
	}
}

func (c *Client) handleDisconnect() {
	if c.hasRoom() {
		c.room.clientUnregister <- c
	}
	c.hub.clientUnregister <- c
	c.conn.Close()
}

func hydrateTypingMessage(m *dto.MessageIn, c *Client) error {
	m.ClientID = c.id
	return nil
}

func (c *Client) handleTyping(m *dto.MessageIn) {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	err := hydrateTypingMessage(m, c)
	if err != nil {
		log.Printf("Error hydrating typing message: %v", err)
		return
	}

	msgSvc := c.hub.svc.Message

	if !msgSvc.PublishChatMessage(c.ctx, c.room.id, m) {
		log.Println("Error sending typing message")
	}
}

func hydrateNotTypingMessage(m *dto.MessageIn, c *Client) error {
	m.ClientID = c.id
	return nil
}

func (c *Client) handleNotTyping(m *dto.MessageIn) {
	if !c.hasRoom() {
		log.Printf("Client is not in a room")
		return
	}

	err := hydrateNotTypingMessage(m, c)
	if err != nil {
		log.Printf("Error hydrating typing message: %v", err)
		return
	}

	msgSvc := c.hub.svc.Message

	if !msgSvc.PublishChatMessage(c.ctx, c.room.id, m) {
		log.Println("Error sending typing message")
	}
}
