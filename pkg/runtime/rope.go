package runtime

import (
	"fmt"
)
type (
	RopeNodeType      int
	RopeBalanceFactor int
)

const (
	MaxLeafLength = 50
)

const (
	RopeNodeInnerType = iota
	RopeNodeLeafType
)

func VisualizeRope(r Rope, level int) {
	node := r.root
	if node == nil {
		return
	}

	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}

	switch n := node.(type) {
	case RopeNodeInner:
		fmt.Printf("%s[Inner W:%d H:%d]\n", indent, n.weight, n.height)
		if n.Left != nil {
			VisualizeRopeTree(n.Left, level+1)
		}
		if n.Right != nil {
			VisualizeRopeTree(n.Right, level+1)
		}
	case RopeNodeLeaf:
		fmt.Printf("%s[Leaf Len:%d]\n", indent, len(n.Data))
	default:
		fmt.Printf("%s[Unknown Node]\n", indent)
	}
}

func VisualizeRopeTree(node RopeNode, level int) {
	if node == nil {
		return
	}

	indent := ""
	for i := 0; i < level; i++ {
		indent += "  " // 2 пробела на каждый уровень
	}

	switch n := node.(type) {
	case RopeNodeInner:
		fmt.Printf("%s[Inner W:%d H:%d]\n", indent, n.weight, n.height)
		if n.Left != nil {
			VisualizeRopeTree(n.Left, level+1)
		}
		if n.Right != nil {
			VisualizeRopeTree(n.Right, level+1)
		}
	case RopeNodeLeaf:
		fmt.Printf("%s[Leaf Len:%d]\n", indent, len(n.Data))
	default:
		fmt.Printf("%s[Unknown Node]\n", indent)
	}
}

type Rope struct {
	root         RopeNode
}

func (r Rope) Len() int {
	if r.root == nil {
		return 0
	}

	return r.root.Weight()
}

func (r Rope) Height() int {
	if r.root == nil {
		return 0
	}

	return r.root.Height()
}

