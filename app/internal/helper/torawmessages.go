package helper

import (
	"encoding/json"
)

func ToRawMessages[T any](arr []T) (json.RawMessage, error) {
	b, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
