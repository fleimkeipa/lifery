package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/lifery/model"
	"github.com/fleimkeipa/lifery/uc"
	"github.com/fleimkeipa/lifery/util"

	"github.com/labstack/echo/v4"
)

type OAuthHandlers struct {
	oauthUC *uc.OAuthUC
}

func NewOAuthHandlers(oauthUC *uc.OAuthUC) *OAuthHandlers {
	return &OAuthHandlers{
		oauthUC: oauthUC,
	}
}

// GoogleAuthURL godoc
//
//	@Summary		Get Google OAuth URL
//	@Description	This endpoint returns the Google OAuth authorization URL.
//	@Tags			oauth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Google OAuth URL"
//	@Failure		500	{object}	FailureResponse		"Internal error"
//	@Router			/oauth/google/url [get]
func (rc *OAuthHandlers) GoogleAuthURL(c echo.Context) error {
	authURL := rc.oauthUC.GetAuthURL(model.GoogleProvider)

	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": authURL,
		"message":  "Google OAuth URL generated successfully",
	})
}

// GoogleCallback godoc
//
//	@Summary		Google OAuth callback
//	@Description	This endpoint handles the Google OAuth callback and creates or logs in the user.
//	@Tags			oauth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.GoogleAuthRequest	true	"Google OAuth code"
//	@Success		200		{object}	AuthResponse			"Successfully authenticated with JWT token"
//	@Failure		400		{object}	FailureResponse			"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse			"Internal error"
//	@Router			/oauth/google/callback [post]
func (rc *OAuthHandlers) GoogleCallback(c echo.Context) error {
	var input model.GoogleAuthRequest

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	user, err := rc.oauthUC.HandleCallback(c.Request().Context(), model.GoogleProvider, input.Code)
	if err != nil {
		return handleEchoError(c, err)
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "Google authentication failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:    jwt,
		Type:     "google",
		Username: user.Username,
		Message:  "Successfully authenticated with Google",
	})
}

// LinkedInAuthURL godoc
//
//	@Summary		Get LinkedIn OAuth URL
//	@Description	This endpoint returns the LinkedIn OAuth authorization URL.
//	@Tags			oauth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"LinkedIn OAuth URL"
//	@Failure		500	{object}	FailureResponse		"Internal error"
//	@Router			/oauth/linkedin/url [get]
func (rc *OAuthHandlers) LinkedInAuthURL(c echo.Context) error {
	authURL := rc.oauthUC.GetAuthURL(model.LinkedInProvider)
	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": authURL,
		"message":  "LinkedIn OAuth URL generated successfully",
	})
}

// LinkedInCallback godoc
//
//	@Summary		LinkedIn OAuth callback
//	@Description	This endpoint handles the LinkedIn OAuth callback and creates or logs in the user.
//	@Tags			oauth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.LinkedInAuthRequest	true	"LinkedIn OAuth code"
//	@Success		200		{object}	AuthResponse				"Successfully authenticated with JWT token"
//	@Failure		400		{object}	FailureResponse				"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse				"Internal error"
//	@Router			/oauth/linkedin/callback [post]
func (rc *OAuthHandlers) LinkedInCallback(c echo.Context) error {
	var input model.LinkedInAuthRequest

	if err := c.Bind(&input); err != nil {
		return handleBindingErrors(c, err)
	}

	if err := c.Validate(&input); err != nil {
		return handleValidatingErrors(c, err)
	}

	user, err := rc.oauthUC.HandleCallback(c.Request().Context(), model.LinkedInProvider, input.Code)
	if err != nil {
		return handleEchoError(c, err)
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "LinkedIn authentication failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:    jwt,
		Type:     "linkedin",
		Username: user.Username,
		Message:  "Successfully authenticated with LinkedIn",
	})
}
