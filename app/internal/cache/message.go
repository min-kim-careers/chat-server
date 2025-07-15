package cache

import (
	"time"

	"github.com/google/uuid"
)

type MessageCache struct {
	Mode      string    `json:"mode"`
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	ClientID  string    `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
	Read      bool      `json:"read"`
	Sent      bool      `json:"sent"`
}
