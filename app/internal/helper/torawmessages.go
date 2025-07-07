package helper

import (
	"encoding/json"
)

func ToRawMessage[T any](arr []T) (json.RawMessage, error) {
	b, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
