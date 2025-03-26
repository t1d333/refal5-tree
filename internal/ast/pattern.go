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

type VariableType int

const (
	SymbolVarType VariableType = iota
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
	Type VariableType
	Name string
}

func (*VarPatternNode) GetPatternType() PatternType {
	return VarPatternType
}

func (v *VarPatternNode) GetVarTypeStr() string {
	switch v.Type {
	case ExprVarType:
		return "e"
	case SymbolVarType:
		return "s"
	case TermVarType:
		return "t"
	}
	return ""
}

func (v *VarPatternNode) ToResultNode() ResultNode {
	return &VarResultNode{
		Type: v.Type,
		Name: v.Name,
	}
}

type GroupedPatternNode struct {
	Patterns []PatternNode
}

func (*GroupedPatternNode) GetPatternType() PatternType {
	return GroupedPatternType
}

func PatternToResult(node PatternNode) ResultNode {
	switch node.GetPatternType() {
	case CharactersPatternType:
		node := node.(*CharactersPatternNode)
		return &CharactersResultNode{
			Value: node.Value,
		}
	case GroupedPatternType:
		node := node.(*GroupedPatternNode)
		result := &GroupedResultNode{}
		for _, n := range node.Patterns {
			result.Results = append(result.Results, PatternToResult(n))
		}
		return result
	case NumberPatternType:
		node := node.(*NumberPatternNode)
		return &NumberResultNode{
			Value: node.Value,
		}
	case StringPatternType:
		node := node.(*NumberPatternNode)
		return &NumberResultNode{
			Value: node.Value,
		}
	case VarPatternType:
		node := node.(*VarPatternNode)
		return node.ToResultNode()
	case WordPatternType:
		node := node.(*WordPatternNode)
		return &WordResultNode{
			Value: node.Value,
		}
	}

	return nil
}
