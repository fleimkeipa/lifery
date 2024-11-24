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
}

func NewConnectHandlers(connectUC *uc.ConnectsUC) *ConnectHandlers {
	return &ConnectHandlers{
		connectUC: connectUC,
	}
}

// Create godoc
//
//	@Summary		Create creates a new connection
//	@Description	This endpoint creates a new connection by binding the incoming JSON request to the ConnectCreateRequest model.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
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
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	user, err := rc.connectUC.Create(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to create user: %v", err),
			Message: "Connect creation failed. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    user.Status,
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
//	@Param			id		path		string						true	"Connection ID to update"
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
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	user, err := rc.connectUC.Update(c.Request().Context(), id, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to update user: %v", err),
			Message: "Connect update failed. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    user.Status,
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
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	err := rc.connectUC.Disconnect(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to disconnect: %v", err),
			Message: "Failed to disconnect. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    "Disconnected",
		Message: "Disconnected successfully.",
	})
}

// List godoc
//
//	@Summary		List all connects
//	@Description	Retrieves a filtered and paginated list of connects based on query parameters.
//	@Tags			connects
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			status			query		string			false	"Filter connects by status"
//	@Param			friend_id		query		string			false	"Filter connects by friend ID"
//	@Param			limit			query		string			false	"Limit the number of connects returned"
//	@Param			skip			query		string			false	"Number of connects to skip for pagination"
//	@Success		200				{object}	SuccessResponse	"Successful response containing the list of connects"
//	@Failure		500				{object}	FailureResponse	"Internal error"
func (rc *ConnectHandlers) List(c echo.Context) error {
	opts := rc.getConnectsFindOpts(c, model.ZeroCreds)

	list, err := rc.connectUC.List(c.Request().Context(), &opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve connect list: %v", err),
			Message: "Unable to retrieve the list of connects. Please check the query parameters and try again.",
		})
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
