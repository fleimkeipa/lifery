package model

type PaginationOpts struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type OrderByOpts struct {
	Column   string
	OrderBy  string
	IsSended bool
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
	OperandIn           Operand = "in"
	OperandNotIn        Operand = "nin"
)

type Filter struct {
	Value    string
	Operand  Operand
	IsSended bool
}
type FieldsOpts struct {
	Fields []string
}
