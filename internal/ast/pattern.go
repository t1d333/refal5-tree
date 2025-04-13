package ast

import "fmt"

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
		node := node.(*StringPatternNode)
		return &StringResultNode{
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

func PatternsToResults(patterns []PatternNode) []ResultNode {
	result := []ResultNode{}

	for _, n := range patterns {
		result = append(result, PatternToResult(n))
	}
	return result
}

func ExtractVars(patterns []PatternNode) []PatternNode {
	result := []PatternNode{}

	for _, n := range patterns {
		if n.GetPatternType() == VarPatternType {
			result = append(result, n)
		}

		if n.GetPatternType() == GroupedPatternType {
			grouped := n.(*GroupedPatternNode)
			result = append(result, ExtractVars(grouped.Patterns)...)
		}
	}

	return result
}

func ReplacePatternVariable(
	patterns []PatternNode,
	target *VarPatternNode,
	replacement []PatternNode,
) []PatternNode {
	result := []PatternNode{}

	for _, curr := range patterns {
		if curr.GetPatternType() == VarPatternType {
			varNode := curr.(*VarPatternNode)
			if varNode.Name != target.Name || varNode.Type != target.Type {
				result = append(result, curr)
			} else {
				result = append(result, replacement...)
			}

		} else if curr.GetPatternType() == GroupedPatternType {
			grouped := curr.(*GroupedPatternNode)
			result = append(result, &GroupedPatternNode{Patterns: ReplacePatternVariable(grouped.Patterns, target, replacement)})
		} else {
			result = append(result, curr)
		}
	}

	return result
}

func PrintPattern(pattern PatternNode) {
	switch pattern.GetPatternType() {
	case NumberPatternType:
		n := pattern.(*NumberPatternNode)
		fmt.Printf("(Number: %d)\n", n.Value)
	case StringPatternType:
		n := pattern.(*StringPatternNode)
		fmt.Printf("(Number: %d)\n", n.Value)
	case CharactersPatternType:
		n := pattern.(*CharactersPatternNode)
		fmt.Printf("(Number: %s)\n", string(n.Value))
	case VarPatternType:
		n := pattern.(*VarPatternNode)
		fmt.Printf("(Var %s %s)\n", n.GetVarTypeStr(), n.Name)
	case GroupedPatternType:
		n := pattern.(*GroupedPatternNode)
		fmt.Printf("(\n")
		for _, g := range n.Patterns {
			PrintPattern(g)
		}
		fmt.Printf(")\n")
	}
}
