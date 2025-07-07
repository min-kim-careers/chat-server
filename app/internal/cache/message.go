package cache

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Mode      string    `json:"mode"`
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	ClientID  string    `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Read      bool      `json:"read"`
	Content   string    `json:"content"`
}

func ToMessageCache(c string) (*Message, error) {
	var m *Message
	err := json.Unmarshal([]byte(c), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
