package interfaces

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
)

type EraRepository interface {
	Create(ctx context.Context, event *model.Era) (*model.Era, error)
	Update(ctx context.Context, eraID string, event *model.Era) (*model.Era, error)
	Delete(ctx context.Context, eraID string) error
	List(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error)
	GetByID(ctx context.Context, eraID string) (*model.Era, error)
}
