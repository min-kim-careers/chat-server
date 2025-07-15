package chat

import (
	"chat-server/internal/dto/messagein"
	"chat-server/internal/dto/messageout"
	"chat-server/internal/service"
	"log"
)

func (c *Client) handleChat(m *messagein.MessageInChat) {
	if !c.hasRoom() {
		return
	}

	msgSvc := c.hub.svc.Message

	cachedMsg, err := msgSvc.CacheChatMessage(c.ctx, service.CacheChatMessageParams{
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
