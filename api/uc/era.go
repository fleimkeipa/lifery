package uc

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/util"
)

type EraUC struct {
	repo interfaces.EraRepository
}

func NewEraUC(repo interfaces.EraRepository) *EraUC {
	return &EraUC{
		repo: repo,
	}
}

func (rc *EraUC) Create(ctx context.Context, req *model.EraCreateRequest) (*model.Era, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	era := model.Era{
		TimeStart: req.TimeStart,
		TimeEnd:   req.TimeEnd,
		Name:      req.Name,
		Color:     req.Color,
		OwnerID:   ownerID,
	}

	return rc.repo.Create(ctx, &era)
}

func (rc *EraUC) Update(ctx context.Context, eraID string, req *model.EraUpdateRequest) (*model.Era, error) {
	// era exist control
	exist, err := rc.GetByID(ctx, eraID)
	if err != nil {
		return nil, err
	}

	era := model.Era{
		TimeStart: req.TimeStart,
		TimeEnd:   req.TimeEnd,
		Name:      req.Name,
		Color:     req.Color,
		OwnerID:   exist.OwnerID,
	}

	return rc.repo.Update(ctx, eraID, &era)
}

func (rc *EraUC) Delete(ctx context.Context, id string) error {
	return rc.repo.Delete(ctx, id)
}

func (rc *EraUC) List(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	if !opts.UserID.IsSended {
		opts.UserID = model.Filter{
			Value:    ownerID,
			IsSended: true,
		}

		return rc.repo.List(ctx, opts)
	}

	return rc.repo.List(ctx, opts)
}

func (rc *EraUC) GetByID(ctx context.Context, id string) (*model.Era, error) {
	return rc.repo.GetByID(ctx, id)
}
