package roomout

import (
	"time"
)

type RoomOut struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"itemId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
