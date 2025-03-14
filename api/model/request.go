package model

type InternalRequest struct {
	Body       interface{}
	Headers    map[string]string
	Method     string
	Paths      []string
	Pagination PaginationOpts
}

type PaginationOpts struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type OrderByOpts struct {
	IsSended bool
	Column   string
	OrderBy  string
}

type Operand string

func (o Operand) String() string {
	return string(o)
}

const (
	OperandEqual        Operand = "eq"
	OperandNot          Operand = "ne"
	OperandGreater      Operand = "gt"
	OperandGreaterEqual Operand = "gte"
	OperandLess         Operand = "lt"
	OperandLessEqual    Operand = "lte"
	OperandLike         Operand = "like"
)

type Filter struct {
	Value    string
	Operand  Operand
	IsSended bool
}
type FieldsOpts struct {
	Fields []string
}
