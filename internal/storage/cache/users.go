package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/event-mgt-api/internal/storage"
)

type UserCache struct {
	rdb *redis.Client
}

const timeExp = time.Minute * 2

func (u *UserCache) Get(ctx context.Context, id int) (*storage.User, error) {
	cacheKey := "user:" + strconv.Itoa(id)

	data, err := u.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user storage.User

	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserCache) Set(ctx context.Context, user *storage.User) error {
	CacheKey := "user:" + strconv.Itoa(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return u.rdb.Set(ctx, CacheKey, data, timeExp).Err()
}

func (u *UserCache) Delete(ctx context.Context, id int) {
	cacheKey := "user:" + strconv.Itoa(id)
	u.rdb.Del(ctx, cacheKey)
}
