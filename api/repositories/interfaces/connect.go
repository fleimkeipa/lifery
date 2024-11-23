package interfaces

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
)

type ConnectInterfaces interface {
	Create(ctx context.Context, connect *model.Connect) (*model.Connect, error)
	Update(ctx context.Context, connectID string, connect *model.Connect) (*model.Connect, error)
	List(ctx context.Context, opts *model.ConnectFindOpts) (*model.ConnectList, error)
	GetByID(ctx context.Context, connectID string) (*model.Connect, error)
	Delete(ctx context.Context, connectID string) error
}
