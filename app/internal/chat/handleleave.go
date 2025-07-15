package chat

func (c *Client) handleLeave() {
	if !c.hasRoom() {
		return
	}

	if c.hasRoom() {
		c.room.clientUnregister <- c
		c.setRoom(nil)
	}
}
