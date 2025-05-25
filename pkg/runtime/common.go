package runtime

import (
	// "fmt"
	"fmt"
	"os"
	// "runtime/pprof"
)

var Buired map[string][]R5Node

var (
	viewFieldLhs            = []ViewFieldNode{}
	viewFieldRhs            = []ViewFieldNode{}
	primaryActiveExpression = []ViewFieldNode{}
)

var StepCounter = 0

var ProfFile *os.File

func RecognitionImpossible() {
	fmt.Fprintf(os.Stderr, "RECOGNITION IMPOSIBLE\n")
	fmt.Fprintln(os.Stderr)

	viewField := append(viewFieldLhs, viewFieldRhs...)

	fmt.Fprintf(os.Stderr, "STEP NUMBER: %d\n", StepCounter)
	fmt.Fprintln(os.Stderr)

	fmt.Fprintf(os.Stderr, "PRIMARY ACTIVE EXPRESSION:\n")
	fmt.Fprintln(os.Stderr)
	printViewfield(primaryActiveExpression)
	fmt.Fprintf(os.Stderr, ">\n")
	fmt.Fprintln(os.Stderr)

	fmt.Fprintf(os.Stderr, "VIEW FIELD:\n")
	fmt.Fprintln(os.Stderr)
	printViewfield(viewField)
	os.Exit(1)
}

func printViewfield(viewField []ViewFieldNode) {
	indent := ""
	for _, node := range viewField {
		switch node.Type() {
		case CloseBracketType:
			indent = indent[:len(indent)-2]
			fmt.Fprintln(os.Stderr, indent+")")
		case CloseCallType:
			indent = indent[:len(indent)-2]
			fmt.Fprintln(os.Stderr, indent+">")
		case OpenBracketType:
			fmt.Fprintln(os.Stderr, indent+"(")
			indent += "  "
		case OpenCallType:
			call := node.(*OpenCallViewFieldNode)
			fmt.Fprintln(os.Stderr, indent+"<"+call.Function.Name)
			indent += "  "
		case RopeType:
			rope := node.(*RopeViewFieldNode).Value
			for i := 0; i < rope.Len(); i++ {
				fmt.Fprintln(os.Stderr, indent+rope.Get(i).String())
			}
		default:
			panic("unexpected runtime.ViewFieldNodeType")
		}
	}
}

func R5tEmpty(i, j int, r *Rope) bool {
	return i+1 >= j
}

func R5tFunctionLeft(i, left, right int, function *R5Function, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	leftNode := r.Get(left)

	if leftNode == nil || leftNode.Type() != R5DatatagFunction {
		return false
	}

	functionNode := leftNode.(*R5NodeFunction)

	if functionNode.Function.Name != function.Name {
		return false
	}

	idxs[i] = left

	return true
}

func R5tFunctionRight(i, left, right int, function *R5Function, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	rightNode := r.Get(right)

	if rightNode == nil || rightNode.Type() != R5DatatagFunction {
		return false
	}

	functionNode := rightNode.(*R5NodeFunction)

	if functionNode.Function.Name != function.Name {
		return false
	}

	idxs[i] = right

	return true
}

func R5tCharLeft(i, left, right int, c byte, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	leftNode := r.Get(left)

	if leftNode == nil || leftNode.Type() != R5DatatagChar {
		return false
	}

	charNode := leftNode.(*R5NodeChar)
	if charNode.Char != c {
		return false
	}

	idxs[i] = left

	return true
}

func R5tCharRight(i, left, right int, c byte, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	node := r.Get(right)

	if node == nil || node.Type() != R5DatatagChar {
		return false
	}

	charNode := node.(*R5NodeChar)
	if charNode.Char != c {
		return false
	}

	idxs[i] = right

	return true
}

