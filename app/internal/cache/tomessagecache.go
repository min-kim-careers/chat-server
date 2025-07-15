package cache

import (
	"encoding/json"
)

func ToMessageCache(c string) (*MessageCache, error) {
	var m *MessageCache
	err := json.Unmarshal([]byte(c), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
