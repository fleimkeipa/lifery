package model

type InternalRequest struct {
	Pagination PaginationOpts
	Method     string
	Paths      []string
	Headers    map[string]string
	Body       interface{}
}

type PaginationOpts struct {
	Limit int
	Skip  int
}

type Filter struct {
	Value    string
	IsSended bool
}

type FieldsOpts struct {
	Fields []string
}
