package repo

import (
	"chat-server/db/gen"
)

type MessageRepo struct {
	queries *gen.Queries
}

func NewMessageRepo(queries *gen.Queries) *MessageRepo {
	return &MessageRepo{queries: queries}
}

func (r *MessageRepo) Queries() *gen.Queries {
	return r.queries
}
