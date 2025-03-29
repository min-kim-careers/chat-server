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

const CACHE_LIMIT = 100

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

func generateKey(roomID string) string {
	return "chat:room:" + roomID
}

func (cache *Cache) PubSub(roomID string) *redis.PubSub {
	key := generateKey(roomID)

	pubsub := cache.client.Subscribe(cacheCtx, key)
	if pubsub == nil {
		log.Println("Error subscribing to room: ", key)
	}
	return pubsub
}

func (cache *Cache) PublishMessage(roomID string, msgJson []byte) error {
	key := generateKey(roomID)

	err := cache.client.Publish(cacheCtx, key, msgJson).Err()
	if err != nil {
		log.Printf("Error publishing following message to room <%s>: %v", key, err)
		return err
	}
	return nil
}

func (cache *Cache) CacheFull(roomID string) bool {
	key := generateKey(roomID)

	res, err := cache.client.LLen(cacheCtx, key).Result()
	if err != nil {
		log.Printf("Error retrieving cache length for room <%s>: %v", roomID, err)
		return false
	}

	return res > CACHE_LIMIT
}

func (cache *Cache) CacheMessage(roomID string, msgJson []byte) error {
	key := generateKey(roomID)

	_, err := cache.client.RPush(cacheCtx, key, msgJson).Result()
	if err != nil {
		log.Printf("Error caching message in room <%s>: %s", roomID, err)
		return err
	}
	return nil
}

func (cache *Cache) RestoreMessages(roomID string, limit int64) []*Message {
	key := generateKey(roomID)

	cachedMsgs, err := cache.client.LRange(cacheCtx, key, -limit, -1).Result()
	if err != nil {
		log.Printf("Error getting cached messages for room <%s>: %v", roomID, err)
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

	return ReverseOrder(msgs)
}
