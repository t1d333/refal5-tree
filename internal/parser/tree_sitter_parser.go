package parser

import (
	"context"
	"fmt"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-tree/internal/ast"
	"github.com/t1d333/refal5-tree/internal/library"
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

func (p *TreeSitterRefal5Parser) CheckErrors(root *sitter.Node) []error {
	errors := []error{}
	iter := sitter.NewIterator(root, sitter.BFSMode)

	for {
		node, err := iter.Next()
		if err != nil {
			break
		}
		if node == nil || !node.HasError() {
			continue
		}
		if node.IsMissing() {
			errors = append(
				errors,
				fmt.Errorf(
					"(Line: %d, Column: %d) Expected '%s', but not found",
					node.Range().StartPoint.Row+1,
					node.Range().StartPoint.Column+1,
					node.Type(),
				),
			)
		} else if node.IsError() {
			errors = append(errors, fmt.Errorf("(Line: %d, Column: %d)-(Line: %d, Column: %d) Unexpected sequence of characters", node.Range().StartPoint.Row+1, node.Range().StartPoint.Column+1, node.Range().EndPoint.Row+1, node.Range().EndPoint.Column+1))
		}
	}
	return errors
}

func (p *TreeSitterRefal5Parser) Parse(source []byte) (*ast.AST, []error) {
	var result *ast.AST
	var cursor *sitter.QueryCursor
	tree, err := p.parser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to parse source code: %w", err)}
	}

	root := tree.RootNode()
	errors := p.CheckErrors(root)

	if len(errors) > 0 {
		return nil, errors
	}

	result = &ast.AST{
		Functions:            []*ast.FunctionNode{},
		ExternalDeclarations: map[string]interface{}{},
	}

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
		funcEntryNode := funcNameNode.Parent().ChildByFieldName("entry")

		if funcEntryNode != nil {
			funcAstNode.Entry = true
		}

		funcAstNode.Name = funcNameNode.Content(source)
		sentences, err := p.walkFunctionBody(funcBodyNode, source)
		if err != nil {
			return nil, []error{fmt.Errorf("failed to build ast: %w", err)}
		}

		funcAstNode.Body = sentences
		result.Functions = append(result.Functions, funcAstNode)
	}

	declarations, err := p.walkExternalDeclarations(root, source)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to walk external declarations: %w", err)}
	}

	for _, declaration := range declarations {
		result.ExternalDeclarations[declaration] = struct{}{}
	}

	errors = append(errors, result.CheckVariableUsage()...)

	if len(errors) > 0 {
		return nil, errors
	}

	result.RebuildBlockSentences()

	return result, errors
}

