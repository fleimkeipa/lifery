package repositories

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
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
		return nil, pkg.NewError(err, "failed to create era", http.StatusInternalServerError)
	}

	return rc.sqlToInternal(sqlEra), nil
}

func (rc *EraRepository) Update(ctx context.Context, eraID string, era *model.Era) (*model.Era, error) {
	if eraID == "" || eraID == "0" {
		return nil, pkg.NewError(nil, "invalid era id "+eraID, http.StatusBadRequest)
	}

	era.ID = eraID

	sqlEra := rc.internalToSQL(era)

	q := rc.db.Model(sqlEra)

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", eraID, ownerID)

	result, err := q.Update()
	if err != nil {
		return nil, pkg.NewError(err, "failed to update era", http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return nil, pkg.NewError(nil, "no era updated", http.StatusBadRequest)
	}

	return rc.sqlToInternal(sqlEra), nil
}

func (rc *EraRepository) Delete(ctx context.Context, id string) error {
	q := rc.db.Model(&era{})

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", id, ownerID)

	result, err := q.Delete()
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pkg.NewError(nil, "no era deleted", http.StatusBadRequest)
	}

	return nil
}

func (rc *EraRepository) List(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error) {
	if opts == nil {
		return nil, pkg.NewError(nil, "opts is nil", http.StatusBadRequest)
	}

	eras := make([]era, 0)

	fields := []string{"*"}
	query := rc.db.Model(&eras).Column(fields...)

	query = applyOrderBy(query, opts.OrderByOpts)

	query = applyStandardQueries(query, opts.PaginationOpts)

	query = rc.fillFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, pkg.NewError(err, "failed to list eras", http.StatusInternalServerError)
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
		return nil, pkg.NewError(nil, "invalid era ID: "+eraID, http.StatusBadRequest)
	}

	resp := new(eraGetResponse)

	err := rc.db.Model(&era{}).
		Relation("Owner").
		Where("era.id = ?", eraID).
		Select(resp)
	if err != nil {
		return nil, pkg.NewError(err, "failed to find era by ID "+eraID, http.StatusInternalServerError)
	}

	return rc.sqlToInternal2(resp), nil
}

func (rc *EraRepository) fillFilter(tx *orm.Query, opts *model.EraFindOpts) *orm.Query {
	if opts.Name.IsSended {
		tx = applyFilterWithOperand(tx, "name", opts.Name)
	}

	if opts.UserID.IsSended {
		tx = applyFilterWithOperand(tx, "owner_id", opts.UserID)
	}

	return tx
}

func (rc *EraRepository) internalToSQL(newEra *model.Era) *era {
	eID, _ := strconv.Atoi(newEra.ID)
	userID, _ := strconv.Atoi(newEra.UserID)
	return &era{
		TimeStart: newEra.TimeStart,
		TimeEnd:   newEra.TimeEnd,
		Name:      newEra.Name,
		Color:     newEra.Color,
		UserID:    userID,
		ID:        eID,
	}
}

func (rc *EraRepository) sqlToInternal(newEra *era) *model.Era {
	eID := strconv.Itoa(newEra.ID)
	userID := strconv.Itoa(newEra.UserID)
	return &model.Era{
		TimeStart: newEra.TimeStart,
		TimeEnd:   newEra.TimeEnd,
		Name:      newEra.Name,
		Color:     newEra.Color,
		UserID:    userID,
		ID:        eID,
	}
}

func (rc *EraRepository) sqlToInternal2(newEra *eraGetResponse) *model.Era {
	eID := strconv.Itoa(newEra.ID)
	userID := strconv.Itoa(newEra.UserID)
	return &model.Era{
		TimeStart: newEra.TimeStart,
		TimeEnd:   newEra.TimeEnd,
		Name:      newEra.Name,
		Color:     newEra.Color,
		UserID:    userID,
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
		return pkg.NewError(err, "failed to create era table", http.StatusInternalServerError)
	}

	return nil
}
