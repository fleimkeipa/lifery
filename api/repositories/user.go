package repositories

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (rc *UserRepository) Create(ctx context.Context, newUser *model.User) (*model.User, error) {
	q := rc.db.Model(newUser)

	if newUser.RoleID <= 0 {
		newUser.RoleID = model.ViewerRole
	}

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (rc *UserRepository) Update(ctx context.Context, userID string, user *model.User) (*model.User, error) {
	if userID == "" || userID == "0" {
		return nil, fmt.Errorf("user id is empty")
	}

	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event id: %w", err)
	}
	user.ID = uID

	q := rc.db.Model(user).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no user updated")
	}

	return user, nil
}

func (rc *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&model.User{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user deleted")
	}

	return nil
}

func (rc *UserRepository) List(ctx context.Context, opts *model.UserFindOpts) (*model.UserList, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}

	users := make([]model.User, 0)

	filter := rc.fillFilter(opts)
	fields := rc.fillFields(opts)
	query := rc.db.Model(&users).Column(fields...)

	if filter != "" {
		query = query.Where(filter)
	}

	query = query.Limit(opts.Limit).Offset(opts.Skip)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if count == 0 {
		return &model.UserList{}, nil
	}

	return &model.UserList{
		Users: users,
		Total: count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *UserRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" || userID == "0" {
		return nil, fmt.Errorf("invalid user ID: %s", userID)
	}

	var user model.User

	query := rc.db.Model(&user).Where("id = ?", userID)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to find user by id [%s]: %w", userID, err)
	}

	return &user, nil
}

func (rc *UserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	if usernameOrEmail == "" {
		return nil, fmt.Errorf("invalid username or email")
	}

	var user model.User

	query := rc.db.Model(&user)

	query = query.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to get user by [%s]: %w", usernameOrEmail, err)
	}

	return &user, nil
}

func (rc *UserRepository) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	if usernameOrEmail == "" {
		return false, fmt.Errorf("invalid username or email")
	}

	query := rc.db.Model(&model.User{})

	query = query.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	exists, err := query.Exists()
	if err != nil {
		return false, fmt.Errorf("failed to get user by [%s]: %w", usernameOrEmail, err)
	}

	return exists, nil
}

func (rc *UserRepository) fillFilter(opts *model.UserFindOpts) string {
	filter := ""

	if opts.Username.IsSended {
		filter = addFilterClause(filter, "username", opts.Username.Value)
	}

	if opts.Email.IsSended {
		filter = addFilterClause(filter, "email", opts.Email.Value)
	}

	if opts.RoleID.IsSended {
		filter = addFilterClause(filter, "role_id", opts.RoleID.Value)
	}

	return filter
}

func (rc *UserRepository) fillFields(opts *model.UserFindOpts) []string {
	fields := opts.Fields

	if len(fields) == 0 {
		return nil
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return []string{
			"id",
			"username",
			"email",
			"role_id",
			"deleted_at",
		}
	}

	return fields
}