func R5tNumberLeft(i, left, right int, n R5Number, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	nodeLeft := r.Get(left)

	if nodeLeft == nil || nodeLeft.Type() != R5DatatagNumber {
		return false
	}

	numberNode := nodeLeft.(*R5NodeNumber)
	if numberNode.Number != n {
		return false
	}

	idxs[i] = left

	return true
}

func R5tNumberRight(i, left, right int, n R5Number, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	nodeRight := r.Get(right)

	if nodeRight == nil || nodeRight.Type() != R5DatatagNumber {
		return false
	}

	numberNode := nodeRight.(*R5NodeNumber)
	if numberNode.Number != n {
		return false
	}

	idxs[i] = right

	return true
}

func R5tStringLeft(i, left, right int, str string, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	leftNode := r.Get(left)

	if leftNode == nil || leftNode.Type() != R5DatatagString {
		return false
	}

	strNode := leftNode.(*R5NodeString)
	if strNode.Value != str {
		return false
	}

	idxs[i] = left

	return true
}

func R5tStringRight(i, left, right int, str string, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	node := r.Get(right)

	if node == nil || node.Type() != R5DatatagString {
		return false
	}

	strNode := node.(*R5NodeString)
	if strNode.Value != str {
		return false
	}

	idxs[i] = right

	return true
}

func R5tBracketsLeft(i, left, right int, r *Rope, idxs []int) bool {
	left += 1
	if left >= right {
		return false
	}

	nodeLeft := r.Get(left)

	if nodeLeft == nil || nodeLeft.Type() != R5DatatagOpenBracket {
		return false
	}

	bracketNode := nodeLeft.(*R5NodeOpenBracket)
	idxs[i] = left
	idxs[i+1] = left + bracketNode.CloseOffset

	return true
}

func R5tBracketsRight(i, left, right int, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	nodeRight := r.Get(right)

	if nodeRight == nil || nodeRight.Type() != R5DatatagCloseBracket {
		return false
	}

	bracketNode := nodeRight.(*R5NodeCloseBracket)
	idxs[i] = right - bracketNode.OpenOffset
	idxs[i+1] = right

	return true
}

func R5tSymbolVarLeft(i, left, right int, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	leftNode := r.Get(left)

	if leftNode == nil || leftNode.Type() == R5DatatagOpenBracket {
		return false
	}

	idxs[i] = left

	return true
}

func R5tSymbolVarRight(i, left, right int, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	rightNode := r.Get(right)

	if rightNode == nil || rightNode.Type() == R5DatatagCloseBracket {
		return false
	}

	idxs[i] = right

	return true
}

func R5tTermVarLeft(i, left, right int, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	leftNode := r.Get(left)

	if leftNode == nil {
		return false
	}

	idxs[i] = left

	if openBracketNode, ok := leftNode.(*R5NodeOpenBracket); ok {
		idxs[i+1] = left + openBracketNode.CloseOffset
	} else {
		idxs[i+1] = left
	}

	return true
}

func R5tTermVarRight(i, left, right int, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	rightNode := r.Get(right)

	if rightNode == nil {
		return false
	}

	idxs[i+1] = right

	if closeBracketNode, ok := rightNode.(*R5NodeCloseBracket); ok {
		idxs[i] = right - closeBracketNode.OpenOffset
	} else {
		idxs[i] = right
	}

	return true
}

func R5tRepeatedSymbolVarLeft(i, left, right, sample int, r *Rope, idxs []int) bool {
	left += 1

	if left >= right {
		return false
	}

	leftNode := r.Get(left)

	if leftNode == nil {
		return false
	}

	sampleNode := r.Get(idxs[sample])

	if !equalNodes(leftNode, sampleNode) {
		return false
	}

	idxs[i] = left

	return true
}

func R5tRepeatedSymbolVarRight(i, left, right int, sample int, r *Rope, idxs []int) bool {
	right -= 1

	if left >= right {
		return false
	}

	rightNode := r.Get(right)

	if rightNode == nil {
		return false
	}

	if !equalNodes(rightNode, r.Get(idxs[sample])) {
		return false
	}

	idxs[i] = right

	return true
}

