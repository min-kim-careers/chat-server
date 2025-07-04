package service

import (
	"chat-server/internal/constant"
	"chat-server/internal/dto"
	"context"
	"log"
)

func (s *MessageService) PublishMessage(ctx context.Context, roomID string, p []byte) bool {
	return s.c.Publish(ctx, roomKey(roomID), p)
}

func (s *MessageService) CacheMessage(ctx context.Context, roomID string, p []byte) bool {
	return s.c.Add(ctx, roomKey(roomID), p)
}

func (s *MessageService) MessageCacheIsFull(ctx context.Context, roomID string) bool {
	return s.c.IsFull(ctx, roomKey(roomID), constant.CACHE_LIMIT)
}

func (s *MessageService) GetMessagesFromCache(ctx context.Context, key string, clientID string) ([]*dto.MessageOut, error) {
	rows := s.c.Range(ctx, roomKey(key), constant.CACHE_LIMIT)
	if len(rows) == 0 {
		return []*dto.MessageOut{}, nil
	}

	dtos := make([]*dto.MessageOut, len(rows))
	for i, r := range rows {
		_m, err := dto.ToMessageOut([]byte(r), clientID)
		if err != nil {
			return nil, err
		}
		dtos[len(rows)-1-i] = _m
	}

	log.Printf("Fetched %d from cache from key <%s>.", len(dtos), key)
	return dtos, nil
}

func (s *MessageService) ClearMessageCache(ctx context.Context, roomID string) {
	s.c.Clear(ctx, roomKey(roomID))
}
