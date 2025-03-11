package runtime

type R5FunctionPtr func(*Rope)

type R5Function struct {
	Name  string
	Entry bool
	Ptr   R5FunctionPtr
}
