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

type CharactersPattern struct {
	Value []byte
}

func (*CharactersPattern) GetPatternType() PatternType {
	return CharactersPatternType
}

type WordPattern struct {
	Value string
}

func (*WordPattern) GetPatternType() PatternType {
	return WordPatternType
}

type NumberPattern struct {
	Value uint
}

func (*NumberPattern) GetPatternType() PatternType {
	return NumberPatternType
}

type StringPattern struct {
	Value string
}

func (*StringPattern) GetPatternType() PatternType {
	return StringPatternType
}

type VarPattern struct {
	Type VaribaleType
	Name string
}

func (*VarPattern) GetPatternType() PatternType {
	return VarPatternType
}

type GroupedPattern struct {
	Patterns []PatternNode
}

func (*GroupedPattern) GetPatternType() PatternType {
	return GroupedPatternType
}
