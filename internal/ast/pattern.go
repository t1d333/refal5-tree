package ast

type PatternType int

const (
	CharactersPatternType PatternType = iota
	VarPatternType
	WordPatternType
	GroupedPatternType
	NumberPatternType
	StringPatternType
)

type VaribaleType int

const (
	SymbolVarType VaribaleType = iota
	TermVarType
	ExprVarType
)

type PatternNode interface {
	GetPatternType() PatternType
}

type CharactersPatternNode struct {
	Value []byte
}

func (*CharactersPatternNode) GetPatternType() PatternType {
	return CharactersPatternType
}

type WordPatternNode struct {
	Value string
}

func (*WordPatternNode) GetPatternType() PatternType {
	return WordPatternType
}

type NumberPatternNode struct {
	Value uint
}

func (*NumberPatternNode) GetPatternType() PatternType {
	return NumberPatternType
}

type StringPatternNode struct {
	Value string
}

func (*StringPatternNode) GetPatternType() PatternType {
	return StringPatternType
}

type VarPatternNode struct {
	Type VaribaleType
	Name string
}

func (*VarPatternNode) GetPatternType() PatternType {
	return VarPatternType
}

type GroupedPatternNode struct {
	Patterns []PatternNode
}

func (*GroupedPatternNode) GetPatternType() PatternType {
	return GroupedPatternType
}
