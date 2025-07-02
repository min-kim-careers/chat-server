package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/db"
	"chat-server/internal/db/gen"
	"chat-server/internal/repo"
	"context"
)

type Services struct {
	Message *MessageService
	Room    *RoomService
}

func NewServices() *Services {
	ctx := context.Background()

	db := db.NewDB(ctx)
	q := gen.New(db.DBPool)
	r := repo.NewRepos(q)
	c := cache.NewCache(ctx)

	return &Services{
		Message: NewMessageService(r.Message, db, c),
		Room:    NewRoomService(r.Room, db, c),
	}
}
