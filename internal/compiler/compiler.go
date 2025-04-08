package compiler

import (
	// "bytes"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/t1d333/refal5-tree/internal/ast"
	"github.com/t1d333/refal5-tree/internal/parser"
)

// Templates

const (
	mainFileTmplString = `
package main

import (
	"github.com/t1d333/refal5-tree/pkg/runtime"
	// "fmt"
)

// Autogenerated functions
{{- range .Functions }}
	 {{ template "r5t-func" . }}
{{- end }}

// Rope with view field
var viewField []runtime.ViewFieldNode

func main() {

	gofunc := &runtime.R5Function{
		Name:	 "GO", 
		Entry: true,
		Ptr:   r5tGO_,
	}
	
	viewField = runtime.InitViewField(gofunc)
	runtime.StartMainLoop(viewField)
}`

	compiledFunctionTmplString = `
func r5t{{.Name}}_ (l, r int, arg *runtime.Rope, viewFieldRhs *[]runtime.ViewFieldNode) {
	{{ range .Body }}
		{{ template "r5t-sentence" . }}
	{{ end }}
  panic("Recognition failed")
}`

	compiledSentenceTmplString = `
  for i := 0; i < 1; i++ {
{{- range $name, $idxs := .VarsToIdxs }}
    /* {{ $name }}: {{- range $place := $idxs }} {{ index $place 0}} {{- end }} */
{{- end }}
    var p []int = make([]int, {{ .VarsArrSize }})
		p[0] = l
		p[1] = r

{{- range $i, $cmd := .Commands}}
    {{ $cmd }}
{{- end}}
		{{ if not .NeedLoopReturn }}
			result := runtime.NewRope([]runtime.R5Node{}) 
			localViewField := &[]runtime.ViewFieldNode{}
			{{ range $cmd := .BuildResultCmds }}
				{{ $cmd }}
			{{ end }}
			if result.Len() > 0 {
				runtime.BuildRopeViewFieldNode(result, localViewField)
			}
			*viewFieldRhs  = append(*localViewField, *viewFieldRhs...)
			return
		{{ end }}
  }`

	elementaryMatchCommandTmplString = `
{{ if eq .NodeType "Empty" }}
if (!runtime.R5t{{.NodeType}}(p[{{.LeftBorder}}], p[{{.RightBorder}}], arg)) {
	continue
}
{{ else }}
	{{if or (eq .NodeType "CloseExprVar") (or (eq .NodeType "SymbolVar") (eq .NodeType "TermVar")) }}
		if (!runtime.R5t{{.NodeType}}{{.Side}}({{.Idx}}, p[{{.LeftBorder}}], p[{{.RightBorder}}], arg, p)) {
			continue
		}
	{{ else }}
		if (!runtime.R5t{{.NodeType}}{{.Side}}({{.Idx}}, p[{{.LeftBorder}}], p[{{.RightBorder}}]{{if ne .NodeType "Brackets" }}, {{.Value}}{{ end }}, arg, p)) {
			continue
		}
	{{ end }}
{{ end }}`
	openExprVarLoopMatchCommandTmplString = `
p[{{ .Idx }}] = p[{{ .Left }}] + 1 
p[{{ .Idx }} + 1] = p[{{ .Left }}]
for end := true; end; end = runtime.R5tOpenEvarAdvance({{ .Idx }}, p[{{ .Right }}], arg, p) {
	{{ range $cmd := .Cmds }}
		{{ $cmd }}
	{{ end }}	
	{{ if .NeedReturn }}
		result := runtime.NewRope([]runtime.R5Node{}) 
		localViewField := &[]runtime.ViewFieldNode{}
		{{ range $cmd := .BuildResultCmds }}
			{{ $cmd }}
		{{ end }}
		
		if result.Len() > 0 {
			runtime.BuildRopeViewFieldNode(result, localViewField)
		}
		
		*viewFieldRhs  = append(*localViewField, *viewFieldRhs...)
		return
	{{ end }}
}`
)

type MatchCmdSideType string

const (
	LeftMatchCmdType  MatchCmdSideType = "Left"
	RightMatchCmdType MatchCmdSideType = "Right"
)

type MatchCmdNodeType string

