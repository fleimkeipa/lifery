package interfaces

import "context"

type CacheRepository interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
}
