package dto

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Message struct {
	Mode      string          `json:"mode"`
	RoomID    uuid.UUID       `json:"roomId"`
	ClientID  string          `json:"clientId"`
	CreatedAt time.Time       `json:"createdAt"`
	Data      json.RawMessage `json:"data"`
	Read      bool            `json:"read"`
}

func validateMessage(m *Message) bool {
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
	case "restore":
		if m.CreatedAt.IsZero() {
			log.Println(m.CreatedAt)
			log.Println("Restore missing created at")
			return false
		}
		return true
	}

	return true
}

func ToMessage(p []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(p, &m)
	if err != nil {
		return nil, err
	}
	if !validateMessage(&m) {
		return nil, errors.New("invalid message format")
	}
	return &m, nil
}