const (
	CharMatchCmdNodeType              MatchCmdNodeType = "Char"
	BracketsMatchCmdNodeType          MatchCmdNodeType = "Brackets"
	NumberMatchCmdNodeType            MatchCmdNodeType = "Number"
	StringMatchCmdNodeType            MatchCmdNodeType = "String"
	FunctionMatchCmdNodeType          MatchCmdNodeType = "Function"
	SymbolVarMatchCmdNodeType         MatchCmdNodeType = "SymbolVar"
	TermVarMatchCmdNodeType           MatchCmdNodeType = "TermVar"
	RepeatedSymbolVarMatchCmdNodeType MatchCmdNodeType = "RepeatedSymbolVar"
	RepeatedTermVarMatchCmdNodeType   MatchCmdNodeType = "RepeatedExprTermVar"
	RepeatedExprVarMatchCmdNodeType   MatchCmdNodeType = "RepeatedExprTermVar"
	CloseExprVarMatchCmdNodeType      MatchCmdNodeType = "CloseExprVar"
	EmptyNodeType                     MatchCmdNodeType = "Empty"
)

type MatchCmdArg struct {
	NodeType    MatchCmdNodeType
	Side        MatchCmdSideType
	Idx         int
	LeftBorder  int
	RightBorder int
	Value       string
}

type ExprVarLoopCmdArg struct {
	Idx             int
	Left            int
	Right           int
	Cmds            []string
	BuildResultCmds []string
	NeedReturn      bool
}

type BuildResultCmdArg struct {
	Cmds []string
}

type CompiledProgram struct {
	Functions []CompiledFunction
}

type CompiledFunction struct {
	Name  string
	Body  []CompiledSentence
	Entry bool
}

type CompiledSentence struct {
	VarsArrSize     int
	VarsToIdxs      map[string][][]int
	Commands        []string
	BuildResultCmds []string
	NeedLoopReturn  bool
}

type Compiler struct {
	parser                     parser.Refal5Parser
	compiledProgramTmpl        *template.Template
	compiledFunctionTmpl       *template.Template
	compiledSentenceTmpl       *template.Template
	compiledMatchCmdTmpl       *template.Template
	compiledOpenExprVarCmdTmpl *template.Template
}

func NewRefal5Compiler() *Compiler {
	mainTmpl, _ := template.New("r5t-main").Parse(mainFileTmplString)
	funcTmpl := template.Must(mainTmpl.New("r5t-func").Parse(compiledFunctionTmplString))
	sentenceTmpl := template.Must(funcTmpl.New("r5t-sentence").Parse(compiledSentenceTmplString))
	matchCmdTmpl, _ := template.New("r5t-match-cmd").Parse(elementaryMatchCommandTmplString)
	openExprVarLoopTmpl, _ := template.New("r5t-open-evar-loop-cmd").
		Parse(openExprVarLoopMatchCommandTmplString)
	// buildResultCmdTmpl, _ := template.New("r5t-build-result-cmd").
	// Parse(buildResultCommandTmplString)

	compiler := &Compiler{
		parser:                     parser.NewTreeSitterRefal5Parser(),
		compiledSentenceTmpl:       sentenceTmpl,
		compiledProgramTmpl:        mainTmpl,
		compiledFunctionTmpl:       funcTmpl,
		compiledMatchCmdTmpl:       matchCmdTmpl,
		compiledOpenExprVarCmdTmpl: openExprVarLoopTmpl,
	}

	return compiler
}

func (c *Compiler) Compile(files []string, options CompilerOptions) {
	sources := [][]byte{}
	trees := []*ast.AST{}
	for _, file := range files {
		code, err := c.readFile(file)
		if err != nil {
			// TODO: wrap error
			return
		}

		sources = append(sources, code)

		for _, source := range sources {
			ast, _ := c.parser.Parse(source)
			trees = append(trees, ast)
		}

		c.Generate(trees)
	}
}

func (c *Compiler) readFile(path string) ([]byte, error) {
	file, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return file, nil
}

