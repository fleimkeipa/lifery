package controller

import (
	"fmt"
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
//	@Description	This endpoint creates a new era by binding the incoming JSON request to the EraCreateRequest model.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Param			Body	body		model.EraCreateRequest	true	"Era creation input"
//	@Success		201		{object}	SuccessResponse			"Era created successfully"
//	@Failure		400		{object}	FailureResponse			"Invalid request data"
//	@Failure		500		{object}	FailureResponse			"Era creation failed"
//	@Router			/eras [post]
func (rc *EraController) Create(c echo.Context) error {
	var request model.EraCreateRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	era, err := rc.EraDBUC.Create(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to create era: %v", err),
			Message: "Era creation failed. Please verify the details and try again.",
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    era.Name,
		Message: "Era created successfully.",
	})
}

// Update handles the update of an existing era.
//
//	@Summary		Update an existing era
//	@Description	This endpoint updates an existing era by binding the incoming JSON request to the EraUpdateRequest model.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Param			Body	body		model.EraUpdateRequest	true	"Era update input"
//	@Success		200		{object}	SuccessResponse			"Era updated successfully"
//	@Failure		400		{object}	FailureResponse			"Invalid request data"
//	@Failure		500		{object}	FailureResponse			"Era update failed"
//	@Router			/eras/{id} [patch]
func (rc *EraController) Update(c echo.Context) error {
	eraID := c.Param("id")
	var request model.EraUpdateRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	era, err := rc.EraDBUC.Update(c.Request().Context(), eraID, &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to update era: %v", err),
			Message: "Era creation failed. Please verify the details and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    era.Name,
		Message: "Era updated successfully.",
	})
}

// Delete handles the deletion of an existing era.
//
//	@Summary		Delete an existing era
//	@Description	This endpoint deletes an existing era by providing era name or UID.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Era name or UID"
//	@Success		200	{object}	SuccessResponse	"Era deleted successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Era delete failed"
//	@Router			/eras/{id} [delete]
func (rc *EraController) Delete(c echo.Context) error {
	nameOrUID := c.Param("id")

	err := rc.EraDBUC.Delete(c.Request().Context(), nameOrUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve era: %v", err),
			Message: "Error fetching the era details. Please verify the era name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Era retrieved successfully.",
	})
}

// List handles the retrieval of a list of eras.
//
//	@Summary		Retrieve a list of eras
//	@Description	This endpoint retrieves a list of eras.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string			false	"Filter eras by name"
//	@Param			user_id	query		string			false	"Filter eras by user id"
//	@Param			limit	query		string			false	"Limit the number of connects returned"
//	@Param			skip	query		string			false	"Number of connects to skip for pagination"
//	@Success		200		{object}	SuccessResponse	"Eras retrieved successfully"
//	@Failure		400		{object}	FailureResponse	"Invalid request data"
//	@Failure		500		{object}	FailureResponse	"Era retrieval failed"
//	@Router			/eras [get]
func (rc *EraController) List(c echo.Context) error {
	opts := rc.getErasFindOpts(c)

	list, err := rc.EraDBUC.List(c.Request().Context(), &opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to list eras: %v", err),
			Message: "There was an issue retrieving eras. Please try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Eras retrieved successfully.",
	})
}

// GetByID handles the retrieval of an era by its name or UID.
//
//	@Summary		Retrieve era by ID
//	@Description	Fetches an era by its unique name or UID from the database.
//	@Tags			eras
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Era name or UID"
//	@Success		200	{object}	SuccessResponse	"Era retrieved successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Era retrieval failed"
//	@Router			/eras/{id} [get]
func (rc *EraController) GetByID(c echo.Context) error {
	nameOrUID := c.Param("id")

	era, err := rc.EraDBUC.GetByID(c.Request().Context(), nameOrUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve era: %v", err),
			Message: "Error fetching the era details. Please verify the era name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    era,
		Message: "Era retrieved successfully.",
	})
}

func (rc *EraController) getErasFindOpts(c echo.Context) model.EraFindOpts {
	return model.EraFindOpts{
		PaginationOpts: getPagination(c),
		Name:           getFilter(c, "name"),
		UserID:         getFilter(c, "user_id"),
	}
}