func (r Rope) traverseLeaves() []RopeNodeLeaf {
	if r.root == nil {
		return nil
	}

	stack := []RopeNode{r.root}
	result := []RopeNodeLeaf{}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch curr.NodeType() {
		case RopeNodeLeafType:
			leaf := curr.(RopeNodeLeaf)
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

func (r *Rope) IsAVLBalanced() bool {
	factor := r.balanceFactorAVL(r.root)
	return factor < 2 && factor > -2
}


func (r *Rope) rotateLeft(node RopeNode) RopeNode {
	if node.IsLeaf() {
		return node
	}

	x := node.(RopeNodeInner)
	y := x.Right

	innerY, ok := y.(RopeNodeInner)

	if !ok {
		panic(123)
	}
	// if inner, ok := y.(RopeNodeInner); ok {
	// 	x.Right = inner.Left
	// 	inner.Left = x
	// }
	//
	// x = r.updateAVLBalanceInfo(x).(RopeNodeInner)
	// y = r.updateAVLBalanceInfo(y)
	newX := NewInner(x.Left, innerY.Left)

	// updatedX := r.updateAVLBalanceInfo(newX)

	// Создаём новый Y, в котором левый ребёнок — это обновлённый X
	newY := NewInner(newX, innerY.Right)

	return newY
}

func (r Rope) rotateRight(node RopeNode) RopeNode {
	if node.NodeType() == RopeNodeLeafType {
		return node
	}

	y := node.(RopeNodeInner)
	x := y.Left

	innerX, ok := x.(RopeNodeInner)
	if !ok {
		panic(123)
	}

	newY := NewInner(innerX.Right, y.Right)

	newX := NewInner(innerX.Left, newY)

	return newX
}

func (r Rope) updateAVLBalanceInfo(x RopeNode) RopeNode {
	if x.NodeType() == RopeNodeLeafType {
		return x
	}

	inner := x.(RopeNodeInner)

	leftHeight := 0
	rightHeight := 0

	inner.weight = 0

	if inner.Left != nil {
		leftHeight = inner.Left.Height()
		inner.weight += inner.Left.Weight()
	}

	if inner.Right != nil {
		rightHeight = inner.Right.Height()
		inner.weight += inner.Right.Weight()
	}

	inner.height = 1 + max(leftHeight, rightHeight)
	return inner
}

func (r *Rope) BalanceFactor() int {
	return r.balanceFactorAVL(r.root)
}

func (r *Rope) balanceFactorAVL(node RopeNode) int {
	if node.NodeType() == RopeNodeLeafType {
		return 0
	}

	root := node.(RopeNodeInner)

	leftHeight := 0
	rightHeight := 0

	if root.Left != nil {
		leftHeight = root.Left.Height()
	}

	if root.Right != nil {
		rightHeight = root.Right.Height()
	}

	return rightHeight - leftHeight
}

func (r Rope) balanceAVL() Rope {
	r.root = r.updateAVLBalanceInfo(r.root)

	balanceFactor := r.balanceFactorAVL(r.root)

	if -2 < balanceFactor && balanceFactor < 2 {
		return r
	}

	root := r.root.(RopeNodeInner)

	if balanceFactor == 2 {

		if right, ok := root.Right.(RopeNodeInner); ok {
			tmp := right
			if r.balanceFactorAVL(tmp) < 0 {
				root.Right = r.rotateRight(tmp)
			}
		}

		return Rope{
			root:         r.rotateLeft(root),
		}

	}

	if balanceFactor == -2 {
		leftTmp := root.Left
		if left, ok := root.Left.(RopeNodeInner); ok {
			tmp := left
			if r.balanceFactorAVL(tmp) > 0 {
				root.Left = r.rotateLeft(leftTmp)
			}
		}

		root := r.rotateRight(root)
		return Rope{
			root:         root,
		}

	}

	return r
}

func (r *Rope) balanceNodeAVL(node RopeNode) RopeNode {
	if node.NodeType() == RopeNodeLeafType {
		return node
	}

	node = r.updateAVLBalanceInfo(node)

	balanceFactor := r.balanceFactorAVL(r.root)

	if -2 < balanceFactor && balanceFactor < 2 {
		return node
	}

	root := node.(RopeNodeInner)

	if balanceFactor == 2 {

		if right, ok := root.Right.(*RopeNodeInner); ok {
			tmp := *right
			if r.balanceFactorAVL(tmp) < 0 {
				root.Right = r.rotateRight(tmp)
			}
		}
		return root

	}

	if balanceFactor == -2 {
		leftTmp := root.Left
		if left, ok := root.Left.(RopeNodeInner); ok {
			tmp := left
			if r.balanceFactorAVL(tmp) > 0 {
				root.Left = r.rotateLeft(leftTmp)
			}
		}

		root := r.rotateRight(root)
		return root

	}
	return node
}

type RopeNode interface {
	NodeType() RopeNodeType
	Weight() int
	Height() int
	IsLeaf() bool
}

type RopeNodeLeaf struct {
	Data   []R5Node
	weight int
}

func NewLeaf(data []R5Node) RopeNodeLeaf {
	return RopeNodeLeaf{
		weight: len(data),
		Data:   data,
	}
}

func NewInner(left, right RopeNode) RopeNodeInner {
	inner := RopeNodeInner{
		weight: left.Weight() + right.Weight(),
		height: max(left.Height(), right.Height()) + 1,
		Left:   left,
		Right:  right,
	}
	return inner
}

func (n RopeNodeLeaf) NodeType() RopeNodeType {
	return RopeNodeLeafType
}

func (n RopeNodeLeaf) Weight() int {
	return n.weight
}

func (n RopeNodeLeaf) Height() int {
	return 1
}

func (n RopeNodeLeaf) IsLeaf() bool {
	return true
}

type RopeNodeInner struct {
	weight int
	height int
	Left   RopeNode
	Right  RopeNode
}

func (n RopeNodeInner) NodeType() RopeNodeType {
	return RopeNodeInnerType
}

func (n RopeNodeInner) Weight() int {
	return n.weight
}

func (n RopeNodeInner) Height() int {
	return n.height
}

func (n RopeNodeInner) IsLeaf() bool {
	return false
}

func NewRope(n []R5Node) Rope {
	if len(n) <= MaxLeafLength {
		return Rope{
			root:         NewLeaf(n),
		}
	}

	left := Rope{
		root:         NewLeaf(n[:MaxLeafLength]),
	}

	return left.ConcatAVL(NewRope(n[MaxLeafLength:]))
}


func (r *Rope) Get(i int) R5Node {
	curr := r.root

	for {
		if curr == nil {
			return nil
		}

		switch curr.NodeType() {
		case RopeNodeInnerType:
			inner := curr.(RopeNodeInner)
			if inner.Left != nil && inner.Left.Weight() > i {
				curr = inner.Left
			} else {
				curr = inner.Right
				if inner.Left != nil {
					i -= inner.Left.Weight()
				}
			}
		case RopeNodeLeafType:
			leaf := curr.(RopeNodeLeaf)
			if len(leaf.Data) > i {
				return leaf.Data[i]
			} else {
				return nil
			}
		}
	}
}

func (r Rope) Split(i int) (Rope, Rope) {
	left, right := r.splitRec2(r.root, i)
	return Rope{
			root:         left,
		}, Rope{
			root:         right,
		}
}

func (r Rope) splitRec2(node RopeNode, idx int) (RopeNode, RopeNode) {
	if node.IsLeaf() {
		leaf := node.(RopeNodeLeaf)
		if idx >= leaf.Weight() {
			return node, nil
		}

		if idx <= 0 {
			return nil, node
		}

		leftLeaf := NewLeaf(leaf.Data[:idx])
		rightLeaf := NewLeaf(leaf.Data[idx:])
		return leftLeaf, rightLeaf
	}

	inner := node.(RopeNodeInner)
	leftLen := inner.Left.Weight()

	if idx < leftLen {
		leftSplit, rightSplit := r.splitRec2(inner.Left, idx)
		newRight := Rope{root: rightSplit}.ConcatAVL(Rope{root: inner.Right})
		return leftSplit, newRight.root
	} else if idx > leftLen {
		leftSplit, rightSplit := r.splitRec2(inner.Right, idx-leftLen)
		newLeft := Rope{root: inner.Left}.ConcatAVL(Rope{root: leftSplit})
		return newLeft.root, rightSplit

	} else {
		return inner.Left, inner.Right
	}
}

func (r Rope) Insert(i int, data []R5Node) Rope {
	if len(data) == 0 {
		return r
	}

	if r.Len() == 0 {
		return NewRope(data)
	}

	if i < 0 || i > r.Len() {
		return NewRope([]R5Node{})
	}

	tmp := NewRope(data)

	if i == 0 {
		return tmp.ConcatAVL(r)
	}

	if i == r.Len() {
		result := r.ConcatAVL(tmp)
		return result
	}

	tmpLhs, tmpRhs := r.Split(i)
	tmp = tmpLhs.ConcatAVL(tmp)
	
	return tmp.ConcatAVL(tmpRhs)
}

func (r Rope) String() string {
	result := ""
	for i := 0; i < r.Len(); i++ {
		node := r.Get(i)
		switch node.Type() {
		case R5DatatagChar:
			charNode := node.(*R5NodeChar)
			result += fmt.Sprintf("%c ", charNode.Char)
		case R5DatatagCloseBracket:
			result += fmt.Sprintf(")")
		case R5DatatagCloseCall:
			result += fmt.Sprintf(">")
		case R5DatatagFunction:
			funcNode := node.(*R5NodeFunction)
			result += fmt.Sprintf("%s ", funcNode.Function.Name)
		case R5DatatagIllegal:
			result += fmt.Sprintf("(Illegal) ")
		case R5DatatagNumber:
			numberNode := node.(*R5NodeNumber)
			result += fmt.Sprintf("%d ", numberNode.Number)
		case R5DatatagString:
			strNode := node.(*R5NodeString)
			result += fmt.Sprintf("%s ", strNode.String())
		case R5DatatagOpenBracket:
			result += fmt.Sprintf("(")
		case R5DatatagOpenCall:
			result += fmt.Sprintf("< ")
		}
	}
	return result
}

func (lhs Rope) ConcatAVL(rhs Rope) Rope {
	if lhs.Len() == 0 {
		return rhs
	}

	if rhs.Len() == 0 {
		return lhs
	}

	if lhs.Height() == 1 && rhs.Height() == 1 {
		lhsLeaf := lhs.root.(RopeNodeLeaf)
		rhsLeaf := rhs.root.(RopeNodeLeaf)

		return Rope{root: lhs.concatLeaves(lhsLeaf, rhsLeaf)}
	}

	diff := rhs.Height() - lhs.Height()
	if -2 < diff && diff < 2 {
		return Rope{
			root: NewInner(lhs.root, rhs.root),
		}
	}

	if diff == 2 || diff == -2 {
		result := Rope{
			root: NewInner(lhs.root, rhs.root),
		}

		return result.balanceAVL()
	}

	if lhs.Height() > rhs.Height() {
		lhsRoot := lhs.root.(RopeNodeInner)
		lhsL := Rope{root: lhsRoot.Left}
		lhsR := Rope{root: lhsRoot.Right}

		return lhsL.ConcatAVL(lhsR.ConcatAVL(rhs))

	} else {
		rhsRoot := rhs.root.(RopeNodeInner)
		lhsL := Rope{root: rhsRoot.Left}
		lhsR := Rope{root: rhsRoot.Right}

		return lhs.ConcatAVL(lhsL).ConcatAVL(lhsR)
	}
}

func (r Rope) concatLeaves(lhs, rhs RopeNodeLeaf) RopeNode {
	buff := make([]R5Node, lhs.weight+rhs.weight)
	copy(buff, lhs.Data)
	copy(buff[len(lhs.Data):], rhs.Data)

	if len(buff) <= MaxLeafLength {
		return NewLeaf(buff)
	}

	return RopeNodeInner{
		weight: lhs.Weight() + rhs.Weight(),
		height: 2,
		Left:   NewLeaf(buff[:MaxLeafLength]),
		Right:  NewLeaf(buff[MaxLeafLength:]),
	}
}
