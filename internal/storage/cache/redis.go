package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func NewRedisClient(addr, pw string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       db, // use default DB
	})

	// Ping the Redis server to check if it's available
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb
}
