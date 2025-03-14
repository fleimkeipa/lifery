package repositories

import (
	"fmt"
	"strings"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg/orm"
)

func addFilterClause(filter string, key string, value string) string {
	if filter == "" {
		return fmt.Sprintf("%s = '%s'", key, value)
	}

	if strings.Contains(value, ",") {
		return fmt.Sprintf("%s AND %s IN (%s)", filter, key, value)
	}

	return fmt.Sprintf("%s AND %s = '%s'", filter, key, value)
}

func applyStandardQueries(tx *orm.Query, pagination model.PaginationOpts) *orm.Query {
	page := 1
	limit := 50
	offset := (page - 1) * limit

	if pagination.Limit > 0 {
		limit = min(pagination.Limit, 200)
	}

	if pagination.Skip > 0 {
		page = max(pagination.Skip, 1)
		offset = (page - 1) * limit
	}

	tx = tx.Limit(limit).Offset(offset)

	return tx
}

func applyFilterWithOperand(tx *orm.Query, key string, filter model.Filter) *orm.Query {
	switch filter.Operand {
	case model.OperandEqual:
		return tx.Where(fmt.Sprintf("%s=?", key), filter.Value)
	case model.OperandNot:
		return tx.Where(fmt.Sprintf("%s!=?", key), filter.Value)
	case model.OperandGreater:
		return tx.Where(fmt.Sprintf("%s>?", key), filter.Value)
	case model.OperandLess:
		return tx.Where(fmt.Sprintf("%s<?", key), filter.Value)
	case model.OperandGreaterEqual:
		return tx.Where(fmt.Sprintf("%s>=?", key), filter.Value)
	case model.OperandLessEqual:
		return tx.Where(fmt.Sprintf("%s<=?", key), filter.Value)
	case model.OperandLike:
		return tx.Where(fmt.Sprintf("%s ILIKE ?", key), "%"+filter.Value+"%")
	default:
		return tx.Where(fmt.Sprintf("%s=?", key), filter.Value)
	}
}

func GetOrderByQuery(orderBy model.OrderByOpts) string {
	if orderBy.IsSended {
		if orderBy.OrderBy == "" {
			orderBy.OrderBy = "asc"
		}
		return fmt.Sprintf("%s %s", orderBy.Column, orderBy.OrderBy)
	}

	return ""
}

func applyOrderBy(tx *orm.Query, orderBy model.OrderByOpts) *orm.Query {
	if orderBy.IsSended {
		if orderBy.OrderBy == "" {
			orderBy.OrderBy = "asc"
		}
		return tx.Order(fmt.Sprintf("%s %s", orderBy.Column, orderBy.OrderBy))
	}

	return tx
}
