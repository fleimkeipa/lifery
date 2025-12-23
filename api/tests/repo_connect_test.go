package tests

import (
	"fmt"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/util"
)

// func TestConnectRepository_ConnectsRequests(t *testing.T) {
// 	testDB, terminateContainer = pkg.GetTestInstance(context.Background())
// 	defer terminateContainer()

// 	repositories.NewUserRepository(testDB)
// 	rc := repositories.NewConnectRepository(testDB)

// 	if err := addConnectTempData(); err != nil {
// 		t.Errorf("ConnectRepository.ConnectsRequests() failed to add connect temp data: %v", err)
// 		return
// 	}

// 	type args struct {
// 		ctx  context.Context
// 		opts *model.ConnectFindOpts
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *model.ConnectList
// 		wantErr bool
// 	}{
// 		{
// 			name: "",
// 			args: args{
// 				ctx: context.TODO(),
// 				opts: &model.ConnectFindOpts{
// 					OrderByOpts: model.OrderByOpts{},
// 					Status:      model.Filter{},
// 					UserID: model.Filter{
// 						Value:    "2",
// 						IsSended: true,
// 					},
// 					FieldsOpts:     model.FieldsOpts{},
// 					PaginationOpts: model.PaginationOpts{},
// 				},
// 			},
// 			want:    &model.ConnectList{},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rc.ConnectsRequests(tt.args.ctx, tt.args.opts)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ConnectRepository.ConnectsRequests() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			for _, v := range got.Connects {
// 				fmt.Printf("id: %v, user id: %v, friend id: %v, user username: %v, friend username: %v\n", v.ID, v.UserID, v.FriendID, v.User.Username, v.Friend.Username)
// 			}
// 		})
// 	}
// }

func addConnectTempData() error {
	if err := addUserTempData(); err != nil {
		return fmt.Errorf("failed to add user temp data: %v", err)
	}

	datas := []model.Connect{
		{
			UserID:   "1",
			FriendID: "2",
			Status:   100, // pending
		},
		{
			UserID:   "1",
			FriendID: "3",
			Status:   100, // pending
		},
		{
			UserID:   "1",
			FriendID: "4",
			Status:   100, // pending
		},
		{
			UserID:   "2",
			FriendID: "4",
			Status:   100, // pending
		},
	}

	for _, v := range datas {
		_, err := testDB.Model(&v).Insert()
		if err != nil {
			return err
		}
	}

	return nil
}

func addUserTempData() error {
	datas := []model.User{
		{
			CreatedAt: util.Now(),
			Username:  "user1",
			Email:     "",
			Password:  "password",
			Connects:  []*model.Connect{},
			RoleID:    7,
		},
		{
			CreatedAt: util.Now(),
			Username:  "user2",
			Email:     "",
			Password:  "password",
			Connects:  []*model.Connect{},
			RoleID:    7,
		},
		{
			CreatedAt: util.Now(),
			Username:  "user3",
			Email:     "",
			Password:  "password",
			Connects:  []*model.Connect{},
			RoleID:    7,
		},
		{
			CreatedAt: util.Now(),
			Username:  "user4",
			Email:     "",
			Password:  "password",
			Connects:  []*model.Connect{},
			RoleID:    7,
		},
	}

	for _, v := range datas {
		_, err := testDB.Model(&v).Insert()
		if err != nil {
			return err
		}
	}

	return nil
}
