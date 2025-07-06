package auth

import (
	"chat-server/internal/service"
	"context"
	"log"

	"github.com/google/uuid"
)

func IsAuthorised(ctx context.Context, roomService *service.RoomService, clientID string, roomID uuid.UUID) error {
	_clientID, err := uuid.Parse(clientID)
	if err != nil {
		return err
	}

	_, err = roomService.GetRoomByIdAndClient(ctx, roomID, _clientID)
	if err != nil {
		return err
	}

	log.Printf("User <%s> authorised", _clientID)
	return nil
}
