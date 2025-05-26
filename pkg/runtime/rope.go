package runtime

import (
	"fmt"
)

type (
	RopeNodeType      int
	RopeBalanceFactor int
)

const (
	MaxLeafLength = 70
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
	fibGenerator func(n int) int
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

func (r *Rope) IsBalanced() bool {
	required := r.fibGenerator(r.root.Height() + 2)
	return r.root.Weight() >= required
}

func (r *Rope) IsAVLBalanced() bool {
	factor := r.balanceFactorAVL(r.root)
	return factor < 2 && factor > -2
}

func (r Rope) Balance() Rope {
	if r.IsBalanced() {
		return r
	}

	leaves := r.traverseLeaves()
	return r.mergeLeaves(leaves, 0, len(leaves))
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
			fibGenerator: r.fibGenerator,
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
			fibGenerator: r.fibGenerator,
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

func (r Rope) mergeLeaves(leaves []RopeNodeLeaf, start, end int) Rope {
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
			fibGenerator: fibonacci(),
			root:         NewLeaf(n),
		}
	}

	left := Rope{
		fibGenerator: fibonacci(),
		root:         NewLeaf(n[:MaxLeafLength]),
	}

	return left.ConcatAVL(NewRope(n[MaxLeafLength:]))
}

func (r Rope) ConcatWithRebalance(other Rope) Rope {
	if other.Len() == 0 {
		return r
	}

	if r.Len() == 0 {
		return other
	}

	res := Rope{
		fibGenerator: r.fibGenerator,
		root: RopeNodeInner{
			weight: r.root.Weight() + other.root.Weight(),
			Left:   r.root,
			Right:  other.root,
			height: max(r.root.Height(), other.root.Height()) + 1,
		},
	}

	if res.IsBalanced() {
		return res
	}

	return res.Balance()
}

