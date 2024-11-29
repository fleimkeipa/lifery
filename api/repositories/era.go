package repositories

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/util"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type EraRepository struct {
	db *pg.DB
}

func NewEraRepository(db *pg.DB) *EraRepository {
	rc := &EraRepository{
		db: db,
	}

	if err := rc.createSchema(db); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	return rc
}

func (rc *EraRepository) Create(ctx context.Context, era *model.Era) (*model.Era, error) {
	sqlEra := rc.internalToSQL(era)

	q := rc.db.Model(sqlEra)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create era [%v]: %w", era.Name, err)
	}

	return rc.sqlToInternal(sqlEra), nil
}

func (rc *EraRepository) Update(ctx context.Context, eraID string, era *model.Era) (*model.Era, error) {
	if eraID == "" || eraID == "0" {
		return nil, fmt.Errorf("invalid era id [%v]", eraID)
	}

	era.ID = eraID

	sqlEra := rc.internalToSQL(era)

	q := rc.db.Model(sqlEra)

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", eraID, ownerID)

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update era [%v]: %w", sqlEra.Name, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no era updated")
	}

	return rc.sqlToInternal(sqlEra), nil
}

func (rc *EraRepository) Delete(ctx context.Context, id string) error {
	q := rc.db.Model(&era{})

	ownerID := util.GetOwnerIDFromCtx(ctx)

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

	eras := make([]era, 0)

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

	internalEras := make([]model.Era, 0)
	for _, v := range eras {
		internalEras = append(internalEras, *rc.sqlToInternal(&v))
	}

	return &model.EraList{
		Eras:  internalEras,
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

	resp := new(eraGetResponse)

	err := rc.db.Model(&era{}).
		Relation("Owner").
		Where("era.id = ?", eraID).
		Select(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to find era by ID [%s]: %w", eraID, err)
	}

	return rc.sqlToInternal2(resp), nil
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

func (rc *EraRepository) internalToSQL(newEra *model.Era) *era {
	eID, _ := strconv.Atoi(newEra.ID)
	ownerID, _ := strconv.Atoi(newEra.OwnerID)
	return &era{
		TimeStart: newEra.TimeStart,
		TimeEnd:   newEra.TimeEnd,
		Name:      newEra.Name,
		Color:     newEra.Color,
		OwnerID:   ownerID,
		ID:        eID,
	}
}

func (rc *EraRepository) sqlToInternal(newEra *era) *model.Era {
	eID := strconv.Itoa(newEra.ID)
	ownerID := strconv.Itoa(newEra.OwnerID)
	return &model.Era{
		TimeStart: newEra.TimeStart,
		TimeEnd:   newEra.TimeEnd,
		Name:      newEra.Name,
		Color:     newEra.Color,
		OwnerID:   ownerID,
		ID:        eID,
	}
}

func (rc *EraRepository) sqlToInternal2(newEra *eraGetResponse) *model.Era {
	eID := strconv.Itoa(newEra.ID)
	ownerID := strconv.Itoa(newEra.OwnerID)
	return &model.Era{
		TimeStart: newEra.TimeStart,
		TimeEnd:   newEra.TimeEnd,
		Name:      newEra.Name,
		Color:     newEra.Color,
		OwnerID:   ownerID,
		ID:        eID,
	}
}

func (rc *EraRepository) createSchema(db *pg.DB) error {
	model := (*era)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return fmt.Errorf("failed to create era table: %w", err)
	}

	return nil
}
