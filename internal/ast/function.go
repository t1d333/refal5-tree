package ast

import "fmt"

type FunctionNode struct {
	Name  string
	Entry bool
	Body  []*SentenceNode
}

func (f *FunctionNode) String() string {
	res := ""
	if f.Entry {
		res += "$ENTRY "
	}

	res += fmt.Sprintf("%s {\n", f.Name)
	for _, s := range f.Body {
		res += fmt.Sprintf("\t%s;\n", s.String())
	}

	res += "}"
	return res
}
