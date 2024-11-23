package uc

import (
	"context"
	"errors"
	"strconv"

	"github.com/fleimkeipa/lifery/model"
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

func (rc *ConnectsUC) Create(ctx context.Context, req model.ConnectCreateRequest) (*model.Connect, error) {
	connect := model.Connect{
		Status:   model.RequestStatusPending,
		UserID:   req.UserID,
		FriendID: req.FriendID,
	}

	if req.UserID == req.FriendID {
		return nil, errors.New("cannot connect to self")
	}

	// owner control
	if !rc.isOwner(ctx, req.UserID) {
		return nil, errors.New("you can't connect to other users")
	}

	// receiver exist control
	receiver, err := rc.userUC.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// sender exist control
	_, err = rc.userUC.GetByID(ctx, req.FriendID)
	if err != nil {
		return nil, err
	}

	connects := receiver.Connects
	for _, v := range connects {
		strID := strconv.Itoa(v)
		if strID == req.FriendID {
			return nil, errors.New("already connected")
		}
	}

	return rc.connectRepo.Create(ctx, &connect)
}

func (rc *ConnectsUC) Update(ctx context.Context, id string, req model.ConnectUpdateRequest) (*model.Connect, error) {
	connect, err := rc.connectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	connect.Status = req.Status

	if req.Status == "" {
		return nil, errors.New("status is required")
	}

	if req.Status == model.RequestStatusPending {
		return nil, errors.New("status is pending already")
	}

	if req.Status != model.RequestStatusApproved && req.Status != model.RequestStatusRejected {
		return nil, errors.New("invalid status")
	}

	// owner control
	if !rc.isOwner(ctx, connect.UserID) {
		return nil, errors.New("you can update only your connects")
	}

	if req.Status == model.RequestStatusRejected {
		if err := rc.Delete(ctx, id); err != nil {
			return nil, err
		}

		return connect, nil
	}

	_, err = rc.userUC.AddConnect(ctx, connect.UserID, connect.FriendID)
	if err != nil {
		return nil, err
	}

	if err := rc.Delete(ctx, id); err != nil {
		return nil, err
	}

	return connect, nil
}

func (rc *ConnectsUC) Disconnect(ctx context.Context, req model.DisconnectRequest) error {
	owner := util.GetOwnerFromCtx(ctx)

	_, err := rc.userUC.DeleteConnect(ctx, strconv.Itoa(int(owner.ID)), req.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func (rc *ConnectsUC) List(ctx context.Context, opts *model.ConnectFindOpts) (*model.ConnectList, error) {
	return rc.connectRepo.List(ctx, opts)
}

func (rc *ConnectsUC) GetByID(ctx context.Context, id string) (*model.Connect, error) {
	return rc.connectRepo.GetByID(ctx, id)
}

func (rc *ConnectsUC) Delete(ctx context.Context, id string) error {
	return rc.connectRepo.Delete(ctx, id)
}

func (rc *ConnectsUC) isOwner(ctx context.Context, id string) bool {
	tokenOwner := util.GetOwnerFromCtx(ctx)
	if id != strconv.Itoa(int(tokenOwner.ID)) {
		return false
	}

	return true
}
