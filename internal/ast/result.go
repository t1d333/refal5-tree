package ast

import "fmt"

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
	String() string
}

func GetResultLengthInRuntimeNodes(r ResultNode) int {
	stack := []ResultNode{r}
	length := 0

	for len(stack) > 0 {
		curr := stack[0]
		stack = stack[1:]

		switch curr.GetResultType() {
		case CharactersResultType:
			node := curr.(*CharactersResultNode)
			length += len(node.Value)
		case WordResultType:
			length += 1
		case StringResultType:
			length += 1
		case FunctionCallResultType:
			// + open call + func + close call
			node := curr.(*FunctionCallResultNode)
			length += 3
			for _, arg := range node.Args {
				length += GetResultLengthInRuntimeNodes(arg)
			}
		case NumberResultType:
			length += 1
		case VarResultType:
			length += 1
		case GroupedResultType:
			node := curr.(*GroupedResultNode)
			length += 2
			for _, r := range node.Results {
				length += GetResultLengthInRuntimeNodes(r)
			}

		}
	}

	return length
}

type CharactersResultNode struct {
	Value []byte
}

func (n *CharactersResultNode) String() string {
	return "'" + string(n.Value) + "'"
}

func (*CharactersResultNode) GetResultType() ResultType {
	return CharactersResultType
}

type WordResultNode struct {
	Value string
}

func (w *WordResultNode) String() string {
	return w.Value
}

func (*WordResultNode) GetResultType() ResultType {
	return WordResultType
}

type FunctionCallResultNode struct {
	Ident string
	Args  []ResultNode
}

func (f *FunctionCallResultNode) String() string {
	res := fmt.Sprintf("<%s", f.Ident)

	for _, arg := range f.Args {
		res += " " + arg.String()
	}

	res += ">"

	return res
}

func (*FunctionCallResultNode) GetResultType() ResultType {
	return FunctionCallResultType
}

type NumberResultNode struct {
	Value uint
}

func (n *NumberResultNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (*NumberResultNode) GetResultType() ResultType {
	return NumberResultType
}

type VarResultNode struct {
	Type VariableType
	Name string
}

func (v *VarResultNode) String() string {
	return fmt.Sprintf("%s.%s", v.GetVarTypeStr(), v.Name)
}

func (v *VarResultNode) GetVarTypeStr() string {
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

func (v *VarResultNode) ToPatternNode() PatternNode {
	return &VarPatternNode{
		Type: v.Type,
		Name: v.Name,
	}
}

type StringResultNode struct {
	Value string
}

func (s *StringResultNode) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
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

func (g *GroupedResultNode) String() string {
	res := "("

	for _, n := range g.Results {
		res += " " + n.String()
	}

	res += " )"
	return res
}

func (*GroupedResultNode) GetResultType() ResultType {
	return GroupedResultType
}

func ReplaceResultVariable(
	results []ResultNode,
	target *VarResultNode,
	replacement []ResultNode,
) []ResultNode {
	result := []ResultNode{}

	for _, curr := range results {
		if curr.GetResultType() == VarResultType {
			varNode := curr.(*VarResultNode)
			if varNode.Name != target.Name || varNode.Type != target.Type {
				result = append(result, curr)
			} else {
				result = append(result, replacement...)
			}

		} else if curr.GetResultType() == GroupedResultType {
			grouped := curr.(*GroupedResultNode)
			newGrouped := &GroupedResultNode{ReplaceResultVariable(grouped.Results, target, replacement)}
			result = append(result, newGrouped)
		} else if curr.GetResultType() == FunctionCallResultType {
			callNode := curr.(*FunctionCallResultNode)
			newCall := &FunctionCallResultNode{Ident: callNode.Ident, Args: ReplaceResultVariable(callNode.Args, target, replacement)}
			result = append(result, newCall)
		} else {
			result = append(result, curr)
		}
	}

	return result
}

func PrintResult(result ResultNode) {
	switch result.GetResultType() {
	case NumberResultType:
		n := result.(*NumberResultNode)
		fmt.Printf("(Number %d)\n", n.Value)
	case StringResultType:
		n := result.(*StringResultNode)
		fmt.Printf("(String %d)\n", n.Value)
	case CharactersResultType:
		n := result.(*CharactersResultNode)
		fmt.Printf("(Char %s)\n", string(n.Value))
	case VarResultType:
		n := result.(*VarResultNode)
		fmt.Printf("(Var %s %s)\n", n.GetVarTypeStr(), n.Name)
	case FunctionCallResultType:
		n := result.(*FunctionCallResultNode)
		fmt.Printf("<%s\n", n.Ident)
		for _, arg := range n.Args {
			PrintResult(arg)
		}
		fmt.Printf(">", n.Ident)
	case GroupedResultType:
		n := result.(*GroupedResultNode)
		fmt.Printf("(\n")
		for _, g := range n.Results {
			PrintResult(g)
		}
		fmt.Printf(")\n")
	}
}
