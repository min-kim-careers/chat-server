package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func (c *Cache) PubSub(ctx context.Context, key string) *redis.PubSub {
	pubsub := c.client.Subscribe(ctx, key)
	if pubsub == nil {
		log.Printf("Error: SUBSCRIBE for key <%s>", key)
	}
	return pubsub
}

func (c *Cache) Publish(ctx context.Context, key string, p []byte) bool {
	err := c.client.Publish(ctx, key, p).Err()
	if err != nil {
		log.Printf("Error: PUBLISH for key <%s>: %v", key, err)
		return false
	}

	log.Printf("Published following data to key <%s>: %s", key, string(p))
	return true
}

func (c *Cache) IsFull(ctx context.Context, key string, limit int64) bool {
	count, err := c.client.LLen(ctx, key).Result()
	if err != nil {
		log.Printf("Error: LLEN for key <%s>: %v", key, err)
		return false
	}

	return count >= limit
}

func (c *Cache) Add(ctx context.Context, key string, p []byte) bool {
	_, err := c.client.RPush(ctx, key, p).Result()
	if err != nil {
		log.Printf("Error: RPUSH for key <%s>: %s", key, err)
		return false
	}

	log.Printf("Cached to key <%s>: %s", key, string(p))
	return true
}

func (c *Cache) Range(ctx context.Context, key string, limit int64) []string {
	res, err := c.client.LRange(ctx, key, -limit, -1).Result()
	if err != nil {
		log.Printf("Error: LRANGE for key <%s>: %v", key, err)
		return nil
	}

	return res
}

func (c *Cache) Clear(ctx context.Context, key string) {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error: DEL for key <%s>: %v", key, err)
	}

	log.Printf("Cleared values for key <%s>", key)
}
