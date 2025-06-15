package interfaces

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) (*model.Notification, error)
	Update(ctx context.Context, notificationID string, notification *model.Notification) (*model.Notification, error)
	List(ctx context.Context, opts *model.NotificationFindOpts) (*model.NotificationList, error)
	GetByID(ctx context.Context, notificationID string) (*model.Notification, error)
	Delete(ctx context.Context, notificationID string) error
}
