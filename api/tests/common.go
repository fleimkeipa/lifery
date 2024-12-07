package tests

import (
	"context"
	"fmt"

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

func addTempData(data interface{}) error {
	_, err := testDB.Model(data).Insert()
	if err != nil {
		return err
	}

	return nil
}

func clearTable(tableName string) error {
	query := fmt.Sprintf("TRUNCATE %s; DELETE FROM %s", tableName, tableName)
	_, err := testDB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
