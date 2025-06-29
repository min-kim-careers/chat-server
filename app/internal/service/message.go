package service

import (
	"chat-server/db/gen"
	"chat-server/internal/dto"
	"chat-server/internal/repo"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type MessageService struct {
	Repo *repo.MessageRepo
}

func NewMessageService(r *repo.MessageRepo) *MessageService {
	return &MessageService{Repo: r}
}

func (s *MessageService) SendMessage(ctx context.Context, input dto.Message) {

}

func ToMessageDto(m gen.Message) dto.Message {
	return dto.Message{
		ID:          int(m.ID),
		MessageType: m.MessageType,
		RoomID:      m.RoomID,
		ClientID:    m.ClientID,
		CreatedAt:   m.CreatedAt.Time,
		Data:        m.Data,
	}
}

type GetPreviousMessagesParams struct {
	RoomID    string    `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	Limit     int       `json:"limit"`
}

func (s *MessageService) GetPreviousMessages(ctx context.Context, arg GetPreviousMessagesParams) ([]dto.Message, error) {
	params := gen.GetPreviousMessagesParams{
		RoomID:    arg.RoomID,
		CreatedAt: pgtype.Timestamp{Time: arg.CreatedAt, Valid: true},
		Limit:     int32(arg.Limit),
	}

	rows, err := s.Repo.Queries().GetPreviousMessages(ctx, params)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.Message, len(rows))
	for i, m := range rows {
		dtos[i] = ToMessageDto(m)
	}

	log.Printf("Fetched %d messages from DB for room <%s>.", len(dtos), arg.RoomID)
	return dtos, nil
}
