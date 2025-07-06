package service

import (
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

func dbRoomToDTO(r gen.Room) *dto.RoomOut {
	return &dto.RoomOut{
		ID:        helper.ToDTOUUID(r.ID),
		ItemID:    r.ItemID,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func dbMessageToDTO(m gen.Message, clientID string) *dto.MessageOutChat {
	return &dto.MessageOutChat{
		CreatedAt: m.CreatedAt.Time,
		Read:      m.Read,
		IsMine:    m.ClientID == clientID,
		Content:   m.Content,
	}
}
