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
	Client *redis.Client
}

func NewCache() *Cache {
	cache := Cache{}
	cache.Client = initCacheClient()
	return &cache
}

func initCacheClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	if redisHost == "" || redisPort == "" || redisPassword == "" || redisDBStr == "" {
		log.Fatal("Missing required Redis environment variables")
	}

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Fatal("Invalid REDIS_DB value:", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	_, err = client.Ping(cacheCtx).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}
	log.Println("Connected to Redis.")

	return client
}

func (cache *Cache) SetMessageHistory(msg *Message) {
	msgJson := SerializeMessage(msg)
	err := cache.Client.Set(cacheCtx, string(msg.RoomID), msgJson, 0).Err()
	if err != nil {
		log.Printf("Failed to cache message history for room <%s>: %v", msg.RoomID, err)
	} else {
		log.Printf("Cached message history for room <%s>.", msg.RoomID)
	}
}

func (cache *Cache) GetMessageHistory(roomID RoomID) ([]*Message, error) {
	val, err := cache.Client.Get(cacheCtx, string(roomID)).Result()
	if err != nil {
		log.Printf("Failed to get message history for room <%s> from cache: %v", roomID, err)
		return nil, err
	}

	var history []*Message

	err = json.Unmarshal([]byte(val), &history)
	if err != nil {
		log.Printf("Failed to unmarshal message history for room <%s>: %v", roomID, err)
		return nil, err
	}

	log.Printf("Fetched message history for room <%s> from cache.", roomID)
	return history, nil
}
