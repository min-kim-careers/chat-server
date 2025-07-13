package repo

import (
	"chat-server/internal/db/gen"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type RoomRepo struct {
	q *gen.Queries
}

func NewRoomRepo(q *gen.Queries) *RoomRepo {
	return &RoomRepo{q: q}
}

func (r *RoomRepo) CreateRoom(ctx context.Context, arg gen.CreateRoomParams) (gen.Room, error) {
	if arg.ItemID == "" || !arg.Client1.Valid || !arg.Client2.Valid {
		return gen.Room{}, errors.New("invalid arg")
	}

	return r.q.CreateRoom(ctx, arg)
}

func (r *RoomRepo) GetRoomByIdAndClient(ctx context.Context, arg gen.GetRoomByIdAndClientParams) (gen.Room, error) {
	if !arg.ID.Valid || !arg.ClientID.Valid {
		return gen.Room{}, errors.New("invalid arg")
	}
	return r.q.GetRoomByIdAndClient(ctx, arg)
}

func (r *RoomRepo) GetRoomByItemAndClients(ctx context.Context, arg gen.GetRoomByItemAndClientsParams) (gen.Room, error) {
	if arg.ItemID == "" || !arg.Client1.Valid || !arg.Client2.Valid {
		return gen.Room{}, errors.New("invalid arg")
	}

	return r.q.GetRoomByItemAndClients(ctx, arg)
}

func (r *RoomRepo) GetAllRoomsByClient(ctx context.Context, clientID pgtype.UUID) ([]gen.Room, error) {
	if !clientID.Valid {
		return []gen.Room{}, errors.New("invalid arg")
	}

	return r.q.GetAllRoomsByClient(ctx, clientID)
}

func (r *RoomRepo) GetRoomById(ctx context.Context, id pgtype.UUID) (gen.Room, error) {
	if !id.Valid {
		return gen.Room{}, errors.New("invalid arg")
	}
	return r.q.GetRoomById(ctx, id)
}

func (r *RoomRepo) DeleteRoomById(ctx context.Context, id pgtype.UUID) error {
	if !id.Valid {
		return errors.New("invalid arg")
	}
	return r.q.DeleteRoomById(ctx, id)
}
