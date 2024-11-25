package repositories

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/util"

	"github.com/go-pg/pg"
)

type EraRepository struct {
	db *pg.DB
}

func NewEraRepository(db *pg.DB) *EraRepository {
	return &EraRepository{
		db: db,
	}
}

func (rc *EraRepository) Create(ctx context.Context, era *model.Era) (*model.Era, error) {
	q := rc.db.Model(era)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create era [%v]: %w", era.Name, err)
	}

	return era, nil
}

func (rc *EraRepository) Update(ctx context.Context, eraID string, era *model.Era) (*model.Era, error) {
	if eraID == "" || eraID == "0" {
		return nil, fmt.Errorf("invalid era id [%v]", eraID)
	}

	eID, err := strconv.ParseInt(eraID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse era id: %w", err)
	}
	era.ID = eID

	q := rc.db.Model(era)

	ownerID := util.GetStrOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", eID, ownerID)

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update era [%v]: %w", era.Name, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no era updated")
	}

	return era, nil
}

func (rc *EraRepository) Delete(ctx context.Context, id string) error {
	q := rc.db.Model(&model.Era{})

	ownerID := util.GetStrOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", id, ownerID)

	result, err := q.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete era: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no era deleted")
	}

	return nil
}

func (rc *EraRepository) List(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}

	eras := make([]model.Era, 0)

	filter := rc.fillFilter(opts)
	fields := []string{"*"}
	query := rc.db.Model(&eras).Column(fields...)

	if filter != "" {
		query = query.Where(filter)
	}

	query = query.Limit(opts.Limit).Offset(opts.Skip)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to list eras: %w", err)
	}

	if count == 0 {
		return &model.EraList{}, nil
	}

	return &model.EraList{
		Eras:  eras,
		Total: count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *EraRepository) GetByID(ctx context.Context, eraID string) (*model.Era, error) {
	if eraID == "" || eraID == "0" {
		return nil, fmt.Errorf("invalid era ID: %s", eraID)
	}

	var era model.Era

	query := rc.db.Model(era).Where("id = ?", eraID)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to find era by ID [%s]: %w", eraID, err)
	}

	return &era, nil
}

func (rc *EraRepository) fillFilter(opts *model.EraFindOpts) string {
	filter := ""

	if opts.Name.IsSended {
		filter = addFilterClause(filter, "name", opts.Name.Value)
	}

	if opts.UserID.IsSended {
		filter = addFilterClause(filter, "owner_id", opts.UserID.Value)
	}

	return filter
}
