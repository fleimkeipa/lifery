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
		repo   interfaces.EventRepository
		userUC *uc.UserUC
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
				repo:   repositories.NewEventRepository(testDB),
				userUC: uc.NewUserUC(repositories.NewUserRepository(testDB)),
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
				repo: repositories.NewEventRepository(testDB),
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
			rc := uc.NewEventUC(tt.fields.repo, tt.fields.userUC)
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
