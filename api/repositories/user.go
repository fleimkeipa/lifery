package repositories

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/pkg/logger"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	rc := &UserRepository{
		db: db,
	}

	if err := rc.createSchema(db); err != nil {
		logger.Log.Fatalf("failed to create schema: %v", err)
	}

	return rc
}

func (rc *UserRepository) Create(ctx context.Context, newUser *model.User) (*model.User, error) {
	sqlUser := rc.internalToSQL(newUser)

	q := rc.db.Model(sqlUser)

	_, err := q.Insert()
	if err != nil {
		return nil, pkg.NewError(err, "failed to create user", http.StatusInternalServerError)
	}

	return rc.sqlToInternal(sqlUser), nil
}

func (rc *UserRepository) Update(ctx context.Context, userID string, user *model.User) (*model.User, error) {
	if userID == "" || userID == "0" {
		return nil, pkg.NewError(nil, "user id is empty", http.StatusBadRequest)
	}

	user.ID = userID

	sqlUser := rc.internalToSQL(user)

	q := rc.db.Model(sqlUser).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, pkg.NewError(err, "failed to update user", http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return nil, pkg.NewError(nil, "no user updated", http.StatusBadRequest)
	}

	return rc.sqlToInternal(sqlUser), nil
}

func (rc *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&user{}).Where("id = ?", id).Delete()
	if err != nil {
		return pkg.NewError(err, "failed to delete user", http.StatusInternalServerError)
	}
	if result.RowsAffected() == 0 {
		return pkg.NewError(nil, "no user deleted: "+id, http.StatusBadRequest)
	}

	return nil
}

func (rc *UserRepository) List(ctx context.Context, opts *model.UserFindOpts) (*model.UserList, error) {
	if opts == nil {
		return nil, pkg.NewError(nil, "opts is nil", http.StatusBadRequest)
	}

	users := make([]user, 0)

	query := rc.db.Model(&users)

	query = applyOrderBy(query, opts.OrderByOpts)

	query = applyStandardQueries(query, opts.PaginationOpts)

	query = rc.fillFields(query, opts)

	query = rc.fillFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, pkg.NewError(err, "failed to list users", http.StatusInternalServerError)
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
		return nil, pkg.NewError(nil, "invalid user ID: "+userID, http.StatusBadRequest)
	}

	var user user

	query := rc.db.Model(&user).Where("id = ?", userID)

	if err := query.Select(); err != nil {
		return nil, pkg.NewError(err, "failed to find user by id "+userID, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(&user), nil
}

func (rc *UserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	if usernameOrEmail == "" {
		return nil, pkg.NewError(nil, "missing username or email", http.StatusBadRequest)
	}

	var user user

	query := rc.db.Model(&user)

	query = query.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	if err := query.Select(); err != nil {
		return nil, pkg.NewError(err, "failed to get user by "+usernameOrEmail, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(&user), nil
}

func (rc *UserRepository) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	if usernameOrEmail == "" {
		return false, pkg.NewError(nil, "invalid username or email", http.StatusBadRequest)
	}

	query := rc.db.Model(&user{})

	query = query.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	exists, err := query.Exists()
	if err != nil {
		return false, pkg.NewError(err, "failed to get user by "+usernameOrEmail, http.StatusInternalServerError)
	}

	return exists, nil
}

func (rc *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, pkg.NewError(nil, "missing email", http.StatusBadRequest)
	}

	var user user

	query := rc.db.
		Model(&user).
		Where("email = ?", email)

	if err := query.Select(); err != nil {
		return nil, pkg.NewError(err, "failed to get user by email "+email, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(&user), nil
}

func (rc *UserRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	if userID == "" || userID == "0" {
		return pkg.NewError(nil, "invalid user ID: "+userID, http.StatusBadRequest)
	}

	if hashedPassword == "" {
		return pkg.NewError(nil, "missing hashed password", http.StatusBadRequest)
	}

	result, err := rc.db.
		Model(&user{}).
		Set("password = ?", hashedPassword).
		Where("id = ?", userID).
		Update()
	if err != nil {
		return pkg.NewError(err, "failed to update user password", http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return pkg.NewError(nil, "no user updated: "+userID, http.StatusBadRequest)
	}

	return nil
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
			"auth_type",
		)
	}

	return tx.Column(fields...)
}

func (rc *UserRepository) internalToSQL(newUser *model.User) *user {
	uID, _ := strconv.Atoi(newUser.ID)
	connects := make([]*connect, 0)
	for _, v := range newUser.Connects {
		connectID, _ := strconv.Atoi(v.ID)
		userID, _ := strconv.Atoi(v.UserID)
		friendID, _ := strconv.Atoi(v.FriendID)
		connects = append(connects, &connect{
			ID:       connectID,
			Status:   int(v.Status),
			UserID:   userID,
			FriendID: friendID,
		})
	}
	return &user{
		CreatedAt: newUser.CreatedAt,
		Connects:  connects,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  newUser.Password,
		ID:        uID,
		RoleID:    UserRole(newUser.RoleID),
		AuthType:  string(newUser.AuthType),
	}
}

func (rc *UserRepository) sqlToInternal(newUser *user) *model.User {
	uID := strconv.Itoa(newUser.ID)
	connects := make([]*model.Connect, 0)
	for _, v := range newUser.Connects {
		connects = append(connects, &model.Connect{
			ID:       strconv.Itoa(v.ID),
			Status:   model.RequestStatus(v.Status),
			UserID:   strconv.Itoa(v.UserID),
			FriendID: strconv.Itoa(v.FriendID),
			User:     *rc.sqlToInternal(v.User),
			Friend:   *rc.sqlToInternal(v.Friend),
		})
	}
	return &model.User{
		CreatedAt: newUser.CreatedAt,
		Connects:  connects,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  newUser.Password,
		ID:        uID,
		RoleID:    model.UserRole(newUser.RoleID),
		AuthType:  model.AuthType(newUser.AuthType),
	}
}

func (rc *UserRepository) createSchema(db *pg.DB) error {
	model := (*user)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return pkg.NewError(err, "failed to create user table", http.StatusInternalServerError)
	}

	return nil
}
