package uc

import (
	"context"
	"net/http"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/util"
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
		return nil, pkg.NewError(nil, "User already exists", http.StatusConflict)
	}

	if req.Password != req.ConfirmPassword {
		return nil, pkg.NewError(nil, "Password and confirm password do not match", http.StatusBadRequest)
	}

	if req.AuthType == "" {
		req.AuthType = string(model.AuthTypeEmail)
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		AuthType: model.AuthType(req.AuthType),
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

	if req.AuthType == "" {
		req.AuthType = string(model.AuthTypeEmail)
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		AuthType: model.AuthType(req.AuthType),
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
	return rc.userRepo.List(ctx, opts)
}

func (rc *UserUC) GetByID(ctx context.Context, id string) (*model.User, error) {
	return rc.userRepo.GetByID(ctx, id)
}

func (rc *UserUC) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	return rc.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
}

func (rc *UserUC) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	return rc.userRepo.Exists(ctx, usernameOrEmail)
}

func (rc *UserUC) Delete(ctx context.Context, userID string) error {
	return rc.userRepo.Delete(ctx, userID)
}

func (rc *UserUC) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return rc.userRepo.GetByEmail(ctx, email)
}

func (rc *UserUC) UpdatePassword(ctx context.Context, userID string, newPassword string) error {
	hashedPassword, err := model.HashPassword(newPassword)
	if err != nil {
		return pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}

	return rc.userRepo.UpdatePassword(ctx, userID, hashedPassword)
}

func (rc *UserUC) UpdateUsername(ctx context.Context, newUsername string) error {
	userID := util.GetOwnerIDFromCtx(ctx)
	if userID == "" {
		return pkg.NewError(nil, "User not authenticated", http.StatusUnauthorized)
	}

	// Check if username already exists
	exists, err := rc.Exists(ctx, newUsername)
	if err != nil {
		return err
	}

	if exists {
		return pkg.NewError(nil, "Username already exists", http.StatusConflict)
	}

	// Get current user to update
	user, err := rc.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Update username
	user.Username = newUsername

	// Update in repository
	_, err = rc.userRepo.Update(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

func (rc *UserUC) UpdatePasswordWithCurrent(ctx context.Context, currentPassword string, newPassword string) error {
	userID := util.GetOwnerIDFromCtx(ctx)
	if userID == "" {
		return pkg.NewError(nil, "User not authenticated", http.StatusUnauthorized)
	}

	// Get current user to verify current password
	user, err := rc.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := model.ValidateUserPassword(user.Password, currentPassword); err != nil {
		return pkg.NewError(err, "Current password is incorrect", http.StatusBadRequest)
	}

	// Update password
	return rc.UpdatePassword(ctx, userID, newPassword)
}
