package dto

import (
	"time"

	"github.com/google/uuid"
)

type RoomOut struct {
	ID        uuid.UUID `json:"id"`
	ItemID    string    `json:"itemId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
