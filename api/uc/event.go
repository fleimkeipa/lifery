package uc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/util"
)

type EventUC struct {
	repo   interfaces.EventRepository
	userUC *UserUC
}

func NewEventUC(repo interfaces.EventRepository, userUC *UserUC) *EventUC {
	return &EventUC{
		repo:   repo,
		userUC: userUC,
	}
}

func (rc *EventUC) Create(ctx context.Context, req *model.EventCreateRequest) (*model.Event, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	date, err := time.Parse(`2006-01-02`, req.Date)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
	}
	timeStart, err := time.Parse(`2006-01-02`, req.TimeStart)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
	}
	timeEnd, err := time.Parse(`2006-01-02`, req.TimeEnd)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeEnd", http.StatusBadRequest)
	}

	event := model.Event{
		Date:        date.Format(`2006-01-02`),
		TimeStart:   timeStart.Format(`2006-01-02`),
		TimeEnd:     timeEnd.Format(`2006-01-02`),
		Name:        req.Name,
		Description: req.Description,
		Items:       req.Items,
		OwnerID:     ownerID,
		Visibility:  req.Visibility,
	}

	newEvent, err := rc.repo.Create(ctx, &event)
	if err != nil {
		return nil, pkg.NewError(err, "failed to create event", http.StatusInternalServerError)
	}

	return newEvent, nil
}

func (rc *EventUC) Update(ctx context.Context, eventID string, req *model.EventUpdateRequest) (*model.Event, error) {
	// event exist control
	exist, err := rc.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	event := model.Event{
		Date:        req.Date,
		TimeStart:   req.TimeStart,
		TimeEnd:     req.TimeEnd,
		Name:        req.Name,
		Description: req.Description,
		Items:       req.Items,
		OwnerID:     exist.OwnerID,
		Visibility:  req.Visibility,
	}

	updatedEvent, err := rc.repo.Update(ctx, eventID, &event)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update event", http.StatusInternalServerError)
	}

	return updatedEvent, nil
}

func (rc *EventUC) Delete(ctx context.Context, id string) error {
	if err := rc.repo.Delete(ctx, id); err != nil {
		return pkg.NewError(err, "failed to delete event", http.StatusInternalServerError)
	}

	return nil
}

func (rc *EventUC) List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	if ownerID == "" {
		if opts.UserID.Value == "" {
			return nil, pkg.NewError(nil, "user id is empty", http.StatusBadRequest)
		}

		opts.UserID = model.Filter{
			Value:    opts.UserID.Value,
			IsSended: true,
		}
		opts.Visibility = model.Filter{
			Value:    fmt.Sprintf("%d", model.EventVisibilityPublic),
			IsSended: true,
		}

		return rc.list(ctx, opts)
	}

	if !opts.UserID.IsSended {
		opts.UserID = model.Filter{
			Value:    ownerID,
			IsSended: true,
		}

		return rc.list(ctx, opts)
	}

	isConnected, err := rc.userUC.IsConnected(ctx, ownerID, opts.UserID.Value)
	if err != nil {
		return nil, err
	}

	if !isConnected {
		opts.Visibility = model.Filter{
			Value:    fmt.Sprintf("%d", model.EventVisibilityPublic),
			IsSended: true,
		}

		return rc.list(ctx, opts)
	}

	opts.Visibility = model.Filter{
		Value:    fmt.Sprintf("%d,%d", model.EventVisibilityPublic, model.EventVisibilityPrivate),
		IsSended: true,
	}

	return rc.list(ctx, opts)
}

func (rc *EventUC) GetByID(ctx context.Context, id string) (*model.Event, error) {
	event, err := rc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, pkg.NewError(err, "failed to get event", http.StatusInternalServerError)
	}

	return event, nil
}

func (rc *EventUC) list(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	list, err := rc.repo.List(ctx, opts)
	if err != nil {
		return nil, pkg.NewError(err, "failed to list events", http.StatusInternalServerError)
	}

	return list, nil
}
