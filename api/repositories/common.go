package repositories

import (
	"fmt"

	"github.com/fleimkeipa/lifery/model"

	"github.com/go-pg/pg/orm"
)

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

func applyOrderBy(tx *orm.Query, orderBy model.OrderByOpts) *orm.Query {
	if orderBy.IsSended {
		if orderBy.OrderBy == "" {
			orderBy.OrderBy = "asc"
		}
		return tx.Order(fmt.Sprintf("%s %s", orderBy.Column, orderBy.OrderBy))
	}

	return tx
}