func (c *Compiler) Generate(trees []*ast.AST) (string, error) {
	// TODO: find Go function and generate code for GO
	// TODO: generate code for another functions

	// Templates
	functions := []CompiledFunction{}

	i := 0
	for i < len(trees) {
		tree := trees[i]
		i += 1
		j := 0
		for j < len(tree.Functions) {
			function := tree.Functions[j]
			j += 1

			generatedBody, err := c.GenerateFunctionBodyCode(tree, function)
			if err != nil {
				return "", fmt.Errorf(
					"failed to generate code for function %s: %w",
					function.Name,
					err,
				)
			}

			compiled := CompiledFunction{
				Name:  function.Name,
				Body:  generatedBody,
				Entry: function.Entry,
			}

			functions = append(functions, compiled)
		}

		c.compiledProgramTmpl.Execute(os.Stdout,
			CompiledProgram{
				Functions: functions,
			},
		)
	}

	// mainTmpl.Execute(os.Stdout, compiledProgram{Functions: generatedFunctions})
	return "", nil
}

func (c *Compiler) GenerateFunctionBodyCode(
	tree *ast.AST,
	function *ast.FunctionNode,
) ([]CompiledSentence, error) {
	body := []CompiledSentence{}

	i := 0
	for i < len(function.Body) {
		// fmt.Println("Lhs")
		// for _, n := range function.Body[i].Lhs {
		// 	ast.PrintPattern(n)
		// }

		// rhs := function.Body[i].Rhs.(*ast.SentenceRhsResultNode)

		// fmt.Println("Rhs")
		// for _, n := range rhs.Result {
		// 	ast.PrintResult(n)
		// }

		compiledSentence := c.GenerateSentence(tree, function, function.Body[i], i)
		i += 1

		body = append(body, compiledSentence)
	}

	return body, nil
}

type patternHole struct {
	patterns []ast.PatternNode
	borders  [][]int
}

type exprVarLoop struct {
	Idx   int
	Left  int
	Right int
	Cmds  []string
}

