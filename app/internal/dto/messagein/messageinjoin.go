package messagein

import (
	"github.com/google/uuid"
)

type MessageInJoin struct {
	Mode     string     `json:"mode"`
	RoomSlug string     `json:"roomSlug,omitempty"`
	RoomID   *uuid.UUID `json:"roomId,omitempty"`
}

func (*MessageInJoin) isMessageIn() {}
