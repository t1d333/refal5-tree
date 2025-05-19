package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/t1d333/refal5-tree/internal/library"
)

type ConditionHelpTemplateType int

const (
	T0TemplateType ConditionHelpTemplateType = iota
	T1TemplateType
	T2TemplateType
	T3TemplateType
	T4TemplateType
	T5TemplateType
	T6TemplateType
	T7TemplateType
)

type AST struct {
	Functions            []*FunctionNode
	ExternalDeclarations map[string]interface{}
}

func (t *AST) AddMuFunction(entryFunctions map[string]int, target int) error {
	muFunction := &FunctionNode{
		Name:  fmt.Sprintf("Mu%d", target),
		Entry: false,
		Body:  []*SentenceNode{},
	}
	// TODO: add stdlib functions
	// TODO: need exclude go funcs?(GO and Go)

	argPatternVar := &VarPatternNode{
		Type: ExprVarType,
		Name: "Arg",
	}

	functionCallTmp := FunctionCallResultNode{
		Args: []ResultNode{
			argPatternVar.ToResultNode(),
		},
	}

	allFunctions := []string{}

	for _, f := range t.Functions {
		allFunctions = append(allFunctions, f.Name)
	}
	//
	// for f := range t.ExternalDeclarations {
	// 	allFunctions = append(allFunctions, f)
	// }

	for f, i := range entryFunctions {
		if i == target {
			continue
		}
		allFunctions = append(allFunctions, f)
	}

	for f := range library.LibraryFunctions {
		allFunctions = append(allFunctions, f)
	}

	for i, f := range allFunctions {
		callName := f
		ident := f
		if i < len(t.Functions) {
			ident, _ = strings.CutSuffix(f, strconv.Itoa(target))
		}

		if srcIdx, ok := entryFunctions[f]; ok && srcIdx != target {
			ident, _ = strings.CutSuffix(f, strconv.Itoa(srcIdx))
			targetIdx := entryFunctions[f]
			callName = fmt.Sprintf("%s%d", f, targetIdx)
		}

		functionCall := functionCallTmp
		functionCall.Ident = callName

		if alias, ok := library.LibraryFuncionOriginToAlias[f]; ok {
			muFunction.Body = append(muFunction.Body, &SentenceNode{
				Lhs: []PatternNode{
					&WordPatternNode{Value: alias},
					argPatternVar,
				},

				Rhs: &SentenceRhsResultNode{Result: []ResultNode{&functionCall}},
			},
				&SentenceNode{
					Lhs: []PatternNode{
						&GroupedPatternNode{
							Patterns: []PatternNode{
								&CharactersPatternNode{Value: []byte(alias)},
							},
						},
						argPatternVar,
					},
					Rhs: &SentenceRhsResultNode{Result: []ResultNode{&functionCall}},
				},
			)
		}
		muFunction.Body = append(muFunction.Body, &SentenceNode{
			Lhs: []PatternNode{
				&WordPatternNode{Value: ident},
				argPatternVar,
			},

			Rhs: &SentenceRhsResultNode{Result: []ResultNode{&functionCall}},
		},
			&SentenceNode{
				Lhs: []PatternNode{
					&GroupedPatternNode{
						Patterns: []PatternNode{
							&CharactersPatternNode{Value: []byte(ident)},
						},
					},
					argPatternVar,
				},
				Rhs: &SentenceRhsResultNode{Result: []ResultNode{&functionCall}},
			},
		)
	}

	// add aliases

	t.Functions = append(t.Functions, muFunction)

	return nil
}

