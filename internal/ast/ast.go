package ast

type AST struct {
	Functions            []*FunctionNode
	ExternalDeclarations []string
}
