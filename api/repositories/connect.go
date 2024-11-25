package repositories

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/fleimkeipa/lifery/model"

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
		return nil, fmt.Errorf("failed to create connect ID [%v]: %w", connect.ID, err)
	}

	return rc.sqlToInternal(sqlConnect), nil
}

func (rc *ConnectRepository) Update(ctx context.Context, connectID string, connect *model.Connect) (*model.Connect, error) {
	if connectID == "" || connectID == "0" {
		return nil, fmt.Errorf("connect id is empty")
	}

	connect.ID = connectID

	sqlConnect := rc.internalToSQL(connect)

	q := rc.db.Model(sqlConnect).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update connect ID [%v]: %w", connectID, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no connect updated")
	}

	return rc.sqlToInternal(sqlConnect), nil
}

func (rc *ConnectRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&connect{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete connect: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no connect deleted")
	}

	return nil
}

func (rc *ConnectRepository) List(ctx context.Context, opts *model.ConnectFindOpts) (*model.ConnectList, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}

	connects := make([]connect, 0)

	filter := rc.fillFilter(opts)
	fields := rc.fillFields(opts)
	query := rc.db.Model(&connects).Column(fields...)

	if filter != "" {
		query = query.Where(filter)
	}

	query = query.Limit(opts.Limit).Offset(opts.Skip)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to list connects: %w", err)
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
		return nil, fmt.Errorf("invalid connect ID: %s", connectID)
	}

	var connect connect

	query := rc.db.Model(&connect).Where("id = ?", connectID)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to find connect by id [%s]: %w", connectID, err)
	}

	return rc.sqlToInternal(&connect), nil
}

func (rc *ConnectRepository) fillFilter(opts *model.ConnectFindOpts) string {
	filter := ""

	if opts.Status.IsSended {
		filter = addFilterClause(filter, "status", opts.Status.Value)
	}

	if opts.FriendID.IsSended {
		filter = addFilterClause(filter, "friend_id", opts.FriendID.Value)
	}

	return filter
}

func (rc *ConnectRepository) fillFields(opts *model.ConnectFindOpts) []string {
	fields := opts.Fields

	if len(fields) == 0 {
		return nil
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return []string{
			"id",
			"status",
			"user_id",
			"friend_id",
		}
	}

	return fields
}

func (rc *ConnectRepository) internalToSQL(newConnect *model.Connect) *connect {
	cID, _ := strconv.Atoi(newConnect.ID)

	return &connect{
		ID:       cID,
		Status:   requestStatus(newConnect.Status),
		UserID:   newConnect.UserID,
		FriendID: newConnect.FriendID,
	}
}

func (rc *ConnectRepository) sqlToInternal(newConnect *connect) *model.Connect {
	cID := strconv.Itoa(newConnect.ID)

	return &model.Connect{
		ID:       cID,
		Status:   model.RequestStatus(newConnect.Status),
		UserID:   newConnect.UserID,
		FriendID: newConnect.FriendID,
	}
}

func (rc *ConnectRepository) createSchema(db *pg.DB) error {
	model := (*connect)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
