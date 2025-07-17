package chat

import (
	"chat-server/internal/auth"
	"chat-server/internal/constant"
	"chat-server/internal/dto/messagein"
	"chat-server/internal/dto/messageout"
	"chat-server/internal/service"
	"log"

	"github.com/google/uuid"
)

func (c *Client) handleChat(m *messagein.MessageInChat) {
	if !c.hasRoom() {
		return
	}

	msgSvc := c.hub.svc.Message

	cachedMsg, err := msgSvc.CacheAndPersistChatMessage(c.ctx, service.CacheChatMessageParams{
		ClientID: c.id,
		RoomID:   c.room.id,
		Content:  m.Content,
	})
	if err != nil {
		log.Printf("error caching chat: %v", err)
		return
	}

	err = msgSvc.PublishMessage(c.ctx, c.room.id, &messageout.MessageOutChat{
		ID:        cachedMsg.ID,
		Mode:      m.Mode,
		TempID:    &m.TempID,
		CreatedAt: cachedMsg.CreatedAt,
		IsMine:    true,
		Content:   cachedMsg.Content,
		Read:      false,
		Sent:      true,
	})
	if err != nil {
		log.Printf("error publish chat: %v", err)
	}
}

func (c *Client) handleDisconnect() {
	if c.hasRoom() {
		c.room.clientUnregister <- c
	}
	c.hub.clientUnregister <- c
	c.conn.Close()
}

func (c *Client) handleJoin(m *messagein.MessageInJoin) {
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

	room, exists := c.hub.getRoom(m.RoomID.String())
	if exists {
		room.clientRegister <- c
		c.setRoom(room)
		return
	}

	newRoom := NewRoom(c.hub, m.RoomID.String())
	c.hub.roomRegister <- newRoom
	newRoom.clientRegister <- c
	c.setRoom(newRoom)
}

func (c *Client) handleLeave() {
	if !c.hasRoom() {
		return
	}

	if c.hasRoom() {
		c.room.clientUnregister <- c
		c.setRoom(nil)
	}
}

func (c *Client) handleNotTyping() {
	if !c.hasRoom() {
		return
	}

	msgSvc := c.hub.svc.Message

	if err := msgSvc.PublishMessage(c.ctx, c.room.id, &messageout.MessageOutEvent{Mode: "not_typing"}); err != nil {
		log.Println("error sending typing message:", err)
		return
	}
}

func (c *Client) handleRestore() {
	if !c.hasRoom() {
		return
	}

	restoredMsgs := []*messageout.MessageOutChat{}

	cachedMsgs, lastCacheID, err := c.hub.svc.Message.GetCachedChatMessages(c.ctx, service.GetCachedChatMessagesParams{
		ClientID:    c.id,
		RoomID:      c.room.id,
		Limit:       constant.RESTORE_LIMIT,
		LastCacheID: c.cursor.LastCacheID,
	})
	if err != nil {
		log.Println("GetCachedChatMessages error:", err)
		return
	}
	if lastCacheID != nil {
		c.cursor.LastCacheID = *lastCacheID
	}

	restoredMsgs = append(restoredMsgs, cachedMsgs...)

	if len(restoredMsgs) != constant.RESTORE_LIMIT {
		_roomID, err := uuid.Parse(c.room.id)
		if err != nil {
			log.Println("error parsing room ID:", err)
			return
		}
		if c.cursor.LastDBID.IsZero() {
			c.cursor.LastDBID = restoredMsgs[len(restoredMsgs)-1].CreatedAt
		}
		dbMsgs, err := c.hub.svc.Message.GetDBMessages(
			c.ctx,

			service.GetDBMessagesParams{
				RoomID:    _roomID,
				CreatedAt: c.cursor.LastDBID,
				Limit:     constant.RESTORE_LIMIT - len(restoredMsgs),
				ClientID:  c.id,
			},
		)
		if err != nil {
			log.Println("GetDBMessages error:", err)
			return
		}

		restoredMsgs = append(restoredMsgs, dbMsgs...)
	}

	p, err := messageout.ToRawMessageOut(&messageout.MessageOutRestored{
		Mode:     "restored",
		Messages: restoredMsgs,
	})
	if err != nil {
		log.Printf("error parsing restore message: %v", err)
		return
	}
	c.channel <- p
}

func (c *Client) handleTyping() {
	if !c.hasRoom() {
		return
	}

	msgSvc := c.hub.svc.Message

	if err := msgSvc.PublishMessage(c.ctx, c.room.id, &messageout.MessageOutEvent{Mode: "typing"}); err != nil {
		log.Println("error sending typing message:", err)
		return
	}
}