func (p *TreeSitterRefal5Parser) walkFunctionBody(
	node *sitter.Node,
	source []byte,
) ([]*ast.SentenceNode, error) {
	if node == nil {
		return nil, fmt.Errorf("got nil node")
	}

	sentences := []*ast.SentenceNode{}

	for i := 0; i < int(node.ChildCount()); i++ {
		inner := node.Child(i)
		if inner.ChildByFieldName("sentence_eq") == nil &&
			inner.ChildByFieldName("sentence_block") == nil {
			continue
		}

		astSentenceNode := &ast.SentenceNode{
			Condtitions: []*ast.ConditionNode{},
			Lhs:         []ast.PatternNode{},
			Rhs:         &ast.SentenceRhsResultNode{},
		}
		sentenceNode := inner.ChildByFieldName("sentence_eq")

		if sentenceNode == nil {
			sentenceNode = inner.ChildByFieldName("sentence_block")
			astSentenceNode.Rhs = &ast.SentenceRhsBlockNode{}
		}

		sentenceLhsNode := sentenceNode.ChildByFieldName("lhs")
		lhs := []ast.PatternNode{}

		// walk lhs
		if sentenceLhsNode != nil {
			for i := 0; i < int(sentenceLhsNode.ChildCount()); i++ {
				lhsPartNode := sentenceLhsNode.Child(i)

				// TODO: check error
				lhsPart, _ := p.walkPattern(lhsPartNode, source)
				lhs = append(lhs, lhsPart)
			}
		}

		astSentenceNode.Lhs = lhs
		// walk conditions

		for j := 0; j < int(sentenceNode.ChildCount()); j++ {
			child := sentenceNode.Child(j)
			if child == nil || child.Type() != "condition" {
				continue
			}

			condition, err := p.walkCondition(child, source)
			if err != nil {
				// TODO: check err
			}

			astSentenceNode.Condtitions = append(astSentenceNode.Condtitions, condition)

		}

		// walk rhs

		switch astSentenceNode.Rhs.GetSentenceRhsType() {
		case ast.SentenceRhsBlockType:
			rhsNode := sentenceNode.ChildByFieldName("block")
			astRhsNode := &ast.SentenceRhsBlockNode{
				Result: []ast.ResultNode{},
			}

			for i := 0; i < int(rhsNode.ChildCount()); i++ {
				if rhsNode.FieldNameForChild(i) != "expr" {
					continue
				}
				resultNode := rhsNode.Child(i)
				tmp, _ := p.walkResult(resultNode, source)
				astRhsNode.Result = append(astRhsNode.Result, tmp)
			}

			bodyNode := rhsNode.ChildByFieldName("body")
			astBody, _ := p.walkFunctionBody(bodyNode, source)
			astRhsNode.Body = astBody

			astSentenceNode.Rhs = astRhsNode
		case ast.SentenceRhsResultType:
			rhsNode := sentenceNode.ChildByFieldName("rhs")

			if rhsNode == nil {

				sentences = append(sentences, astSentenceNode)
				continue
			}

			astRhsNode := &ast.SentenceRhsResultNode{
				Result: []ast.ResultNode{},
			}
			for i := 0; i < int(rhsNode.ChildCount()); i++ {
				child := rhsNode.Child(i)
				tmp, _ := p.walkResult(child, source)
				astRhsNode.Result = append(astRhsNode.Result, tmp)
			}
			astSentenceNode.Rhs = astRhsNode
		}

		sentences = append(sentences, astSentenceNode)

	}

	return sentences, nil
}

func (p *TreeSitterRefal5Parser) walkPattern(
	node *sitter.Node,
	source []byte,
) (ast.PatternNode, error) {
	// TODO: check node == nil
	// var result ast.PatternNode

	switch node.Type() {
	case "ident":
		return &ast.WordPatternNode{Value: node.Content(source)}, nil
	case "string":
		return &ast.StringPatternNode{Value: node.Content(source)}, nil
	case "number":
		val, _ := strconv.Atoi(node.Content(source))
		return &ast.NumberPatternNode{Value: uint(val)}, nil
	case "variable":
		varStrType := node.ChildByFieldName("type").Content(source)
		varType := ast.SymbolVarType
		switch varStrType {
		case "t":
			varType = ast.TermVarType
		case "e":
			varType = ast.ExprVarType
		}
		return &ast.VarPatternNode{
			Name: node.ChildByFieldName("name").Content(source),
			Type: varType,
		}, nil
	case "grouped_pattern":
		pattern := &ast.GroupedPatternNode{Patterns: []ast.PatternNode{}}
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
		chars := []byte(node.Content(source))[1:]
		chars = chars[:len(chars)-1]
		return &ast.CharactersPatternNode{Value: chars}, nil
	}

	return nil, fmt.Errorf("undefined pattern")
}

func (p *TreeSitterRefal5Parser) walkResult(
	node *sitter.Node,
	source []byte,
) (ast.ResultNode, error) {
	// TODO: check node == nil
	// var result ast.PatternNode

	switch node.Type() {
	case "ident":
		return &ast.WordResultNode{Value: node.Content(source)}, nil
	case "string":
		return &ast.StringResultNode{Value: node.Content(source)}, nil
	case "number":
		val, _ := strconv.Atoi(node.Content(source))
		return &ast.NumberResultNode{Value: uint(val)}, nil
	case "variable":
		varStrType := node.ChildByFieldName("type").Content(source)
		varType := ast.SymbolVarType
		switch varStrType {
		case "t":
			varType = ast.TermVarType
		case "e":
			varType = ast.ExprVarType
		}
		return &ast.VarResultNode{
			Name: node.ChildByFieldName("name").Content(source),
			Type: varType,
		}, nil
	case "grouped_expr":
		pattern := &ast.GroupedResultNode{Results: []ast.ResultNode{}}
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if !child.IsNamed() {
				continue
			}

			// TODO: check error
			nestedResult, _ := p.walkResult(child, source)
			pattern.Results = append(pattern.Results, nestedResult)

		}
		return pattern, nil
	case "symbols":
		chars := []byte(node.Content(source))[1:]
		chars = chars[:len(chars)-1]
		return &ast.CharactersResultNode{Value: chars}, nil
	case "function_call":
		functionCallNode := &ast.FunctionCallResultNode{
			Ident: "",
			Args:  []ast.ResultNode{},
		}
		nameNode := node.ChildByFieldName("name")
		functionCallNode.Ident = nameNode.Content(source)

		for i := 0; i < int(node.ChildCount()); i++ {
			if node.FieldNameForChild(i) == "param" {
				child := node.Child(i)
				arg, _ := p.walkResult(child, source)
				functionCallNode.Args = append(functionCallNode.Args, arg)
			}
		}

		return functionCallNode, nil
	}

	return nil, fmt.Errorf("undefined result")
}

