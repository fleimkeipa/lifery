package controller

import (
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/uc"

	"github.com/labstack/echo/v4"
)

type EraController struct {
	EraDBUC *uc.EraUC
}

func NewEraController(eraUC *uc.EraUC) *EraController {
	return &EraController{EraDBUC: eraUC}
}

// Create handles the creation of a new era.
//
//	@Summary		Create a new era
//	@Description	This endpoint creates a new era by binding the incoming JSON request to the EraCreateInput model.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.EraCreateInput	true	"Era creation input"
//	@Success		201		{object}	SuccessResponse			"Era created successfully"
//	@Failure		400		{object}	FailureResponse			"Invalid request data"
//	@Failure		500		{object}	FailureResponse			"Era creation failed"
//	@Router			/eras [post]
func (rc *EraController) Create(c echo.Context) error {
	var input model.EraCreateInput

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	_, err := rc.EraDBUC.Create(c.Request().Context(), &input)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Era created successfully",
	})
}

// Update handles the update of an existing era.
//
//	@Summary		Update an existing era
//	@Description	This endpoint updates an existing era by binding the incoming JSON request to the EraUpdateInput model.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.EraUpdateInput	true	"Era update input"
//	@Success		200		{object}	SuccessResponse			"Era updated successfully"
//	@Failure		400		{object}	FailureResponse			"Invalid request data"
//	@Failure		500		{object}	FailureResponse			"Era update failed"
//	@Router			/eras/{id} [patch]
func (rc *EraController) Update(c echo.Context) error {
	eraID := c.Param("id")
	var input model.EraUpdateInput

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	_, err := rc.EraDBUC.Update(c.Request().Context(), eraID, &input)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Era updated successfully",
	})
}

// Delete handles the deletion of an existing era.
//
//	@Summary		Delete an existing era
//	@Description	This endpoint deletes an existing era by providing era name or UID.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string			true	"Era name or UID"
//	@Success		200	{object}	SuccessResponse	"Era deleted successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Era delete failed"
//	@Router			/eras/{id} [delete]
func (rc *EraController) Delete(c echo.Context) error {
	eraID := c.Param("id")

	if err := rc.EraDBUC.Delete(c.Request().Context(), eraID); err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Era deleted successfully",
	})
}

// List handles the retrieval of a list of eras.
//
//	@Summary		Retrieve a list of eras
//	@Description	This endpoint retrieves a list of eras.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			false	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			name			query		string			false	"Filter eras by name"		example(eq:test)
//	@Param			user_id			query		string			false	"Filter eras by user id"	example(eq:1)
//	@Param			limit			query		string			false	"Limit the number of connects returned"
//	@Param			skip			query		string			false	"Number of connects to skip for pagination"
//	@Param			order			query		string			false	"Order by column (prefix with asc: or desc:)"	example(desc:created_at)
//	@Success		200				{object}	SuccessListResponse	"Eras retrieved successfully"
//	@Failure		400				{object}	FailureResponse	"Invalid request data"
//	@Failure		500				{object}	FailureResponse	"Era retrieval failed"
//	@Router			/eras [get]
func (rc *EraController) List(c echo.Context) error {
	opts := rc.getErasFindOpts(c)

	list, err := rc.EraDBUC.List(c.Request().Context(), &opts)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessListResponse{
		Data:  list.Eras,
		Total: list.Total,
		Limit: list.Limit,
		Skip:  list.Skip,
	})
}

// GetByID handles the retrieval of an era by its name or UID.
//
//	@Summary		Retrieve era by ID
//	@Description	Fetches an era by its unique name or UID from the database.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string			true	"Era name or UID"
//	@Success		200	{object}	SuccessResponse	"Era retrieved successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Era retrieval failed"
//	@Router			/eras/{id} [get]
func (rc *EraController) GetByID(c echo.Context) error {
	eraID := c.Param("id")

	era, err := rc.EraDBUC.GetByID(c.Request().Context(), eraID)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessListResponse{
		Data: era,
	})
}

func (rc *EraController) getErasFindOpts(c echo.Context) model.EraFindOpts {
	return model.EraFindOpts{
		OrderByOpts:    getOrder(c),
		PaginationOpts: getPagination(c),
		Name:           getFilter(c, "name"),
		UserID:         getFilter(c, "user_id"),
	}
}
