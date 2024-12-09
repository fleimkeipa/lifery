package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
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
		return nil, fmt.Errorf("failed to create event [%v]: %w", event.Name, err)
	}

	return rc.sqlToInternal(sqlEvent), nil
}

func (rc *EventRepository) Update(ctx context.Context, eventID string, event *model.Event) (*model.Event, error) {
	if eventID == "" || eventID == "0" {
		return nil, fmt.Errorf("event id is empty")
	}

	event.ID = eventID

	sqlEvent := rc.internalToSQL(event)

	q := rc.db.Model(sqlEvent)

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", eventID, ownerID)

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update event [%v]: %w", sqlEvent.Name, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no event updated: [%s]", eventID)
	}

	return rc.sqlToInternal(sqlEvent), nil
}

func (rc *EventRepository) Delete(ctx context.Context, eventID string) error {
	q := rc.db.Model(&event{})

	ownerID := util.GetOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", eventID, ownerID)

	result, err := q.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete event [%s]: %w", eventID, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no event deleted: [%s]", eventID)
	}

	return nil
}

func (rc *EventRepository) List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	if opts == nil {
		return nil, errors.New("opts is nil")
	}

	events := make([]event, 0)

	filter := rc.fillFilter(opts)
	fields := []string{"*"}
	query := rc.db.Model(&events).Column(fields...)

	if filter != "" {
		query = query.Where(filter)
	}

	query = query.Limit(opts.Limit).Offset(opts.Skip)

	count, err := query.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	if count == 0 {
		return &model.EventList{}, nil
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
		return nil, fmt.Errorf("invalid event ID: %s", eventID)
	}

	event := new(event)

	query := rc.db.Model(event).Where("id = ?", eventID)

	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("failed to find event by ID [%s]: %w", eventID, err)
	}

	return rc.sqlToInternal(event), nil
}

func (rc *EventRepository) fillFilter(opts *model.EventFindOpts) string {
	filter := ""

	if opts.UserID.IsSended {
		filter = addFilterClause(filter, "owner_id", opts.UserID.Value)
	}

	if opts.Visibility.IsSended {
		filter = addFilterClause(filter, "visibility", opts.Visibility.Value)
	}

	return filter
}

func (rc *EventRepository) internalToSQL(newEvent *model.Event) *event {
	eID, _ := strconv.Atoi(newEvent.ID)
	ownerID, _ := strconv.Atoi(newEvent.OwnerID)

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
		OwnerID:     ownerID,
		Visibility:  int(newEvent.Visibility),
	}
}

func (rc *EventRepository) sqlToInternal(newEvent *event) *model.Event {
	eID := strconv.Itoa(newEvent.ID)
	ownerID := strconv.Itoa(newEvent.OwnerID)

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
		OwnerID:     ownerID,
		Visibility:  model.Visibility(newEvent.Visibility),
	}
}

func (rc *EventRepository) createSchema(db *pg.DB) error {
	model := (*event)(nil)

	opts := &orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}

	if err := db.Model(model).CreateTable(opts); err != nil {
		return fmt.Errorf("failed to create event table: %w", err)
	}

	return nil
}
