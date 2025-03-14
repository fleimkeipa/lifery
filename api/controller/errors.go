package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fleimkeipa/lifery/pkg"

	"github.com/labstack/echo/v4"
)

// HandleEchoError handles errors that occur within the Echo framework.
func HandleEchoError(c echo.Context, err error) error {
	var pe *pkg.Error

	if errors.As(err, &pe) {
		message := pe.Message()
		errorMessage := func() string {
			errMessage := pe.Error()
			if errMessage == "" {
				return ""
			}
			if strings.HasPrefix(errMessage, "ERROR:") {
				return errMessage
			}
			return fmt.Sprintf("error: %s", errMessage)
		}()

		if message != "" && errorMessage != "" {
			message = fmt.Sprintf("%s | %s", message, errorMessage)
		}
		if message == "" {
			message = errorMessage
		}

		return c.JSON(pe.StatusCode(), FailureResponse{
			Error:   err.Error(),
			Message: message,
		})
	}

	return c.JSON(http.StatusInternalServerError, FailureResponse{
		Error:   err.Error(),
		Message: "Internal Server Error",
	})
}
