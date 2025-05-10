package controller

import (
	"strconv"
	"strings"

	"github.com/fleimkeipa/lifery/model"

	"github.com/labstack/echo/v4"
)

type FailureResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessListResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
	Limit int         `json:"limit"`
	Skip  int         `json:"skip"`
}

type SuccessResponse struct {
	Message string `json:"message"`
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
	paramQ := c.QueryParam(query)
	if paramQ == "" {
		return model.Filter{}
	}

	basicQuery := model.Filter{
		IsSended: true,
		Operand:  model.OperandEqual,
		Value:    paramQ,
	}

	if !strings.Contains(paramQ, ":") {
		return basicQuery
	}

	// example: is_ocr_checked=eq:true
	splitted := strings.Split(paramQ, ":")
	if len(splitted) == 1 {
		return basicQuery
	}

	operand := model.OperandEqual
	switch splitted[0] {
	case model.OperandNot.String():
		operand = model.OperandNot
	case model.OperandGreater.String():
		operand = model.OperandGreater
	case model.OperandGreaterEqual.String():
		operand = model.OperandGreaterEqual
	case model.OperandLess.String():
		operand = model.OperandLess
	case model.OperandLessEqual.String():
		operand = model.OperandLessEqual
	case model.OperandLike.String():
		operand = model.OperandLike
	}

	return model.Filter{
		IsSended: true,
		Operand:  operand,
		Value:    splitted[1],
	}
}

func getOrder(c echo.Context) model.OrderByOpts {
	orderQ := c.QueryParam("order")
	if orderQ == "" {
		return model.OrderByOpts{}
	}

	// example: created_at:desc
	splitted := strings.Split(orderQ, ":")
	if len(splitted) != 2 {
		return model.OrderByOpts{
			IsSended: true,
			Column:   orderQ,
			OrderBy:  "asc",
		}
	}

	return model.OrderByOpts{
		IsSended: true,
		Column:   splitted[1],
		OrderBy:  splitted[0],
	}
}
