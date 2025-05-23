package uc

import (
	"context"
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/repositories/interfaces"
	"github.com/fleimkeipa/lifery/util"
)

type ConnectsUC struct {
	userUC      *UserUC
	connectRepo interfaces.ConnectInterfaces
}

func NewConnectsUC(userUC *UserUC, connectRepo interfaces.ConnectInterfaces) *ConnectsUC {
	return &ConnectsUC{
		userUC:      userUC,
		connectRepo: connectRepo,
	}
}

func (rc *ConnectsUC) Create(ctx context.Context, req model.ConnectCreateInput) (*model.Connect, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)
	connect := model.Connect{
		Status:   model.RequestStatusPending,
		UserID:   ownerID,
		FriendID: req.FriendID,
	}

	if ownerID == req.FriendID {
		return nil, pkg.NewError(nil, "cannot connect to self", http.StatusBadRequest)
	}

	// sender exist control
	_, err := rc.userUC.GetByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	// receiver exist control
	_, err = rc.userUC.GetByID(ctx, req.FriendID)
	if err != nil {
		return nil, err
	}

	connectsFindOpts := model.ConnectFindOpts{
		UserID: model.Filter{
			Value:    ownerID,
			IsSended: true,
		},
	}
	connectsList, err := rc.connectRepo.ConnectsRequests(ctx, &connectsFindOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range connectsList.Connects {
		if v.FriendID == req.FriendID {
			switch v.Status {
			case model.RequestStatusApproved:
				return nil, pkg.NewError(nil, "already connected", http.StatusBadRequest)
			case model.RequestStatusPending:
				return nil, pkg.NewError(nil, "already sended request", http.StatusBadRequest)
			}
		}
	}

	return rc.connectRepo.Create(ctx, &connect)
}

func (rc *ConnectsUC) Update(ctx context.Context, id string, req model.ConnectUpdateInput) error {
	existConnect, err := rc.connectRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existConnect.Status = req.Status

	if req.Status == 0 {
		return pkg.NewError(nil, "status is required", http.StatusBadRequest)
	}

	if req.Status == model.RequestStatusPending {
		return pkg.NewError(nil, "status is pending already", http.StatusBadRequest)
	}

	if req.Status != model.RequestStatusApproved && req.Status != model.RequestStatusRejected {
		return pkg.NewError(nil, "invalid status", http.StatusBadRequest)
	}

	// users can only update their own connects
	if !rc.isOwner(ctx, existConnect.FriendID) {
		return pkg.NewError(nil, "you can update only your connects", http.StatusBadRequest)
	}

	if req.Status == model.RequestStatusRejected {
		return rc.Delete(ctx, id)
	}

	_, err = rc.connectRepo.Update(ctx, id, existConnect)
	if err != nil {
		return err
	}

	return nil
}

func (rc *ConnectsUC) ConnectsRequests(ctx context.Context, opts *model.ConnectFindOpts) (*model.ConnectList, error) {
	if err := rc.checkOwner(ctx, opts); err != nil {
		return nil, err
	}

	return rc.connectRepo.ConnectsRequests(ctx, opts)
}

func (rc *ConnectsUC) GetByID(ctx context.Context, id string) (*model.Connect, error) {
	return rc.connectRepo.GetByID(ctx, id)
}

func (rc *ConnectsUC) Delete(ctx context.Context, id string) error {
	return rc.connectRepo.Delete(ctx, id)
}

func (rc *ConnectsUC) AddConnectOnUser(ctx context.Context, connectID string, userID, friendID string) error {
	// receiver exist control
	receiver, err := rc.userUC.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// sender exist control
	sender, err := rc.userUC.GetByID(ctx, friendID)
	if err != nil {
		return err
	}

	err = rc.userUC.UpdateConnects(ctx, receiver, sender.ID, receiver.ID)
	if err != nil {
		return err
	}

	return rc.userUC.UpdateConnects(ctx, sender, receiver.ID, sender.ID)
}

func (rc *ConnectsUC) DeleteConnectOnUser(ctx context.Context, userID, friendID string) error {
	// receiver exist control
	receiver, err := rc.userUC.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// sender exist control
	sender, err := rc.userUC.GetByID(ctx, friendID)
	if err != nil {
		return err
	}

	err = rc.userUC.DeleteUserConnect(ctx, receiver, sender.ID)
	if err != nil {
		return err
	}

	return rc.userUC.DeleteUserConnect(ctx, sender, receiver.ID)
}

func (rc *ConnectsUC) isOwner(ctx context.Context, id string) bool {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	return id == ownerID
}

func (rc *ConnectsUC) checkOwner(ctx context.Context, opts *model.ConnectFindOpts) error {
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
		return pkg.NewError(nil, "you cannot get another users connects", http.StatusForbidden)
	}

	return nil
}
