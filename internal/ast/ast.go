package ast

import "fmt"

type AST struct {
	Functions            []*FunctionNode
	ExternalDeclarations []string
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
				Lhs:         append(variables, firstConditon.Pattern...),
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
				Lhs: append(variables, &VarPatternNode{
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
				Lhs: append(variables, &VarPatternNode{
					Type: ExprVarType,
					Name: "Other",
				}),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sForward1", f.Name),
							Args:  t.BuildTemplateT0(sentence.Lhs, openEvarMap),
						},
					},
				},
			},
		)

		// generate forward and next functions

		for i := range openEvarList {
			// build i forward func
			forwardFunc := &FunctionNode{
				Name:  fmt.Sprintf("%sForward%d", f.Name, i),
				Entry: false,
			}

			// build i next func
			nextFunc := &FunctionNode{
				Name:  fmt.Sprintf("%sNext%d", f.Name, i),
				Entry: false,
			}

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

	//
	// for i := 0; i <=len(openEvarList); i++ {
	// 	curr :=
	// }

	// TODO: if openevar not emty build forward and next functions
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
	return &FunctionNode{
		Name:  fmt.Sprintf("%sForward%d", originFunc.Name, k),
		Entry: false,
		Body: []*SentenceNode{
			//  T1(Pat, K) = <F_next_K T2(Pat, K)>;
			{
				Lhs: t.BuildTemplateT1(k, originSentence.Lhs, openEvars),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sNext%d", originFunc.Name, k),
							Args:  t.BuildTemplateT2(k, originSentence.Lhs, openEvars),
						},
					},
				},
			},
			// T3(Pat, K) = <F_forward_K+1 T4(Pat, K)>;
			{
				Lhs: t.BuildTemplateT3(k, originSentence.Lhs, openEvars),
				Rhs: &SentenceRhsResultNode{
					Result: []ResultNode{
						&FunctionCallResultNode{
							Ident: fmt.Sprintf("%sNext%d", originFunc.Name, k+1),
							Args:  t.BuildTemplateT4(k, originSentence.Lhs, openEvars),
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
	openEvars []*VarPatternNode,
) *FunctionNode {
	return &FunctionNode{
		Name:  fmt.Sprintf("%sNext%d", originFunc, k),
		Entry: false,
		Body: []*SentenceNode{
			{},
			{},
		},
	}
}

func (a *AST) BuildTemplateT1(
	k int,
	patterns []PatternNode,
	openEvars []*VarPatternNode,
) []PatternNode {
	result := []PatternNode{}
	// queue := patterns

	return result
}

func (a *AST) BuildTemplateT2(
	k int,
	patterns []PatternNode,
	openEvars []*VarPatternNode,
) []ResultNode {
	result := []ResultNode{}
	// queue := patterns

	return result
}

func (a *AST) BuildTemplateT3(
	k int,
	patterns []PatternNode,
	openEvars []*VarPatternNode,
) []PatternNode {
	result := []PatternNode{}
	// queue := patterns

	return result
}

func (a *AST) BuildTemplateT4(
	k int,
	patterns []PatternNode,
	openEvar []*VarPatternNode,
) []ResultNode {
	result := []ResultNode{}
	// queue := patterns

	return result
}

func (a *AST) BuildTemplateT5(
	patterns []PatternNode,
	openExprVarMap map[string]bool,
) []PatternNode {
	result := []PatternNode{}
	// queue := patterns

	return result
}

func (a *AST) BuildTemplateT6(
	patterns []PatternNode,
	openExprVarMap map[string]bool,
) []PatternNode {
	result := []PatternNode{}
	// queue := patterns

	return result
}

func (a *AST) BuildTemplateT7(
	k int,
	patterns []PatternNode,
	openEvar []*VarPatternNode,
) []ResultNode {
	result := []ResultNode{}
	// queue := patterns

	return result
}
