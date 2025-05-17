package runtime

import (
	"fmt"
)

type (
	RopeNodeType      int
	RopeBalanceFactor int
)

const (
	RopeBalanceFactorFibonacci = iota
	RopeBalanceFactorAVL
)

const (
	RopeNodeInnerType = iota
	RopeNodeLeafType
)

const (
	RopeLeafLengthMax = 10
)

type Rope struct {
	root         RopeNode
	fibGenerator func(n int) int
}

func (r *Rope) Len() int {
	if r.root == nil {
		return 0
	}

	return r.root.GetWeight()
}

func (r *Rope) Height() int {
	if r.root == nil {
		return 0
	}

	return r.root.GetHeight()
}

func (r *Rope) traverseLeaves() []*RopeNodeLeaf {
	if r.root == nil {
		return nil
	}

	stack := []RopeNode{r.root}
	result := []*RopeNodeLeaf{}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch curr.NodeType() {
		case RopeNodeLeafType:
			leaf := curr.(*RopeNodeLeaf)
			result = append(result, leaf)
		case RopeNodeInnerType:
			inner := curr.(*RopeNodeInner)

			if inner.Right != nil {
				stack = append(stack, inner.Right)
			}

			if inner.Left != nil {
				stack = append(stack, inner.Left)
			}
		}
	}

	return result
}

func (r *Rope) IsBalanced() bool {
	required := r.fibGenerator(r.root.GetHeight() + 2)
	return r.root.GetWeight() >= required
}

func (r *Rope) IsAVLBalanced() bool {
	factor := r.balanceFactorAVL(r.root)
	return factor < 2 && factor > -2
}

func (r *Rope) findSlot(length int) int {
	for i := 0; ; i++ {
		if length >= r.fibGenerator(i) && length < r.fibGenerator(i+1) {
			return i
		}
	}
}

func (r *Rope) Balance() *Rope {
	if r.IsBalanced() {
		return r
	}

	leaves := r.traverseLeaves()
	return r.mergeLeaves(leaves, 0, len(leaves))
}

func (r *Rope) rotateLeft(node RopeNode) RopeNode {
	if node.NodeType() == RopeNodeLeafType {
		return node
	}

	x := node.(*RopeNodeInner)
	y := x.Right

	if inner, ok := y.(*RopeNodeInner); ok {
		x.Right = inner.Left
		inner.Left = x
	}

	r.updateAVLBalanceInfo(x)
	r.updateAVLBalanceInfo(y)

	return y
}

func (r *Rope) rotateRight(node RopeNode) RopeNode {
	if node.NodeType() == RopeNodeLeafType {
		return node
	}

	y := node.(*RopeNodeInner)
	x := y.Left

	if inner, ok := x.(*RopeNodeInner); ok {
		y.Left = inner.Right
		inner.Right = y
	}

	r.updateAVLBalanceInfo(y)
	r.updateAVLBalanceInfo(x)

	return x
}

func (r *Rope) updateAVLBalanceInfo(x RopeNode) {
	if x.NodeType() == RopeNodeLeafType {
		return
	}

	inner := x.(*RopeNodeInner)

	leftHeight := 0
	rightHeight := 0

	inner.Weight = 0

	if inner.Left != nil {
		leftHeight = inner.Left.GetHeight()
		inner.Weight += inner.Left.GetWeight()
	}

	if inner.Right != nil {
		rightHeight = inner.Right.GetHeight()
		inner.Weight += inner.Right.GetWeight()
	}

	inner.Height = 1 + max(leftHeight, rightHeight)
}

func (r *Rope) balanceFactorAVL(node RopeNode) int {
	if node.NodeType() == RopeNodeLeafType {
		return 0
	}

	root := node.(*RopeNodeInner)

	leftHeight := 0
	rightHeight := 0

	if root.Left != nil {
		leftHeight = root.Left.GetHeight()
	}

	if root.Right != nil {
		rightHeight = root.Right.GetHeight()
	}

	return rightHeight - leftHeight
}

