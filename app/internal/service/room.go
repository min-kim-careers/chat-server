package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/db"
	"chat-server/internal/db/gen"
	"chat-server/internal/dto/roomout"
	"chat-server/internal/helper"
	"chat-server/internal/repo"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoomService struct {
	r  *repo.RoomRepo
	db *db.DB
	c  *cache.Cache
}

func NewRoomService(r *repo.RoomRepo, db *db.DB, c *cache.Cache) *RoomService {
	return &RoomService{
		r:  r,
		db: db,
		c:  c,
	}
}

func (s *RoomService) RegisterRoom(ctx context.Context, itemID string, client1 uuid.UUID, client2 uuid.UUID) (*roomout.RoomOut, error) {
	if itemID == "" || client1 == uuid.Nil || client2 == uuid.Nil || client1 == client2 {
		return nil, errors.New("invalid params")
	}

	tx, err := s.db.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	client1, client2 = sortClientIds(client1, client2)

	row, err := s.r.GetRoomByItemAndClients(ctx, gen.GetRoomByItemAndClientsParams{
		ItemID:  itemID,
		Client1: helper.ToDBUUID(client1),
		Client2: helper.ToDBUUID(client2),
	})
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		dto := dbToRoomOut(row)
		log.Printf("Existing room found: <%s>", dto.ID)
		return dto, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	row, err = s.r.CreateRoom(ctx, gen.CreateRoomParams{
		ItemID:  itemID,
		Client1: helper.ToDBUUID(client1),
		Client2: helper.ToDBUUID(client2),
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	dto := dbToRoomOut(row)

	log.Printf("Registered room ID: <%s>", dto.ID)
	return dto, nil
}

func (s *RoomService) GetRoomByIdAndClient(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) (*roomout.RoomOut, error) {
	row, err := s.r.GetRoomByIdAndClient(ctx, gen.GetRoomByIdAndClientParams{
		ID:       helper.ToDBUUID(roomID),
		ClientID: helper.ToDBUUID(clientID),
	})
	if err != nil {
		return nil, err
	}
	dto := dbToRoomOut(row)
	return dto, nil
}

func (s *RoomService) GetAllRoomsByClient(ctx context.Context, roomID uuid.UUID) ([]*roomout.RoomOut, error) {
	rows, err := s.r.GetAllRoomsByClient(ctx, helper.ToDBUUID(roomID))
	if err != nil {
		return nil, err
	}
	dtos := make([]*roomout.RoomOut, len(rows))
	for i, r := range rows {
		dtos[i] = dbToRoomOut(r)
	}
	return dtos, nil
}

func (s *RoomService) GetRoomById(ctx context.Context, roomID uuid.UUID) (*roomout.RoomOut, error) {
	room, err := s.r.GetRoomById(ctx, helper.ToDBUUID(roomID))
	if err != nil {
		return nil, err
	}

	return dbToRoomOut(room), nil
}

func (s *RoomService) DeleteRoomById(ctx context.Context, roomID uuid.UUID) error {
	err := s.r.DeleteRoomById(ctx, helper.ToDBUUID(roomID))
	if err != nil {
		return err
	}

	return nil
}
