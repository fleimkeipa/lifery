package interfaces

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) (*model.Event, error)
	Update(ctx context.Context, eventID string, event *model.Event) (*model.Event, error)
	Delete(ctx context.Context, eventID string) error
	List(ctx context.Context, opts *model.EventFindOpts) (*model.EventList, error)
	GetByID(ctx context.Context, eventID string) (*model.Event, error)
}