func (t *AST) RebuildBlockSentences() {
	i := 0

	for ; i < len(t.Functions); i++ {
		function := t.Functions[i]
		for idx, sentence := range function.Body {
			if sentence.Rhs.GetSentenceRhsType() != SentenceRhsBlockType {
				continue
			}
			sentenceBlock := sentence.Rhs.(*SentenceRhsBlockNode)

			lhsVars := t.ExtractVariables(sentence.Lhs, map[string]interface{}{})
			groupedLhsVars := t.GroupExprPatternVars(lhsVars)
			rhsFunction := &FunctionNode{
				Name:  fmt.Sprintf("%s_%d", function.Name, idx),
				Entry: false,
				Body:  []*SentenceNode{},
			}

			for _, innerSentence := range sentenceBlock.Body {
				updatedSentence := innerSentence
				updatedSentence.Lhs = append(groupedLhsVars, updatedSentence.Lhs...)
				rhsFunction.Body = append(rhsFunction.Body, updatedSentence)
			}

			t.Functions = append(t.Functions, rhsFunction)

			function.Body[idx].Rhs = &SentenceRhsResultNode{
				Result: []ResultNode{
					&FunctionCallResultNode{
						Ident: fmt.Sprintf("%s_%d", function.Name, idx),
						Args:  append(PatternsToResults(groupedLhsVars), sentenceBlock.Result...),
					},
				},
			}
		}
	}
}

func (t *AST) BuildHelpFunctionsForSentenceConditions(
	f *FunctionNode,
	sentenceIdx int,
	variables []PatternNode,
	openEvarList []*VarPatternNode,
) {
	sentence := f.Body[sentenceIdx]

	extendetVariables := []PatternNode{}
	rhsConvertedVariables := []ResultNode{}

	for _, v := range variables {
		varNode := v.(*VarPatternNode)
		if varNode.Type == ExprVarType {
			groupedNode := &GroupedPatternNode{Patterns: []PatternNode{varNode}}
			extendetVariables = append(
				extendetVariables,
				groupedNode,
			)
			rhsConvertedVariables = append(rhsConvertedVariables, PatternToResult(groupedNode))
		} else {
			extendetVariables = append(extendetVariables, varNode)
			rhsConvertedVariables = append(rhsConvertedVariables, PatternToResult(varNode))
		}

	}

	if len(sentence.Condtitions) == 0 {
		return
	}

	firstConditon := sentence.Condtitions[0]
	otherConditions := sentence.Condtitions[1:]

	contFunction := &FunctionNode{
		Name:  fmt.Sprintf("%sCont", f.Name),
		Entry: false,
		Body:  []*SentenceNode{},
	}

	for _, t := range f.Body[sentenceIdx+1:] {
		contFunction.Body = append(contFunction.Body, t)
	}

	checkFunction := &FunctionNode{
		Name:  fmt.Sprintf("%sCheck", f.Name),
		Entry: false,
		Body: []*SentenceNode{
			{
				Lhs:         append(t.GroupExprPatternVars(variables), firstConditon.Pattern...),
				Condtitions: otherConditions,
				Rhs:         sentence.Rhs,
			},
		},
	}

	if len(openEvarList) == 0 {

		lhsResults := []ResultNode{}
		for _, n := range sentence.Lhs {
			lhsResults = append(lhsResults, PatternToResult(n))
		}
		checkFunction.Body = append(checkFunction.Body,
			//   [перем] e.Other = <F_cont Pat1>;
			&SentenceNode{
				Lhs: append(t.GroupExprPatternVars(variables), &VarPatternNode{
					Type: ExprVarType,
					Name: "Other",
				}),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sCont", f.Name),
							Args:  lhsResults,
						},
					},
				},
			},
		)
	} else {
		openEvarMap := map[string]bool{}
		for _, v := range openEvarList {
			openEvarMap[v.Name] = false
		}
		checkFunction.Body = append(checkFunction.Body,

			// [перем] e.Other = <F_forward_1 T0(Pat)>;
			&SentenceNode{
				Lhs: append(t.GroupExprPatternVars(variables), &VarPatternNode{
					Type: ExprVarType,
					Name: "Other",
				}),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sForward0", f.Name),
							Args:  PatternsToResults(t.BuildConditionTemplate(-1, T0TemplateType, sentence.Lhs, openEvarList, map[string]interface{}{})),
						},
					},
				},
			},
		)

		// generate forward and next functions

		for i := range openEvarList {
			// build i forward func
			forwardFunc := t.BuildForwardFunction(i, f, sentence, openEvarList)
			// build i next func
			nextFunc := t.BuildNextFunction(i, f, sentence, firstConditon, variables, openEvarList)

			t.Functions = append(t.Functions, forwardFunc)
			t.Functions = append(t.Functions, nextFunc)
		}
	}

	tmpBody := []*SentenceNode{
		// Pat1 = <F_check [перем] ResC1>;
		{
			Lhs: append(sentence.Lhs),
			Rhs: &SentenceRhsResultNode{
				Result: []ResultNode{
					&FunctionCallResultNode{
						Ident: fmt.Sprintf("%sCheck", f.Name),
						Args:  append(rhsConvertedVariables, firstConditon.Result...),
					},
				},
			},
		},

		// e.X = <F_cont e.X>;
		{
			Lhs: []PatternNode{
				&VarPatternNode{
					Type: ExprVarType,
					Name: "X",
				},
			},
			Rhs: &SentenceRhsResultNode{
				Result: []ResultNode{
					&FunctionCallResultNode{
						Ident: fmt.Sprintf("%sCont", f.Name),
						Args: []ResultNode{
							&VarResultNode{
								Type: ExprVarType,
								Name: "X",
							},
						},
					},
				},
			},
		},
	}

	f.Body = append(f.Body[:sentenceIdx], tmpBody...)

	t.Functions = append(t.Functions, checkFunction)
	t.Functions = append(t.Functions, contFunction)
}

