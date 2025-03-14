package repositories

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type ConnectRepository struct {
	db *pg.DB
}

func NewConnectRepository(db *pg.DB) *ConnectRepository {
	rc := &ConnectRepository{
		db: db,
	}

	if err := rc.createSchema(db); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	return rc
}

func (rc *ConnectRepository) Create(ctx context.Context, connect *model.Connect) (*model.Connect, error) {
	sqlConnect := rc.internalToSQL(connect)

	q := rc.db.Model(sqlConnect)

	_, err := q.Insert()
	if err != nil {
		return nil, pkg.NewError(err, "failed to create connect ID "+connect.ID, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(sqlConnect), nil
}

func (rc *ConnectRepository) Update(ctx context.Context, connectID string, connect *model.Connect) (*model.Connect, error) {
	if connectID == "" || connectID == "0" {
		return nil, pkg.NewError(nil, "connect id is empty", http.StatusBadRequest)
	}

	connect.ID = connectID

	sqlConnect := rc.internalToSQL(connect)

	q := rc.db.Model(sqlConnect).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, pkg.NewError(err, "failed to update connect ID "+connectID, http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return nil, pkg.NewError(nil, "no connect updated", http.StatusBadRequest)
	}

	return rc.sqlToInternal(sqlConnect), nil
}

func (rc *ConnectRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&connect{}).Where("id = ?", id).Delete()
	if err != nil {
		return pkg.NewError(err, "failed to delete connect", http.StatusInternalServerError)
	}
	if result.RowsAffected() == 0 {
		return pkg.NewError(nil, "no connect deleted", http.StatusBadRequest)
	}

	return nil
}

func (rc *ConnectRepository) ConnectsRequests(ctx context.Context, opts *model.ConnectFindOpts) (*model.ConnectList, error) {
	if opts == nil {
		return nil, pkg.NewError(nil, "opts is nil", http.StatusBadRequest)
	}

	connects := make([]connect, 0)

	query := rc.db.Model(&connects)

	query = applyOrderBy(query, opts.OrderByOpts)

	query = applyStandardQueries(query, opts.PaginationOpts)

	query = rc.fillFields(query, opts)

	query = rc.fillConnectsRequestsFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, pkg.NewError(err, "failed to list connects", http.StatusInternalServerError)
	}

	if count == 0 {
		return &model.ConnectList{}, nil
	}

	internalConnects := make([]model.Connect, 0)
	for _, v := range connects {
		internalConnects = append(internalConnects, *rc.sqlToInternal(&v))
	}

	return &model.ConnectList{
		Connects: internalConnects,
		Total:    count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *ConnectRepository) GetByID(ctx context.Context, connectID string) (*model.Connect, error) {
	if connectID == "" || connectID == "0" {
		return nil, pkg.NewError(nil, "invalid connect ID: "+connectID, http.StatusBadRequest)
	}

	connect := new(connect)

	query := rc.db.Model(connect).Where("id = ?", connectID)

	if err := query.Select(); err != nil {
		return nil, pkg.NewError(err, "failed to find connect by id "+connectID, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(connect), nil
}

func (rc *ConnectRepository) fillConnectsRequestsFilter(tx *orm.Query, opts *model.ConnectFindOpts) *orm.Query {
	if opts.Status.IsSended {
		tx = applyFilterWithOperand(tx, "status", opts.Status)
	}

	if opts.FriendID.IsSended {
		tx = applyFilterWithOperand(tx, "friend_id", opts.FriendID)
	}

	return tx
}

func (rc *ConnectRepository) fillFields(tx *orm.Query, opts *model.ConnectFindOpts) *orm.Query {
	fields := opts.Fields

	if len(fields) == 0 {
		return tx
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return tx.Column(
			"id",
			"status",
			"user_id",
			"friend_id",
		)
	}

	return tx.Column(fields...)
}

func (rc *ConnectRepository) internalToSQL(newConnect *model.Connect) *connect {
	cID, _ := strconv.Atoi(newConnect.ID)
	userID, _ := strconv.Atoi(newConnect.UserID)
	friendID, _ := strconv.Atoi(newConnect.FriendID)

	return &connect{
		ID:       cID,
		Status:   int(newConnect.Status),
		UserID:   userID,
		FriendID: friendID,
	}
}

func (rc *ConnectRepository) sqlToInternal(newConnect *connect) *model.Connect {
	cID := strconv.Itoa(newConnect.ID)
	userID := strconv.Itoa(newConnect.UserID)
	friendID := strconv.Itoa(newConnect.FriendID)

	return &model.Connect{
		ID:       cID,
		Status:   model.RequestStatus(newConnect.Status),
		UserID:   userID,
		FriendID: friendID,
	}
}

func (rc *ConnectRepository) createSchema(db *pg.DB) error {
	model := (*connect)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return pkg.NewError(err, "failed to create connect table", http.StatusInternalServerError)
	}

	return nil
}
