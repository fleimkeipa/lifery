package uc

import (
	"context"
	"net/http"
	"time"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/util"
)

type NotificationUC struct {
	repo interfaces.NotificationRepository
}

func NewNotificationUC(repo interfaces.NotificationRepository) *NotificationUC {
	return &NotificationUC{
		repo: repo,
	}
}

func (rc *NotificationUC) Create(ctx context.Context, req model.NotificationCreateInput) (*model.Notification, error) {
	createdAt, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse start time", http.StatusBadRequest)
	}

	notification := model.Notification{
		UserID:    req.UserID,
		Type:      req.Type,
		Message:   req.Message,
		Read:      model.NotificationStatusUnread,
		CreatedAt: createdAt,
	}

	return rc.repo.Create(ctx, &notification)
}

func (rc *NotificationUC) Update(ctx context.Context, id string, req model.NotificationUpdateInput) error {
	existNotification, err := rc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !rc.isOwner(ctx, existNotification.UserID) {
		return pkg.NewError(nil, "you can update only your notifications", http.StatusForbidden)
	}

	existNotification.Read = req.Read

	_, err = rc.repo.Update(ctx, id, existNotification)
	if err != nil {
		return err
	}

	return nil
}

func (rc *NotificationUC) List(ctx context.Context, opts *model.NotificationFindOpts) (*model.NotificationList, error) {
	if err := rc.checkOwner(ctx, opts); err != nil {
		return nil, err
	}

	return rc.repo.List(ctx, opts)
}

func (rc *NotificationUC) GetByID(ctx context.Context, id string) (*model.Notification, error) {
	return rc.repo.GetByID(ctx, id)
}

func (rc *NotificationUC) Delete(ctx context.Context, id string) error {
	notification, err := rc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !rc.isOwner(ctx, notification.UserID) {
		return pkg.NewError(nil, "you can delete only your notifications", http.StatusForbidden)
	}

	return rc.repo.Delete(ctx, id)
}

func (rc *NotificationUC) isOwner(ctx context.Context, id string) bool {
	ownerID := util.GetOwnerIDFromCtx(ctx)
	return id == ownerID
}

func (rc *NotificationUC) checkOwner(ctx context.Context, opts *model.NotificationFindOpts) error {
	ownerID := util.GetOwnerIDFromCtx(ctx)
	if !opts.UserID.IsSended {
		opts.UserID = model.Filter{
			Value:    ownerID,
			IsSended: true,
		}

		return nil
	}

	if opts.UserID.Value == ownerID {
		return nil
	}

	owner := util.GetOwnerFromCtx(ctx)
	if owner.RoleID != model.AdminRole {
		return pkg.NewError(nil, "you cannot get another users notifications", http.StatusForbidden)
	}

	return nil
}
