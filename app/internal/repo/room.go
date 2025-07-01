package repo

import (
	"chat-server/db/gen"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type RoomRepo struct {
	q *gen.Queries
}

func NewRoomRepo(q *gen.Queries) *RoomRepo {
	return &RoomRepo{q: q}
}

func (r *RoomRepo) CreateRoom(ctx context.Context, args gen.CreateRoomParams) (gen.Room, error) {
	return r.q.CreateRoom(ctx, args)
}

func (r *RoomRepo) GetRoomById(ctx context.Context, id pgtype.UUID) (gen.Room, error) {
	return r.q.GetRoomById(ctx, id)
}

func (r *RoomRepo) GetRoomByItemAndClients(ctx context.Context, arg gen.GetRoomByItemAndClientsParams) (gen.Room, error) {
	return r.q.GetRoomByItemAndClients(ctx, arg)
}