func (a *AST) BuildTemplateT0(
	patterns []PatternNode,
	openExprVarMap map[string]bool,
) []ResultNode {
	result := []ResultNode{}
	queue := patterns

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		switch curr.GetPatternType() {
		case GroupedPatternType:
			grouped := curr.(*GroupedPatternNode)
			result = append(result, a.BuildTemplateT0(grouped.Patterns, openExprVarMap)...)
			// queue = append(grouped.Patterns, queue...)
		case VarPatternType:
			varNode := curr.(*VarPatternNode)
			if varNode.Type == ExprVarType && !openExprVarMap[varNode.Name] {
				openExprVarMap[varNode.Name] = true

				result = append(
					result,
					&GroupedResultNode{Results: []ResultNode{varNode.ToResultNode()}},
				)
			} else {
				result = append(result, varNode.ToResultNode())
			}
		default:
			result = append(result, PatternToResult(curr))
		}
	}

	return result
}

func (t *AST) BuildForwardFunction(
	k int,
	originFunc *FunctionNode,
	originSentence *SentenceNode,
	openEvars []*VarPatternNode,
) *FunctionNode {
	nextForwardIdent := fmt.Sprintf("%sForward%d", originFunc.Name, k+1)
	if k+1 == len(openEvars) {
		nextForwardIdent = fmt.Sprintf("%sCont", originFunc.Name)
	}
	return &FunctionNode{
		Name:  fmt.Sprintf("%sForward%d", originFunc.Name, k),
		Entry: false,
		Body: []*SentenceNode{
			//  T1(Pat, K) = <F_next_K T2(Pat, K)>;
			{
				Lhs: t.BuildConditionTemplate(
					k,
					T1TemplateType,
					originSentence.Lhs,
					openEvars,
					map[string]interface{}{},
				),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sNext%d", originFunc.Name, k),
							Args: PatternsToResults(t.BuildConditionTemplate(
								k,
								T2TemplateType,
								originSentence.Lhs,
								openEvars,
								map[string]interface{}{},
							)),
						},
					},
				},
			},
			// T3(Pat, K) = <F_forward_K+1 T4(Pat, K)>;
			{
				Lhs: t.BuildConditionTemplate(
					k,
					T3TemplateType,
					originSentence.Lhs,
					openEvars,
					map[string]interface{}{},
				),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: nextForwardIdent,
							Args: PatternsToResults(
								t.BuildConditionTemplate(
									k,
									T4TemplateType,
									originSentence.Lhs,
									openEvars,
									map[string]interface{}{},
								),
							),
						},
					},
				},
			},
		},
	}
}

