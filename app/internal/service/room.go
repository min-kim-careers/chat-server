package service

import (
	"chat-server/db"
	"chat-server/db/gen"
	"chat-server/internal/dto"
	"chat-server/internal/repo"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomService struct {
	r  *repo.RoomRepo
	db *db.DB
}

func NewRoomService(r *repo.RoomRepo, db *db.DB) *RoomService {
	return &RoomService{
		r:  r,
		db: db,
	}
}

func (s *RoomService) RegisterRoom(ctx context.Context, itemID string, client1 uuid.UUID, client2 uuid.UUID) (*dto.Room, error) {
	tx, err := s.db.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	client1, client2 = sortClientIds(client1, client2)

	row, err := s.r.GetRoomByItemAndClients(ctx, gen.GetRoomByItemAndClientsParams{
		ItemID:  itemID,
		Client1: pgtype.UUID{Bytes: client1, Valid: true},
		Client2: pgtype.UUID{Bytes: client2, Valid: true},
	})
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		dto := toRoomDTO(row)
		return dto, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	row, err = s.r.CreateRoom(ctx, gen.CreateRoomParams{
		ItemID:  itemID,
		Client1: pgtype.UUID{Bytes: client1, Valid: true},
		Client2: pgtype.UUID{Bytes: client2, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	dto := toRoomDTO(row)

	log.Printf("Registered room ID: <%s>", dto.ID)
	return dto, nil
}

func (s *RoomService) GetRoomById(ctx context.Context, roomID uuid.UUID) (*dto.Room, error) {
	row, err := s.r.GetRoomById(ctx, pgtype.UUID{Bytes: roomID, Valid: true})
	if err != nil {
		return nil, err
	}
	dto := toRoomDTO(row)
	return dto, nil
}
