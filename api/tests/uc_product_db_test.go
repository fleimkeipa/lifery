package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/uc"
)

func TestEventDBUC_Create(t *testing.T) {
	testCache, _ = pkg.GetCacheTestInstance(context.Background())
	testDB, terminateContainer = pkg.GetTestInstance(context.Background())
	defer terminateContainer()

	type fields struct {
		repo  interfaces.EventRepository
		cache *uc.EventCacheUC
	}
	type tempData struct {
		barcode string
		brandID int
	}
	type args struct {
		ctx   context.Context
		event *model.EventCreateRequest
	}
	tests := []struct {
		fields    fields
		args      args
		want      *model.Event
		name      string
		tempDatas []tempData
		wantErr   bool
	}{
		{
			name: "correct - not exist",
			fields: fields{
				repo:  repositories.NewEventRepository(testDB),
				cache: uc.NewEventCacheUC(repositories.NewCacheRepository(testCache)),
			},
			args: args{
				ctx: context.TODO(),
				event: &model.EventCreateRequest{
					Name: "1234",
				},
			},
			want: &model.Event{
				Name: "1234",
			},
			wantErr: false,
		},
		{
			name: "correct - exist",
			fields: fields{
				repo:  repositories.NewEventRepository(testDB),
				cache: uc.NewEventCacheUC(repositories.NewCacheRepository(testCache)),
			},
			tempDatas: []tempData{
				{
					barcode: "1234",
				},
			},
			args: args{
				ctx: context.TODO(),
				event: &model.EventCreateRequest{
					Name: "1234",
				},
			},
			want: &model.Event{
				Name: "1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.tempDatas {
				cacheID := uc.EventCacheID(v.brandID, v.barcode)
				addTestCacheData(tt.args.ctx, cacheID, v.barcode)
			}
			rc := uc.NewEventUC(tt.fields.repo, tt.fields.cache)
			got, err := rc.Create(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventDBUC.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventDBUC.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
