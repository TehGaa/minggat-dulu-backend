package database

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	syncOnce    sync.Once
)

func GetRedisClient(redisUri string) *redis.Client {
	syncOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opt, err := redis.ParseURL(redisUri)
		if err != nil {
			log.Fatalf("Failed to parse Redis URI: %v", err)
		}

		redisClient = redis.NewClient(opt)

		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Println("Successfully connected to Redis")
	})

	return redisClient
}
