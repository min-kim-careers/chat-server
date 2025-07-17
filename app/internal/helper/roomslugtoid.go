package helper

import (
	"log"

	"github.com/google/uuid"
)

func RoomSlugToID(roomSlug string) *uuid.UUID {
	if len(roomSlug) == 0 {
		return nil
	}

	roomID, err := DecodeSlug(roomSlug)
	if err != nil {
		log.Println("error:", err)
		return nil
	}

	_roomID, err := uuid.FromBytes(roomID)
	if err != nil {
		log.Println("error:", err)
		return nil
	}

	return &_roomID
}
