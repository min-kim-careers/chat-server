package chat

func (c *Client) handleDisconnect() {
	if c.hasRoom() {
		c.room.clientUnregister <- c
	}
	c.hub.clientUnregister <- c
	c.conn.Close()
}