func R5tRepeatedExprTermVarLeft(i, left, right, sample int, r *Rope, idxs []int) bool {
	curr := left + 1
	limit := right

	curr_sample := idxs[sample]
	limit_sample := idxs[sample+1] + 1

	for curr != limit && curr_sample != limit_sample && equalNodes(r.Get(curr), r.Get(curr_sample)) {
		curr += 1
		curr_sample += 1
	}

	if curr_sample == limit_sample {
		idxs[i] = left + 1
		idxs[i+1] = curr - 1
		return true
	}

	return false
}

func R5tRepeatedExprTermVarRight(i, left, right, sample int, r *Rope, idxs []int) bool {
	curr := right - 1
	limit := left

	curr_sample := idxs[sample+1]
	limit_sample := idxs[sample] - 1

	for curr != limit && curr_sample != limit_sample && equalNodes(r.Get(curr), r.Get(curr_sample)) {
		curr -= 1
		curr_sample -= 1
	}

	if curr_sample == limit_sample {
		idxs[i] = curr + 1
		idxs[i+1] = right - 1
		return true
	}

	return false
}

func equalNodes(lhs, rhs R5Node) bool {
	if lhs == nil && rhs == nil {
		return true
	}

	if (lhs == nil || rhs == nil) && (lhs != rhs) {
		return false
	}

	if lhs.Type() != rhs.Type() {
		return false
	}

	switch lhs.Type() {
	case R5DatatagChar:
		lhsCharNode := lhs.(*R5NodeChar)
		rhsCharNode := rhs.(*R5NodeChar)
		return lhsCharNode.Char == rhsCharNode.Char
	case R5DatatagOpenBracket:
		return true
	case R5DatatagCloseBracket:
		return true
	case R5DatatagFunction:
		lhsFunctionNode := lhs.(*R5NodeFunction)
		rhsFunctionNode := rhs.(*R5NodeFunction)
		return lhsFunctionNode.Function.Name == rhsFunctionNode.Function.Name
	case R5DatatagNumber:
		lhsNumberNode := lhs.(*R5NodeNumber)
		rhsNumberNode := rhs.(*R5NodeNumber)
		return lhsNumberNode.Number == rhsNumberNode.Number
	case R5DatatagString:
		lhsNumberNode := lhs.(*R5NodeString)
		rhsNumberNode := rhs.(*R5NodeString)
		return lhsNumberNode.Value == rhsNumberNode.Value
	default:
		// TODO: panic
	}
	return false
}

func R5tCloseExprVar(i, left, right int, r *Rope, idxs []int) bool {
	idxs[i] = left + 1
	idxs[i+1] = right - 1
	return true
}

func R5tOpenEvarAdvance(i, right int, r *Rope, idxs []int) bool {
	term := make([]int, 2)

	if R5tTermVarLeft(0, idxs[i+1], right, r, term) {
		idxs[i+1] = term[1]
		return true
	}

	return false
}

