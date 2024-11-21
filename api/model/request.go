package model

type InternalRequest struct {
	Body       interface{}
	Headers    map[string]string
	Method     string
	Paths      []string
	Pagination PaginationOpts
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
