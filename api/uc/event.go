package uc

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/util"
)

type EventUC struct {
	repo   interfaces.EventRepository
	userUC *UserUC
	cache  *EventCacheUC
}

func NewEventUC(repo interfaces.EventRepository, cache *EventCacheUC, userUC *UserUC) *EventUC {
	return &EventUC{
		repo:   repo,
		userUC: userUC,
		cache:  cache,
	}
}

func (rc *EventUC) Create(ctx context.Context, req *model.EventCreateRequest) (*model.Event, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	event := model.Event{
		Date:       req.Date,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Name:       req.Name,
		Items:      req.Items,
		OwnerID:    ownerID,
		Visibility: req.Visibility,
	}

	return rc.repo.Create(ctx, &event)
}

func (rc *EventUC) Update(ctx context.Context, eventID string, req *model.EventUpdateRequest) (*model.Event, error) {
	// event exist control
	exist, err := rc.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	event := model.Event{
		Date:       req.Date,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Name:       req.Name,
		Items:      req.Items,
		OwnerID:    exist.OwnerID,
		Visibility: req.Visibility,
	}

	return rc.repo.Update(ctx, eventID, &event)
}

func (rc *EventUC) Delete(ctx context.Context, id string) error {
	return rc.repo.Delete(ctx, id)
}

func (rc *EventUC) List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	if !opts.UserID.IsSended {
		opts.UserID = model.Filter{
			Value:    ownerID,
			IsSended: true,
		}

		return rc.repo.List(ctx, opts)
	}

	isConnected, err := rc.userUC.IsConnected(ctx, ownerID, opts.UserID.Value)
	if err != nil {
		return nil, err
	}

	if !isConnected {
		opts.Private = model.Filter{
			Value:    "false",
			IsSended: true,
		}

		return rc.repo.List(ctx, opts)
	}

	return rc.repo.List(ctx, opts)
}

func (rc *EventUC) GetByID(ctx context.Context, id string) (*model.Event, error) {
	return rc.repo.GetByID(ctx, id)
}
