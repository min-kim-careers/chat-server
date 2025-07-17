package messageout

import (
	"time"

	"github.com/google/uuid"
)

type MessageOut interface {
	isMessageOut()
}

type MessageOutBase struct {
	Mode string `json:"mode"`
}

func (*MessageOutBase) isMessageOut() {}

type MessageOutChat struct {
	Mode      string    `json:"mode"`
	TempID    *string   `json:"tempId,omitempty"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	IsMine    bool      `json:"isMine"`
	Content   string    `json:"content"`
	Read      bool      `json:"read"`
	Sent      bool      `json:"sent"`
}

func (*MessageOutChat) isMessageOut() {}

type MessageOutEvent struct {
	Mode string `json:"mode"`
}

func (*MessageOutEvent) isMessageOut() {}

type MessageOutRestored struct {
	Mode     string            `json:"mode"`
	Messages []*MessageOutChat `json:"messages"`
}

func (*MessageOutRestored) isMessageOut() {}
