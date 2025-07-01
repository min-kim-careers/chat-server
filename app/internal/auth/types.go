package auth

import "github.com/google/uuid"

type AuthResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type VerifyTokenRequest struct {
	IDToken string `json:"idToken"`
}

type VerifyTokenResponse struct {
	UserID uuid.UUID `json:"userId"`
}
