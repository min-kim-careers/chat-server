package handler

import (
	"chat-server/internal/auth"
	"chat-server/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

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

func isAuthorised(ctx context.Context, roomService *service.RoomService, userID uuid.UUID, roomID uuid.UUID) error {
	_, err := roomService.GetRoomByIdAndClient(ctx, roomID, userID)
	if err != nil {
		return err
	}

	log.Printf("User <%s> authorised", userID)
	return nil
}
