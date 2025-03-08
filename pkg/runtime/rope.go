package runtime

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
				i -= inner.Left.GetWeight()
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
