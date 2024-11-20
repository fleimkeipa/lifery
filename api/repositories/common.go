package repositories

import "fmt"

func addFilterClause(filter string, key string, value string) string {
	if filter == "" {
		return fmt.Sprintf("%s = '%s'", key, value)
	}

	return fmt.Sprintf("%s AND %s = '%s'", filter, key, value)
}