func (a *AST) BuildNextFunction(
	k int,
	originFunc *FunctionNode,
	originSentence *SentenceNode,
	condition *ConditionNode,
	lhsVariables []PatternNode,
	openEvars []*VarPatternNode,
) *FunctionNode {
	targetVariable := openEvars[k]
	// vars(Pat) | e.K → e.K_fix e.K_var
	if originFunc.Name == "Evaluate0" {
		for _, e := range openEvars {
			fmt.Printf("OPENEVARS: %#v\n", e)
			fmt.Printf("LAST : %#v\n", originSentence.Lhs[len(originSentence.Lhs)-1])
		}

		for _, e := range originSentence.Lhs {
			fmt.Printf("LHS : %#v\n", e)
		}

		fmt.Printf("RHS : %#v\n", originSentence.Rhs)

	}

	replacementPatternVariables := []PatternNode{
		&VarPatternNode{
			Type: ExprVarType,
			Name: fmt.Sprintf("%sFix", targetVariable.Name),
		},
		&VarPatternNode{
			Type: ExprVarType,
			Name: fmt.Sprintf("%sVar", targetVariable.Name),
		},
	}

	lhsTmp := []PatternNode{}

	for _, n := range lhsVariables {
		varNode := n.(*VarPatternNode)
		if varNode.Type == ExprVarType {
			lhsTmp = append(lhsTmp, &GroupedPatternNode{
				Patterns: []PatternNode{
					n,
				},
			})
		} else {
			lhsTmp = append(lhsTmp, n)
		}
	}

	replacedPatternVariables := ReplacePatternVariable(
		lhsTmp,
		targetVariable,
		replacementPatternVariables,
	)

	nextForwardIdent := fmt.Sprintf("%sForward%d", originFunc.Name, k+1)
	if k+1 == len(openEvars) {
		nextForwardIdent = fmt.Sprintf("%sCont", originFunc.Name)
	}

	if originFunc.Name == "Evaluate0" && k == 0 {
		for _, a := range originSentence.Lhs {
			fmt.Printf("%#v\n", a)
		}

		fmt.Printf("-------------------\n")
		t := a.BuildConditionTemplate(
			k,
			T5TemplateType,
			originSentence.Lhs,
			openEvars,
			map[string]interface{}{})

		for _, a := range t {
			fmt.Printf("%#v\n", a)
			g, ok := a.(*GroupedPatternNode)
			if ok {
				for _, b := range g.Patterns {
					fmt.Printf("%#v\n", b)
				}
			}
		}
	}

	return &FunctionNode{
		Name:  fmt.Sprintf("%sNext%d", originFunc.Name, k),
		Entry: false,
		Body: []*SentenceNode{
			// T5(Pat, K)
			// = <F_check
			// vars(Pat) | e.K → e.K_fix e.K_var
			// ResC | e.K → e.K_fix e.K_var
			// >;

			{
				Lhs: a.BuildConditionTemplate(
					k,
					T5TemplateType,
					originSentence.Lhs,
					openEvars,
					map[string]interface{}{},
				),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sCheck", originFunc.Name),
							Args: append(
								PatternsToResults(replacedPatternVariables),
								ReplaceResultVariable(
									a.GroupExprResultVars(condition.Result),
									PatternToResult(targetVariable).(*VarResultNode),
									PatternsToResults(replacementPatternVariables),
								)...),
						},
					},
				},
			},

			// T6(Pat, K) = <F_forward_K+1 T7(Pat, K)>;
			{
				Lhs: a.BuildConditionTemplate(
					k,
					T6TemplateType,
					originSentence.Lhs,
					openEvars,
					map[string]interface{}{},
				),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: nextForwardIdent,
							Args: PatternsToResults(a.BuildConditionTemplate(
								k,
								T7TemplateType,
								originSentence.Lhs,
								openEvars,
								map[string]interface{}{},
							)),
						},
					},
				},
			},
		},
	}
}

