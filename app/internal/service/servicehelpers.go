package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/db/gen"
	"chat-server/internal/dto/messageout"
	"chat-server/internal/dto/roomout"
	"chat-server/internal/helper"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func sortClientIds(client1 uuid.UUID, client2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if client2.String() < client1.String() {
		return client2, client1
	}
	return client1, client2
}

func dbToRoomOut(r gen.Room) *roomout.RoomOut {
	return &roomout.RoomOut{
		ID:        helper.EncodeSlug(r.ID.Bytes[:]),
		ItemID:    r.ItemID,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func dbToMessageOutChat(m gen.Message, clientID string) *messageout.MessageOutChat {
	return &messageout.MessageOutChat{
		Mode:      "chat",
		ID:        *helper.ToDTOUUID(m.ID),
		CreatedAt: m.CreatedAt.Time,
		Read:      m.Read,
		IsMine:    m.ClientID == clientID,
		Content:   m.Content,
	}
}

func cacheToMessageOutChat(s map[string]any, clientID string) (*messageout.MessageOutChat, error) {
	p, err := json.Marshal(s)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var c *cache.CacheMessage
	err = json.Unmarshal(p, &c)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	m := &messageout.MessageOutChat{
		Mode:      "chat",
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		Content:   c.Content,
		IsMine:    c.ClientID == clientID,
	}
	return m, nil
}
