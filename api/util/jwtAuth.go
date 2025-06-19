package util

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type contextKey string

const userContextKey contextKey = "user"

// check for valid admin token
func JWTAuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := ValidateJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Authentication required",
				"error":   err.Error(),
			})
		}

		if err := ValidateAdminRoleJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Only Administrators are allowed to perform this action",
				"error":   err.Error(),
			})
		}

		setOwnerOnCtx(c)

		return next(c)
	}
}

// check for valid user token
func JWTAuthEditor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := ValidateJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Authentication required",
				"error":   err.Error(),
			})
		}

		if err := ValidateEditorRoleJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Only registered Editors are allowed to perform this action",
				"error":   err.Error(),
			})
		}

		setOwnerOnCtx(c)

		return next(c)
	}
}

// check for valid viewer token
func JWTAuthViewer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if IsUserPublic(c) {
			return next(c)
		}

		if err := ValidateJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Authentication required",
				"error":   err.Error(),
			})
		}

		if err := ValidateViewerRoleJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Only registered Viewers are allowed to perform this action",
				"error":   err.Error(),
			})
		}

		setOwnerOnCtx(c)

		return next(c)
	}
}

func setOwnerOnCtx(c echo.Context) {
	user, _ := GetOwnerFromToken(c)

	ctx := context.WithValue(c.Request().Context(), userContextKey, user)

	c.SetRequest(c.Request().WithContext(ctx))
}
