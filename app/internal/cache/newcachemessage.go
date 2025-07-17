package cache

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func NewCacheMessage(roomID string, clientID string, content string) (*CacheMessage, error) {
	newID, err := uuid.NewUUID()
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	_roomID, err := uuid.Parse(roomID)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	m := &CacheMessage{
		ID:        newID,
		RoomID:    _roomID,
		ClientID:  clientID,
		CreatedAt: time.Now(),
		Content:   content,
	}
	return m, nil
}
