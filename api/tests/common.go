package tests

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/redis/go-redis/v9"
)

var (
	testDB             *pg.DB
	testCache          *redis.Client
	terminateContainer = func() {}
)

func addTestCacheData(ctx context.Context, key string, value string) {
	err := testCache.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
