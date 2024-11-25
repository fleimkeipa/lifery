package uc

import (
	"context"
	"strconv"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
)

type UserUC struct {
	userRepo interfaces.UserInterfaces
}

func NewUserUC(repo interfaces.UserInterfaces) *UserUC {
	return &UserUC{
		userRepo: repo,
	}
}

func (rc *UserUC) Create(ctx context.Context, req model.UserCreateRequest) (*model.User, error) {
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()

	return rc.userRepo.Create(ctx, &user)
}

func (rc *UserUC) Update(ctx context.Context, userID string, req model.UserCreateRequest) (*model.User, error) {
	// user exist control
	_, err := rc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return rc.userRepo.Update(ctx, userID, &user)
}

func (rc *UserUC) IsConnected(ctx context.Context, userID, friendID string) (bool, error) {
	// receiver exist control
	receiver, err := rc.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	connects := receiver.Connects
	if connects == nil || len(connects) == 0 {
		return false, nil
	}

	for _, v := range connects {
		if strconv.Itoa(v) == friendID {
			return true, nil
		}
	}

	return false, nil
}

func (rc *UserUC) AddConnect(ctx context.Context, userID, friendID string) (*model.User, error) {
	// receiver exist control
	receiver, err := rc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// sender exist control
	sender, err := rc.GetByID(ctx, friendID)
	if err != nil {
		return nil, err
	}

	receiver, err = rc.addConnect(ctx, receiver, sender.ID, receiver.ID)
	if err != nil {
		return nil, err
	}

	return rc.addConnect(ctx, sender, receiver.ID, sender.ID)
}

func (rc *UserUC) DeleteConnect(ctx context.Context, userID, friendID string) (*model.User, error) {
	// receiver exist control
	receiver, err := rc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// sender exist control
	sender, err := rc.GetByID(ctx, friendID)
	if err != nil {
		return nil, err
	}

	receiver, err = rc.deleteConnect(ctx, receiver, sender.ID)
	if err != nil {
		return nil, err
	}

	return rc.deleteConnect(ctx, sender, receiver.ID)
}

func (rc *UserUC) List(ctx context.Context, opts *model.UserFindOpts) (*model.UserList, error) {
	return rc.userRepo.List(ctx, opts)
}

func (rc *UserUC) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := rc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (rc *UserUC) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	user, err := rc.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (rc *UserUC) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	return rc.userRepo.Exists(ctx, usernameOrEmail)
}

func (rc *UserUC) Delete(ctx context.Context, id string) error {
	return rc.userRepo.Delete(ctx, id)
}

func (rc *UserUC) addConnect(ctx context.Context, user *model.User, senderID, receiverID string) (*model.User, error) {
	sID, err := strconv.Atoi(senderID)
	if err != nil {
		return nil, err
	}

	user.Connects = append(user.Connects, int(sID))

	return rc.userRepo.Update(ctx, receiverID, user)
}

func (rc *UserUC) deleteConnect(ctx context.Context, user *model.User, userID string) (*model.User, error) {
	for i, v := range user.Connects {
		if strconv.Itoa(v) == userID {
			user.Connects = append(user.Connects[:i], user.Connects[i+1:]...)
			break
		}
	}

	return rc.userRepo.Update(ctx, userID, user)
}