func (r *Rope) Concat(other *Rope) *Rope {
	if r.root.NodeType() == RopeNodeLeafType && other.root.NodeType() == RopeNodeLeafType {

		rLeaf := r.root.(RopeNodeLeaf)
		otherLeaf := other.root.(RopeNodeLeaf)

		if r.Len()+other.Len() <= MaxLeafLength {
			return &Rope{
				fibGenerator: r.fibGenerator,
				root:         NewLeaf(append(rLeaf.Data, otherLeaf.Data...)),
			}
		}

		tmp := append(rLeaf.Data, otherLeaf.Data...)

		return &Rope{
			fibGenerator: r.fibGenerator,
			root: RopeNodeInner{
				weight: len(tmp),
				height: 1,
				Left:   NewLeaf(tmp[:MaxLeafLength]),
				Right:  NewLeaf(tmp[MaxLeafLength:]),
			},
		}
	}

	res := &Rope{
		fibGenerator: r.fibGenerator,
		root: RopeNodeInner{
			weight: r.root.Weight() + other.root.Weight(),
			Left:   r.root,
			Right:  other.root,
			height: max(r.root.Height(), other.root.Height()) + 1,
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

func (r *Rope) Set(i int, data R5Node) {
	curr := r.root

	for {
		if curr == nil {
			return
		}

		switch curr.NodeType() {
		case RopeNodeInnerType:
			inner := curr.(*RopeNodeInner)
			if inner.Left != nil && inner.Left.Weight() > i {
				curr = inner.Left
			} else {
				curr = inner.Right
				i -= inner.Left.Weight()
			}
		case RopeNodeLeafType:
			leaf := curr.(RopeNodeLeaf)
			if len(leaf.Data) > i {
				leaf.Data[i] = data
			} else {
				return
			}
		}
	}
}

func (r Rope) Split(i int) (Rope, Rope) {
	left, right := r.splitRec2(r.root, i)
	return Rope{
			root:         left,
			fibGenerator: r.fibGenerator,
		}, Rope{
			root:         right,
			fibGenerator: r.fibGenerator,
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

// func (r *Rope) split(node RopeNode, index int) (RopeNode, RopeNode) {
// 	if node.IsLeaf() {
// 		leaf := node.(RopeNodeLeaf)
// 		if index >= len(leaf.Text) {
// 			return node, nil
// 		}
// 		if index <= 0 {
// 			return nil, node
// 		}
//
// 		leftLeaf := RopeNodeLeaf{Text: leaf.Text[:index]}
// 		rightLeaf := RopeNodeLeaf{Text: leaf.Text[index:]}
// 		return leftLeaf, rightLeaf
// 	}
//
// 	inner := node.(RopeNodeInner)
// 	leftLen := r.length(inner.Left)
//
// 	if index < leftLen {
// 		// Весь split уходит в левое поддерево
// 		leftSplit, rightSplit := r.split(inner.Left, index)
// 		newRight := r.concat(rightSplit, inner.Right)
// 		return leftSplit, newRight
// 	} else if index > leftLen {
// 		// Split делится между левым и правым поддеревом
// 		leftSplit, rightSplit := r.split(inner.Right, index-leftLen)
// 		newLeft := r.concat(inner.Left, leftSplit)
// 		return newLeft, rightSplit
// 	} else {
// 		// Точно на границе
// 		return inner.Left, inner.Right
// 	}
// }

func (r *Rope) splitRec(n RopeNode, i int) (RopeNode, RopeNode) {
	if n == nil {
		return nil, nil
	}

	if n.NodeType() == RopeNodeLeafType {

		leaf := n.(RopeNodeLeaf)
		// buff := make([]R5Node, len(leaf.Data))
		// copy(buff, leaf.Data)

		if i <= 0 {
			return NewLeaf([]R5Node{}), NewLeaf(leaf.Data)
		}

		if i >= leaf.Weight() {
			return NewLeaf(leaf.Data), NewLeaf([]R5Node{})
		}

		left := NewLeaf(leaf.Data[:i])
		right := NewLeaf(leaf.Data[i:])
		return left, right

	}

	inner := n.(*RopeNodeInner)

	if inner.Left != nil && inner.Left.Weight() > i {
		l1, l2 := r.splitRec(inner.Left, i)
		right := NewInner(l2, inner.Right)

		return r.balanceNodeAVL(l1), r.balanceNodeAVL(right)
	}

	if inner.Left != nil {
		i -= inner.Left.Weight()
	}

	r1, r2 := r.splitRec(inner.Right, i)
	left := RopeNodeInner{
		Left:  inner.Left,
		Right: r1,
	}

	left = r.updateAVLBalanceInfo(left).(RopeNodeInner)

	return r.balanceNodeAVL(left), r.balanceNodeAVL(r2)
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
		// return tmp.ConcatWithRebalance(r)
	}

	if i == r.Len() {
		result := r.ConcatAVL(tmp)
		// result := r.ConcatWithRebalance(tmp)
		return result
	}

	tmpLhs, tmpRhs := r.Split(i)
	tmp = tmpLhs.ConcatAVL(tmp)
	// tmp = tmpLhs.ConcatWithRebalance(tmp)
	// return tmp.ConcatWithRebalance(tmpRhs)
	return tmp.ConcatAVL(tmpRhs)
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

	tmp := tmpLhs.ConcatAVL(tmpRhs)
	r.root = tmp.root
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
			// closeBrNode := node.(*R5NodeCloseBracket)
			result += fmt.Sprintf(")")
		case R5DatatagCloseCall:
			// closeCallNode := node.(*R5NodeCloseCall)
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
			// openBrNode := node.(*R5NodeOpenBracket)
			result += fmt.Sprintf("(")
		case R5DatatagOpenCall:
			// openCallNode := node.(*R5NodeOpenCall)
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
			// fibGenerator: r.fibGenerator,
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
	// buff := append(lhs.Data, rhs.Data...)

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

// func (r *Rope) concatRight(lhs, rhs RopeNode) RopeNode {
// 	if lhs.NodeType() == RopeNodeLeafType && rhs.NodeType() == RopeNodeLeafType {
// 		lhsLeaf := lhs.(*RopeNodeLeaf)
// 		rhsLeaf := rhs.(*RopeNodeLeaf)
//
// 		return r.concatLeaves(lhsLeaf, rhsLeaf)
// 	}
//
// 	diff := lhs.Height() - rhs.Height()
//
// 	if diff == 2 || diff == -2 {
// 		result := &RopeNodeInner{
// 			weight: lhs.Weight() + rhs.Weight(),
// 			height: max(lhs.Height(), rhs.Height()) + 1,
// 			Left:   lhs,
// 			Right:  rhs,
// 		}
//
// 		return r.balanceNodeAVL(result)
// 	}
//
// 	root, ok := lhs.(*RopeNodeInner)
//
// 	if !ok {
// 		result := &RopeNodeInner{
// 			Weight: lhs.Weight() + rhs.Weight(),
// 			Height: max(lhs.Height(), rhs.Height()) + 1,
// 			Left:   lhs,
// 			Right:  rhs,
// 		}
// 		return result
// 	}
//
// 	diff = root.Right.GetHeight() - rhs.GetHeight()
//
// 	if root.Right == nil || (-1 <= diff && diff <= 1) {
//
// 		newNode := &RopeNodeInner{
// 			Left:  root.Right,
// 			Right: rhs,
// 		}
//
// 		r.updateAVLBalanceInfo(newNode)
//
// 		newRoot := &RopeNodeInner{Left: root.Left, Right: newNode}
// 		r.updateAVLBalanceInfo(newRoot)
//
// 		return r.balanceNodeAVL(newRoot)
// 	}
//
// 	root.Right = r.concatRight(root.Right, rhs)
//
// 	return r.balanceNodeAVL(root)
// }
//
// func (r *Rope) concatLeft(lhs, rhs RopeNode) RopeNode {
// 	if lhs.NodeType() == RopeNodeLeafType && rhs.NodeType() == RopeNodeLeafType {
// 		lhsLeaf := lhs.(*RopeNodeLeaf)
// 		rhsLeaf := rhs.(*RopeNodeLeaf)
//
// 		return r.concatLeaves(lhsLeaf, rhsLeaf)
// 	}
//
// 	diff := lhs.Height() - rhs.Height()
//
// 	if diff == 2 || diff == -2 {
// 		result := &RopeNodeInner{
// 			weight: lhs.Weight() + rhs.Weight(),
// 			height: max(lhs.Height(), rhs.Height()) + 1,
// 			Left:   lhs,
// 			Right:  rhs,
// 		}
//
// 		return r.balanceNodeAVL(result)
// 	}
//
// 	root, ok := rhs.(*RopeNodeInner)
//
// 	if !ok {
// 		// panic(123)
// 		result := &RopeNodeInner{
// 			weight: lhs.Weight() + rhs.Weight(),
// 			height: max(lhs.Height(), rhs.Height()) + 1,
// 			Left:   lhs,
// 			Right:  rhs,
// 		}
// 		return result
// 	}
//
// 	diff = root.Left.Height() - lhs.Height()
//
// 	if root.Left == nil || (-1 <= diff && diff <= 1) {
//
// 		newNode := &RopeNodeInner{
// 			Left:  lhs,
// 			Right: root.Left,
// 		}
//
// 		r.updateAVLBalanceInfo(newNode)
//
// 		newRoot := &RopeNodeInner{Left: newNode, Right: root.Right}
// 		r.updateAVLBalanceInfo(newRoot)
//
// 		return r.balanceNodeAVL(newRoot)
// 	}
//
// 	root.Left = r.concatLeft(lhs, root.Left)
//
// 	return r.balanceNodeAVL(root)
// }

// // concatRight: A taller than B
// func concatRight(A, B *RopeNode) *RopeNode {
//     // if right subtree of A is short enough
//     if A.right == nil || abs(A.right.height-B.height) <= 1 {
//         newNode := &RopeNode{left: A.right, right: B}
//         newNode.update()
//         return balance(&RopeNode{left: A.left, right: balance(newNode)})
//     }
//     // recurse
//     A.right = concatRight(A.right, B)
//     return balance(A)
// }
//
// // concatLeft: B taller than A
// func concatLeft(A, B *RopeNode) *RopeNode {
//     if B.left == nil || abs(B.left.height-A.height) <= 1 {
//         newNode := &RopeNode{left: A, right: B.left}
//         newNode.update()
//         return balance(&RopeNode{left: balance(newNode), right: B.right})
//     }
//     B.left = concatLeft(A, B.left)
//     return balance(B)
// }
