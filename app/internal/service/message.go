package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/constant"
	"chat-server/internal/db"
	"chat-server/internal/db/gen"
	"chat-server/internal/dto"
	"chat-server/internal/repo"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageService struct {
	r  *repo.MessageRepo
	db *db.DB
	c  *cache.Cache
}

func NewMessageService(r *repo.MessageRepo, db *db.DB, c *cache.Cache) *MessageService {
	return &MessageService{
		r:  r,
		db: db,
		c:  c,
	}
}

func (s *MessageService) GetDBMessages(ctx context.Context, roomID uuid.UUID, createdAt time.Time, limit int) ([]*dto.MessagePayload, error) {
	if roomID == uuid.Nil || createdAt.IsZero() || limit < 1 {
		return nil, errors.New("invalid params")
	}

	rows, err := s.r.GetMessages(ctx, gen.GetAllMessagesBeforeCreatedAtParams{
		RoomID:    pgtype.UUID{Bytes: roomID, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.MessagePayload, len(rows))
	for i, m := range rows {
		dtos[i] = ToMessagePayload(m)
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(dtos), roomID)
	return dtos, nil
}

func (s *MessageService) PublishMessage(ctx context.Context, roomID string, p []byte) bool {
	return s.c.Publish(ctx, roomKey(roomID), p)
}

func (s *MessageService) CacheMessage(ctx context.Context, roomID string, p []byte) bool {
	return s.c.Add(ctx, roomKey(roomID), p)
}

func (s *MessageService) MessageCacheIsFull(ctx context.Context, roomID string) bool {
	return s.c.IsFull(ctx, roomKey(roomID), constant.CACHE_LIMIT)
}

func (s *MessageService) GetCachedMessages(ctx context.Context, key string) []*dto.MessagePayload {
	cache := s.c.Range(ctx, roomKey(key), constant.CACHE_LIMIT)
	if len(cache) == 0 {
		return []*dto.MessagePayload{}
	}

	msgs := make([]*dto.MessagePayload, len(cache))
	for i, c := range cache {
		var msg dto.MessagePayload
		err := json.Unmarshal([]byte(c), &msg)
		if err != nil {
			log.Printf("Error unmarshalling message from key <%s>: %v", key, err)
			return nil
		}
		msgs[i] = &msg
	}

	log.Printf("Fetched %d from cache from key <%s>.", len(msgs), key)
	return msgs
}

func (s *MessageService) ClearMessageCache(ctx context.Context, roomID string) {
	s.c.Clear(ctx, roomKey(roomID))
}
