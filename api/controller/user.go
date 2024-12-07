package controller

import (
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/pkg"
	"github.com/fleimkeipa/lifery/uc"

	"github.com/labstack/echo/v4"
)

type UserHandlers struct {
	userUC *uc.UserUC
}

func NewUserHandlers(uc *uc.UserUC) *UserHandlers {
	return &UserHandlers{
		userUC: uc,
	}
}

// Create godoc
//
//	@Summary		Create creates a new user
//	@Description	This endpoint creates a new user by providing username, email, password, and role ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		model.UserCreateRequest	true	"User creation input"
//	@Success		201				{object}	SuccessResponse			"user username"
//	@Failure		400				{object}	FailureResponse			"Error message including details on failure"
//	@Failure		500				{object}	FailureResponse			"Interval error"
//	@Router			/users [post]
func (rc *UserHandlers) Create(c echo.Context) error {
	var input model.UserCreateRequest

	if err := c.Bind(&input); err != nil {
		err := pkg.NewError(
			err,
			"Invalid request format. Please check the input data and try again.",
			http.StatusBadRequest,
		)
		return HandleEchoError(c, err)
	}

	user, err := rc.userUC.Create(c.Request().Context(), input)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    user.Username,
		Message: "User created successfully.",
	})
}

// UpdateUser godoc
//
//	@Summary		UpdateUser updates an existing user
//	@Description	This endpoint updates a user by providing username, email, password, and role ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		model.UserCreateRequest	true	"User update input"
//	@Success		200				{object}	SuccessResponse			"user username"
//	@Failure		400				{object}	FailureResponse			"Error message including details on failure"
//	@Failure		500				{object}	FailureResponse			"Interval error"
//	@Router			/users/{id} [patch]
func (rc *UserHandlers) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var input model.UserCreateRequest

	if err := c.Bind(&input); err != nil {
		err := pkg.NewError(
			err,
			"Invalid request format. Please check the input data and try again.",
			http.StatusBadRequest,
		)
		return HandleEchoError(c, err)
	}

	user, err := rc.userUC.Update(c.Request().Context(), id, input)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    user.Username,
		Message: "User updated successfully.",
	})
}

// List godoc
//
//	@Summary		List all users
//	@Description	Retrieves a filtered and paginated list of users from the database based on query parameters.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			username		query		string			false	"Filter users by username"
//	@Param			email			query		string			false	"Filter users by email"
//	@Param			role_id			query		string			false	"Filter users by role ID"
//	@Param			limit			query		string			false	"Limit the number of users returned"
//	@Param			skip			query		string			false	"Number of users to skip for pagination"
//	@Success		200				{object}	SuccessResponse	"Successful response containing the list of users"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/users [get]
func (rc *UserHandlers) List(c echo.Context) error {
	opts := rc.getUsersFindOpts(c, model.ZeroCreds)

	list, err := rc.userUC.List(c.Request().Context(), &opts)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Users retrieved successfully.",
	})
}

// GetByID godoc
//
//	@Summary		Retrieve user by ID
//	@Description	Fetches a user by their unique ID from the database.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			id				path		string			true	"User ID"
//	@Success		200				{object}	SuccessResponse	"Successful response containing the user information"
//	@Failure		500				{object}	FailureResponse	"Internal server error"
//	@Router			/users/{id} [get]
func (rc *UserHandlers) GetByID(c echo.Context) error {
	id := c.Param("id") // Extract the user ID from the path parameters

	user, err := rc.userUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return HandleEchoError(c, err)
	}

	// Remove the password from the user object before returning it
	user.Password = ""

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    user,
		Message: "User retrieved successfully.",
	})
}

// DeleteUser godoc
//
//	@Summary		DeleteUser deletes an existing user
//	@Description	This endpoint deletes a user by providing user id.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	SuccessResponse	"user username"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/users/{id} [delete]
func (rc *UserHandlers) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if err := rc.userUC.Delete(c.Request().Context(), id); err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "User deleted successfully.",
	})
}

func (rc *UserHandlers) getUsersFindOpts(c echo.Context, fields ...string) model.UserFindOpts {
	return model.UserFindOpts{
		PaginationOpts: getPagination(c),
		FieldsOpts: model.FieldsOpts{
			Fields: fields,
		},
		Username: getFilter(c, "username"),
		Email:    getFilter(c, "email"),
		RoleID:   getFilter(c, "role_id"),
	}
}
