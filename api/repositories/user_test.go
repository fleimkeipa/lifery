package repositories

import (
	"reflect"
	"testing"

	"github.com/fleimkeipa/lifery/pkg"

	"github.com/go-pg/pg"
)

func TestUserRepository_getConnects(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				db: pkg.NewPSQLClient(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &UserRepository{
				db: tt.fields.db,
			}
			rc.getConnects()
		})
	}
}

func TestUserRepository_GetMutualConnections(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []connect
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				db: pkg.NewPSQLClient(),
			},
			args: args{
				userID: 3,
			},
			want:    []connect{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &UserRepository{
				db: tt.fields.db,
			}
			got, err := rc.GetMutualConnections(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetMutualConnections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetMutualConnections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_getUsers(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				db: pkg.NewPSQLClient(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &UserRepository{
				db: tt.fields.db,
			}
			rc.getUsers()
		})
	}
}
