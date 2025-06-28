package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var authUrl = os.Getenv("AUTH_URL")
var authKey = os.Getenv("AUTH_KEY")
var authKeyHeader = os.Getenv("AUTH_KEY_HEADER")

type AuthResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type VerifyTokenRequest struct {
	IDToken string `json:"idToken"`
}

type VerifyTokenResponse struct {
	UserID string `json:"userId"`
}

func VerifyToken(idToken string) (*VerifyTokenResponse, bool) {
	client := &http.Client{Timeout: 5 * time.Second}

	reqBody, err := json.Marshal(VerifyTokenRequest{IDToken: idToken})
	if err != nil {
		return nil, false
	}

	req, err := http.NewRequest("POST", authUrl+"/auth/firebase/verify-token", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authKeyHeader, authKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Token verification error:", err)
		return nil, false
	}
	defer resp.Body.Close()

	var result AuthResponse[VerifyTokenResponse]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding response:", err)
		return nil, false
	}

	log.Printf("Token verified: %+v", result)
	return &result.Data, true
}
