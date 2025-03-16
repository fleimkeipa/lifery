package uc

import (
	"context"
	"net/http"
	"strconv"

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

func (rc *ConnectsUC) Create(ctx context.Context, req model.ConnectCreateRequest) (*model.Connect, error) {
	ownerID := util.GetOwnerIDFromCtx(ctx)
	connect := model.Connect{
		Status:   model.RequestStatusPending,
		UserID:   ownerID,
		FriendID: req.FriendID,
	}

	if ownerID == req.FriendID {
		return nil, pkg.NewError(nil, "cannot connect to self", http.StatusBadRequest)
	}

	// owner control
	if !rc.isOwner(ctx, ownerID) {
		return nil, pkg.NewError(nil, "you can't connect to other users", http.StatusBadRequest)
	}

	// receiver exist control
	receiver, err := rc.userUC.GetByID(ctx, ownerID)
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
			return nil, pkg.NewError(nil, "already connected", http.StatusBadRequest)
		}
	}

	return rc.connectRepo.Create(ctx, &connect)
}

func (rc *ConnectsUC) Update(ctx context.Context, id string, req model.ConnectUpdateRequest) error {
	connect, err := rc.connectRepo.GetByID(ctx, id)
	if err != nil {
		return pkg.NewError(nil, "connect not found", http.StatusNotFound)
	}

	connect.Status = req.Status

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
	if !rc.isOwner(ctx, connect.FriendID) {
		return pkg.NewError(nil, "you can update only your connects", http.StatusBadRequest)
	}

	if req.Status == model.RequestStatusRejected {
		if err := rc.Delete(ctx, id); err != nil {
			return err
		}

		return nil
	}

	_, err = rc.AddConnectOnUser(ctx, id, connect.UserID, connect.FriendID)
	if err != nil {
		return err
	}

	return rc.Delete(ctx, id)
}

func (rc *ConnectsUC) Disconnect(ctx context.Context, req model.DisconnectRequest) error {
	ownerID := util.GetOwnerIDFromCtx(ctx)

	_, err := rc.DeleteConnectOnUser(ctx, ownerID, req.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func (rc *ConnectsUC) ConnectsRequests(ctx context.Context, opts *model.ConnectFindOpts) (*model.ConnectList, error) {
	return rc.connectRepo.ConnectsRequests(ctx, opts)
}

func (rc *ConnectsUC) GetByID(ctx context.Context, id string) (*model.Connect, error) {
	return rc.connectRepo.GetByID(ctx, id)
}

func (rc *ConnectsUC) Delete(ctx context.Context, id string) error {
	return rc.connectRepo.Delete(ctx, id)
}

func (rc *ConnectsUC) AddConnectOnUser(ctx context.Context, connectID string, userID, friendID string) (*model.User, error) {
	// receiver exist control
	receiver, err := rc.userUC.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// sender exist control
	sender, err := rc.userUC.GetByID(ctx, friendID)
	if err != nil {
		return nil, err
	}

	receiver, err = rc.userUC.UpdateConnects(ctx, receiver, sender.ID, receiver.ID)
	if err != nil {
		return nil, err
	}

	sender, err = rc.userUC.UpdateConnects(ctx, sender, receiver.ID, sender.ID)
	if err != nil {
		return nil, err
	}

	return sender, nil
}

func (rc *ConnectsUC) DeleteConnectOnUser(ctx context.Context, userID, friendID string) (*model.User, error) {
	// receiver exist control
	receiver, err := rc.userUC.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// sender exist control
	sender, err := rc.userUC.GetByID(ctx, friendID)
	if err != nil {
		return nil, err
	}

	receiver, err = rc.userUC.DeleteUserConnect(ctx, receiver, sender.ID)
	if err != nil {
		return nil, err
	}

	return rc.userUC.DeleteUserConnect(ctx, sender, receiver.ID)
}

func (rc *ConnectsUC) isOwner(ctx context.Context, id string) bool {
	ownerID := util.GetOwnerIDFromCtx(ctx)
	if id != ownerID {
		return false
	}

	return true
}
