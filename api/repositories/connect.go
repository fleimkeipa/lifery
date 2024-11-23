package repositories

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg"
)

type ConnectRepository struct {
	db *pg.DB
}

func NewConnectRepository(db *pg.DB) *ConnectRepository {
	return &ConnectRepository{
		db: db,
	}
}

func (rc *ConnectRepository) Create(ctx context.Context, connect *model.Connect) (*model.Connect, error) {
	q := rc.db.Model(connect)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create connect: %w", err)
	}

	return connect, nil
}

func (rc *ConnectRepository) Update(ctx context.Context, connectID string, connect *model.Connect) (*model.Connect, error) {
	if connectID == "" || connectID == "0" {
		return nil, fmt.Errorf("connect id is empty")
	}

	uID, err := strconv.ParseInt(connectID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event id: %w", err)
	}
	connect.ID = uID

	q := rc.db.Model(connect).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update connect: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no connect updated")
	}

	return connect, nil
}

func (rc *ConnectRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&model.Connect{}).Where("id = ?", id).Delete()
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

	connects := make([]model.Connect, 0)

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

	return &model.ConnectList{
		Connects: connects,
		Total:    count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *ConnectRepository) GetByID(ctx context.Context, friendRequestID string) (*model.Connect, error) {
	if friendRequestID == "" || friendRequestID == "0" {
		return nil, fmt.Errorf("invalid connect ID: %s", friendRequestID)
	}

	var friendRequest model.Connect

	query := rc.db.Model(&friendRequest).Where("id = ?", friendRequestID)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to find connect by id [%s]: %w", friendRequestID, err)
	}

	return &friendRequest, nil
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
