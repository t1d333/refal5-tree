package runtime

import "fmt"

type ViewFieldNodeType int

const (
	OpenCallType ViewFieldNodeType = iota
	CloseCallType
	RopeType
)

type OpenCallViewFieldNode struct {
	Function R5Function
}

func (*OpenCallViewFieldNode) Type() ViewFieldNodeType {
	return OpenCallType
}

type CloseCallViewFieldNode struct{}

func (*CloseCallViewFieldNode) Type() ViewFieldNodeType {
	return CloseCallType
}

type RopeViewFieldNode struct {
	Value *Rope
}

func (*RopeViewFieldNode) Type() ViewFieldNodeType {
	return RopeType
}

type ViewFieldNode interface {
	Type() ViewFieldNodeType
}

func PrintViewField(viewField []ViewFieldNode) {
	for _, n := range viewField {
		switch n.Type() {
		case OpenCallType:
			openCall := n.(*OpenCallViewFieldNode)
			fmt.Printf("< %s", openCall.Function.Name)
		case CloseCallType:
			fmt.Printf("> ")
		case RopeType:
			rope := n.(*RopeViewFieldNode)
			fmt.Printf(rope.Value.String())
		}
	}
	fmt.Println()
}
