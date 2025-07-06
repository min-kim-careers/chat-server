package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/helper"

	"chat-server/internal/db"
	"chat-server/internal/db/gen"
	"chat-server/internal/dto"
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

func (s *MessageService) GetMessagesDB(ctx context.Context, roomID uuid.UUID, createdAt time.Time, limit int, clientID string) ([]*dto.MessageOutChat, error) {
	if roomID == uuid.Nil || createdAt.IsZero() || limit < 1 {
		return nil, errors.New("invalid params")
	}

	rows, err := s.r.GetMessages(ctx, gen.GetAllMessagesBeforeCreatedAtParams{
		RoomID:    helper.ToDBUUID(roomID),
		CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.MessageOutChat, len(rows))
	for i, r := range rows {
		dtos[len(dtos)-1-i] = dbMessageToDTO(r, clientID)
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(dtos), roomID)
	return dtos, nil
}

func (s *MessageService) BulkInsertMessagesDB(ctx context.Context, chats []*dto.MessageOutChat) error {
	if len(chats) == 0 {
		return errors.New("no messages to insert")
	}

	argMsgs := make([]gen.BulkInsertMessagesParams, len(chats))
	for i, c := range chats {
		argMsgs[i] = gen.BulkInsertMessagesParams{
			RoomID:    helper.ToDBUUID(c.RoomID),
			ClientID:  c.ClientID,
			CreatedAt: helper.ToDBTimestamp(c.CreatedAt),
			Content:   c.Content,
		}
	}

	count, err := s.r.BulkInsertMessages(ctx, argMsgs)
	if err != nil {
		return err
	}
	if int(count) != len(chats) {
		log.Printf("%d messages given but %d inserted", len(chats), count)
	}

	log.Printf("Persisted %d messages", count)
	return nil
}
