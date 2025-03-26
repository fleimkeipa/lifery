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

	var err error
	var date, timeStart, timeEnd time.Time
	if req.Date != "" {
		date, err = util.ParseTime(req.Date)
		if err != nil {
			return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
		}
	}
	if req.TimeStart != "" && req.TimeEnd != "" {
		timeStart, err = util.ParseTime(req.TimeStart)
		if err != nil {
			return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
		}
		timeEnd, err = util.ParseTime(req.TimeEnd)
		if err != nil {
			return nil, pkg.NewError(err, "failed to parse timeEnd", http.StatusBadRequest)
		}
	}

	event := model.Event{
		Date:        date,
		TimeStart:   timeStart,
		TimeEnd:     timeEnd,
		Name:        req.Name,
		Description: req.Description,
		Items:       req.Items,
		UserID:      ownerID,
		Visibility:  req.Visibility,
		CreatedAt:   util.Now(),
	}

	newEvent, err := rc.repo.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (rc *EventUC) Update(ctx context.Context, eventID string, req *model.EventUpdateRequest) (*model.Event, error) {
	// event exist control
	exist, err := rc.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	date, err := util.ParseTime(req.Date)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse date", http.StatusBadRequest)
	}
	timeStart, err := util.ParseTime(req.TimeStart)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
	}
	timeEnd, err := util.ParseTime(req.TimeEnd)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeEnd", http.StatusBadRequest)
	}

	event := model.Event{
		Date:        date,
		TimeStart:   timeStart,
		TimeEnd:     timeEnd,
		Name:        req.Name,
		Description: req.Description,
		Items:       req.Items,
		UserID:      exist.UserID,
		Visibility:  req.Visibility,
		CreatedAt:   exist.CreatedAt,
		UpdatedAt:   util.Now(),
	}

	updatedEvent, err := rc.repo.Update(ctx, eventID, &event)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (rc *EventUC) Delete(ctx context.Context, id string) error {
	return rc.repo.Delete(ctx, id)
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
	return rc.repo.GetByID(ctx, id)
}

func (rc *EventUC) list(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error) {
	return rc.repo.List(ctx, opts)
}
