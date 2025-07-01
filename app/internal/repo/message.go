package repo

import (
	"chat-server/db/gen"
	"context"
)

type MessageRepo struct {
	q *gen.Queries
}

func NewMessageRepo(q *gen.Queries) *MessageRepo {
	return &MessageRepo{q: q}
}

func (r *MessageRepo) GetPreviousMessages(ctx context.Context, arg gen.GetAllMessagesBeforeCreatedAtParams) ([]gen.Message, error) {
	return r.q.GetAllMessagesBeforeCreatedAt(ctx, arg)
}
