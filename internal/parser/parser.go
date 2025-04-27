package parser

import "github.com/t1d333/refal5-tree/internal/ast"

type Refal5Parser interface {
	Parse(prog []byte) (*ast.AST, error)
	ParseFiles(progs [][]byte) ([]*ast.AST, *ast.FunctionNode, error)
}
