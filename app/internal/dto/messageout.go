package dto

import (
	"encoding/json"
	"errors"
	"log"
)

type MessageOut struct {
	Mode string `json:"mode"`
}

func ValidateMessageOut(m *MessageOut) bool {
	_, valid := MessageModes[m.Mode]
	if !valid {
		log.Println("Invalid message mode:", m.Mode)
		return false
	}

	return true
}

func ToRawMessageOut(m *MessageOut) ([]byte, error) {
	if !ValidateMessageOut(m) {
		return nil, errors.New("invalid message out")
	}
	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return p, nil
}
