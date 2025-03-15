package chat

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var cacheCtx = context.Background()

type Cache struct {
	Client *redis.Client
}

func NewCache() *Cache {
	c := Cache{}
	c.Client = newRedisClient()
	return &c
}

func newRedisClient() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	if redisAddr == "" || redisPassword == "" || redisDBStr == "" {
		log.Fatal("Missing required Redis environment variables")
	}

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Fatal("Invalid REDIS_DB value:", err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
}

func (cache *Cache) CacheMessageHistory(msg *Message) {
	msgJson := SerializeMessage(msg)
	err := cache.Client.Set(cacheCtx, string(msg.RoomID), msgJson, 0)
	if err != nil {
		log.Printf("Failed to cache message history for room <%s>: %v", msg.RoomID, err)
	} else {
		log.Printf("Cached message history for room <%s>.", msg.RoomID)
	}
}

func (cache *Cache) GetMessageHistory(roomID RoomID) (msgHistory []*Message) {
	err := cache.Client.Get(cacheCtx, string(roomID))
	if err != nil {
		log.Printf("Failed to get message history for room <%s> from cache: %v", roomID, err)
	} else {
		log.Printf("Fetched message history for room <%s> from cache.", roomID)
	}
	return make([]*Message, 0)
}
