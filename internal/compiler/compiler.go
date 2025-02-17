package compiler

import (
	"fmt"
	"os"

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
			c.parser.Parse(source)
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
