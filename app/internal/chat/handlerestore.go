package chat

import (
	"chat-server/internal/constant"
	"chat-server/internal/dto/messagein"
	"chat-server/internal/dto/messageout"
	"chat-server/internal/service"
	"log"
)

func (c *Client) handleRestore(m *messagein.MessageInRestore) {
	if !c.hasRoom() {
		return
	}

	restoredMsgs := []*messageout.MessageOutChat{}

	cachedMsgs, err := c.hub.svc.Message.GetCachedMessages(c.ctx, service.GetCachedChatMessagesParams{
		ClientID: c.id,
		RoomID:   c.room.id,
		Before:   m.CreatedAt,
		Limit:    constant.CACHE_LIMIT,
	})
	if err != nil {
		log.Println("GetCachedChatMessages error:", err)
		return
	}
	restoredMsgs = append(restoredMsgs, cachedMsgs...)

	// if len(cachedMsgs) != constant.RESTORE_LIMIT {
	// 	_roomID, err := uuid.Parse(c.room.id)
	// 	if err != nil {
	// 		log.Println("error parsing room ID:", err)
	// 		return
	// 	}
	// 	dbMsgs, err := c.hub.svc.Message.GetDBMessages(
	// 		c.ctx,
	// 		_roomID,
	// 		m.CreatedAt,
	// 		constant.RESTORE_LIMIT-len(cachedMsgs),
	// 		c.id,
	// 	)
	// 	if err != nil {
	// 		log.Printf("error fetching messages from db for client <%s> in room <%s>: %v", c.id, c.room.id, err)
	// 		return
	// 	}
	// 	restoredMsgs = append(restoredMsgs, dbMsgs...)
	// }

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
