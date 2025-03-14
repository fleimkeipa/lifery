package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	rc := &UserRepository{
		db: db,
	}

	if err := rc.createSchema(db); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	return rc
}

func (rc *UserRepository) Create(ctx context.Context, newUser *model.User) (*model.User, error) {
	if newUser.RoleID <= 0 {
		newUser.RoleID = model.ViewerRole
	}

	sqlUser := rc.internalToSQL(newUser)

	q := rc.db.Model(sqlUser)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return rc.sqlToInternal(sqlUser), nil
}

func (rc *UserRepository) Update(ctx context.Context, userID string, user *model.User) (*model.User, error) {
	if userID == "" || userID == "0" {
		return nil, fmt.Errorf("user id is empty")
	}

	user.ID = userID

	sqlUser := rc.internalToSQL(user)

	q := rc.db.Model(sqlUser).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no user updated")
	}

	return rc.sqlToInternal(sqlUser), nil
}

func (rc *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&user{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user deleted: %s", id)
	}

	return nil
}

func (rc *UserRepository) List(ctx context.Context, opts *model.UserFindOpts) (*model.UserList, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}

	users := make([]user, 0)

	query := rc.db.Model(&users)

	query = applyOrderBy(query, opts.OrderByOpts)

	query = applyStandardQueries(query, opts.PaginationOpts)

	query = rc.fillFields(query, opts)

	query = rc.fillFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if count == 0 {
		return &model.UserList{}, nil
	}

	internalUsers := make([]model.User, 0)
	for _, v := range users {
		internalUsers = append(internalUsers, *rc.sqlToInternal(&v))
	}

	return &model.UserList{
		Users: internalUsers,
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

	var user user

	query := rc.db.Model(&user).Where("id = ?", userID)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to find user by id [%s]: %w", userID, err)
	}

	return rc.sqlToInternal(&user), nil
}

func (rc *UserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	if usernameOrEmail == "" {
		return nil, errors.New("missing username or email")
	}

	var user user

	query := rc.db.Model(&user)

	query = query.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to get user by [%s]: %w", usernameOrEmail, err)
	}

	return rc.sqlToInternal(&user), nil
}

func (rc *UserRepository) GetConnects(ctx context.Context, opts *model.UserConnectsFindOpts) (*model.UserConnects, error) {
	if !opts.UserID.IsSended {
		return nil, errors.New("missing user id")
	}

	var connects userConnects

	query := rc.db.Model(&connects)

	query = applyOrderBy(query, opts.OrderByOpts)

	query = applyStandardQueries(query, opts.PaginationOpts)

	query = rc.fillConnectsFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get user connects: %w", err)
	}

	if len(connects.Connects) == 0 {
		return &model.UserConnects{}, nil
	}

	respConnects := new(model.UserConnects)
	for _, v := range connects.Connects {
		respConnects.Connects = append(respConnects.Connects, *rc.sqlToInternal(&v))
	}

	return &model.UserConnects{
		Connects: respConnects.Connects,
		Total:    count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *UserRepository) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	if usernameOrEmail == "" {
		return false, fmt.Errorf("invalid username or email")
	}

	query := rc.db.Model(&user{})

	query = query.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	exists, err := query.Exists()
	if err != nil {
		return false, fmt.Errorf("failed to get user by [%s]: %w", usernameOrEmail, err)
	}

	return exists, nil
}

func (rc *UserRepository) fillFilter(tx *orm.Query, opts *model.UserFindOpts) *orm.Query {
	if opts.Username.IsSended {
		tx = applyFilterWithOperand(tx, "username", opts.Username)
	}

	if opts.Email.IsSended {
		tx = applyFilterWithOperand(tx, "email", opts.Email)
	}

	if opts.RoleID.IsSended {
		tx = applyFilterWithOperand(tx, "role_id", opts.RoleID)
	}

	return tx
}

func (rc *UserRepository) fillFields(tx *orm.Query, opts *model.UserFindOpts) *orm.Query {
	fields := opts.Fields

	if len(fields) == 0 {
		return tx
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return tx.Column(
			"id",
			"username",
			"email",
			"role_id",
			"deleted_at",
		)
	}

	return tx.Column(fields...)
}

func (rc *UserRepository) internalToSQL(newUser *model.User) *user {
	uID, _ := strconv.Atoi(newUser.ID)
	return &user{
		DeletedAt: newUser.DeletedAt,
		CreatedAt: newUser.CreatedAt,
		Connects:  newUser.Connects,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  newUser.Password,
		ID:        uID,
		RoleID:    newUser.RoleID,
	}
}

func (rc *UserRepository) sqlToInternal(newUser *user) *model.User {
	uID := strconv.Itoa(newUser.ID)
	return &model.User{
		DeletedAt: newUser.DeletedAt,
		CreatedAt: newUser.CreatedAt,
		Connects:  newUser.Connects,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  newUser.Password,
		ID:        uID,
		RoleID:    newUser.RoleID,
	}
}

func (rc *UserRepository) createSchema(db *pg.DB) error {
	model := (*user)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	return nil
}

func (rc *UserRepository) fillConnectsFilter(tx *orm.Query, opts *model.UserConnectsFindOpts) *orm.Query {
	filter := tx

	if opts.UserID.IsSended {
		filter = applyFilterWithOperand(tx, "id", opts.UserID)
	}

	return filter
}
