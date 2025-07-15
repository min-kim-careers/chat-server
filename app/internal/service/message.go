package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/dto/messageout"
	"chat-server/internal/helper"

	"chat-server/internal/db"
	"chat-server/internal/db/gen"
	"chat-server/internal/repo"
	"context"
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

func (s *MessageService) GetDBMessages(ctx context.Context, roomID uuid.UUID, createdAt time.Time, limit int, clientID string) ([]*messageout.MessageOutChat, error) {
	if roomID == uuid.Nil || createdAt.IsZero() || limit < 1 {
		return nil, errors.New("invalid params")
	}

	rows, err := s.r.GetMessages(ctx, gen.GetMessagesBeforeCreatedAtParams{
		RoomID:    helper.ToDBUUID(roomID),
		CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]*messageout.MessageOutChat, len(rows))
	for i, r := range rows {
		dtos[len(dtos)-1-i] = dbToMessageOut(r, clientID)
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(dtos), roomID)
	return dtos, nil
}

func (s *MessageService) FlushCachedMessagesToDB(ctx context.Context, roomID string, clientID string, cacheSize int64) error {
	// rows := s.c.Range(ctx, roomMsgKey(roomID), cacheSize)
	rows := []string{}
	if len(rows) == 0 {
		return nil
	}

	cached := make([]gen.BulkInsertMessagesParams, len(rows))
	for i, r := range rows {
		c, err := cache.ToMessageCache(r)
		if err != nil {
			log.Printf("error parsing cache rows")
			return err
		}
		cached[i] = gen.BulkInsertMessagesParams{
			ID:        helper.ToDBUUID(c.ID),
			RoomID:    helper.ToDBUUID(c.RoomID),
			ClientID:  c.ClientID,
			CreatedAt: helper.ToDBTimestamp(c.CreatedAt),
			Content:   c.Content,
			Read:      c.Read,
		}
	}

	count, err := s.r.BulkInsertMessages(ctx, cached)
	if err != nil {
		return err
	}
	if int(count) != len(cached) {
		log.Printf("%d messages given but %d inserted", len(cached), count)
	}

	log.Printf("Persisted %d messages", count)
	return nil
}
