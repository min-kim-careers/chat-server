package messagein

import (
	"github.com/google/uuid"
)

type MessageIn interface {
	isMessageIn()
}

type MessageInBase struct {
	Mode string `json:"mode"`
}

func (*MessageInBase) isMessageIn() {}

type MessageInChat struct {
	Mode    string `json:"mode"`
	TempID  string `json:"tempId"`
	Content string `json:"content"`
}

func (*MessageInChat) isMessageIn() {}

type MessageInEvent struct {
	Mode     string `json:"mode"`
	ClientID string `json:"clientId,omitempty"`
}

func (*MessageInEvent) isMessageIn() {}

type MessageInJoin struct {
	Mode     string     `json:"mode"`
	RoomSlug string     `json:"roomSlug,omitempty"`
	RoomID   *uuid.UUID `json:"roomId,omitempty"`
}

func (*MessageInJoin) isMessageIn() {}
