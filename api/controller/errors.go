package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fleimkeipa/lifery/pkg"
	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

// handleEchoError handles errors that occur within the Echo framework.
func handleEchoError(c echo.Context, err error) error {
	var pe *pkg.Error

	if errors.As(err, &pe) {
		message := pe.Message()
		errorMessage := func() string {
			errMessage := pe.Error()
			if errMessage == "" {
				return ""
			}
			if strings.HasPrefix(errMessage, "pg:") {
				return errMessage
			}
			return fmt.Sprintf("error: %s", errMessage)
		}()

		return c.JSON(pe.StatusCode(), FailureResponse{
			Error:   errorMessage,
			Message: message,
		})
	}

	return c.JSON(http.StatusInternalServerError, FailureResponse{
		Error:   err.Error(),
		Message: "Internal Server Error",
	})
}

func handleBindingErrors(c echo.Context, err error) error {
	var errorMessage string
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		validationError := validationErrors[0]
		if validationError.Tag() == "required" {
			errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
		}
	}

	return c.JSON(http.StatusBadRequest, FailureResponse{
		Error:   fmt.Sprintf("Failed to bind request: %v", errorMessage),
		Message: "Invalid request data. Please check your input and try again.",
	})
}
