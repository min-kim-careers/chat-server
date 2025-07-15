package chat

import (
	"chat-server/internal/auth"
	"chat-server/internal/dto/messagein"
	"log"
)

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