func StartMainLoop(initViewField []ViewFieldNode) error {
	// viewFieldLhs := []ViewFieldNode{}
	viewFieldRhs = initViewField

	for len(viewFieldRhs) > 0 {
		curr := viewFieldRhs[0]
		viewFieldRhs = viewFieldRhs[1:]

		switch curr.Type() {
		case CloseCallType:
			if argNode, ok := viewFieldLhs[len(viewFieldLhs)-1].(*RopeViewFieldNode); ok {
				argRope := argNode.Value
				openCall := viewFieldLhs[len(viewFieldLhs)-2].(*OpenCallViewFieldNode)
				viewFieldLhs = viewFieldLhs[:len(viewFieldLhs)-2]
				primaryActiveExpression = []ViewFieldNode{openCall, argNode}
				openCall.Function.Ptr(-1, argRope.Len(), argRope, &viewFieldRhs)
			} else {
				argRope := NewRope([]R5Node{})
				openCall := viewFieldLhs[len(viewFieldLhs)-1].(*OpenCallViewFieldNode)
				primaryActiveExpression = []ViewFieldNode{openCall}
				viewFieldLhs = viewFieldLhs[:len(viewFieldLhs)-1]
				openCall.Function.Ptr(-1, argRope.Len(), argRope, &viewFieldRhs)
			}
			StepCounter += 1
		case OpenCallType:
			viewFieldLhs = append(viewFieldLhs, curr)
		case OpenBracketType:
			viewFieldLhs = append(viewFieldLhs, curr)
		case CloseBracketType:
			if inner, ok := viewFieldLhs[len(viewFieldLhs)-1].(*RopeViewFieldNode); ok {
				viewFieldLhs = viewFieldLhs[:len(viewFieldLhs)-2]
				offset := inner.Value.Len() + 1
				grouped := NewRope([]R5Node{&R5NodeOpenBracket{CloseOffset: offset}})
				grouped = grouped.ConcatAVL(inner.Value)
				grouped = grouped.ConcatAVL(
					NewRope([]R5Node{&R5NodeCloseBracket{OpenOffset: offset}}),
				)

				viewFieldRhs = append(
					[]ViewFieldNode{&RopeViewFieldNode{Value: grouped}},
					viewFieldRhs...)

			} else if _, ok := viewFieldLhs[len(viewFieldLhs)-1].(*OpenBracketViewFieldNode); ok {
				viewFieldLhs = viewFieldLhs[:len(viewFieldLhs)-1]
				grouped := NewRope([]R5Node{
					&R5NodeOpenBracket{CloseOffset: 1},
					&R5NodeCloseBracket{OpenOffset: 1},
				})

				viewFieldRhs = append(
					[]ViewFieldNode{&RopeViewFieldNode{Value: grouped}},
					viewFieldRhs...)
			}
		case RopeType:
			if len(viewFieldLhs) == 0 {
				viewFieldLhs = append(viewFieldLhs, curr)
				continue
			}

			lastLhs := viewFieldLhs[len(viewFieldLhs)-1]
			if lastLhs.Type() == RopeType {
				viewFieldLhs = viewFieldLhs[:len(viewFieldLhs)-1]
				lhsRope := lastLhs.(*RopeViewFieldNode)

				rhsRope := curr.(*RopeViewFieldNode)
				viewFieldLhs = append(
					viewFieldLhs,
					&RopeViewFieldNode{Value: lhsRope.Value.ConcatAVL(rhsRope.Value)},
				)
			} else {
				viewFieldLhs = append(viewFieldLhs, curr)
			}

		default:
			panic("unexpected runtime.ViewFieldNodeType")
		}
	}

	return nil
}

func UpdateOffsets(l, r, diff int, viewField *Rope) {
	openBrStack := []int{}
	openCallStack := []int{}

	for i := 0; i < l; i++ {
		node := viewField.Get(i)
		if node.Type() == R5DatatagOpenBracket {
			openBrStack = append(openBrStack, i)
		} else if node.Type() == R5DatatagCloseBracket {
			openBrStack = openBrStack[:len(openBrStack)-1]
		} else if node.Type() == R5DatatagOpenCall {
			openCallStack = append(openCallStack, i)
		} else if node.Type() == R5DatatagCloseCall {
			openCallStack = openCallStack[:len(openCallStack)-1]
		}
	}

	for _, i := range openBrStack {
		node := viewField.Get(i).(*R5NodeOpenBracket)
		node.CloseOffset += diff
	}

	for _, i := range openCallStack {
		node := viewField.Get(i).(*R5NodeOpenCall)
		node.CloseOffset += diff
	}

	openBrStack = []int{}
	openCallStack = []int{}
	unpairedCloseCall := []int{}
	unpairedCloseBr := []int{}
	for i := l + 1; i < viewField.Len(); i++ {
		node := viewField.Get(i)
		if node.Type() == R5DatatagOpenBracket {
			openBrStack = append(openBrStack, i)
		} else if node.Type() == R5DatatagCloseBracket {
			if len(openBrStack) == 0 {
				unpairedCloseBr = append(unpairedCloseBr, i)
			} else {
				openBrStack = openBrStack[:len(openBrStack)-1]
			}
		} else if node.Type() == R5DatatagOpenCall {
			openCallStack = append(openCallStack, i)
		} else if node.Type() == R5DatatagCloseCall {
			if len(openCallStack) == 0 {
				unpairedCloseCall = append(unpairedCloseCall, i)
			} else {
				openCallStack = openCallStack[:len(openCallStack)-1]
			}
		}
	}

	for _, i := range unpairedCloseBr {
		node := viewField.Get(i).(*R5NodeCloseBracket)
		node.OpenOffset += diff
	}

	for _, i := range unpairedCloseCall {
		node := viewField.Get(i).(*R5NodeCloseCall)
		node.OpenOffset += diff
	}
}

