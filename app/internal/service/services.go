package service

import (
	"chat-server/db"
	"chat-server/internal/repo"
)

type Services struct {
	Message *MessageService
	Room    *RoomService
}

func NewServices(r *repo.Repos, db *db.DB) *Services {
	return &Services{
		Message: NewMessageService(r.Message),
		Room:    NewRoomService(r.Room, db),
	}
}