func (p *TreeSitterRefal5Parser) walkCondition(
	condition *sitter.Node,
	source []byte,
) (*ast.ConditionNode, error) {
	astConditionNode := &ast.ConditionNode{
		Result:  []ast.ResultNode{},
		Pattern: []ast.PatternNode{},
	}

	for i := 0; i < int(condition.ChildCount()); i++ {
		conditionChild := condition.Child(i)
		switch condition.FieldNameForChild(i) {
		case "result":
			result, err := p.walkResult(conditionChild, source)
			if err != nil {
				return nil, fmt.Errorf("failed to walk result in walkCondition: %w", err)
			}
			astConditionNode.Result = append(astConditionNode.Result, result)
		case "pattern":
			pattern, err := p.walkPattern(conditionChild, source)
			if err != nil {
				return nil, fmt.Errorf("failed to walk pattern in walkCondition: %w", err)
			}
			astConditionNode.Pattern = append(astConditionNode.Pattern, pattern)
		}
	}

	return astConditionNode, nil
}

func (p *TreeSitterRefal5Parser) walkExternalDeclarations(
	root *sitter.Node,
	source []byte,
) ([]string, error) {
	cursor := sitter.NewQueryCursor()
	query, err := sitter.NewQuery([]byte(`
	(external_declaration
		(external_modifier)
		func_name_list: (function_name_list) @func_name_list
	)`), tree_sitter_refal5.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to build sitter query in walkExternalDeclarations: %w", err)
	}

	cursor.Exec(query, root)
	externals := []string{}

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		funcNameListNode := match.Captures[0].Node

		for i := 0; i < int(funcNameListNode.ChildCount()); i++ {
			child := funcNameListNode.Child(i)
			if !child.IsNamed() {
				continue
			}
			externals = append(externals, child.Content(source))
		}
	}

	return externals, nil
}

func (p *TreeSitterRefal5Parser) ParseFiles(
	progs [][]byte,
) ([]*ast.AST, *ast.FunctionNode, [][]error) {
	var goFunctPtr *ast.FunctionNode = nil
	trees := []*ast.AST{}
	errors := make([][]error, len(progs))
	foundErrors := false

	for i, prog := range progs {
		tree, fileErrors := p.Parse(prog)
		foundErrors = len(fileErrors) > 0 || foundErrors
		errors[i] = append(errors[i], fileErrors...)
		trees = append(trees, tree)
	}

	if foundErrors {
		return nil, nil, errors
	}

	entryFuncMapping := map[string]*ast.FunctionNode{}
	funcToSourceMapping := map[string]int{}
	newFuncMapping := map[string]*ast.FunctionNode{}

	for idx := range trees {
		funcMapping, fileErrors := p.UpdateFunctionsForManyFilesCompilation(idx, trees)
		foundErrors = len(fileErrors) > 0 || foundErrors
		errors[idx] = append(errors[idx], fileErrors...)
		for name, function := range funcMapping {
			newFuncMapping[function.Name] = function
			if function.Entry {
				if j, ok := funcToSourceMapping[name]; ok {
					foundErrors = true
					err := fmt.Errorf("Entry function %s is multiple defined", name)
					errors[idx] = append(errors[idx], err)
					errors[j] = append(errors[j], err)
				} else {
					funcToSourceMapping[name] = idx
					entryFuncMapping[name] = function
				}
			}
		}
	}

	for idx := range trees {
		trees[idx].AddMuFunction(funcToSourceMapping, idx)
		muFunc := trees[idx].Functions[len(trees[idx].Functions)-1]
		muMapping := map[string]*ast.FunctionNode{}
		muMapping["Mu"] = muFunc
		newFuncMapping[fmt.Sprintf("Mu%d", idx)] = muFunc

		p.UpdateFunctionsCallsForManyFilesCompilation(
			muMapping,
			map[string]*ast.FunctionNode{},
			trees[idx],
			false,
			false,
		)
	}

	for idx := range trees {
		fileErrors := p.UpdateFunctionsCallsForManyFilesCompilation(
			entryFuncMapping,
			newFuncMapping,
			trees[idx],
			true,
			true,
		)

		foundErrors = len(fileErrors) > 0 || foundErrors
		errors[idx] = append(errors[idx], fileErrors...)
	}

	if f, ok := entryFuncMapping["GO"]; ok {
		goFunctPtr = f
	} else {
		goFunctPtr = entryFuncMapping["Go"]
	}

	return trees, goFunctPtr, errors
}