func (a *AST) BuildConditionTemplate(
	target int,
	templateType ConditionHelpTemplateType,
	patterns []PatternNode,
	openEvars []*VarPatternNode,
	varsSeen map[string]interface{},
) []PatternNode {
	resultLhs := []PatternNode{}
	resultRhs := []PatternNode{}
	queue := patterns

	for len(queue) > 0 {

		result := &[]PatternNode{}
		var curr PatternNode = nil
		currStart := queue[0]
		currEnd := queue[len(queue)-1]

		if currStart.GetPatternType() != VarPatternType {
			curr = currStart
			queue = queue[1:]
			result = &resultLhs
		} else {
			varNode := currStart.(*VarPatternNode)
			// TODO: check if not expr var
			if varNode.Type != ExprVarType {
				curr = currStart
				queue = queue[1:]
				result = &resultLhs
				// TODO: check if var already seen
			} else if _, ok := varsSeen[varNode.Name]; ok {
				curr = currStart
				queue = queue[1:]
				result = &resultLhs
			}
		}

		if curr == nil {
			if currEnd.GetPatternType() != VarPatternType {
				curr = currEnd
				queue = queue[:len(queue)-1]
				result = &resultRhs
			} else {
				varNode := currEnd.(*VarPatternNode)
				// TODO: check if not expr var
				if varNode.Type != ExprVarType {
					curr = currEnd
					queue = queue[:len(queue)-1]
					result = &resultRhs
					// TODO: check if var already seen
				} else if _, ok := varsSeen[varNode.Name]; ok {
					curr = currEnd
					queue = queue[:len(queue)-1]
					result = &resultRhs
				}
			}
		}

		if curr == nil {
			curr = currStart
			queue = queue[1:]
			result = &resultLhs
		}

		if curr.GetPatternType() == GroupedPatternType {
			grouped := curr.(*GroupedPatternNode)
			*result = append(
				*result,
				&GroupedPatternNode{
					Patterns: a.BuildConditionTemplate(
						target,
						templateType,
						grouped.Patterns,
						openEvars,
						varsSeen,
					),
				},
			)
			continue
		}

		if curr.GetPatternType() != VarPatternType {
			*result = append(*result, curr)
			continue
		}

		varNode := curr.(*VarPatternNode)

		if varNode.Type != ExprVarType || len(queue) == 0 {
			*result = append(*result, curr)
			continue
		}

		if _, ok := varsSeen[varNode.Name]; !ok && templateType == T0TemplateType {
			*result = append(*result, &GroupedPatternNode{
				Patterns: []PatternNode{
					curr,
				},
			})
			varsSeen[varNode.Name] = struct{}{}
			continue
		} else if templateType == T0TemplateType {
			*result = append(*result, curr)
			continue
		}

		if _, ok := varsSeen[varNode.Name]; !ok && templateType == T5TemplateType &&
			varNode.Name == openEvars[target].Name {
			varsSeen[varNode.Name] = struct{}{}
			// (e.X_fix) e.X_var
			*result = append(*result, &GroupedPatternNode{[]PatternNode{&VarPatternNode{
				Name: fmt.Sprintf("%sFix", openEvars[target].Name),
				Type: ExprVarType,
			}}}, &VarPatternNode{Name: fmt.Sprintf("%sVar", openEvars[target].Name), Type: ExprVarType})
			continue
		} else if templateType == T5TemplateType && varNode.Name == openEvars[target].Name {
			// e.X_fix e.X_var
			*result = append(*result, &VarPatternNode{
				Name: fmt.Sprintf("%sFix", openEvars[target].Name),
				Type: ExprVarType,
			}, &VarPatternNode{Name: fmt.Sprintf("%sVar", openEvars[target].Name), Type: ExprVarType})
			continue
		}

		if varNode.Name == openEvars[target].Name {
			if templateType == T1TemplateType {
				// (e.X_fix) t.X_next e.X_rest
				*result = append(*result,
					&GroupedPatternNode{
						Patterns: []PatternNode{
							&VarPatternNode{
								Name: fmt.Sprintf("%sFix", varNode.Name),
								Type: ExprVarType,
							},
						},
					},
					&VarPatternNode{
						Name: fmt.Sprintf("%sNext", varNode.Name),
						Type: TermVarType,
					},
					&VarPatternNode{
						Name: fmt.Sprintf("%sRest", varNode.Name),
						Type: ExprVarType,
					},
				)
				break
			}
			if templateType == T2TemplateType {
				// (e.X_fix t.X_next) e.X_rest
				*result = append(*result,
					&GroupedPatternNode{
						Patterns: []PatternNode{
							&VarPatternNode{
								Name: fmt.Sprintf("%sFix", varNode.Name),
								Type: ExprVarType,
							},
							&VarPatternNode{
								Name: fmt.Sprintf("%sNext", varNode.Name),
								Type: TermVarType,
							},
						},
					},
					&VarPatternNode{
						Name: fmt.Sprintf("%sRest", varNode.Name),
						Type: ExprVarType,
					},
				)
				break
			}
			if templateType == T3TemplateType {
				// (e.X_fix)
				*result = append(*result,
					&GroupedPatternNode{
						Patterns: []PatternNode{
							&VarPatternNode{
								Name: fmt.Sprintf("%sFix", varNode.Name),
								Type: ExprVarType,
							},
						},
					},
				)
				break
			}
			if templateType == T4TemplateType {
				// e.X_fix
				*result = append(*result,
					&VarPatternNode{
						Name: fmt.Sprintf("%sFix", varNode.Name),
						Type: ExprVarType,
					},
				)
				break
			}
			if templateType == T6TemplateType {
				// (e.X_fix) e.X_rest
				*result = append(*result,
					&GroupedPatternNode{
						Patterns: []PatternNode{
							&VarPatternNode{
								Name: fmt.Sprintf("%sFix", varNode.Name),
								Type: ExprVarType,
							},
						},
					},
					&VarPatternNode{
						Name: fmt.Sprintf("%sRest", varNode.Name),
						Type: ExprVarType,
					},
				)
				break
			}
			if templateType == T7TemplateType {
				// e.X_fix e.X_rest
				*result = append(*result,
					&VarPatternNode{
						Name: fmt.Sprintf("%sFix", varNode.Name),
						Type: ExprVarType,
					},
					&VarPatternNode{
						Name: fmt.Sprintf("%sRest", varNode.Name),
						Type: ExprVarType,
					},
				)
				break
			}
		} else if templateType == T5TemplateType {
			isInactive := false
			for i := target + 1; i < len(openEvars); i++ {
				openVar := openEvars[i]
				if varNode.Name == openVar.Name {
					*result = append(*result, &GroupedPatternNode{Patterns: []PatternNode{
						varNode,
					}})
					isInactive = true
				}
			}

			if !isInactive {
				*result = append(*result, varNode)
			}
		} else {
			isInactive := false
			for i := 0; i < target; i++ {
				openVar := openEvars[i]
				if varNode.Name == openVar.Name {
					*result = append(*result, &VarPatternNode{
						Name: fmt.Sprintf("%sRest", varNode.Name),
						Type: ExprVarType,
					})
					queue = []PatternNode{}
					isInactive = true
					break
				}
			}

			if !isInactive {
				*result = append(*result, &GroupedPatternNode{
					Patterns: []PatternNode{varNode},
				})
			}
		}

	}

	for i := len(resultRhs) - 1; i >= 0; i-- {
		resultLhs = append(resultLhs, resultRhs[i])
	}

	return resultLhs
}

