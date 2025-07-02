package repo

import (
	"chat-server/internal/db/gen"
	"context"
	"errors"
)

type MessageRepo struct {
	q *gen.Queries
}

func NewMessageRepo(q *gen.Queries) *MessageRepo {
	return &MessageRepo{q: q}
}

func (r *MessageRepo) GetMessages(ctx context.Context, arg gen.GetAllMessagesBeforeCreatedAtParams) ([]gen.Message, error) {
	if !arg.RoomID.Valid || !arg.CreatedAt.Valid || arg.Limit < 1 {
		return []gen.Message{}, errors.New("invalid params")
	}
	return r.q.GetAllMessagesBeforeCreatedAt(ctx, arg)
}
