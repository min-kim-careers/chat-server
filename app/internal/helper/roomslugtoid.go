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
		log.Printf("error decoding room slug: %s", roomSlug)
		return nil
	}

	_roomID, err := uuid.FromBytes(roomID)
	if err != nil {
		log.Printf("error parsing room slug to UUID")
		return nil
	}

	return &_roomID
}
