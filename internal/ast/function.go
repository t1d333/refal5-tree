package ast

type FunctionNode struct {
	Name  string
	Entry bool
	Body  []*SentenceNode
}
