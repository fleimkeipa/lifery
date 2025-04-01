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
	return c.JSON(http.StatusBadRequest, FailureResponse{
		Error:   fmt.Sprintf("Failed to bind request: %v", err),
		Message: "Invalid request data. Please check your input and try again.",
	})
}

func handleValidatingErrors(c echo.Context, err error) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err.Error()),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	var allErrMessages string
	for _, vErr := range validationErrors {
		var errorMessage string
		switch vErr.Tag() {
		case "required":
			errorMessage = fmt.Sprintf("%s field not provided", vErr.Field())
		case "iscolor":
			errorMessage = fmt.Sprintf("%s field is not color(hexcolor|rgb|rgba|hsl|hsla)", vErr.Field())
		default:
			errorMessage = vErr.Error()
		}

		if allErrMessages == "" {
			allErrMessages = errorMessage
		} else {
			allErrMessages = fmt.Sprintf("%s | %s", allErrMessages, errorMessage)
		}
	}

	if allErrMessages == "" {
		allErrMessages = err.Error()
	}

	return c.JSON(http.StatusBadRequest, FailureResponse{
		Error:   fmt.Sprintf("Failed to bind request: %v", allErrMessages),
		Message: "Invalid request data. Please check your input and try again.",
	})
}
