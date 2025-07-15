package messagein

import (
	"time"
)

type MessageInRestore struct {
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

func (*MessageInRestore) isMessageIn() {}
