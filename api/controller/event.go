package controller

import (
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/uc"

	"github.com/labstack/echo/v4"
)

type EventController struct {
	EventDBUC *uc.EventUC
}

func NewEventController(eventUC *uc.EventUC) *EventController {
	return &EventController{EventDBUC: eventUC}
}

// Create handles the creation of a new event.
//
//	@Summary		Create a new event
//	@Description	This endpoint creates a new event by binding the incoming JSON request to the EventCreateRequest model.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.EventCreateRequest	true	"Event creation input"
//	@Success		201		{object}	SuccessResponse				"Event created successfully"
//	@Failure		400		{object}	FailureResponse				"Invalid request data"
//	@Failure		500		{object}	FailureResponse				"Event creation failed"
//	@Router			/events [post]
func (rc *EventController) Create(c echo.Context) error {
	var request model.EventCreateRequest

	if err := c.Bind(&request); err != nil {
		err := pkg.NewError(
			err,
			"Invalid request format. Please check the input data and try again.",
			http.StatusBadRequest,
		)
		return HandleEchoError(c, err)
	}

	event, err := rc.EventDBUC.Create(c.Request().Context(), &request)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    event.Name,
		Message: "Event created successfully.",
	})
}

// Update handles the update of an existing event.
//
//	@Summary		Update an existing event
//	@Description	This endpoint updates an existing event by binding the incoming JSON request to the EventUpdateRequest model.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Body	body		model.EventUpdateRequest	true	"Event update input"
//	@Success		200		{object}	SuccessResponse				"Event updated successfully"
//	@Failure		400		{object}	FailureResponse				"Invalid request data"
//	@Failure		500		{object}	FailureResponse				"Event update failed"
//	@Router			/events/{id} [patch]
func (rc *EventController) Update(c echo.Context) error {
	eventID := c.Param("id")
	var request model.EventUpdateRequest

	if err := c.Bind(&request); err != nil {
		err := pkg.NewError(
			err,
			"Invalid request format. Please check the input data and try again.",
			http.StatusBadRequest,
		)
		return HandleEchoError(c, err)
	}

	event, err := rc.EventDBUC.Update(c.Request().Context(), eventID, &request)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    event.Name,
		Message: "Event updated successfully.",
	})
}

// Delete handles the deletion of an existing event.
//
//	@Summary		Delete an existing event
//	@Description	This endpoint deletes an existing event by providing event name or UID.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string			true	"Event name or UID"
//	@Success		200	{object}	SuccessResponse	"Event deleted successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Event delete failed"
//	@Router			/events/{id} [delete]
func (rc *EventController) Delete(c echo.Context) error {
	eventID := c.Param("id")

	if err := rc.EventDBUC.Delete(c.Request().Context(), eventID); err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Event deleted successfully.",
	})
}

// List handles the retrieval of a list of events.
//
//	@Summary		Retrieve a list of events
//	@Description	This endpoint retrieves a list of events.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			false	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			user_id			query		string			false	"Filter events by user id, returns owners events if not provided"
//	@Param			visibility		query		string			false	"Filter events by visibility status (public:1, private:2, just me:3)"
//	@Param			limit			query		string			false	"Limit the number of events returned"
//	@Param			skip			query		string			false	"Number of events to skip for pagination"
//	@Success		200				{object}	SuccessResponse	"Events retrieved successfully"
//	@Failure		400				{object}	FailureResponse	"Invalid request data"
//	@Failure		500				{object}	FailureResponse	"Event retrieval failed"
//	@Router			/events [get]
func (rc *EventController) List(c echo.Context) error {
	opts := rc.getEventsFindOpts(c)

	list, err := rc.EventDBUC.List(c.Request().Context(), &opts)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Events retrieved successfully.",
	})
}

// GetByID godoc
//
//	@Summary		Retrieve event by ID
//	@Description	Fetches an event by its unique name or UID from the database.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string			true	"Event name or UID"
//	@Success		200	{object}	SuccessResponse	"Event retrieved successfully"
//	@Failure		400	{object}	FailureResponse	"Invalid request data"
//	@Failure		500	{object}	FailureResponse	"Event retrieval failed"
//	@Router			/events/{id} [get]
func (rc *EventController) GetByID(c echo.Context) error {
	eventID := c.Param("id")

	event, err := rc.EventDBUC.GetByID(c.Request().Context(), eventID)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    event,
		Message: "Event retrieved successfully.",
	})
}

func (rc *EventController) getEventsFindOpts(c echo.Context) model.EventFindOpts {
	return model.EventFindOpts{
		PaginationOpts: getPagination(c),
		UserID:         getFilter(c, "user_id"),
		Visibility:     getFilter(c, "visibility"),
	}
}
