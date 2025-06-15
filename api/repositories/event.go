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

type EventRepository struct {
	db *pg.DB
}

func NewEventRepository(db *pg.DB) *EventRepository {
	rc := &EventRepository{
		db: db,
	}

	if err := rc.createSchema(db); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	return rc
}

func (rc *EventRepository) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	if event.Visibility == 0 {
		event.Visibility = model.EventVisibilityPublic
	}

	sqlEvent := rc.internalToSQL(event)

	q := rc.db.Model(sqlEvent)

	_, err := q.Insert()
	if err != nil {
		return nil, pkg.NewError(err, "failed to create event "+event.Name, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(sqlEvent), nil
}

func (rc *EventRepository) Update(ctx context.Context, eventID string, event *model.Event) (*model.Event, error) {
	if eventID == "" || eventID == "0" {
		return nil, pkg.NewError(nil, "event id is empty", http.StatusBadRequest)
	}

	event.ID = eventID

	sqlEvent := rc.internalToSQL(event)

	q := rc.db.Model(sqlEvent)

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND user_id = ?", eventID, ownerID)

	result, err := q.Update()
	if err != nil {
		return nil, pkg.NewError(err, "failed to update event "+sqlEvent.Name, http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return nil, pkg.NewError(nil, "no event updated: "+eventID, http.StatusBadRequest)
	}

	return rc.sqlToInternal(sqlEvent), nil
}

func (rc *EventRepository) Delete(ctx context.Context, eventID string) error {
	q := rc.db.Model(&event{})

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND user_id = ?", eventID, ownerID)

	result, err := q.Delete()
	if err != nil {
		return pkg.NewError(err, "failed to delete event "+eventID, http.StatusInternalServerError)
	}

	if result.RowsAffected() == 0 {
		return pkg.NewError(nil, "no event deleted: "+eventID, http.StatusBadRequest)
	}

	return nil
}

func (rc *EventRepository) List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	if opts == nil {
		return nil, pkg.NewError(nil, "opts is nil", http.StatusBadRequest)
	}

	events := make([]event, 0)

	fields := []string{"*"}
	query := rc.db.Model(&events).Column(fields...)

	query = applyOrderBy(query, opts.OrderByOpts)

	query = applyStandardQueries(query, opts.PaginationOpts)

	query = rc.fillFilter(query, opts)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, pkg.NewError(err, "failed to list events", http.StatusInternalServerError)
	}

	internalEvents := make([]model.Event, 0)
	for _, v := range events {
		internalEvents = append(internalEvents, *rc.sqlToInternal(&v))
	}

	return &model.EventList{
		Events: internalEvents,
		Total:  count,
		PaginationOpts: model.PaginationOpts{
			Limit: opts.Limit,
			Skip:  opts.Skip,
		},
	}, nil
}

func (rc *EventRepository) GetByID(ctx context.Context, eventID string) (*model.Event, error) {
	if eventID == "" || eventID == "0" {
		return nil, pkg.NewError(nil, "invalid event ID: "+eventID, http.StatusBadRequest)
	}

	event := new(event)

	query := rc.db.Model(event).Where("id = ?", eventID)

	if err := query.Select(); err != nil {
		return nil, pkg.NewError(err, "failed to find event by ID "+eventID, http.StatusInternalServerError)
	}

	return rc.sqlToInternal(event), nil
}

func (rc *EventRepository) fillFilter(tx *orm.Query, opts *model.EventFindOpts) *orm.Query {
	if opts.UserID.IsSended {
		tx = applyFilterWithOperand(tx, "user_id", opts.UserID)
	}

	if opts.Visibility.IsSended {
		tx = applyFilterWithOperand(tx, "visibility", opts.Visibility)
	}

	if opts.Name.IsSended {
		tx = applyFilterWithOperand(tx, "name", opts.Name)
	}

	return tx
}

func (rc *EventRepository) internalToSQL(newEvent *model.Event) *event {
	eID, _ := strconv.Atoi(newEvent.ID)
	ownerID, _ := strconv.Atoi(newEvent.UserID)

	items := []eventItem{}
	for _, v := range newEvent.Items {
		items = append(items, eventItem{
			Data: v.Data,
			Type: int(v.Type),
		})
	}

	return &event{
		Date:        newEvent.Date,
		TimeStart:   newEvent.TimeStart,
		TimeEnd:     newEvent.TimeEnd,
		Name:        newEvent.Name,
		Description: newEvent.Description,
		Items:       items,
		ID:          eID,
		UserID:      ownerID,
		Visibility:  int(newEvent.Visibility),
		CreatedAt:   newEvent.CreatedAt,
		UpdatedAt:   newEvent.UpdatedAt,
	}
}

func (rc *EventRepository) sqlToInternal(newEvent *event) *model.Event {
	eID := strconv.Itoa(newEvent.ID)
	ownerID := strconv.Itoa(newEvent.UserID)

	items := []model.EventItem{}
	for _, v := range newEvent.Items {
		items = append(items, model.EventItem{
			Data: v.Data,
			Type: model.EventType(v.Type),
		})
	}

	return &model.Event{
		Date:        newEvent.Date,
		TimeStart:   newEvent.TimeStart,
		TimeEnd:     newEvent.TimeEnd,
		Name:        newEvent.Name,
		Description: newEvent.Description,
		Items:       items,
		ID:          eID,
		UserID:      ownerID,
		Visibility:  model.Visibility(newEvent.Visibility),
		CreatedAt:   newEvent.CreatedAt,
		UpdatedAt:   newEvent.UpdatedAt,
	}
}

func (rc *EventRepository) createSchema(db *pg.DB) error {
	model := (*event)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return pkg.NewError(err, "failed to create event table", http.StatusInternalServerError)
	}

	return nil
}
