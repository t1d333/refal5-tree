package runtime

func R5tEmptyHole(i, j int, re *Rope) bool {
	return i >= j
}

func R5tFunctionLeft(i, lhs, rhs int, function *R5Function, r *Rope) {
}

func R5tFunctionRight(i, lhs, rhs int, function *R5Function, r *Rope) {
}

func R5tCharLeft(i, lhs, rhs int, c byte, r *Rope) {
}

func R5tCharRight(i, lhs, rhs int, c byte, r *Rope) {
}

func R5tNumberLeft(i, lhs, rhs int, n uint32, r *Rope) {
}

func R5tNumberRight(i, lhs, rhs int, n uint32, r *Rope) {
}

func R5tSymbolVarLeft(i, lhs, rhs int, r *Rope) {
}

func R5tSymbolRight(i, lhs, rhs int, r *Rope) {
}

func R5tTermVarLeft(i, lhs, rhs int, r *Rope) {
}

func R5tTermVarRight(i, lhs, rhs int, r *Rope) {
}

func R5tRepeatedSymbolVarLeft(i, lhs, rhs int, r *Rope) {
}

func R5tRepeatedSymbolRight(i, lhs, rhs int, r *Rope) {
}

func R5tRepeatedTermVarLeft(i, lhs, rhs int, r *Rope) {
}

func R5tRepeatedTermVarRight(i, lhs, rhs int, r *Rope) {
}
