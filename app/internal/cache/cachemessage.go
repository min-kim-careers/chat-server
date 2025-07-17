package cache

import (
	"time"

	"github.com/google/uuid"
)

type CacheMessage struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	ClientID  string    `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
}
