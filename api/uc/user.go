package uc

import (
	"context"
	"net/http"
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

func (rc *UserUC) Create(ctx context.Context, req model.UserCreateInput) (*model.User, error) {
	exists, err := rc.Exists(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, pkg.NewError(nil, "User already exists", http.StatusBadRequest)
	}

	if req.Password != req.ConfirmPassword {
		return nil, pkg.NewError(nil, "Password and confirm password do not match", http.StatusBadRequest)
	}

	user := model.User{
		Username: req.Username,
		// Email:    req.Email,
		Password: req.Password,
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		return nil, pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()

	if user.RoleID <= 0 {
		user.RoleID = model.EditorRole
	}

	newUser, err := rc.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (rc *UserUC) Update(ctx context.Context, userID string, req model.UserCreateInput) (*model.User, error) {
	// user exist control
	_, err := rc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		return nil, pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	updatedUser, err := rc.userRepo.Update(ctx, userID, &user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (rc *UserUC) List(ctx context.Context, opts *model.UserFindOpts) (*model.UserList, error) {
	list, err := rc.userRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list, nil
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
	exists, err := rc.userRepo.Exists(ctx, usernameOrEmail)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (rc *UserUC) Delete(ctx context.Context, userID string) error {
	if err := rc.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	return nil
}
