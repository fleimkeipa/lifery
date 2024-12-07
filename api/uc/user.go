package uc

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
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
		return nil, pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()

	newUser, err := rc.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, pkg.NewError(err, "failed to create user", http.StatusInternalServerError)
	}

	return newUser, nil
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
		return nil, pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	updatedUser, err := rc.userRepo.Update(ctx, userID, &user)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update user", http.StatusInternalServerError)
	}

	return updatedUser, nil
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
	list, err := rc.userRepo.List(ctx, opts)
	if err != nil {
		return nil, pkg.NewError(err, "failed to list users", http.StatusInternalServerError)
	}

	return list, nil
}

func (rc *UserUC) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := rc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, pkg.NewError(err, "user not found", http.StatusNotFound)
	}

	return user, nil
}

func (rc *UserUC) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	user, err := rc.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return nil, pkg.NewError(err, "user not found", http.StatusNotFound)
	}

	return user, nil
}

func (rc *UserUC) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	exists, err := rc.userRepo.Exists(ctx, usernameOrEmail)
	if err != nil {
		return false, pkg.NewError(err, "failed to get user by username or email", http.StatusInternalServerError)
	}

	return exists, nil
}

func (rc *UserUC) Delete(ctx context.Context, userID string) error {
	if err := rc.userRepo.Delete(ctx, userID); err != nil {
		return pkg.NewError(err, "failed to delete user", http.StatusInternalServerError)
	}

	return nil
}

func (rc *UserUC) addConnect(ctx context.Context, user *model.User, senderID, receiverID string) (*model.User, error) {
	sID, err := strconv.Atoi(senderID)
	if err != nil {
		return nil, pkg.NewError(err, "failed to convert string to int", http.StatusBadRequest)
	}

	user.Connects = append(user.Connects, int(sID))

	updatedUser, err := rc.userRepo.Update(ctx, receiverID, user)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update user", http.StatusInternalServerError)
	}

	return updatedUser, nil
}

func (rc *UserUC) deleteConnect(ctx context.Context, user *model.User, userID string) (*model.User, error) {
	for i, v := range user.Connects {
		if strconv.Itoa(v) == userID {
			user.Connects = append(user.Connects[:i], user.Connects[i+1:]...)
			break
		}
	}

	updatedUser, err := rc.userRepo.Update(ctx, userID, user)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update user", http.StatusInternalServerError)
	}

	return updatedUser, nil
}
