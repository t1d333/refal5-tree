package library

import (
	"bufio"
	// "fmt"
	"os"

	"github.com/t1d333/refal5-tree/pkg/runtime"
)

// TODO: implement
func R5tCard(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text() // Возвращает строку без символа конца строки
		lineBytes := []byte(line)
		rope := runtime.NewRope([]runtime.R5Node{})
		for _, b := range lineBytes {
			charNode := &runtime.R5NodeChar{Char: b}
			rope = rope.Insert(rope.Len(), []runtime.R5Node{charNode})

		}

		*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: rope,
		}}, *rhsStack...)
		return
	}

	endCode := &runtime.R5NodeNumber{Number: 0}
	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope([]runtime.R5Node{endCode}),
	}}, *rhsStack...)
}

func R5tExistsFile(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	filename := ""
	for i := l; i < r; i++ {
		curr := arg.Get(i)
		if curr.Type() != runtime.R5DatatagChar {
			// Panic?
			return
		}

		charNode := curr.(*runtime.R5NodeChar)
		filename += string(charNode.Char)
	}

	_, err := os.Stat(filename)

	var result *runtime.Rope
	if os.IsNotExist(err) {
		result = runtime.NewRope([]runtime.R5Node{&runtime.R5NodeFunction{
			Function: &runtime.R5Function{
				Name: "False",
			},
		}})
	} else {
		result = runtime.NewRope([]runtime.R5Node{&runtime.R5NodeFunction{
			Function: &runtime.R5Function{
				Name: "True",
			},
		}})
	}
	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{Value: result}},
		*rhsStack...)
}

func R5tPrint(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tProut(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}
