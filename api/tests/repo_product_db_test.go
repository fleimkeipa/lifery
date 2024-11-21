package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories"

	"github.com/go-pg/pg"
)

func TestEventDBRepository_Create(t *testing.T) {
	testDB, terminateContainer = pkg.GetTestInstance(context.Background())
	defer terminateContainer()

	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx   context.Context
		event *model.Event
	}
	tests := []struct {
		args    args
		fields  fields
		want    *model.Event
		name    string
		wantErr bool
	}{
		{
			name: "correct",
			fields: fields{
				db: testDB,
			},
			args: args{
				ctx: context.TODO(),
				event: &model.Event{
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
			rc := repositories.NewEventRepository(testDB)
			got, err := rc.Create(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventDBRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventDBRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
