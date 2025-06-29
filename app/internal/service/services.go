package service

import (
	"chat-server/internal/repo"
)

type Services struct {
	MessageService *MessageService
}

func NewServices(r *repo.Repos) *Services {
	return &Services{
		MessageService: NewMessageService(r.Message),
	}
}
