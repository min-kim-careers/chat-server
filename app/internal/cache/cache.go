package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(ctx context.Context) *Cache {
	return &Cache{
		Client: initCacheClient(ctx),
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

	createPersistStreamGroup(client, ctx)

	return client
}

func createPersistStreamGroup(c *redis.Client, ctx context.Context) {
	c.XGroupCreateMkStream(ctx, PersistStreamKey(), PersistGroupKey(), "$")
}
