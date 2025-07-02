package dto

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
)

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

func ToMessageDTO(data []byte) (*Message, error) {
	var msg Message

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	if !validateMessage(&msg) {
		return nil, errors.New("invalid message format")
	}

	return &msg, nil
}

func ToMessagePayload(p []byte) ([]byte, error) {
	var m MessagePayload

	err := json.Unmarshal(p, &m)
	if err != nil {
		return nil, err
	}

	_p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return _p, nil
}

func ToRawMessages[T any](arr []T) (json.RawMessage, error) {
	b, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
