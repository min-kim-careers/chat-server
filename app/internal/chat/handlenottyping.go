package chat

import (
	"chat-server/internal/dto/messageout"
	"log"
)

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
