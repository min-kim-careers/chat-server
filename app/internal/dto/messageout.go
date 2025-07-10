package dto

import (
	"chat-server/internal/constant"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type MessageOut struct {
	Mode string `json:"mode"`
}

type MessageOutRestore struct {
	Mode     string            `json:"mode"`
	Messages []*MessageOutChat `json:"messages"`
}

type MessageOutChat struct {
	Mode      string    `json:"mode"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Read      bool      `json:"read"`
	IsMine    bool      `json:"isMine"`
	Content   string    `json:"content"`
}

func validateMessageOut(s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	mode := v.FieldByName("Mode")
	if !mode.IsValid() {
		return fmt.Errorf("field 'Mode' not found")
	}

	_, valid := constant.ChatModes[mode.String()]
	if !valid {
		return fmt.Errorf("invalid message mode")
	}

	return nil
}

func ToRawMessageOut(m any) ([]byte, error) {
	if err := validateMessageOut(m); err != nil {
		return nil, err
	}
	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return p, nil
}
