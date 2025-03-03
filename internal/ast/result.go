package ast

type ResultType int

const (
	CharactersResultType = iota
	WordResultType
	StringResultType
	FunctionCallResultType
	NumberResultType
	VarResultType
	GroupedResultType
)

type ResultNode interface {
	GetResultType() ResultType
}

type CharactersResultNode struct {
	Value []byte
}

func (*CharactersResultNode) GetResultType() ResultType {
	return CharactersResultType
}

type WordResultNode struct {
	Value string
}

func (*WordResultNode) GetResultType() ResultType {
	return WordResultType
}

type FunctionCallResultNode struct {
	Ident string
	Args  []ResultNode
}

func (*FunctionCallResultNode) GetResultType() ResultType {
	return FunctionCallResultType
}

type NumberResultNode struct {
	Value uint
}

func (*NumberResultNode) GetResultType() ResultType {
	return FunctionCallResultType
}

type VarResultNode struct {
	Type VaribaleType
	Name string
}

type StringResultNode struct {
	Value string
}

func (*StringResultNode) GetResultType() ResultType {
	return StringResultType
}

func (*VarResultNode) GetResultType() ResultType {
	return VarResultType
}

type GroupedResultNode struct {
	Results []ResultNode
}

func (*GroupedResultNode) GetResultType() ResultType {
	return GroupedResultType
}
