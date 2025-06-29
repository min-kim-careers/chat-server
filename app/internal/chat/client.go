package chat

import (
	"chat-server/internal/cache"
	"chat-server/internal/constant"
	"chat-server/internal/dto"
	"chat-server/internal/service"
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ctx  context.Context
	conn *websocket.Conn
	id   string
	lock sync.Mutex
	room *Room
}

func NewClient(ctx context.Context, conn *websocket.Conn, id string, room *Room) *Client {
	return &Client{
		ctx:  ctx,
		conn: conn,
		id:   id,
		room: room,
	}
}

func (c *Client) services() *service.Services {
	return c.room.hub.deps.Services
}

func (c *Client) cache() *cache.Cache {
	return c.room.hub.deps.Cache
}

func (c *Client) initializeMessages(room *Room) {
	// c.lock.Lock()
	// defer c.lock.Unlock()

	// msgs := helper.ReverseOrder(room.hub.cache.Restore(room.id, constant.RESTORE_LIMIT))

	// if len(msgs) < constant.RESTORE_LIMIT {
	// 	var lastTimestamp time.Time
	// 	if len(msgs) == 0 {
	// 		lastTimestamp = time.Now().Format(constant.TIMESTAMP_FORMAT)
	// 	} else {
	// 		lastTimestamp = msgs[len(msgs)-1].CreatedAt
	// 	}

	// 	delta := constant.RESTORE_LIMIT - len(msgs)
	// 	dbMsgs := c.services().RestoreMessages(context.Background())
	// 	if dbMsgs != nil {
	// 		msgs = append(msgs, dbMsgs...)
	// 	}
	// }

	// msg := &dto.MessageDTO{
	// 	MessageType: "restore",
	// 	RoomID:      room.id,
	// 	ClientID:    c.id,
	// 	Data:        msgs,
	// }

	// msgJson := dto.SerializeMessage(msg)
	// if msgJson == nil {
	// 	log.Printf("Error restoring cached messages for client <%s> in room <%s>.", c.id, room.id)
	// 	return
	// }

	// c.Send(msgJson)
}

func (c *Client) handleChatMessage(roomID string, msgJson []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.cache().Publish(c.ctx, roomID, msgJson) {
		return
	}

	if !c.cache().Add(c.ctx, roomID, msgJson) {
		return
	}

	if !c.cache().IsFull(c.ctx, roomID, constant.CACHE_LIMIT) {
		return
	}

	cachedMsgs := c.cache().Restore(c.ctx, roomID, constant.RESTORE_LIMIT)
	if len(cachedMsgs) == 0 {
		return
	}

	// if !c.services().BulkAddMessages(cachedMsgs) {
	// 	return
	// }

	c.cache().Clear(c.ctx, roomID)
}

func (c *Client) handleRestoreMessage(roomID string, timestamp time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// dbMsgs := c.services().RestoreMessages(roomID, timestamp, helper.RESTORE_LIMIT)
	dbMsgs := []byte{}

	newMsg := &dto.Message{
		MessageType: "restore",
		RoomID:      roomID,
		ClientID:    c.id,
		Data:        dbMsgs,
	}

	msgJson := dto.SerializeMessage(newMsg)
	if msgJson == nil {
		log.Printf("Error restoring cached messages for client <%s> in room <%s>.", c.id, roomID)
		return
	}

	c.Send(msgJson)
}

// Send to client/browser
func (c *Client) Send(msgJson []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, msgJson)
	if err != nil {
		log.Printf("Error sending message to client <%s>: %v", c.id, err)
	} else {
		log.Printf("Sent message to client <%s>: %v", c.id, string(msgJson))
	}
}

// Read from client/browser
func (c *Client) Read(room *Room) {
	roomID := room.id

	for {
		_, msgJson, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %v", c.id, err)
			room.unregister <- c
			return
		}

		msg := dto.DeserializeMessage(msgJson)
		if msg == nil {
			continue
		}

		switch msg.MessageType {
		case "chat":
			c.handleChatMessage(roomID, msgJson)
		case "restore":
			c.handleRestoreMessage(roomID, msg.CreatedAt)
		}

	}
}

func (c *Client) Run(room *Room) {
	c.initializeMessages(room)

	go c.Read(room)
}
