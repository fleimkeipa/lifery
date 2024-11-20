package uc

import (
	"context"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
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
	era := model.Era{
		Name:      req.Name,
		TimeRange: req.TimeRange,
	}

	return rc.repo.Create(ctx, &era)
}

func (rc *EraUC) Update(ctx context.Context, eraID string, req *model.EraUpdateRequest) (*model.Era, error) {
	// era exist control
	_, err := rc.GetByID(ctx, eraID)
	if err != nil {
		return nil, err
	}

	era := model.Era{
		Name:      req.Name,
		TimeRange: req.TimeRange,
	}

	return rc.repo.Update(ctx, eraID, &era)
}

func (rc *EraUC) Delete(ctx context.Context, id string) error {
	return rc.repo.Delete(ctx, id)
}

func (rc *EraUC) List(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error) {
	return rc.repo.List(ctx, opts)
}

func (rc *EraUC) GetByID(ctx context.Context, id string) (*model.Era, error) {
	return rc.repo.GetByID(ctx, id)
}
