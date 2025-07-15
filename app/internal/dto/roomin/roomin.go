package roomin

import (
	"time"

	"github.com/google/uuid"
)

type RoomIn struct {
	ID        uuid.UUID `json:"id"`
	ItemID    string    `json:"itemId"`
	Client1   uuid.UUID `json:"client1"`
	Client2   uuid.UUID `json:"client2"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
