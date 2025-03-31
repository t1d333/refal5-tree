package runtime

import "fmt"

type RopeNodeType int

const (
	RopeNodeInnerType = iota
	RopeNodeLeafType
)

type Rope struct {
	root RopeNode
}

func (r *Rope) Len() int {
	if r.root == nil {
		return 0
	}

	return r.root.GetWeight()
}

type RopeNode interface {
	NodeType() RopeNodeType
	GetWeight() int
}

type RopeNodeLeaf struct {
	Data []R5Node
}

func (n *RopeNodeLeaf) NodeType() RopeNodeType {
	return RopeNodeLeafType
}

func (n *RopeNodeLeaf) GetWeight() int {
	return len(n.Data)
}

type RopeNodeInner struct {
	Weight int
	Left   RopeNode
	Right  RopeNode
}

func (n *RopeNodeInner) NodeType() RopeNodeType {
	return RopeNodeInnerType
}

func (n *RopeNodeInner) GetWeight() int {
	return n.Weight
}

func NewRope(n []R5Node) *Rope {
	return &Rope{
		root: &RopeNodeLeaf{
			Data: n,
		},
	}
}

func (r *Rope) Concat(other *Rope) *Rope {
	newRoot := &RopeNodeInner{
		Weight: r.root.GetWeight() + other.root.GetWeight(),
		Left:   r.root,
		Right:  other.root,
	}
	newRoot.Left = r.root
	newRoot.Right = other.root
	res := Rope{root: newRoot}
	return &res
}

func (r *Rope) Get(i int) R5Node {
	curr := r.root

	for {
		if curr == nil {
			return nil
		}

		switch curr.NodeType() {
		case RopeNodeInnerType:
			inner := curr.(*RopeNodeInner)
			if inner.Left != nil && inner.Left.GetWeight() > i {
				curr = inner.Left
			} else {
				curr = inner.Right
				if inner.Left != nil {
					i -= inner.Left.GetWeight()
				}
			}
		case RopeNodeLeafType:
			leaf := curr.(*RopeNodeLeaf)
			if len(leaf.Data) > i {
				return leaf.Data[i]
			} else {
				return nil
			}
		}
	}
}

func (r *Rope) Set(i int, data R5Node) {
	curr := r.root

	for {
		if curr == nil {
			return
		}

		switch curr.NodeType() {
		case RopeNodeInnerType:
			inner := curr.(*RopeNodeInner)
			if inner.Left != nil && inner.Left.GetWeight() > i {
				curr = inner.Left
			} else {
				curr = inner.Right
				i -= inner.Left.GetWeight()
			}
		case RopeNodeLeafType:
			leaf := curr.(*RopeNodeLeaf)
			if len(leaf.Data) > i {
				leaf.Data[i] = data
			} else {
				return
			}
		}
	}
}

func (r *Rope) Split(i int) (*Rope, *Rope) {
	if i < 0 || i > r.Len() {
		return nil, nil
	}

	path := []RopeNode{}
	curr := r.root

	for {
		switch curr.NodeType() {
		case RopeNodeInnerType:
			path = append(path, curr)
			inner := curr.(*RopeNodeInner)
			if inner.Left != nil && inner.Left.GetWeight() > i {
				curr = inner.Left
			} else {
				curr = inner.Right
				if inner.Left != nil {
					i -= inner.Left.GetWeight()
				}
			}
		case RopeNodeLeafType:
			leaf := curr.(*RopeNodeLeaf)
			var prevLhsNode RopeNode = &RopeNodeLeaf{
				Data: leaf.Data[:i],
			}
			var prevRhsNode RopeNode = &RopeNodeLeaf{
				Data: leaf.Data[i:],
			}

			for i := len(path) - 1; i >= 0; i -= 1 {
				currParent := path[i].(*RopeNodeInner)

				newLhs := &RopeNodeInner{
					Weight: prevLhsNode.GetWeight(),
					Left:   nil,
					Right:  prevLhsNode,
				}

				if currParent.Left != nil && currParent.Left != curr {
					newLhs.Weight += currParent.Left.GetWeight()
					newLhs.Left = currParent.Left
				}

				newRhs := &RopeNodeInner{
					Weight: prevRhsNode.GetWeight(),
					Left:   prevRhsNode,
					Right:  nil,
				}

				if currParent.Right != nil && currParent.Right != curr {
					newRhs.Right = currParent.Right
					newRhs.Weight += currParent.Right.GetWeight()
				}

				prevLhsNode = newLhs
				prevRhsNode = newRhs
			}

			return &Rope{
				root: prevLhsNode,
			}, &Rope{root: prevRhsNode}
		}
	}
}

func (r *Rope) Insert(i int, data []R5Node) {
	if i < 0 || i > r.Len() {
		return
	}

	tmp := NewRope(data)

	if i == 0 {
		tmp = tmp.Concat(r)
		r.root = tmp.root
		return
	}

	if i == r.Len() {
		tmp = r.Concat(tmp)
		r.root = tmp.root
		return
	}

	tmpLhs, tmpRhs := r.Split(i)
	tmp = tmpLhs.Concat(tmp)
	tmp = tmp.Concat(tmpRhs)
	r.root = tmp.root
}

func (r *Rope) Delete(i int) {
	if i < 0 || i > r.Len() {
		return
	}

	tmpLhs, _ := r.Split(i)

	if i == r.Len()-1 {
		r.root = tmpLhs.root
		return
	}

	_, tmpRhs := r.Split(i + 1)

	tmp := tmpLhs.Concat(tmpRhs)
	r.root = tmp.root
}

func (r *Rope) String() string {
	result := ""
	for i := 0; i < r.Len(); i++ {
		node := r.Get(i)
		switch node.Type() {
		case R5DatatagChar:
			charNode := node.(*R5NodeChar)
			result += fmt.Sprintf("(Char: %c) ", charNode.Char)
		case R5DatatagCloseBracket:
			closeBrNode := node.(*R5NodeCloseBracket)
			result += fmt.Sprintf("(CloseBracket, OpenOffset: %d) ", closeBrNode.OpenOffset)
		case R5DatatagCloseCall:
			closeCallNode := node.(*R5NodeCloseCall)
			result += fmt.Sprintf("(CloseCall, OpenOffset: %d) ", closeCallNode.OpenOffset)
		case R5DatatagFunction:
			funcNode := node.(*R5NodeFunction)
			result += fmt.Sprintf("(Function: %s) ", funcNode.Function.Name)
		case R5DatatagIllegal:
			result += fmt.Sprintf("(Illegal) ")
		case R5DatatagNumber:
			numberNode := node.(*R5NodeNumber)
			result += fmt.Sprintf("(Number: %d) ", numberNode.Number)
		case R5DatatagString:
			strNode := node.(*R5NodeString)
			result += fmt.Sprintf("(String: %s) ", strNode.String)
		case R5DatatagOpenBracket:
			openBrNode := node.(*R5NodeOpenBracket)
			result += fmt.Sprintf("(OpenBracket: CloseOffset: %d) ", openBrNode.CloseOffset)
		case R5DatatagOpenCall:
			openCallNode := node.(*R5NodeOpenCall)
			result += fmt.Sprintf("(OpenCall: CloseOffset: %d) ", openCallNode.CloseOffset)
		}
	}
	return result
}
