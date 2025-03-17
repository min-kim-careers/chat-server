package chat

import (
	"context"
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

func (cache *Cache) AddToStream(roomID string, msgJson string) {
	key := "chat:room:" + roomID

	_, err := cache.client.XAdd(cacheCtx, &redis.XAddArgs{
		Stream: key,
		Values: map[string]string{
			"message": msgJson,
		},
	}).Result()
	if err != nil {
		log.Printf("Error adding message to stream <%s>: %v", key, err)
		return
	}

	log.Printf("Message sent to stream <%s>.", key)
}

func (cache *Cache) GetReadStreams(roomID string) []redis.XStream {
	key := "chat:room:" + roomID

	streams, err := cache.client.XRead(cacheCtx, &redis.XReadArgs{
		Streams: []string{key, "$"},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		log.Printf("Error reading message from stream <%s>: %v", key, err)
		return nil
	}

	return streams
}
