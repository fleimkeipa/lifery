package repositories

import (
	"fmt"
	"strings"
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
