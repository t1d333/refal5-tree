package parser

import (
	"context"
	"fmt"
	"strconv"

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
	(function_definition
		name: (ident) @function_name
		body: (body) @body
	)`), tree_sitter_refal5.GetLanguage())
	cursor.Exec(query, root)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		funcAstNode := &ast.FunctionNode{}

		funcNameNode := match.Captures[0].Node
		funcBodyNode := match.Captures[1].Node

		funcAstNode.Name = funcNameNode.Content(source)
		fmt.Println("Found function", funcAstNode.Name)
		sentences, err := p.walkFunctionBody(funcBodyNode, source)
		if err != nil {
			return nil, fmt.Errorf("failed to build ast: %v", err)
		}

		funcAstNode.Body = sentences

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

	return result, nil
}

func (p *TreeSitterRefal5Parser) walkFunctionBody(
	node *sitter.Node,
	source []byte,
) ([]*ast.SentenceNode, error) {
	if node == nil {
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		inner := node.Child(i)
		if inner.ChildByFieldName("sentence_eq") == nil &&
			inner.ChildByFieldName("sentence_block") == nil {
			continue
		}

		sentenceNode := inner.ChildByFieldName("sentence_eq")
		if sentenceNode == nil {
			sentenceNode = inner.ChildByFieldName("sentence_block")
		}

		sentenceLhsNode := sentenceNode.ChildByFieldName("lhs")
		lhs := []ast.PatternNode{}

		// walk lhs
		for i := 0; i < int(sentenceLhsNode.ChildCount()); i++ {
			lhsPartNode := sentenceLhsNode.Child(i)

			// TODO: check error
			lhsPart, _ := p.walkPattern(lhsPartNode, source)
			lhs = append(lhs, lhsPart)
		}

		for j := 0; j < int(sentenceNode.ChildCount()); j++ {
			child := sentenceNode.Child(j)
			if child == nil || child.Type() != "condition" {
				continue
			}

			// isPatternStart := "False"
			for k := 0; k < int(child.ChildCount()); k++ {
				child.NamedChildCount
			}
			
		}

		// walk conditions?

		// walk rhs

	}

	return nil, nil
}

func (p *TreeSitterRefal5Parser) walkPattern(
	node *sitter.Node,
	source []byte,
) (ast.PatternNode, error) {
	// TODO: check node == nil
	// var result ast.PatternNode

	switch node.Type() {
	case "ident":
		return &ast.WordPattern{Value: node.Content(source)}, nil
	case "string":
		return &ast.StringPattern{Value: node.Content(source)}, nil
	case "number":
		val, _ := strconv.Atoi(node.Content(source))
		return &ast.NumberPattern{Value: uint(val)}, nil
	case "variable":
		varStrType := node.ChildByFieldName("type").Content(source)
		varType := ast.SymbolVarType
		switch varStrType {
		case "t":
			varType = ast.TermVarType
		case "e":
			varType = ast.ExprVarType
		}
		return &ast.VarPattern{
			Name: node.ChildByFieldName("name").Content(source),
			Type: varType,
		}, nil
	case "grouped_pattern":
		pattern := &ast.GroupedPattern{Patterns: []ast.PatternNode{}}
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if !child.IsNamed() {
				continue
			}

			// TODO: check error
			nestedPattern, _ := p.walkPattern(child, source)
			pattern.Patterns = append(pattern.Patterns, nestedPattern)

		}
		return pattern, nil
	case "symbols":
		return &ast.CharactersPattern{Value: []byte(node.Content(source))}, nil
	}

	return nil, fmt.Errorf("undefined pattern")
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
