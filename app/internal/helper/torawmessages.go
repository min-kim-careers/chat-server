package helper

import (
	"encoding/json"
	"log"
)

func ToRawMessage[T any](arr []T) (json.RawMessage, error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	return json.RawMessage(b), nil
}
