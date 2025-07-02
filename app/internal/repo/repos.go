package repo

import (
	"chat-server/internal/db/gen"
)

type Repos struct {
	Message *MessageRepo
	Room    *RoomRepo
}

func NewRepos(queries *gen.Queries) *Repos {
	return &Repos{
		Message: NewMessageRepo(queries),
		Room:    NewRoomRepo(queries),
	}
}
