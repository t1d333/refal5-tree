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

	if len(sentence.Condtitions) == 0 {
		return
	}

	resultVariables := []ResultNode{}

	for _, v := range variables {
		varNode := v.(*VarPatternNode)
		resultVariables = append(resultVariables, varNode.ToResultNode())
	}

	firstConditon := sentence.Condtitions[0]
	otherConditions := sentence.Condtitions[1:]

	contFunction := &FunctionNode{
		Name:  fmt.Sprintf("%sCont", f.Name),
		Entry: false,
		Body:  f.Body[sentenceIdx+1:],
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
	}

	tmpBody := []*SentenceNode{
		// Pat1 = <F_check [перем] ResC1>;
		{
			Lhs: append(sentence.Lhs),
			Rhs: &SentenceRhsResultNode{
				Result: []ResultNode{
					&FunctionCallResultNode{
						Ident: fmt.Sprintf("%sCheck", f.Name),
						Args:  firstConditon.Result,
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
