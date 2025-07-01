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

func attachServiceKey(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authKeyHeader, authKey)
}

func VerifyToken(idToken string) (*VerifyTokenResponse, bool) {
	client := &http.Client{Timeout: 5 * time.Second}

	reqBody, err := json.Marshal(VerifyTokenRequest{IDToken: idToken})
	if err != nil {
		return nil, false
	}

	body := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest(
		"POST",
		authUrl+"/auth/firebase/verify-token",
		body,
	)
	if err != nil {
		return nil, false
	}

	attachServiceKey(req)

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
