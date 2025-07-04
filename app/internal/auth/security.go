package auth

import (
	"chat-server/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func VerifyClient(r *http.Request) *uuid.UUID {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("Verification token missing")
		return nil
	}

	resp, ok := VerifyToken(token)
	if !ok {
		return nil
	}

	return &resp.UserID
}

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
