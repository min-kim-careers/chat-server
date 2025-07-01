package chat

import (
	"chat-server/internal/cache"
	"chat-server/internal/constant"
	"chat-server/internal/deps"
	"chat-server/internal/dto"
	"chat-server/internal/service"
	"context"
	"log"

	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(ctx context.Context, conn *websocket.Conn, id string, deps *deps.Container) *Client {
	ctx, cancel := context.WithCancel(ctx)
	return &Client{
		id:     id,
		deps:   deps,
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
	}
}

func (c *Client) services() *service.Services {
	return c.deps.Services
}

func (c *Client) cache() *cache.Cache {
	return c.deps.Cache
}

func (c *Client) handleChatMessage(roomID uuid.UUID, msgJson []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.cache().Publish(c.ctx, roomID.String(), msgJson) {
		return
	}

	if !c.cache().Add(c.ctx, roomID.String(), msgJson) {
		return
	}

	if !c.cache().IsFull(c.ctx, roomID.String(), constant.CACHE_LIMIT) {
		return
	}

	cachedMsgs := c.cache().Restore(c.ctx, roomID.String(), constant.RESTORE_LIMIT)
	if len(cachedMsgs) == 0 {
		return
	}

	// if !c.services().BulkAddMessages(cachedMsgs) {
	// 	return
	// }

	c.cache().Clear(c.ctx, roomID.String())
}

func (c *Client) handleRestoreMessage(roomID uuid.UUID, createdAt time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	dbMsgs, err := c.services().Message.GetPreviousMessages(
		c.ctx,
		roomID,
		createdAt,
		constant.RESTORE_LIMIT,
	)
	if err != nil {
		log.Printf("Error fetching old messages for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}

	data, err := dto.EncodeRaw(dbMsgs)
	if err != nil {
		log.Printf("Error encoding old messages for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}

	newMsg := &dto.Message{
		Mode:     "restore",
		RoomID:   roomID,
		ClientID: c.id,
		Data:     data,
	}

	msgJson := dto.SerializeMessage(newMsg)
	if msgJson == nil {
		log.Printf("Error restoring old messages for client <%s> in room <%s>: %v", c.id, roomID, err)
		return
	}

	c.Send(msgJson)
}

// Send to client/browser
func (c *Client) Send(msgJson []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.conn.WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		log.Printf("Error sending message to client <%s>: %v", c.id, err)
	} else {
		log.Printf("Sent message to client <%s>: %v", c.id, string(msgJson))
	}
}

// Read from client/browser
func (c *Client) Read(room *Room) {
	roomID := room.ID

	for {
		_, msgJson, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %v", c.id, err)
			room.Unregister <- c
			return
		}

		msg := dto.DeserializeMessage(msgJson)
		if msg == nil {
			continue
		}

		switch msg.Mode {
		case "chat":
			c.handleChatMessage(roomID, msgJson)
		case "restore":
			c.handleRestoreMessage(roomID, msg.CreatedAt)
		}

	}
}

func (c *Client) Run(room *Room) {
	go c.Read(room)
}