func BuildOpenCallViewFieldNode(f R5Function, viewField *[]ViewFieldNode) {
	*viewField = append(*viewField, &OpenCallViewFieldNode{Function: f})
}

func BuildCloseCallViewFieldNode(viewField *[]ViewFieldNode) {
	*viewField = append(*viewField, &CloseCallViewFieldNode{})
}

func BuildOpenBracketViewFieldNode(viewField *[]ViewFieldNode) {
	*viewField = append(*viewField, &OpenBracketViewFieldNode{})
}

func BuildCloseBracketViewFieldNode(viewField *[]ViewFieldNode) {
	*viewField = append(*viewField, &CloseBracketViewFieldNode{})
}

func BuildRopeViewFieldNode(r *Rope, viewField *[]ViewFieldNode) {
	if r.Len() == 0 {
		return
	}

	*viewField = append(*viewField, &RopeViewFieldNode{Value: r})
}

func MoveExprTermVar(l, r int, src, dst *Rope) {
	if r < l {
		return
	}

	_, tmp := src.Split(l)
	tmp, _ = tmp.Split(r - l + 1)
	//
	// tmp2 := NewRope([]R5Node{})
	//
	// for i := l; i <= r; i++ {
	// 	// buff[i - l] = src.Get(i)
	// 	tmp2 = tmp2.Insert(tmp2.Len(), []R5Node{src.Get(i)})
	// 	// *dst = *dst.Insert(dst.Len(), []R5Node{src.Get(i)})
	// }
	//
	// fmt.Println("------------")
	// VisualizeRope(tmp, 0)
	// fmt.Println("------------")
	// fmt.Println("L", l, "R", r,  tmp.IsAVLBalanced(), tmp2.String() == tmp.String(), tmp.Len() == tmp2.Len())
	// *dst = *dst.ConcatAVL(tmp)

	// tmp := NewRope([]R5Node{})
	//
	// for i := l; i <= r; i++ {
	// 	// buff[i - l] = src.Get(i)
	// 	tmp = tmp.Insert(tmp.Len(), []R5Node{src.Get(i)})
	// 	// *dst = *dst.Insert(dst.Len(), []R5Node{src.Get(i)})
	// }

	*dst = *dst.ConcatAVL(tmp)
}

func CopyExprTermVar(l, r int, src, dst *Rope) {
	// buff := make([]R5Node, r - l + 1)
	tmp := NewRope([]R5Node{})

	for i := l; i <= r; i++ {
		// buff[i - l] = src.Get(i)
		tmp = tmp.Insert(tmp.Len(), []R5Node{src.Get(i)})
		// *dst = *dst.Insert(dst.Len(), []R5Node{src.Get(i)})
	}

	*dst = *dst.ConcatAVL(tmp)
}

func CopySymbolVar(i int, src, dst *Rope) {
	*dst = *dst.Insert(dst.Len(), []R5Node{src.Get(i)})
}