func (r *Rope) balanceAVL() *Rope {
	r.updateAVLBalanceInfo(r.root)

	balanceFactor := r.balanceFactorAVL(r.root)

	if balanceFactor < 2 && balanceFactor > -2 {
		return r
	}

	root := *r.root.(*RopeNodeInner)

	if balanceFactor == 2 {

		if right, ok := root.Right.(*RopeNodeInner); ok {
			tmp := *right
			if r.balanceFactorAVL(&tmp) < 0 {
				root.Right = r.rotateRight(&tmp)
			}
		}
		return &Rope{
			root:         r.rotateLeft(&root),
			fibGenerator: r.fibGenerator,
		}

	}

	if balanceFactor == -2 {
		leftTmp := root.Left
		if left, ok := root.Left.(*RopeNodeInner); ok {
			tmp := *left
			if r.balanceFactorAVL(&tmp) > 0 {
				root.Left = r.rotateLeft(leftTmp)
			}
		}

		root := r.rotateRight(&root)
		return &Rope{
			root:         root,
			fibGenerator: r.fibGenerator,
		}

	}

	return r
}

// function balance(node: RopeNode) -> RopeNode:
//     update(node)
//     if balanceFactor(node) == 2:
//         if balanceFactor(node.right) < 0:
//             node.right = rotateRight(node.right)
//         return rotateLeft(node)
//     if balanceFactor(node) == -2:
//         if balanceFactor(node.left) > 0:
//             node.left = rotateLeft(node.left)
//         return rotateRight(node)
//     return node

// function update(node: RopeNode):
//     node.height = 1 + max(height(node.left), height(node.right))
//     node.weight = length(node.left)
//
// function balanceFactor(node: RopeNode) -> int:
//     return height(node.right) - height(node.left)
//
// function rotateLeft(x: RopeNode) -> RopeNode:
//     y = x.right
//     x.right = y.left
//     y.left = x
//     update(x)
//     update(y)
//     return y
//
// function rotateRight(y: RopeNode) -> RopeNode:
//     x = y.left
//     y.left = x.right
//     x.right = y
//     update(y)
//     update(x)
//     return x
//
// function balance(node: RopeNode) -> RopeNode:
//     update(node)
//     if balanceFactor(node) == 2:
//         if balanceFactor(node.right) < 0:
//             node.right = rotateRight(node.right)
//         return rotateLeft(node)
//     if balanceFactor(node) == -2:
//         if balanceFactor(node.left) > 0:
//             node.left = rotateLeft(node.left)
//         return rotateRight(node)
//     return node

func (r *Rope) mergeLeaves(leaves []*RopeNodeLeaf, start, end int) *Rope {
	rng := end - start

	if rng == 1 {
		return NewRope(leaves[start].Data)
	}

	if rng == 2 {
		return NewRope(leaves[start].Data).ConcatWithRebalance(NewRope(leaves[start+1].Data))
	}

	mid := start + (rng / 2)

	return r.mergeLeaves(leaves, start, mid).ConcatWithRebalance(r.mergeLeaves(leaves, mid, end))
}

