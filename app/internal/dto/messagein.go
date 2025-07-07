package dto

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MessageIn struct {
	Mode      string    `json:"mode"`
	RoomID    uuid.UUID `json:"roomId"`
	ClientID  string    `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
	Read      bool      `json:"read"`
}

func validateMessageIn(m *MessageIn) bool {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		log.Println("Invalid message:", err)
		return false
	}

	mode, valid := MessageModes[m.Mode]
	if !valid {
		log.Println("Invalid message mode:", m.Mode)
		return false
	}

	switch mode {
	case "chat":
		if len(m.Content) == 0 {
			log.Println("blank chat message")
			return false
		}
	case "restore":
		if m.CreatedAt.IsZero() {
			log.Println("missing CreatedAt")
			return false
		}
	case "join":
		if m.RoomID == uuid.Nil {
			log.Println("missing RoomID")
			return false
		}
	}

	return true
}

func ToMessageIn(p []byte) (*MessageIn, error) {
	var m MessageIn
	err := json.Unmarshal(p, &m)
	if err != nil {
		return nil, err
	}
	if !validateMessageIn(&m) {
		return nil, errors.New("invalid message format")
	}
	return &m, nil
}
