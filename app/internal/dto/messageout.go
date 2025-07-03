package dto

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type MessageOut struct {
	Mode      string          `json:"mode"`
	CreatedAt time.Time       `json:"createdAt"`
	Data      json.RawMessage `json:"data"`
	Read      bool            `json:"read"`
	IsMine    bool            `json:"isMine"`
}

func validateMessageOut(m *MessageOut) bool {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		log.Println("Invalid message out:", err)
		return false
	}

	mode, valid := MessageModes[m.Mode]
	if !valid {
		log.Println("Invalid message mode:", m.Mode)
		return false
	}

	switch mode {
	case "restore":

	}

	return true
}

func ToMessageOut(p []byte, clientID string) (*MessageOut, error) {
	var m Message
	err := json.Unmarshal(p, &m)
	if err != nil {
		log.Printf("Error unmarshalling cached message: %v", err)
		return nil, err
	}

	_m := MessageOut{
		Mode:      m.Mode,
		CreatedAt: m.CreatedAt,
		Data:      m.Data,
		Read:      m.Read,
		IsMine:    m.ClientID == clientID,
	}

	if !validateMessageOut(&_m) {
		return nil, errors.New("invalid message out format")
	}
	return &_m, nil
}
