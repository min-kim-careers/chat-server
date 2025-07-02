package dto

import (
	"encoding/json"
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

var MessageModes = map[string]string{
	"connect":    "connect",
	"disconnect": "disconnect",
	"chat":       "true",
	"restore":    "restore",
	"empty":      "empty",
}

type Message struct {
	Mode      string          `json:"mode"`
	RoomID    uuid.UUID       `json:"roomId"`
	ClientID  string          `json:"clientId"`
	CreatedAt time.Time       `json:"createdAt"`
	Data      json.RawMessage `json:"data"`
}

type MessagePayload struct {
	Mode      string          `json:"mode"`
	CreatedAt time.Time       `json:"createdAt"`
	Data      json.RawMessage `json:"data"`
}
