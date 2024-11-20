package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const expireTime = 24 * time.Hour

type CacheRepository struct {
	Client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		Client: client,
	}
}

func (rc *CacheRepository) Set(ctx context.Context, key string, value string) error {
	return rc.Client.Set(ctx, key, value, expireTime).Err()
}

func (rc *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return rc.Client.Get(ctx, key).Result()
}

func (rc *CacheRepository) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rc.Client.Exists(ctx, keys...).Result()
}
