package dto

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `json:"id"`
	ItemID    string    `json:"itemId"`
	Client1   uuid.UUID `json:"client1"`
	Client2   uuid.UUID `json:"client2"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoomOut struct {
	ID        uuid.UUID `json:"id"`
	ItemID    string    `json:"itemId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
