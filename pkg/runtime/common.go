package runtime

import "fmt"

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
	idxs[i] = bracketNode.OpenOffset
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
		lhsNumberNode := lhs.(*R5NodeChar)
		rhsNumberNode := rhs.(*R5NodeChar)
		return lhsNumberNode.Char == rhsNumberNode.Char
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

//	int r05_open_evar_advance(struct r05_node **evar, struct r05_node *right) {
//	  struct r05_node *term[2];
//
//	  if (r05_tvar_left(term, evar[1], right)) {
//	    evar[1] = term[1];
//	    return 1;
//	  } else {
//	    return 0;
//	  }
//	}
func R5tOpenEvarAdvance(i, right int, r *Rope, idxs []int) bool {
	term := make([]int, 2)

	if R5tTermVarLeft(0, idxs[i+1], right, r, term) {
		idxs[i+1] = term[1]
		return true
	}

	return false
}

func StartMainLoop(viewField *Rope) error {
	callStack := [][]int{{1, viewField.Len() - 2}}

	for len(callStack) > 0 {
		tmp := callStack[0]
		callStack = callStack[1:]
		begin := tmp[0]
		end := tmp[1]

		functionNode := viewField.Get(begin + 1)

		if f, ok := functionNode.(*R5NodeFunction); ok {
			f.Function.Ptr(begin+1, end, viewField)
		} else {
			panic("Recognition Imposible")
		}
	}
	return nil
}

func PrintViewField(viewField *Rope) {
	fmt.Print("ViewField{")
	for i := 0; i < viewField.Len(); i++ {
		node := viewField.Get(i)
		switch node.Type() {
		case R5DatatagChar:
			charNode := node.(*R5NodeChar)
			fmt.Printf("(Char: %c), ", charNode.Char)
		case R5DatatagCloseBracket:
			closeBrNode := node.(*R5NodeCloseBracket)
			fmt.Printf("(CloseBracket, OpenOffset: %d), ", closeBrNode.OpenOffset)
		case R5DatatagCloseCall:
			closeCallNode := node.(*R5NodeCloseCall)
			fmt.Printf("(CloseCall, OpenOffset: %d), ", closeCallNode.OpenOffset)
		case R5DatatagFunction:
			funcNode := node.(*R5NodeFunction)
			fmt.Printf("(Function: %s), ", funcNode.Function.Name)
		case R5DatatagIllegal:
			fmt.Printf("(Illegal), ")
		case R5DatatagNumber:
			numberNode := node.(*R5NodeNumber)
			fmt.Printf("(Number: %d), ", numberNode.Number)
		case R5DatatagOpenBracket:
			openBrNode := node.(*R5NodeOpenBracket)
			fmt.Printf("(OpenBracket, CloseOffset: %d), ", openBrNode.CloseOffset)
		case R5DatatagOpenCall:
			openCallNode := node.(*R5NodeOpenCall)
			fmt.Printf("(OpenCall, CloseOffset: %d), ", openCallNode.CloseOffset)
		}
	}
	fmt.Print("}")
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
