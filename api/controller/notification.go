package controller

import (
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/uc"
	"github.com/fleimkeipa/lifery/util"

	"github.com/labstack/echo/v4"
)

type NotificationHandlers struct {
	notificationUC *uc.NotificationUC
}

func NewNotificationHandlers(notificationUC *uc.NotificationUC) *NotificationHandlers {
	return &NotificationHandlers{
		notificationUC: notificationUC,
	}
}

// Update godoc
//
//	@Summary		Update updates an existing notification
//	@Description	This endpoint updates a notification by binding the incoming JSON request to the NotificationUpdateInput model.
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string							true	"Notification ID to update"
//	@Param			Body	body		model.NotificationUpdateInput	true	"Notification update input"
//	@Success		200		{object}	SuccessResponse					"Notification updated successfully"
//	@Failure		400		{object}	FailureResponse					"Invalid request data"
//	@Failure		500		{object}	FailureResponse					"Notification update failed"
//	@Router			/notifications/{id} [patch]
func (rc *NotificationHandlers) Update(c echo.Context) error {
	id := c.Param("id")
	var input model.NotificationUpdateInput

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	err := rc.notificationUC.Update(c.Request().Context(), id, input)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notification updated successfully",
	})
}

// List godoc
//
//	@Summary		List lists all notifications
//	@Description	Retrieves a filtered and paginated list of notifications based on query parameters.
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			user_id	query		string			false	"Filter notifications by user id if you are admin"
//	@Param			read	query		string			false	"Filter notifications by read status"
//	@Param			limit	query		string			false	"Limit the number of notifications returned"
//	@Param			skip	query		string			false	"Number of notifications to skip for pagination"
//	@Success		200		{object}	SuccessResponse	"Successful response containing the list of notifications"
//	@Failure		500		{object}	FailureResponse	"Internal error"
//	@Router			/notifications [get]
func (rc *NotificationHandlers) List(c echo.Context) error {
	opts := rc.getNotificationFindOpts(c, model.ZeroCreds)

	list, err := rc.notificationUC.List(c.Request().Context(), &opts)
	if err != nil {
		return handleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessListResponse{
		Data:  list.Notifications,
		Total: list.Total,
		Limit: list.Limit,
		Skip:  list.Skip,
	})
}

func (rc *NotificationHandlers) getNotificationFindOpts(c echo.Context, fields ...string) model.NotificationFindOpts {
	defaultFilter := model.NotificationFindOpts{
		OrderByOpts:    getOrder(c),
		PaginationOpts: getPagination(c),
		FieldsOpts: model.FieldsOpts{
			Fields: fields,
		},
		UserID: getFilter(c, "user_id"),
		Read:   getFilter(c, "read"),
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
