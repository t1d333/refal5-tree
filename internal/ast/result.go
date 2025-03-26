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
	return NumberResultType
}

type VarResultNode struct {
	Type VariableType
	Name string
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
