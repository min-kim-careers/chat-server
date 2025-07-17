package cache

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func BulkStreamToCacheMessages(s []redis.XMessage) ([]CacheMessage, []string, error) {
	cachedMsgs := make([]CacheMessage, len(s))
	streamIDs := make([]string, len(s))
	for i, msg := range s {
		c, err := StreamToCacheMessage(msg.Values)
		if err != nil {
			log.Println("error:", err)
			return nil, nil, err
		}
		cachedMsgs[i] = *c
		streamIDs[i] = msg.ID
	}
	return cachedMsgs, streamIDs, nil
}
