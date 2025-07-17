package cache

import (
	"encoding/json"
)

func StreamToCacheMessage(s map[string]any) (*CacheMessage, error) {
	p, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var c CacheMessage
	err = json.Unmarshal(p, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
