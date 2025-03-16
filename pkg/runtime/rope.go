package runtime

import "fmt"

// import (
// 	"fmt"
// )

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

func NewRope(n []R5Node) Rope {
	return Rope{
		root: &RopeNodeLeaf{
			Data: n,
		},
	}
}

func (r *Rope) Concat(other Rope) Rope {
	newRoot := &RopeNodeInner{
		Weight: r.root.GetWeight() + other.root.GetWeight(),
		Left:   r.root,
		Right:  other.root,
	}
	newRoot.Left = r.root
	newRoot.Right = other.root
	res := Rope{root: newRoot}
	return res
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

	fmt.Println("-------", i)

	for {
		switch curr.NodeType() {
		case RopeNodeInnerType:
			path = append(path, curr)
			inner := curr.(*RopeNodeInner)
			if inner.Left != nil && inner.Left.GetWeight() > i {
				curr = inner.Left
			} else {
				curr = inner.Right
				i -= inner.Left.GetWeight()
			}
		case RopeNodeLeafType:
			leaf := curr.(*RopeNodeLeaf)
			fmt.Println(leaf.Data[:i], leaf.Data[i:])
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
		tmp = tmp.Concat(*r)
		r.root = tmp.root
		fmt.Println("Root in split", r.Len())
		return
	}

	if i == r.Len()-1 {
		tmp = r.Concat(tmp)
		r.root = tmp.root
		return
	}

	tmpLhs, tmpRhs := r.Split(i)
	tmp = tmpLhs.Concat(tmp)
	tmp = tmp.Concat(*tmpRhs)
	r.root = tmp.root
}

func (r *Rope) Delete(i int) {
	if i < 0 || i > r.Len() {
		return
	}

	fmt.Println("++++++++++++", i)
	tmpLhs, _ := r.Split(i)

	if i == r.Len()-1 {
		r.root = tmpLhs.root
		return
	}

	_, tmpRhs := r.Split(i + 1)
	fmt.Println("Left: ", tmpLhs.Len(), "Right: ", tmpRhs.Len())

	tmp := tmpLhs.Concat(*tmpRhs)
	r.root = tmp.root
}
