package controller

import (
	"strconv"

	"github.com/fleimkeipa/lifery/model"

	"github.com/labstack/echo/v4"
)

type FailureResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type AuthResponse struct {
	Type     string `json:"type" example:"basic,oauth2"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func getPagination(c echo.Context) model.PaginationOpts {
	limitQ := c.QueryParam("limit")
	skipQ := c.QueryParam("skip")

	limit, _ := strconv.Atoi(limitQ)
	if limit <= 0 {
		limit = 30
	}

	skip, _ := strconv.Atoi(skipQ)
	if skip < 0 {
		skip = 0
	}

	return model.PaginationOpts{
		Limit: limit,
		Skip:  skip,
	}
}

func getFilter(c echo.Context, query string) model.Filter {
	param := c.QueryParam(query)
	if param == "" {
		return model.Filter{}
	}

	return model.Filter{
		IsSended: true,
		Value:    param,
	}
}
