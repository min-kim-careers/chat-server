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

type GetDBMessagesParams struct {
	RoomID    uuid.UUID
	CreatedAt time.Time
	Limit     int
	ClientID  string
}

func (s *MessageService) GetChatMessagesFromDB(ctx context.Context, arg GetDBMessagesParams) ([]*messageout.MessageOutChat, error) {
	if arg.RoomID == uuid.Nil || arg.CreatedAt.IsZero() || arg.Limit < 1 {
		return nil, errors.New("invalid params")
	}

	rows, err := s.r.GetMessages(ctx, gen.GetMessagesBeforeCreatedAtParams{
		RoomID:    helper.ToDBUUID(arg.RoomID),
		CreatedAt: pgtype.Timestamp{Time: arg.CreatedAt, Valid: true},
		Limit:     int32(arg.Limit),
	})
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	dtos := make([]*messageout.MessageOutChat, len(rows))
	for i, r := range rows {
		dtos[len(dtos)-1-i] = dbToMessageOutChat(r, arg.ClientID)
	}

	log.Printf("fetched %d from db", len(dtos))
	return dtos, nil
}

func (s *MessageService) FlushCacheBatchMessagesToDB(ctx context.Context, cachedMsgs []cache.CacheMessage) error {
	if len(cachedMsgs) == 0 {
		return nil
	}

	cached := make([]gen.BulkInsertMessagesParams, len(cachedMsgs))
	for i, c := range cachedMsgs {
		cached[i] = gen.BulkInsertMessagesParams{
			ID:        helper.ToDBUUID(c.ID),
			RoomID:    helper.ToDBUUID(c.RoomID),
			ClientID:  c.ClientID,
			CreatedAt: helper.ToDBTimestamp(c.CreatedAt),
			Content:   c.Content,
			Read:      false,
		}
	}

	count, err := s.r.BulkInsertMessages(ctx, cached)
	if err != nil {
		log.Println("error:", err)
		return err
	}
	if int(count) != len(cached) {
		log.Printf("%d messages given but %d inserted", len(cached), count)
	}

	log.Printf("persisted %d messages", count)
	return nil
}
