package uc

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
)

type EventUC struct {
	repo  interfaces.EventRepository
	cache *EventCacheUC
}

func NewEventUC(repo interfaces.EventRepository, cache *EventCacheUC) *EventUC {
	return &EventUC{
		repo:  repo,
		cache: cache,
	}
}

func (rc *EventUC) Create(ctx context.Context, req *model.EventCreateRequest) (*model.Event, error) {
	event := model.Event{
		Name:      req.Name,
		Items:     req.Items,
		Date:      req.Date,
		TimeStart: req.TimeStart,
		TimeEnd:   req.TimeEnd,
	}

	return rc.repo.Create(ctx, &event)
}

func (rc *EventUC) Update(ctx context.Context, eventID string, req *model.EventUpdateRequest) (*model.Event, error) {
	// event exist control
	_, err := rc.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	event := model.Event{
		Name:      req.Name,
		Items:     req.Items,
		Date:      req.Date,
		TimeStart: req.TimeStart,
		TimeEnd:   req.TimeEnd,
	}

	return rc.repo.Update(ctx, eventID, &event)
}

func (rc *EventUC) Delete(ctx context.Context, id string) error {
	return rc.repo.Delete(ctx, id)
}

func (rc *EventUC) List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	return rc.repo.List(ctx, opts)
}

func (rc *EventUC) GetByID(ctx context.Context, id string) (*model.Event, error) {
	return rc.repo.GetByID(ctx, id)
}
