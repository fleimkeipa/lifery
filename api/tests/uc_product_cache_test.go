package tests

import (
	"context"
	"testing"

	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/uc"
)

func TestEventCacheUC_IsExist(t *testing.T) {
	testCache, terminateContainer = pkg.GetCacheTestInstance(context.Background())
	defer terminateContainer()

	type fields struct {
		repo interfaces.CacheRepository
	}
	type tempData struct {
		barcode string
		brandID int
	}
	type args struct {
		ctx     context.Context
		barcode string
		brandID int
	}
	tests := []struct {
		fields    fields
		args      args
		name      string
		tempDatas []tempData
		want      bool
	}{
		{
			name: "correct - not exist",
			fields: fields{
				repo: repositories.NewCacheRepository(testCache),
			},
			args: args{
				ctx:     context.Background(),
				brandID: 123,
				barcode: "abc",
			},
			want: false,
		},
		{
			name: "correct - exist",
			fields: fields{
				repo: repositories.NewCacheRepository(testCache),
			},
			tempDatas: []tempData{
				{
					brandID: 123,
					barcode: "abc",
				},
			},
			args: args{
				ctx:     context.Background(),
				brandID: 123,
				barcode: "abc",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tempData := range tt.tempDatas {
				cacheID := uc.EventCacheID(tempData.brandID, tempData.barcode)
				addTestCacheData(tt.args.ctx, cacheID, tempData.barcode)
			}
			rc := uc.NewEventCacheUC(tt.fields.repo)
			if got := rc.IsExist(tt.args.ctx, tt.args.brandID, tt.args.barcode); got != tt.want {
				t.Errorf("EventCacheUC.IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
