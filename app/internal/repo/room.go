package repo

import (
	"chat-server/db/gen"
)

type RoomRepo struct {
	queries *gen.Queries
}

func NewRoomRepo(queries *gen.Queries) *RoomRepo {
	return &RoomRepo{queries: queries}
}
