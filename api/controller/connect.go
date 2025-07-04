package controller

import (
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/uc"
	"github.com/fleimkeipa/lifery/util"

	"github.com/labstack/echo/v4"
)

type ConnectHandlers struct {
	connectUC *uc.ConnectsUC
	userUC    *uc.UserUC
}

func NewConnectHandlers(connectUC *uc.ConnectsUC, userUC *uc.UserUC) *ConnectHandlers {
	return &ConnectHandlers{
		connectUC: connectUC,
		userUC:    userUC,
	}
}

// Create godoc
//
//	@Summary		Create creates a new connection
//	@Description	This endpoint creates a new connection by binding the incoming JSON request to the ConnectCreateInput model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.ConnectCreateInput	true	"Connect creation input"
//	@Success		201		{object}	SuccessResponse				"Connect created successfully"
//	@Failure		400		{object}	FailureResponse				"Invalid request data"
//	@Failure		500		{object}	FailureResponse				"Connect creation failed"
//	@Router			/connects [post]
func (rc *ConnectHandlers) Create(c echo.Context) error {
	var input model.ConnectCreateInput

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	_, err := rc.connectUC.Create(c.Request().Context(), input)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Connect created successfully",
	})
}

// Update godoc
//
//	@Summary		Update updates an existing connection
//	@Description	This endpoint updates a connection by binding the incoming JSON request to the ConnectUpdateInput model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string						true	"Connection ID to update,approved:101, rejected:102"
//	@Param			Body	body		model.ConnectUpdateInput	true	"Connect update input"
//	@Success		200		{object}	SuccessResponse				"Connect updated successfully"
//	@Failure		400		{object}	FailureResponse				"Invalid request data"
//	@Failure		500		{object}	FailureResponse				"Connect update failed"
//	@Router			/connects/{id} [patch]
func (rc *ConnectHandlers) Update(c echo.Context) error {
	id := c.Param("id")
	var input model.ConnectUpdateInput

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	err := rc.connectUC.Update(c.Request().Context(), id, input)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Connect updated successfully",
	})
}

// Delete godoc
//
//	@Summary		Delete deletes an existing connection
//	@Description	This endpoint deletes a connection by binding the incoming JSON request to the ConnectUpdateInput model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string			true	"Connection ID to delete"
//	@Success		200	{object}	SuccessResponse	"Connect deleted successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Connect update failed"
//	@Router			/connects/{id} [delete]
func (rc *ConnectHandlers) Delete(c echo.Context) error {
	id := c.Param("id")

	err := rc.connectUC.Delete(c.Request().Context(), id)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Connect deleted successfully",
	})
}

// ConnectsRequests godoc
//
//	@Summary		ConnectsRequests list all connects requests
//	@Description	Retrieves a filtered and paginated list of connects requests based on query parameters.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status	query		string			false	"Filter connects by status"
//	@Param			user_id	query		string			false	"Filter connects by user id if you are admin"
//	@Param			limit	query		string			false	"Limit the number of connects returned"
//	@Param			skip	query		string			false	"Number of connects to skip for pagination"
//	@Success		200		{object}	SuccessResponse	"Successful response containing the list of connects"
//	@Failure		500		{object}	FailureResponse	"Internal error"
//	@Router			/connects [get]
func (rc *ConnectHandlers) ConnectsRequests(c echo.Context) error {
	opts := rc.getConnectsFindOpts(c, model.ZeroCreds)

	list, err := rc.connectUC.ConnectsRequests(c.Request().Context(), &opts)
	if err != nil {
		return handleEchoError(c, err)
	}

	rc.relocateFriend(c, list)

	return c.JSON(http.StatusOK, SuccessListResponse{
		Data:  list.Connects,
		Total: list.Total,
		Limit: list.Limit,
		Skip:  list.Skip,
	})
}

func (rc *ConnectHandlers) relocateFriend(c echo.Context, list *model.ConnectList) {
	ownerID := util.GetOwnerIDFromCtx(c.Request().Context())

	for i, connect := range list.Connects {
		if connect.FriendID == ownerID {
			friend := connect.Friend
			connect.Friend = connect.User
			connect.User = friend
			list.Connects[i] = connect
		}
	}
}

func (rc *ConnectHandlers) getConnectsFindOpts(c echo.Context, fields ...string) model.ConnectFindOpts {
	defaultFilter := model.ConnectFindOpts{
		OrderByOpts:    getOrder(c),
		PaginationOpts: getPagination(c),
		FieldsOpts: model.FieldsOpts{
			Fields: fields,
		},
		Status:   getFilter(c, "status"),
		Username: getFilter(c, "username"),
	}

	owner, err := util.GetOwnerFromToken(c)
	if err != nil {
		return defaultFilter
	}

	defaultFilter.UserID = getFilter(c, "user_id")

	if !defaultFilter.UserID.IsSended {
		defaultFilter.UserID = model.Filter{
			Value:    owner.ID,
			IsSended: owner.ID != "",
		}

		return defaultFilter
	}

	if owner.RoleID == model.AdminRole {
		return defaultFilter
	}

	defaultFilter.UserID = model.Filter{
		Value:    owner.ID,
		IsSended: owner.ID != "",
	}

	return defaultFilter
}