func (c *Compiler) GenerateSentence(
	tree *ast.AST,
	f *ast.FunctionNode,
	sentence *ast.SentenceNode,
	sentenceIdx int,
) CompiledSentence {
	// fmt.Println("Generate sentence for func", f.Name, sentenceIdx)
	compiledSentence := CompiledSentence{
		VarsArrSize:    0,
		VarsToIdxs:     map[string][][]int{},
		Commands:       []string{},
		NeedLoopReturn: false,
	}

	// TODO: build pattern matching

	patternHoles := []patternHole{{
		patterns: sentence.Lhs,
		borders:  [][]int{{0, 1}},
	}}

	cmds := []string{}
	nextBorder := 2

	exprVarLoops := []exprVarLoop{}
	openEvars := []*ast.VarPatternNode{}
	allVars := []ast.PatternNode{}

	for len(patternHoles) > 0 {
		hole := patternHoles[0]
		patternHoles = patternHoles[1:]
		patterns := hole.patterns

		borders := hole.borders

		for len(borders) > 0 {

			left, right := borders[0][0], borders[0][1]
			borders = borders[1:]

			cmdArg := MatchCmdArg{
				Idx:         nextBorder,
				LeftBorder:  left,
				RightBorder: right,
			}

			// TODO: check empty hole
			if len(patterns) == 0 {
				cmdArg.NodeType = EmptyNodeType
				cmd := c.generateMatchCmd(cmdArg)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}

				break
			}

			// check if left hole is symbol
			if charNode, ok := patterns[0].(*ast.CharactersPatternNode); ok {
				cmdArg.NodeType = CharMatchCmdNodeType
				cmdArg.Side = LeftMatchCmdType
				cmdArg.Value = fmt.Sprintf("'%c'", charNode.Value[0])
				cmd := c.generateMatchCmd(cmdArg)

				charNode.Value = charNode.Value[1:]

				if len(charNode.Value) == 0 {
					patterns = patterns[1:]
				}

				borders = append([][]int{{nextBorder, right}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 1
				continue
			}

			// check if right hole is symbol
			if charNode, ok := patterns[len(patterns)-1].(*ast.CharactersPatternNode); ok {
				cmdArg.NodeType = CharMatchCmdNodeType
				cmdArg.Side = RightMatchCmdType
				cmdArg.Value = fmt.Sprintf("'%c'", charNode.Value[0])
				cmd := c.generateMatchCmd(cmdArg)

				charNode.Value = charNode.Value[1:]

				if len(charNode.Value) == 0 {
					patterns = patterns[:len(patterns)-1]
				}

				borders = append([][]int{{left, nextBorder}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 1
				continue
			}

			// check if left hole is number
			if numberNode, ok := patterns[0].(*ast.NumberPatternNode); ok {
				cmdArg.NodeType = NumberMatchCmdNodeType
				cmdArg.Side = LeftMatchCmdType
				cmdArg.Value = fmt.Sprintf("%d", numberNode.Value)
				cmd := c.generateMatchCmd(cmdArg)

				patterns = patterns[1:]
				borders = append([][]int{{nextBorder, right}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 1
				continue
			}

			// check if right hole is number
			if numberNode, ok := patterns[len(patterns)-1].(*ast.NumberPatternNode); ok {
				cmdArg.NodeType = NumberMatchCmdNodeType
				cmdArg.Side = RightMatchCmdType
				cmdArg.Value = fmt.Sprintf("%d", numberNode.Value)

				cmd := c.generateMatchCmd(cmdArg)

				patterns = patterns[:len(patterns)-1]
				borders = append([][]int{{left, nextBorder}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 1
				continue
			}

			// check if left hole is string
			if strNode, ok := patterns[0].(*ast.StringPatternNode); ok {
				cmdArg.NodeType = StringMatchCmdNodeType
				cmdArg.Side = LeftMatchCmdType
				cmdArg.Value = fmt.Sprintf("%s", strNode.Value)
				cmd := c.generateMatchCmd(cmdArg)

				patterns = patterns[1:]
				borders = append([][]int{{nextBorder, right}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 1
				continue
			}

			// check if right hole is string
			if strNode, ok := patterns[len(patterns)-1].(*ast.StringPatternNode); ok {
				cmdArg.NodeType = StringMatchCmdNodeType
				cmdArg.Side = RightMatchCmdType
				cmdArg.Value = fmt.Sprintf("%s", strNode.Value)

				cmd := c.generateMatchCmd(cmdArg)

				patterns = patterns[:len(patterns)-1]
				borders = append([][]int{{left, nextBorder}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 1
				continue
			}

			// TODO: check if left hole is bracket
			if grouped, ok := patterns[0].(*ast.GroupedPatternNode); ok {
				cmdArg.NodeType = BracketsMatchCmdNodeType
				cmdArg.Side = LeftMatchCmdType
				cmd := c.generateMatchCmd(cmdArg)

				patterns = patterns[1:]
				patternHoles = append(patternHoles, patternHole{
					patterns: grouped.Patterns,
					borders:  [][]int{{nextBorder, nextBorder + 1}},
				})
				borders = append([][]int{{nextBorder + 1, right}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 2
				continue
			}

			// TODO: check if right hole is bracket
			if grouped, ok := patterns[len(patterns)-1].(*ast.GroupedPatternNode); ok {
				cmdArg.NodeType = BracketsMatchCmdNodeType
				cmdArg.Side = RightMatchCmdType
				cmd := c.generateMatchCmd(cmdArg)

				patterns = patterns[:len(patterns)-1]
				patternHoles = append(patternHoles, patternHole{
					patterns: grouped.Patterns,
					borders:  [][]int{{nextBorder, nextBorder + 1}},
				})
				borders = append([][]int{{left, nextBorder + 1}}, borders...)
				if len(exprVarLoops) > 0 {
					exprVarLoops[len(exprVarLoops)-1].Cmds = append(
						exprVarLoops[len(exprVarLoops)-1].Cmds,
						cmd,
					)
				} else {
					cmds = append(cmds, cmd)
				}
				nextBorder += 2
				continue
			}

			var varNode *ast.VarPatternNode

			leftVarNode := patterns[0].(*ast.VarPatternNode)
			rightVarNode := patterns[len(patterns)-1].(*ast.VarPatternNode)

			needLeft := false
			needRight := false
			if leftVarNode != nil {
				if leftVarNode.Type == ast.SymbolVarType || leftVarNode.Type == ast.TermVarType {
					needLeft = true
				} else if _, ok := compiledSentence.VarsToIdxs[fmt.Sprintf("%s.%s", leftVarNode.GetVarTypeStr(), leftVarNode.Name)]; ok {
					needLeft = true
				}
			}

			if rightVarNode != nil {
				if rightVarNode.Type == ast.SymbolVarType || leftVarNode.Type == ast.TermVarType {
					needRight = true
				} else if _, ok := compiledSentence.VarsToIdxs[fmt.Sprintf("%s.%s", rightVarNode.GetVarTypeStr(), rightVarNode.Name)]; ok {
					needRight = true
				}
			}

			if needLeft || !needRight {
				cmdArg.Side = LeftMatchCmdType
				varNode = leftVarNode
				patterns = patterns[1:]
			} else if needRight {
				cmdArg.Side = RightMatchCmdType
				varNode = rightVarNode
				patterns = patterns[:len(patterns)-1]
			}

			if varNode != nil {
				allVars = append(allVars, varNode)
			}

			if varNode != nil &&
				(varNode.Type == ast.SymbolVarType || varNode.Type == ast.TermVarType) {

				ident := fmt.Sprintf("%s.%s", varNode.GetVarTypeStr(), varNode.Name)
				// check repeated var
				if varIdxs, ok := compiledSentence.VarsToIdxs[ident]; ok {

					// TODO: check repeated svar
					cmdArg.Value = fmt.Sprintf("%d", (varIdxs[0][0]))

					switch varNode.Type {
					case ast.SymbolVarType:
						cmdArg.NodeType = RepeatedSymbolVarMatchCmdNodeType
						compiledSentence.VarsToIdxs[ident] = append(
							compiledSentence.VarsToIdxs[ident],
							[]int{nextBorder},
						)
						switch cmdArg.Side {
						case LeftMatchCmdType:
							borders = append([][]int{{nextBorder, right}}, borders...)
						case RightMatchCmdType:
							borders = append([][]int{{left, nextBorder}}, borders...)
						}
						nextBorder += 1
					case ast.TermVarType:
						cmdArg.NodeType = RepeatedTermVarMatchCmdNodeType

						compiledSentence.VarsToIdxs[ident] = append(
							compiledSentence.VarsToIdxs[ident],
							[]int{nextBorder, nextBorder + 1},
						)
						switch cmdArg.Side {
						case LeftMatchCmdType:
							borders = append([][]int{{nextBorder + 1, right}}, borders...)
						case RightMatchCmdType:
							borders = append([][]int{{left, nextBorder + 1}}, borders...)
						}
						nextBorder += 2
					}

					cmd := c.generateMatchCmd(cmdArg)
					if len(exprVarLoops) > 0 {
						exprVarLoops[len(exprVarLoops)-1].Cmds = append(
							exprVarLoops[len(exprVarLoops)-1].Cmds,
							cmd,
						)
					} else {
						cmds = append(cmds, cmd)
					}
					continue
				} else {

					switch varNode.Type {
					case ast.SymbolVarType:
						cmdArg.NodeType = SymbolVarMatchCmdNodeType
						compiledSentence.VarsToIdxs[ident] = append(varIdxs, []int{nextBorder})
						switch cmdArg.Side {
						case LeftMatchCmdType:
							borders = append([][]int{{nextBorder, right}}, borders...)
						case RightMatchCmdType:
							borders = append([][]int{{left, nextBorder}}, borders...)
						}
						nextBorder += 1
					case ast.TermVarType:
						cmdArg.NodeType = TermVarMatchCmdNodeType
						compiledSentence.VarsToIdxs[ident] = append(varIdxs, []int{nextBorder})
						switch cmdArg.Side {
						case LeftMatchCmdType:
							borders = append([][]int{{nextBorder + 1, right}}, borders...)
						case RightMatchCmdType:
							borders = append([][]int{{left, nextBorder + 1}}, borders...)
						}
						nextBorder += 2
					}
					cmd := c.generateMatchCmd(cmdArg)
					if len(exprVarLoops) > 0 {
						exprVarLoops[len(exprVarLoops)-1].Cmds = append(exprVarLoops[len(exprVarLoops)-1].Cmds, cmd)
					} else {
						cmds = append(cmds, cmd)
					}
					continue
				}
			}

			if varNode != nil && varNode.Type == ast.ExprVarType {
				ident := fmt.Sprintf("%s.%s", varNode.GetVarTypeStr(), varNode.Name)
				// check repeated var
				if varIdxs, ok := compiledSentence.VarsToIdxs[ident]; ok {
					cmdArg.Value = fmt.Sprintf("%d", (varIdxs[0][0]))
					cmdArg.NodeType = RepeatedExprVarMatchCmdNodeType

					compiledSentence.VarsToIdxs[ident] = append(
						compiledSentence.VarsToIdxs[ident],
						[]int{nextBorder, nextBorder + 1},
					)
					switch cmdArg.Side {
					case LeftMatchCmdType:
						borders = append([][]int{{nextBorder + 1, right}}, borders...)
					case RightMatchCmdType:
						borders = append([][]int{{left, nextBorder + 1}}, borders...)
					}
					cmd := c.generateMatchCmd(cmdArg)

					if len(exprVarLoops) > 0 {
						exprVarLoops[len(exprVarLoops)-1].Cmds = append(
							exprVarLoops[len(exprVarLoops)-1].Cmds,
							cmd,
						)
					} else {
						cmds = append(cmds, cmd)
					}
					// borders = append([][]int{{nextBorder + 1, right}}, borders...)
					nextBorder += 2
					continue
				} else if len(patterns) == 0 {
					cmdArg := MatchCmdArg{
						Idx:         nextBorder,
						LeftBorder:  left,
						RightBorder: right,
						NodeType:    CloseExprVarMatchCmdNodeType,
					}

					compiledSentence.VarsToIdxs[ident] = append(varIdxs, []int{nextBorder})
					cmd := c.generateMatchCmd(cmdArg)
					if len(exprVarLoops) > 0 {
						exprVarLoops[len(exprVarLoops)-1].Cmds = append(
							exprVarLoops[len(exprVarLoops)-1].Cmds,
							cmd,
						)
					} else {
						cmds = append(cmds, cmd)
					}
					nextBorder += 2
					continue
				} else {
					// TODO: Open evar
					// hole.pattern
					openEvars = append([]*ast.VarPatternNode{varNode}, openEvars...)
					compiledSentence.NeedLoopReturn = true
					compiledSentence.VarsToIdxs[ident] = [][]int{{nextBorder}}
					exprVarLoops = append(exprVarLoops, exprVarLoop{
						Idx:   nextBorder,
						Left:  left,
						Right: right,
						Cmds:  []string{},
					})
					borders = append([][]int{{nextBorder + 1, right}}, borders...)
					nextBorder += 2
					continue
				}

			}
			panic("Uknown pattern")
		}

	}

	tree.BuildHelpFunctionsForSentenceConditions(f, sentenceIdx, allVars, openEvars)

	// TODO: build result

	sentenceRhs := f.Body[sentenceIdx].Rhs.(*ast.SentenceRhsResultNode)

	buildResultCmds := []string{}
	for _, r := range sentenceRhs.Result {
		tmp := c.buildResultCmds(r, compiledSentence.VarsToIdxs)
		buildResultCmds = append(buildResultCmds, tmp...)
	}

	compiledSentence.BuildResultCmds = buildResultCmds
	prevLoopCmd := ""

	for i := len(exprVarLoops) - 1; i >= 0; i -= 1 {
		buff := bytes.Buffer{}
		loop := exprVarLoops[i]

		if prevLoopCmd != "" {
			loop.Cmds = append(loop.Cmds, prevLoopCmd)
		}

		tmplArg := ExprVarLoopCmdArg{
			Idx:             loop.Idx,
			Left:            loop.Left,
			Right:           loop.Right,
			Cmds:            loop.Cmds,
			NeedReturn:      i == (len(exprVarLoops) - 1),
			BuildResultCmds: compiledSentence.BuildResultCmds,
		}

		c.compiledOpenExprVarCmdTmpl.Execute(&buff, tmplArg)

		prevLoopCmd = buff.String()
	}

	cmds = append(cmds, prevLoopCmd)

	compiledSentence.VarsArrSize = nextBorder
	compiledSentence.Commands = cmds
	return compiledSentence
}

func (c *Compiler) buildResultCmds(node ast.ResultNode, varsToIdxs map[string][][]int) []string {
	switch node.GetResultType() {
	case ast.CharactersResultType:
		cNode := node.(*ast.CharactersResultNode)
		cmd := "result = result.Insert(result.Len(), []runtime.R5Node{"
		for _, c := range cNode.Value {
			cmd += fmt.Sprintf("&runtime.R5NodeChar{Char: %d}, ", c)
		}
		return []string{cmd + "})\n"}
	case ast.FunctionCallResultType:
		fNode := node.(*ast.FunctionCallResultNode)
		fCmds := []string{
			"runtime.BuildRopeViewFieldNode(result, localViewField)\n",
			"result = runtime.NewRope([]runtime.R5Node{})\n",
		}

		fCmds = append(
			fCmds,
			fmt.Sprintf(
				"runtime.BuildOpenCallViewFieldNode(runtime.R5Function{Name: \"%s\", Ptr: r5t%s_}, localViewField)\n",
				fNode.Ident,
				fNode.Ident,
			),
		)

		for _, arg := range fNode.Args {
			fCmds = append(fCmds, c.buildResultCmds(arg, varsToIdxs)...)
		}
		fCmds = append(
			fCmds,
			"runtime.BuildRopeViewFieldNode(result,localViewField)\n",
			"runtime.BuildCloseCallViewFieldNode(localViewField)\n",
			"result = runtime.NewRope([]runtime.R5Node{})\n",
		)

		return fCmds
	case ast.StringResultType:
		sNode := node.(*ast.StringResultNode)

		return []string{
			fmt.Sprintf(
				"result = result.Insert(result.Len(), []runtime.R5Node{&runtime.R5NodeString{String: %s}})",
				sNode.Value,
			),
		}
	case ast.WordResultType:
		wNode := node.(*ast.WordResultNode)

		return []string{
			fmt.Sprintf(
				"result = result.Insert(result.Len(), []runtime.R5Node{&runtime.R5NodeFunction{Name: %s}})",
				wNode.Value,
			),
		}
	case ast.NumberResultType:
		nNode := node.(*ast.NumberResultNode)

		return []string{
			fmt.Sprintf(
				"result = result.Insert(result.Len(), []runtime.R5Node{&runtime.R5NodeNumber{Number: %d}})",
				nNode.Value,
			),
		}
	case ast.VarResultType:
		vNode := node.(*ast.VarResultNode)
		if vNode.Type == ast.SymbolVarType {
			idxs := varsToIdxs[fmt.Sprintf("%s.%s", vNode.GetVarTypeStr(), vNode.Name)]
			return []string{fmt.Sprintf("runtime.CopySymbolVar(p[%d], arg, result)", idxs[0][0])}
		} else {
			idxs := varsToIdxs[fmt.Sprintf("%s.%s", vNode.GetVarTypeStr(), vNode.Name)]
			return []string{
				fmt.Sprintf("runtime.CopyExprTermVar(p[%d], p[%d], arg, result)", idxs[0][0], idxs[0][0]+1),
			}
		}
	case ast.GroupedResultType:
		gNode := node.(*ast.GroupedResultNode)
		gCmds := []string{}

		for _, r := range gNode.Results {
			gCmds = append(gCmds, c.buildResultCmds(r, varsToIdxs)...)
		}

		gCmds = append(gCmds,
			"runtime.BuildRopeViewFieldNode(result, localViewField)\n",
		)

		tmp := []string{
			"runtime.BuildOpenBracketViewFieldNode(localViewField)\n",
		}

		tmp = append(tmp, gCmds...)

		tmp = append(tmp,
			"result = runtime.NewRope([]runtime.R5Node{})\n",
			"runtime.BuildCloseBracketViewFieldNode(localViewField)\n",
		)

		return tmp
	}
	return []string{}
}

func (c *Compiler) generateMatchCmd(arg MatchCmdArg) string {
	buff := bytes.Buffer{}
	c.compiledMatchCmdTmpl.Execute(&buff, arg)

	return buff.String()
}
