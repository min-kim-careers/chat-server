package chat

import (
	"chat-server/internal/constants"
	"chat-server/internal/dto/messagein"
	"chat-server/internal/dto/messageout"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) handleOutbound() {
	for p := range c.outbound {
		err := c.conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			log.Printf("error sending message to client <%s>: %v", c.id, err)
			c.ctxCancel()
			return
		}
		log.Println("sent message:", string(p))
	}
}

func (c *Client) handleInbound() {
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			c.ctxCancel()
			return
		}
		c.inbound <- p
		log.Println("read message:", string(p))
	}
}

func (c *Client) handleMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.inbound:
			c.touch()

			m, err := messagein.ToMessageIn(data)
			if err != nil {
				log.Printf("error parsing message in: %v", err)
				continue
			}

			switch v := m.(type) {
			case *messagein.MessageInChat:
				c.handleChat(v)
			case *messagein.MessageInJoin:
				c.handleJoin(v)
			case *messagein.MessageInEvent:
				switch v.Mode {
				case "restore":
					switch c.hasNoMessages() {
					case true:
						c.handleNoMessages()
					case false:
						c.handleRestore()
					}
				case "leave":
					c.handleLeave()
				case "typing":
					c.handleTyping()
				case "not_typing":
					c.handleNotTyping()
				}
			}
		}
	}
}

func (c *Client) handleIdleTimeout() {
	ticker := time.NewTicker(constants.CLIENT_IDLE_TIMEOUT_CHECK_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()

			timeoutAt := c.lastActivity.Add(constants.CLIENT_IDLE_TIMEOUT)
			if time.Now().After(timeoutAt) {
				p, err := messageout.ToRawMessageOut(&messageout.MessageOutEvent{
					Mode: "timeout",
				})
				if err != nil {
					log.Println("error sending timeout:", err)
					continue
				}
				c.outbound <- p
			}

			c.mu.Unlock()
		}
	}
}

func (c *Client) handleClose() {
	<-c.ctx.Done()
	if c.hasRoom() {
		c.room.clientUnregister <- c
	}
	c.hub.clientUnregister <- c
	close(c.outbound)
	c.outbound = nil
	close(c.inbound)
	c.inbound = nil
	c.conn.Close()
}
