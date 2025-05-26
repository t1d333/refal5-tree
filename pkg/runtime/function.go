package runtime

type R5FunctionPtr func(l, r int, arg Rope, rhsStack *[]ViewFieldNode)

type R5Function struct {
	Name  string
	Entry bool
	Ptr   R5FunctionPtr
}