func (t *AST) GroupExprPatternVars(patterns []PatternNode) []PatternNode {
	result := []PatternNode{}
	for _, pattern := range patterns {
		varNode, ok := pattern.(*VarPatternNode)
		if !ok {
			result = append(result, pattern)
			continue
		}

		if varNode.Type == ExprVarType {
			result = append(result, &GroupedPatternNode{
				Patterns: []PatternNode{
					pattern,
				},
			})
		} else {
			result = append(result, pattern)
		}
	}
	return result
}

func (t *AST) GroupExprResultVars(results []ResultNode) []ResultNode {
	result := []ResultNode{}
	for _, r := range results {
		varNode, ok := r.(*VarResultNode)
		if !ok {
			result = append(result, r)
			continue
		}

		if varNode.Type == ExprVarType {
			result = append(result, &GroupedResultNode{
				Results: []ResultNode{
					r,
				},
			})
		} else {
			result = append(result, r)
		}
	}
	return result
}

func (t *AST) ExtractVariables(
	patterns []PatternNode,
	varsSeen map[string]interface{},
) []PatternNode {
	result := []PatternNode{}

	for _, pattern := range patterns {
		switch pattern.GetPatternType() {
		case GroupedPatternType:
			groupedNode := pattern.(*GroupedPatternNode)
			result = append(result, t.ExtractVariables(groupedNode.Patterns, varsSeen)...)
		case VarPatternType:
			varNode := pattern.(*VarPatternNode)
			ident := fmt.Sprintf("%s.%s", varNode.GetVarTypeStr(), varNode.Name)
			if _, ok := varsSeen[ident]; !ok {
				varsSeen[ident] = struct{}{}
				result = append(result, varNode)
			}
		default:
			continue

		}
	}

	return result
}

