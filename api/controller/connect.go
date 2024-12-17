package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/uc"

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
//	@Description	This endpoint creates a new connection by binding the incoming JSON request to the ConnectCreateRequest model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.ConnectCreateRequest	true	"Connect creation input"
//	@Success		201		{object}	SuccessResponse				"Connect created successfully"
//	@Failure		400		{object}	FailureResponse				"Invalid request data"
//	@Failure		500		{object}	FailureResponse				"Connect creation failed"
//	@Router			/connects [post]
func (rc *ConnectHandlers) Create(c echo.Context) error {
	var input model.ConnectCreateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	connect, err := rc.connectUC.Create(c.Request().Context(), input)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    connect.Status,
		Message: "Connect created successfully.",
	})
}

// Update godoc
//
//	@Summary		Update updates an existing connection
//	@Description	This endpoint updates a connection by binding the incoming JSON request to the ConnectUpdateRequest model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string						true	"Connection ID to update,approved:101, rejected:102"
//	@Param			Body	body		model.ConnectUpdateRequest	true	"Connect update input"
//	@Success		200		{object}	SuccessResponse				"Connect updated successfully"
//	@Failure		400		{object}	FailureResponse				"Invalid request data"
//	@Failure		500		{object}	FailureResponse				"Connect update failed"
//	@Router			/connects/{id} [patch]
func (rc *ConnectHandlers) Update(c echo.Context) error {
	id := c.Param("id")
	var input model.ConnectUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	connect, err := rc.connectUC.Update(c.Request().Context(), id, input)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    connect.Status,
		Message: "Connect updated successfully.",
	})
}

// Disconnect godoc
//
//	@Summary		Disconnects an existing connection
//	@Description	This endpoint disconnects an existing connection by binding the incoming JSON request to the DisconnectRequest model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.DisconnectRequest	true	"Disconnect input"
//	@Success		200		{object}	SuccessResponse			"Disconnected successfully"
//	@Failure		400		{object}	FailureResponse			"Invalid request data"
//	@Failure		500		{object}	FailureResponse			"Disconnect failed"
//	@Router			/connects/disconnect [patch]
func (rc *ConnectHandlers) Disconnect(c echo.Context) error {
	var input model.DisconnectRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	if err := rc.connectUC.Disconnect(c.Request().Context(), input); err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    "Disconnected",
		Message: "Disconnected successfully.",
	})
}

// ConnectsRequests godoc
//
//	@Summary		ConnectsRequests all connects
//	@Description	Retrieves a filtered and paginated list of connects based on query parameters.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status		query		string			false	"Filter connects by status"
//	@Param			friend_id	query		string			false	"Filter connects by friend ID"
//	@Param			limit		query		string			false	"Limit the number of connects returned"
//	@Param			skip		query		string			false	"Number of connects to skip for pagination"
//	@Success		200			{object}	SuccessResponse	"Successful response containing the list of connects"
//	@Failure		500			{object}	FailureResponse	"Internal error"
func (rc *ConnectHandlers) ConnectsRequests(c echo.Context) error {
	opts := rc.getConnectsFindOpts(c, model.ZeroCreds)

	list, err := rc.connectUC.List(c.Request().Context(), &opts)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Connects requests retrieved successfully.",
	})
}

// GetConnects godoc
//
//	@Summary		List all connects
//	@Description	Retrieves a filtered and paginated list of connects based on query parameters.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			limit	query		string			false	"Limit the number of connects returned"
//	@Param			skip	query		string			false	"Number of connects to skip for pagination"
//	@Success		200		{object}	SuccessResponse	"Successful response containing the list of connects"
//	@Failure		500		{object}	FailureResponse	"Internal error"
//	@Router			/connects/{user_id} [get]
func (rc *ConnectHandlers) GetConnects(c echo.Context) error {
	id := c.Param("user_id")

	list, err := rc.userUC.GetConnects(c.Request().Context(), id)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Connects retrieved successfully.",
	})
}

func (rc *ConnectHandlers) getConnectsFindOpts(c echo.Context, fields ...string) model.ConnectFindOpts {
	return model.ConnectFindOpts{
		PaginationOpts: getPagination(c),
		FieldsOpts: model.FieldsOpts{
			Fields: fields,
		},
		Status:   getFilter(c, "status"),
		FriendID: getFilter(c, "friend_id"),
	}
}
