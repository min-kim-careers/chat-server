package service

import (
	"chat-server/db/gen"
	"chat-server/internal/dto"

	"github.com/google/uuid"
)

func sortClientIds(client1 uuid.UUID, client2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if client2.String() < client1.String() {
		return client2, client1
	}
	return client1, client2
}

func toRoomDTO(r gen.Room) *dto.Room {
	return &dto.Room{
		ID:        r.ID.Bytes,
		ItemID:    r.ItemID,
		Client1:   r.Client1.Bytes,
		Client2:   r.Client2.Bytes,
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}
