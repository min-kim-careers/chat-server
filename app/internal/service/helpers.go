package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/db/gen"
	"chat-server/internal/dto"
	"chat-server/internal/helper"

	"github.com/google/uuid"
)

func sortClientIds(client1 uuid.UUID, client2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if client2.String() < client1.String() {
		return client2, client1
	}
	return client1, client2
}

func roomDBToDTO(r gen.Room) *dto.RoomOut {
	return &dto.RoomOut{
		ID:        helper.ToDTOUUID(r.ID),
		ItemID:    r.ItemID,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func messageDBToDTO(m gen.Message, clientID string) *dto.MessageOutChat {
	return &dto.MessageOutChat{
		Mode:      "chat",
		ID:        *helper.ToDTOUUID(m.ID),
		CreatedAt: m.CreatedAt.Time,
		Read:      m.Read,
		IsMine:    m.ClientID == clientID,
		Content:   m.Content,
	}
}

func messageCacheToDTO(p string, clientID string) (*dto.MessageOutChat, error) {
	c, err := cache.ToMessageCache(p)
	if err != nil {
		return nil, err
	}
	return &dto.MessageOutChat{
		Mode:      "chat",
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		Read:      c.Read,
		Content:   c.Content,
		IsMine:    c.ClientID == clientID,
	}, nil
}