type RopeNode interface {
	NodeType() RopeNodeType
	GetWeight() int
	GetHeight() int
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

func (n *RopeNodeLeaf) GetHeight() int {
	return 0
}

type RopeNodeInner struct {
	Weight int
	Height int
	Left   RopeNode
	Right  RopeNode
}

func (n *RopeNodeInner) NodeType() RopeNodeType {
	return RopeNodeInnerType
}

func (n *RopeNodeInner) GetWeight() int {
	return n.Weight
}

func (n *RopeNodeInner) GetHeight() int {
	return n.Height
}

func NewRope(n []R5Node) *Rope {
	return &Rope{
		fibGenerator: fibonacci(),
		root: &RopeNodeLeaf{
			Data: n,
		},
	}
}

func (r *Rope) ConcatWithRebalance(other *Rope) *Rope {
	if other.Len() == 0 {
		return r
	}

	if r.Len() == 0 {
		return other
	}

	res := &Rope{
		fibGenerator: r.fibGenerator,
		root: &RopeNodeInner{
			Weight: r.root.GetWeight() + other.root.GetWeight(),
			Left:   r.root,
			Right:  other.root,
			Height: max(r.root.GetHeight(), other.root.GetHeight()) + 1,
		},
	}

	if res.IsBalanced() {
		return res
	}

	return res.Balance()
}

func (r *Rope) Concat(other *Rope) *Rope {
	res := &Rope{
		fibGenerator: r.fibGenerator,
		root: &RopeNodeInner{
			Weight: r.root.GetWeight() + other.root.GetWeight(),
			Left:   r.root,
			Right:  other.root,
			Height: max(r.root.GetHeight(), other.root.GetHeight()) + 1,
		},
	}

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
	left, right := r.splitRec(r.root, i)
	return &Rope{
			root:         left,
			fibGenerator: r.fibGenerator,
		}, &Rope{
			root:         right,
			fibGenerator: r.fibGenerator,
		}
}

func (r *Rope) splitRec(n RopeNode, i int) (RopeNode, RopeNode) {
	if n == nil {
		return nil, nil
	}

	if n.NodeType() == RopeNodeLeafType {

		leaf := n.(*RopeNodeLeaf)
		
		if i <= 0 {
			return &RopeNodeLeaf{} , &RopeNodeLeaf{leaf.Data}
		}

		if i >= leaf.GetWeight() {
			return &RopeNodeLeaf{leaf.Data}, &RopeNodeLeaf{}
		}

		left := &RopeNodeLeaf{leaf.Data[i:]}
		right := &RopeNodeLeaf{leaf.Data[:i]}
		return left, right

	}

	inner := n.(*RopeNodeInner)

	if inner.Left != nil && inner.Left.GetWeight() > i {
		l1, l2 := r.splitRec(inner.Left, i)
		right := &RopeNodeInner{
			Weight: l2.GetWeight(),
			Height: l2.GetHeight() + 1,
			Left:   l2,
			Right:  inner.Right,
		}

		if inner.Right != nil {
			right.Weight += inner.Right.GetWeight()
			right.Height = max(inner.Right.GetHeight(), l2.GetHeight()) + 1
		}

		return l1, right
	}

	if inner.Left != nil {
		i -= inner.Left.GetWeight()
	}

	r1, r2 := r.splitRec(inner.Right, i)
	left := &RopeNodeInner{
		Weight: r1.GetWeight(),
		Height: r1.GetHeight() + 1,
		Left:   inner.Left,
		Right:  r1,
	}

	if inner.Left != nil {
		left.Weight += inner.Left.GetWeight()
		left.Height = max(inner.Left.GetHeight(), r1.GetHeight()) + 1
	}

	return left, r2
}

// func Split(n *RopeNode, i int) (*RopeNode, *RopeNode) {
//     if n == nil {
//         return nil, nil
//     }
//
//     if isLeaf(n) {
//         if i <= 0 {
//             return nil, n
//         } else if i >= len(n.value) {
//             return n, nil
//         } else {
//             left := &RopeNode{value: n.value[:i]}
//             right := &RopeNode{value: n.value[i:]}
//             return left, right
//         }
//     }
//
//     if i < n.weight {
//         // Идём влево
//         l1, l2 := Split(n.left, i)
//         right := makeNode(l2, n.right)
//         return l1, right
//     } else {
//         // Идём вправо
//         r1, r2 := Split(n.right, i-n.weight)
//         left := makeNode(n.left, r1)
//         return left, r2
//     }
// }

func (r *Rope) Insert(i int, data []R5Node) *Rope {
	if i < 0 || i > r.Len() {
		return nil
	}

	tmp := NewRope(data)

	if i == 0 {
		return tmp.ConcatWithRebalance(r)
	}

	if i == r.Len() {
		result := r.ConcatWithRebalance(tmp)
		return result
	}

	tmpLhs, tmpRhs := r.Split(i)
	tmp = tmpLhs.ConcatWithRebalance(tmp)
	return tmp.ConcatWithRebalance(tmpRhs)
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

	tmp := tmpLhs.ConcatWithRebalance(tmpRhs)
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
