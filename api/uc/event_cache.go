package uc

import (
	"context"
	"fmt"

	"github.com/fleimkeipa/lifery/repositories/interfaces"
)

type EventCacheUC struct {
	repo interfaces.CacheRepository
}

func NewEventCacheUC(repo interfaces.CacheRepository) *EventCacheUC {
	return &EventCacheUC{
		repo: repo,
	}
}

func (rc *EventCacheUC) IsExist(ctx context.Context, brandID int, barcode string) bool {
	cacheID := EventCacheID(brandID, barcode)

	count, err := rc.repo.Exists(ctx, cacheID)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	}

	go func(ctx context.Context, cacheID string, barcode string) {
		if err := rc.repo.Set(ctx, cacheID, barcode); err != nil {
			return
		}
	}(ctx, cacheID, barcode)

	return false
}

func EventCacheID(brandID int, barcode string) string {
	return fmt.Sprintf("event:%d:%v", brandID, barcode)
}
