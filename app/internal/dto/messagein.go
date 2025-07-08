package dto

import (
	"chat-server/internal/helper"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MessageIn struct {
	Mode      string     `json:"mode"`
	RoomSlug  string     `json:"roomSlug,omitempty"`
	RoomID    *uuid.UUID `json:"roomId,omitempty"`
	ClientID  string     `json:"clientId,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	Content   string     `json:"content,omitempty"`
	Read      bool       `json:"read,omitempty"`
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
		if len(m.RoomSlug) == 0 {
			log.Println("missing RoomSlug")
			return false
		}
	}

	return true
}

func parseSlug(m *MessageIn) error {
	if len(m.RoomSlug) == 0 {
		return nil
	}

	roomID, err := helper.DecodeSlug(m.RoomSlug)
	if err != nil {
		log.Printf("Error decoding room slug: %s", m.RoomSlug)
		return err
	}

	_roomID, err := uuid.FromBytes(roomID)
	if err != nil {
		log.Printf("Error parsing room slug")
		return err
	}

	m.RoomID = &_roomID
	return nil
}

func ToMessageIn(p []byte) (*MessageIn, error) {
	var m *MessageIn
	err := json.Unmarshal(p, &m)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", &p)
		return nil, err
	}

	err = parseSlug(m)
	if err != nil {
		return nil, err
	}

	if !validateMessageIn(m) {
		return nil, errors.New("invalid message format")
	}
	return m, nil
}
