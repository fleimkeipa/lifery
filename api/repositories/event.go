package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/util"

	"github.com/go-pg/pg"
)

type EventRepository struct {
	db *pg.DB
}

func NewEventRepository(db *pg.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (rc *EventRepository) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	items := []EventItem{}
	for _, v := range event.Items {
		items = append(items, EventItem{
			Data: v.Data,
			Type: EventType(v.Type),
		})
	}

	sqlEvent := Event{
		Date:      event.Date,
		TimeStart: event.TimeStart,
		TimeEnd:   event.TimeEnd,
		Name:      event.Name,
		Items:     items,
		ID:        event.ID,
		OwnerID:   event.OwnerID,
		Private: sql.NullBool{
			Bool:  event.Private,
			Valid: true,
		},
	}
	q := rc.db.Model(&sqlEvent)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create event [%v]: %w", event.Name, err)
	}

	return event, nil
}

func (rc *EventRepository) Update(ctx context.Context, eventID string, event *model.Event) (*model.Event, error) {
	if eventID == "" || eventID == "0" {
		return nil, fmt.Errorf("event id is empty")
	}

	eID, err := strconv.ParseInt(eventID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event id: %w", err)
	}
	event.ID = eID

	q := rc.db.Model(event)

	ownerID := util.GetStrOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", eID, ownerID)

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update event [%v]: %w", event.Name, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no event updated")
	}

	return event, nil
}

func (rc *EventRepository) Delete(ctx context.Context, id string) error {
	q := rc.db.Model(&model.Event{})

	ownerID := util.GetStrOwnerIDFromCtx(ctx)

	q = q.Where("id = ? AND owner_id = ?", id, ownerID)

	result, err := q.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no event deleted")
	}

	return nil
}

func (rc *EventRepository) List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}

	events := make([]model.Event, 0)

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

	return &model.EventList{
		Events: events,
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

	var event model.Event

	query := rc.db.Model(&event).Where("id = ?", eventID)

	if err := query.Select(&event); err != nil {
		return nil, fmt.Errorf("failed to find event by ID [%s]: %w", eventID, err)
	}

	return &event, nil
}

func (rc *EventRepository) fillFilter(opts *model.EventFindOpts) string {
	filter := ""

	if opts.UserID.IsSended {
		filter = addFilterClause(filter, "owner_id", opts.UserID.Value)
	}

	if opts.Private.IsSended {
		filter = addFilterClause(filter, "private", opts.Private.Value)
	}

	return filter
}
