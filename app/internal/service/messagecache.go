package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/constant"
	"chat-server/internal/dto"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func (s *MessageService) PublishChatMessage(ctx context.Context, roomID string, m *dto.MessageIn) bool {
	p, err := dto.ToRawMessageOut(m)
	if err != nil {
		log.Printf("%v", err)
		return false
	}
	return s.c.Publish(ctx, roomKey(roomID), p)
}

func (s *MessageService) CacheChatMessage(ctx context.Context, roomID string, m *dto.MessageIn) bool {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return false
	}
	p, err := json.Marshal(cache.Message{
		ID:        id,
		Mode:      m.Mode,
		RoomID:    m.RoomID,
		ClientID:  m.ClientID,
		CreatedAt: m.CreatedAt,
		Read:      m.Read,
		Content:   m.Content,
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return s.c.Add(ctx, roomKey(roomID), p)
}

func (s *MessageService) GetCacheSize(ctx context.Context, roomID string) int64 {
	return s.c.Size(ctx, roomKey(roomID), constant.CACHE_LIMIT)
}

func (s *MessageService) ClearChatMessageCache(ctx context.Context, roomID string) {
	s.c.Clear(ctx, roomKey(roomID))
}

func (s *MessageService) GetCachedChatMessages(ctx context.Context, key string, clientID string) ([]*dto.MessageOutChat, error) {
	rows := s.c.Range(ctx, roomKey(key), constant.CACHE_LIMIT)
	if len(rows) == 0 {
		return []*dto.MessageOutChat{}, nil
	}

	dtos := make([]*dto.MessageOutChat, len(rows))
	for i, r := range rows {
		dto, err := messageCacheToDTO(r, clientID)
		if err != nil {
			return nil, err
		}
		dtos[len(dtos)-1-i] = dto
	}

	log.Printf("Fetched %d from cache from key <%s>.", len(dtos), key)
	return dtos, nil
}
