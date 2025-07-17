package cache

import (
	"encoding/json"
	"log"
)

func StreamToCacheMessage(s map[string]any) (*CacheMessage, error) {
	p, err := json.Marshal(s)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	var c CacheMessage
	err = json.Unmarshal(p, &c)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	return &c, nil
}
