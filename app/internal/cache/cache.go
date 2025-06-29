package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"chat-server/internal/dto"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(ctx context.Context) *Cache {
	return &Cache{
		client: initCacheClient(ctx),
	}
}

func initCacheClient(ctx context.Context) *redis.Client {
	cacheHost := os.Getenv("CACHE_HOST")
	cachePort := os.Getenv("CACHE_PORT")
	cachePassword := os.Getenv("CACHE_PASSWORD")
	cacheDBStr := os.Getenv("CACHE_DB")

	if cacheHost == "" || cachePort == "" || cachePassword == "" || cacheDBStr == "" {
		log.Fatal("Missing required Redis environment variables")
	}

	redisDB, err := strconv.Atoi(cacheDBStr)
	if err != nil {
		log.Fatalf("Invalid %s value: %v", cacheDBStr, err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cacheHost, cachePort),
		Password: cachePassword,
		DB:       redisDB,
	})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not connect to cache server:", err)
	}
	log.Println("Connected to cache server.")

	return client
}

func roomKey(roomID string) string {
	return "chat:room:" + roomID
}

func (c *Cache) PubSub(ctx context.Context, roomID string) *redis.PubSub {
	key := roomKey(roomID)

	pubsub := c.client.Subscribe(ctx, key)
	if pubsub == nil {
		log.Println("Error subscribing to room: ", key)
	}
	return pubsub
}

func (c *Cache) Publish(ctx context.Context, roomID string, msgJson []byte) bool {
	key := roomKey(roomID)

	err := c.client.Publish(ctx, key, msgJson).Err()
	if err != nil {
		log.Printf("Error publishing following message to room <%s>: %v", key, err)
		return false
	}

	log.Printf("Published following message to room <%s>: %v", key, err)
	return true
}

func (c *Cache) IsFull(ctx context.Context, roomID string, cacheLimit int64) bool {
	key := roomKey(roomID)

	count, err := c.client.LLen(ctx, key).Result()
	if err != nil {
		log.Printf("Error checking if cache is full <%s>: %v", roomID, err)
		return false
	}

	return count >= cacheLimit
}

func (c *Cache) Add(ctx context.Context, roomID string, msgJson []byte) bool {
	key := roomKey(roomID)

	_, err := c.client.RPush(ctx, key, msgJson).Result()
	if err != nil {
		log.Printf("Error caching message in room <%s>: %s", roomID, err)
		return false
	}

	log.Printf("Cached message in room <%s>.", roomID)
	return true
}

func (c *Cache) Restore(ctx context.Context, roomID string, limit int64) []*dto.Message {
	key := roomKey(roomID)

	cachedMsgs, err := c.client.LRange(ctx, key, -limit, -1).Result()
	if err != nil {
		log.Printf("Error restoring cached messages for room <%s>: %v", roomID, err)
		return nil
	}

	msgs := []*dto.Message{}

	for _, cachedMsg := range cachedMsgs {
		var msg dto.Message
		err = json.Unmarshal([]byte(cachedMsg), &msg)
		if err != nil {
			log.Println("Error unmarshalling cached messages:", err)
			return nil
		}
		msgs = append(msgs, &msg)
	}

	log.Printf("Fetched %d messages from cache for room <%s>.", len(msgs), roomID)
	return msgs
}

func (c *Cache) Clear(ctx context.Context, roomID string) {
	key := roomKey(roomID)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error clearing messages for room <%s>: %v", roomID, err)
	}

	log.Printf("Cleared messages for room <%s>: %v", roomID, err)
}
