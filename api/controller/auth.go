package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/uc"
	"github.com/fleimkeipa/lifery/util"

	"github.com/labstack/echo/v4"
)

type AuthHandlers struct {
	userUC *uc.UserUC
}

func NewAuthHandlers(uc *uc.UserUC) *AuthHandlers {
	return &AuthHandlers{
		userUC: uc,
	}
}

// Register godoc
//
//	@Summary		User register
//	@Description	This endpoint allows a user to log in by providing a valid username and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Register	true	"User register input"
//	@Success		200		{object}	AuthResponse	"Successfully registered in with JWT token"
//	@Failure		400		{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse	"Interval error"
//	@Router			/auth/register [post]
func (rc *AuthHandlers) Register(c echo.Context) error {
	var input model.Register

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	newUser := model.UserCreateRequest{
		Username:        input.Username,
		Email:           input.Email,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	}

	user, err := rc.userUC.Create(c.Request().Context(), newUser)
	if err != nil {
		return handleEchoError(c, err)
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "Register failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token:    jwt,
		Type:     "basic",
		Username: input.Username,
		Message:  "Successfully registered in",
	})
}

// Login godoc
//
//	@Summary		User login
//	@Description	This endpoint allows a user to log in by providing a valid username and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Login		true	"User login input"
//	@Success		200		{object}	AuthResponse	"Successfully logged in with JWT token"
//	@Failure		400		{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse	"Interval error"
//	@Router			/auth/login [post]
func (rc *AuthHandlers) Login(c echo.Context) error {
	var input model.Login

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	user, err := rc.userUC.GetByUsernameOrEmail(c.Request().Context(), input.Username)
	if err != nil {
		return handleEchoError(c, err)
	}

	if err := model.ValidateUserPassword(user.Password, input.Password); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Invalid password: %v", err),
			Message: "Invalid password. Please check the password and try again.",
		})
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "Login failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:    jwt,
		Type:     "basic",
		Username: input.Username,
		Message:  "Successfully logged in",
	})
}
