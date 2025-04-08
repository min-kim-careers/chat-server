package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var cacheCtx = context.Background()

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	return &Cache{
		client: initCacheClient(),
	}
}

func initCacheClient() *redis.Client {
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

	_, err = client.Ping(cacheCtx).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}
	log.Println("Connected to Redis.")

	return client
}

func generateCacheKey(roomID string) string {
	return "chat:room:" + roomID
}

func (cache *Cache) PubSub(roomID string) *redis.PubSub {
	key := generateCacheKey(roomID)

	pubsub := cache.client.Subscribe(cacheCtx, key)
	if pubsub == nil {
		log.Println("Error subscribing to room: ", key)
	}
	return pubsub
}

func (cache *Cache) Publish(roomID string, msgJson []byte) error {
	key := generateCacheKey(roomID)

	err := cache.client.Publish(cacheCtx, key, msgJson).Err()
	if err != nil {
		log.Printf("Error publishing following message to room <%s>: %v", key, err)
		return err
	}
	return nil
}

func (cache *Cache) IsFull(roomID string) bool {
	key := generateCacheKey(roomID)

	res, err := cache.client.LLen(cacheCtx, key).Result()
	if err != nil {
		log.Printf("Error retrieving cache length for room <%s>: %v", roomID, err)
		return false
	}

	return res > CACHE_LIMIT
}

func (cache *Cache) Add(roomID string, msgJson []byte) error {
	key := generateCacheKey(roomID)

	_, err := cache.client.RPush(cacheCtx, key, msgJson).Result()
	if err != nil {
		log.Printf("Error caching message in room <%s>: %s", roomID, err)
		return err
	}
	return nil
}

func (cache *Cache) Restore(roomID string, limit int64) []*Message {
	key := generateCacheKey(roomID)

	cachedMsgs, err := cache.client.LRange(cacheCtx, key, -limit, -1).Result()
	if err != nil {
		log.Printf("Error restoring cached messages for room <%s>: %v", roomID, err)
		return nil
	}

	var msgs []*Message

	for _, cachedMsg := range cachedMsgs {
		var msg Message
		err = json.Unmarshal([]byte(cachedMsg), &msg)
		if err != nil {
			log.Println("Error unmarshalling cached messages:", err)
			return nil
		}
		msgs = append(msgs, &msg)
	}

	log.Printf("Fetched %d messages from cache for room <%s>.", len(msgs), roomID)
	return ReverseOrder(msgs)
}

func (cache *Cache) Clear(roomID string) error {
	key := generateCacheKey(roomID)

	err := cache.client.Del(cacheCtx, key).Err()
	if err != nil {
		log.Printf("Error clearing messages for room <%s>: %v", roomID, err)
		return err
	}

	return nil
}
