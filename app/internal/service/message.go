package service

import (
	"chat-server/db/gen"
	"chat-server/internal/dto"
	"chat-server/internal/repo"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageService struct {
	r *repo.MessageRepo
}

func NewMessageService(r *repo.MessageRepo) *MessageService {
	return &MessageService{r: r}
}

func (s *MessageService) SendMessage(ctx context.Context, input dto.Message) {

}

func ToMessageDto(m gen.Message) dto.Message {
	return dto.Message{
		ID:        int(m.ID),
		Mode:      m.Mode,
		RoomID:    m.RoomID.Bytes,
		ClientID:  m.ClientID,
		CreatedAt: m.CreatedAt.Time,
		Data:      m.Data,
	}
}

func (s *MessageService) GetPreviousMessages(ctx context.Context, roomID uuid.UUID, createdAt time.Time, limit int) ([]dto.Message, error) {
	rows, err := s.r.GetPreviousMessages(ctx, gen.GetAllMessagesBeforeCreatedAtParams{
		RoomID:    pgtype.UUID{Bytes: roomID, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.Message, len(rows))
	for i, m := range rows {
		dtos[i] = ToMessageDto(m)
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(dtos), roomID)
	return dtos, nil
}
