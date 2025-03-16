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

type EraUC struct {
	repo interfaces.EraRepository
}

func NewEraUC(repo interfaces.EraRepository) *EraUC {
	return &EraUC{
		repo: repo,
	}
}

func (rc *EraUC) Create(ctx context.Context, req *model.EraCreateRequest) (*model.Era, error) {
	userID := util.GetOwnerIDFromCtx(ctx)

	timeStart, err := time.Parse(`2006-01-02`, req.TimeStart)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
	}
	timeEnd, err := time.Parse(`2006-01-02`, req.TimeEnd)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeEnd", http.StatusBadRequest)
	}

	era := model.Era{
		TimeStart: timeStart.Format(`2006-01-02`),
		TimeEnd:   timeEnd.Format(`2006-01-02`),
		Name:      req.Name,
		Color:     req.Color,
		UserID:    userID,
	}

	newEra, err := rc.repo.Create(ctx, &era)
	if err != nil {
		return nil, err
	}

	return newEra, nil
}

func (rc *EraUC) Update(ctx context.Context, eraID string, req *model.EraUpdateRequest) (*model.Era, error) {
	// era exist control
	exist, err := rc.GetByID(ctx, eraID)
	if err != nil {
		return nil, err
	}

	timeStart, err := time.Parse(`2006-01-02`, req.TimeStart)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeStart", http.StatusBadRequest)
	}
	timeEnd, err := time.Parse(`2006-01-02`, req.TimeEnd)
	if err != nil {
		return nil, pkg.NewError(err, "failed to parse timeEnd", http.StatusBadRequest)
	}

	era := model.Era{
		TimeStart: timeStart.Format(`2006-01-02`),
		TimeEnd:   timeEnd.Format(`2006-01-02`),
		Name:      req.Name,
		Color:     req.Color,
		UserID:    exist.UserID,
	}

	updatedEra, err := rc.repo.Update(ctx, eraID, &era)
	if err != nil {
		return nil, err
	}

	return updatedEra, nil
}

func (rc *EraUC) Delete(ctx context.Context, id string) error {
	return rc.repo.Delete(ctx, id)
}

func (rc *EraUC) List(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	if ownerID == "" {
		if opts.UserID.Value == "" {
			return nil, pkg.NewError(nil, "user id is empty", http.StatusBadRequest)
		}

		opts.UserID = model.Filter{
			Value:    opts.UserID.Value,
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

	return rc.list(ctx, opts)
}

func (rc *EraUC) GetByID(ctx context.Context, id string) (*model.Era, error) {
	era, err := rc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return era, nil
}

func (rc *EraUC) list(ctx context.Context, opts *model.EraFindOpts) (*model.EraList, error) {
	list, err := rc.repo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list, nil
}
