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

func (r *MessageRepo) GetMessages(ctx context.Context, arg gen.GetMessagesBeforeCreatedAtParams) ([]gen.Message, error) {
	if !arg.RoomID.Valid || !arg.CreatedAt.Valid || arg.Limit < 1 {
		return []gen.Message{}, errors.New("invalid params")
	}
	return r.q.GetMessagesBeforeCreatedAt(ctx, arg)
}

func (r *MessageRepo) BulkInsertMessages(ctx context.Context, arg []gen.BulkInsertMessagesParams) (int64, error) {
	if len(arg) == 0 {
		return 0, errors.New("invalid params")
	}
	return r.q.BulkInsertMessages(ctx, arg)
}

// if arg.Mode == "" || !arg.RoomID.Valid || arg.ClientID == "" || !arg.CreatedAt.Valid || arg.Data != nil {
// 	return []gen.Message{}, errors.New("invalid params")
// }