func (t *AST) collectVariables(p []PatternNode) map[string]interface{} {
	variables := map[string]interface{}{}
	queue := p

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.GetPatternType() == VarPatternType {
			variable := curr.(*VarPatternNode)
			variables[fmt.Sprintf("%s.%s", variable.GetVarTypeStr(), variable.Name)] = struct{}{}

		} else if curr.GetPatternType() == GroupedPatternType {
			grouped := curr.(*GroupedPatternNode)
			queue = append(queue, grouped.Patterns...)
		}
	}

	return variables
}

func (t *AST) checkResultVarUsage(result []ResultNode, variables map[string]interface{}) []error {
	errors := []error{}
	queue := result

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == nil {
			fmt.Println(curr)
		}

		if curr.GetResultType() == VarResultType {
			variable := curr.(*VarResultNode)
			ident := fmt.Sprintf("%s.%s", variable.GetVarTypeStr(), variable.Name)
			if _, ok := variables[ident]; !ok {
				errors = append(errors, fmt.Errorf("Variable %s not found", ident))
			}
		} else if curr.GetResultType() == GroupedResultType {
			grouped := curr.(*GroupedResultNode)
			queue = append(queue, grouped.Results...)
		} else if curr.GetResultType() == FunctionCallResultType {
			call := curr.(*FunctionCallResultNode)
			queue = append(queue, call.Args...)
		}
	}

	return errors
}

func (t *AST) checkSentenceVarUsage(
	sentence *SentenceNode,
	variables map[string]interface{},
) []error {
	errors := []error{}

	lhsVariables := t.collectVariables(sentence.Lhs)

	for v := range lhsVariables {
		variables[v] = struct{}{}
	}

	for _, condition := range sentence.Condtitions {
		conditionVariables := t.collectVariables(condition.Pattern)

		for v := range conditionVariables {
			variables[v] = struct{}{}
		}

		condErrors := t.checkResultVarUsage(condition.Result, variables)
		errors = append(errors, condErrors...)
	}

	if sentence.Rhs.GetSentenceRhsType() == SentenceRhsResultType {
		sentenceRhs := sentence.Rhs.(*SentenceRhsResultNode)
		errors = append(errors, t.checkResultVarUsage(sentenceRhs.Result, variables)...)
	} else {
		sentenceRhs := sentence.Rhs.(*SentenceRhsBlockNode)
		errors = append(errors, t.checkResultVarUsage(sentenceRhs.Result, variables)...)

		for _, sentence := range sentenceRhs.Body {
			errors = append(errors, t.checkSentenceVarUsage(sentence, variables)...)
		}
	}

	return errors
}

func (t *AST) CheckVariableUsage() []error {
	errors := []error{}
	for _, f := range t.Functions {
		for _, sentence := range f.Body {
			errors = append(errors, t.checkSentenceVarUsage(sentence, map[string]interface{}{})...)
		}
	}
	return errors
}
