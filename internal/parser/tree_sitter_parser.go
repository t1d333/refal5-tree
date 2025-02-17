package parser

import (
	"context"
	"fmt"
	"math"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-tree/internal/ast"
	"github.com/t1d333/refal5-tree/internal/tree_sitter_refal5"
)

type TreeSitterRefal5Parser struct {
	parser *sitter.Parser
}

func NewTreeSitterRefal5Parser() Refal5Parser {
	parser := sitter.NewParser()
	parser.SetLanguage(tree_sitter_refal5.GetLanguage())
	return &TreeSitterRefal5Parser{
		parser: parser,
	}
}

func (p *TreeSitterRefal5Parser) GetSymbolTable() {
}

func (p *TreeSitterRefal5Parser) Parse(source []byte) (*ast.AST, error) {
	var result *ast.AST
	var cursor *sitter.QueryCursor
	tree, err := p.parser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return result, fmt.Errorf("failed to parse source code: %v", err)
	}

	root := tree.RootNode()

	// TODO: build AST and symbol table
	// TODO: walk functions

	cursor = sitter.NewQueryCursor()
	query, _ := sitter.NewQuery([]byte(`
	// (function_definition
	// 	name: (ident) @function_name
	// 	body: (body) @body
	// )`), tree_sitter_refal5.GetLanguage())
	cursor.Exec(query, root)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		funcAstNode := &ast.FunctionNode{}

		funcNameNode := match.Captures[0].Node
		funcBodyNode := match.Captures[0].Node

		funcAstNode.Name = funcNameNode.Content(source)
		sentences, err := p.walkFunctionBody(funcBodyNode)
		if err != nil {
			return nil, fmt.Errorf("failed to build ast: %v", err)
		}

		funcAstNode.Body = sentences

		fmt.Println(match.Captures)
	}

	// TODO: walk external declarations
	// cursor = sitter.NewQueryCursor()
	// query, _ = sitter.NewQuery([]byte(`
	// (external_declaration
	// 	(
	// 		function_name_list
	// 		(ident) @external_function_name
	// 	)
	// )`), tree_sitter_refal5.GetLanguage())
	// cursor.Exec(query, root)

	fmt.Println(tree.RootNode().String())

	return result, nil
}

func (p *TreeSitterRefal5Parser) walkFunctionBody(node *sitter.Node) ([]*ast.SentenceNode, error) {
	// TODO: walk sentences
	return nil, nil
}

func (p *TreeSitterRefal5Parser) walkExternalDeclarations(
	root *sitter.Node,
) (*ast.FunctionNode, error) {
	return nil, nil
}

func (p *TreeSitterRefal5Parser) walkSentence(node *sitter.Node) (*ast.SentenceNode, error) {
	return nil, nil
}

func (p *TreeSitterRefal5Parser) ParseFiles(progs [][]byte) ([]*ast.AST, error) {
	return nil, nil
}
