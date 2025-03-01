package compiler

import (
	"fmt"
	"os"

	astp "github.com/t1d333/refal5-tree/internal/ast"
	"github.com/t1d333/refal5-tree/internal/parser"
)

type Compiler struct {
	parser parser.Refal5Parser
}

func NewRefal5Compiler() *Compiler {
	return &Compiler{
		parser: parser.NewTreeSitterRefal5Parser(),
	}
}

func (c *Compiler) Compile(files []string, options CompilerOptions) {
	sources := [][]byte{}
	for _, file := range files {
		code, err := c.readFile(file)
		if err != nil {
			// TODO: wrap error
			return
		}

		sources = append(sources, code)

		for _, source := range sources {
			ast, _ := c.parser.Parse(source)
			for _, f := range ast.Functions {
				fmt.Println("Function ", f.Name, "Entry ", f.Entry)
				for _, s := range f.Body {
					for _, l := range s.Lhs {
						fmt.Println(l)
					}

					for _, c := range s.Condtitions {
						fmt.Println("Condition", *c)
					}

					fmt.Println("Rhs", s.Rhs, s.Rhs.GetSentenceRhsType())
					if rhs, ok := s.Rhs.(*astp.SentenceRhsBlockNode); ok {
						fmt.Println("Rhs result", rhs.Result)
						fmt.Println("Rhs body", rhs.Body)
						for _, s1 := range rhs.Body {
							fmt.Println("Block sentance", s1)
						}
					} else {
						fmt.Println("not ok")
					}
				}
			}
		}
	}
}

func (c *Compiler) readFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", path, err)
	}

	return file, nil
}
