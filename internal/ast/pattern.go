package ast

type PatternType int

const (
	SymbolPatternType = iota
	VarPatternType
	WordPatternType
	GroupedPatternType
	NumberPatternType
)

type PatternNode interface {
	GetPatternType() PatternType
}

type SymbolPattern struct{}

func (*SymbolPattern) GetPatternType() PatternType {
	return SymbolPatternType
}

type WordPattern struct{}

func (*WordPattern) GetPatternType() PatternType {
	return SymbolPatternType
}

type NumberPattern struct{}

func (*NumberPattern) GetPatternType() PatternType {
	return SymbolPatternType
}

type VarPattern struct{}

func (*VarPattern) GetPatternType() PatternType {
	return SymbolPatternType
}

type GroupedPattern struct{}

func (*GroupedPattern) GetPatternType() PatternType {
	return SymbolPatternType
}
