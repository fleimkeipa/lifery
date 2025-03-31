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

func TestEventRepository_List(t *testing.T) {
	testDB, terminateContainer = pkg.GetTestInstance(context.Background())
	defer terminateContainer()

	type tempDatas struct {
		events []model.Event
	}
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx  context.Context
		opts *model.EventFindOpts
	}
	tests := []struct {
		args      args
		fields    fields
		want      *model.EventList
		name      string
		tempDatas tempDatas
		wantErr   bool
	}{
		{
			name: "correct",
			tempDatas: tempDatas{
				events: []model.Event{
					{
						Name:       "1234",
						Visibility: model.Visibility(model.EventVisibilityPublic),
					},
				},
			},
			fields: fields{
				db: testDB,
			},
			args: args{
				ctx:  context.TODO(),
				opts: &model.EventFindOpts{},
			},
			want: &model.EventList{
				Events: []model.Event{
					{
						Name:       "1234",
						Visibility: model.Visibility(model.EventVisibilityPublic),
					},
				},
				Total:          1,
				PaginationOpts: model.PaginationOpts{},
			},
			wantErr: false,
		},
		{
			name: "correct - with opts",
			tempDatas: tempDatas{
				events: []model.Event{
					{
						Name:       "1234",
						Visibility: model.Visibility(model.EventVisibilityPublic),
					},
					{
						Name:       "1234",
						Visibility: model.Visibility(model.EventVisibilityPrivate),
					},
				},
			},
			fields: fields{
				db: testDB,
			},
			args: args{
				ctx:  context.TODO(),
				opts: &model.EventFindOpts{},
			},
			want: &model.EventList{
				Events: []model.Event{
					{
						Name:       "1234",
						Visibility: model.Visibility(model.EventVisibilityPublic),
					},
					{
						Name:       "1234",
						Visibility: model.Visibility(model.EventVisibilityPrivate),
					},
				},
				Total:          1,
				PaginationOpts: model.PaginationOpts{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositories.NewUserRepository(testDB)
			rc := repositories.NewEventRepository(testDB)
			if err := addEventTempData(tt.tempDatas.events); err != nil {
				t.Errorf("addTempData() error = %v", err)
				return
			}
			got, err := rc.List(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Events[0].Name != tt.want.Events[0].Name {
				t.Errorf("EventRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func addEventTempData(data []model.Event) error {
	for _, v := range data {
		_, err := testDB.Model(&v).Insert()
		if err != nil {
			return err
		}
	}

	return nil
}
