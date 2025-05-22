package cache

import "github.com/go-redis/redis/v8"

type CacheStorage struct {
	Users  UserCache
	Events EventCache
}

func NewCacheStorage(rdb *redis.Client) *CacheStorage {
	return &CacheStorage{
		Users:  UserCache{rdb},
		Events: EventCache{rdb},
	}
}