func (p *TreeSitterRefal5Parser) UpdateFunctionsForManyFilesCompilation(
	target int,
	trees []*ast.AST,
) (map[string]*ast.FunctionNode, []error) {
	targetTree := trees[target]
	errors := []error{}

	funcMapping := map[string]*ast.FunctionNode{}
	newToOldFuncMapping := map[string]*ast.FunctionNode{}

	for idx, function := range targetTree.Functions {
		updatedFunction := &ast.FunctionNode{
			Name:  fmt.Sprintf("%s%d", function.Name, target),
			Entry: function.Entry,
			Body:  function.Body,
		}

		if _, ok := funcMapping[function.Name]; ok {
			errors = append(errors,
				fmt.Errorf("Function %s is multiple defined", function.Name))
		}

		funcMapping[function.Name] = updatedFunction
		newToOldFuncMapping[updatedFunction.Name] = updatedFunction

		targetTree.Functions[idx] = updatedFunction
	}

	p.UpdateFunctionsCallsForManyFilesCompilation(
		funcMapping,
		newToOldFuncMapping,
		targetTree,
		false,
		false,
	)

	return funcMapping, errors
}

func (p *TreeSitterRefal5Parser) UpdateFunctionsCallsForManyFilesCompilation(
	funcMapping map[string]*ast.FunctionNode,
	newToOldFuncMapping map[string]*ast.FunctionNode,
	tree *ast.AST,
	onlyExternals bool,
	triggerUndefinedCalls bool,
) []error {
	sentences := []*ast.SentenceNode{}
	queue := []ast.ResultNode{}
	errors := []error{}

	for _, function := range tree.Functions {
		sentences = append(sentences, function.Body...)
	}

	for len(sentences) > 0 {
		sentence := sentences[0]
		sentences = sentences[1:]
		for _, cond := range sentence.Condtitions {
			queue = append(queue, cond.Result...)
		}

		sentenceRhs := sentence.Rhs

		if sentenceRhs.GetSentenceRhsType() == ast.SentenceRhsBlockType {
			blockRhs := sentenceRhs.(*ast.SentenceRhsBlockNode)
			sentences = append(sentences, blockRhs.Body...)
			queue = append(queue, blockRhs.Result...)
		} else {
			resultRhs := sentenceRhs.(*ast.SentenceRhsResultNode)
			queue = append(queue, resultRhs.Result...)
		}
	}

	for len(queue) > 0 {
		result := queue[0]
		queue = queue[1:]

		if result.GetResultType() == ast.GroupedResultType {
			groupedNode := result.(*ast.GroupedResultNode)
			queue = append(queue, groupedNode.Results...)
			continue
		}

		if result.GetResultType() != ast.FunctionCallResultType {
			continue
		}

		functionCall := result.(*ast.FunctionCallResultNode)
		queue = append(queue, functionCall.Args...)
		if function, ok := funcMapping[functionCall.Ident]; ok {
			if _, ok := tree.ExternalDeclarations[functionCall.Ident]; (ok && onlyExternals) ||
				!onlyExternals {
				functionCall.Ident = function.Name
			}
		} else if _, ok := newToOldFuncMapping[functionCall.Ident]; !ok && triggerUndefinedCalls {
			if _, ok := library.LibraryFunctions[functionCall.Ident]; ok && triggerUndefinedCalls {
				continue
			}
			errors = append(errors, fmt.Errorf("Function %s is not defined", functionCall.Ident))
		}
	}

	return errors
}
