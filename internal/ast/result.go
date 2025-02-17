package ast

type ResultType int

const (
	SymbolResultType = iota
	WordResultType
	FunctionCallResultType
	NumberResultType
	VarResultType
	GroupedResultType
)

type ResultNode interface {
	GetResultType() ResultType
}

type SymbolResultNode struct{}

func (*SymbolResultNode) GetResultType() ResultType {
	return SymbolResultType
}

type WordResultNode struct{}

func (*WordResultNode) GetResultType() ResultType {
	return WordResultType
}

type FunctionCallResultNode struct{}

func (*FunctionCallResultNode) GetResultType() ResultType {
	return FunctionCallResultType
}

type NumberResultNode struct{}

func (*NumberResultNode) GetResultType() ResultType {
	return FunctionCallResultType
}

type VarResultNode struct{}

func (*VarResultNode) GetResultType() ResultType {
	return VarResultType
}

type GroupedResultNode struct{}

func (*GroupedResultNode) GetResultType() ResultType {
	return GroupedResultType
}
