package dto

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type MessageOutChat struct {
	Mode      string    `json:"mode"`
	RoomID    uuid.UUID `json:"roomId"`
	ClientID  string    `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Read      bool      `json:"read"`
	IsMine    bool      `json:"isMine"`
	Content   string    `json:"content"`
}

func ToChatMessageOut(p []byte, clientID string) (*MessageOutChat, error) {
	var c MessageOutChat
	err := json.Unmarshal(p, &c)
	if err != nil {
		log.Printf("Error unmarshalling message out: %v", err)
		return nil, err
	}
	return &c, nil
}
