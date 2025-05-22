package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/event-mgt-api/internal/storage"
)

type EventCache struct {
	rdb *redis.Client
}

func (u *EventCache) Get(ctx context.Context, id int) (*storage.Event, error) {
	cacheKey := "event:" + strconv.Itoa(id)

	data, err := u.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var event storage.Event

	if data != "" {
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, err
		}
	}

	return &event, nil
}

func (u *EventCache) Set(ctx context.Context, event *storage.Event) error {
	CacheKey := "event:" + strconv.Itoa(event.ID)
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return u.rdb.Set(ctx, CacheKey, data, timeExp).Err()
}

func (u *EventCache) Delete(ctx context.Context, id int) {
	cacheKey := "user:" + strconv.Itoa(id)
	u.rdb.Del(ctx, cacheKey)
}
