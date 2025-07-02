package handler

import (
	"chat-server/internal/auth"
	"chat-server/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func mockVerifyClient() *uuid.UUID {
	id, err := uuid.Parse("dd1dca96-45c6-4c39-b652-7cc66cc59e09")
	if err != nil {
		return nil
	}
	return &id

}

func verifyClient(r *http.Request) *uuid.UUID {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("Verification token missing")
		return nil
	}

	resp, ok := auth.VerifyToken(token)
	if !ok {
		return nil
	}

	return &resp.UserID
}

func isAuthorised(ctx context.Context, roomService *service.RoomService, userID uuid.UUID, roomID uuid.UUID) bool {
	room, err := roomService.GetRoomById(ctx, roomID)
	if err != nil {
		return false
	}
	if room.Client1 != userID && room.Client2 != userID {
		log.Printf("Unauthorised user <%s>", userID)
		return false
	}

	log.Printf("User <%s> authorised", userID)
	return true
}
