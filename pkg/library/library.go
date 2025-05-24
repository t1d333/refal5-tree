package library

import (
	"bufio"
	"fmt"
	"unicode"

	// "io"
	"math/big"
	"math/rand/v2"
	"strings"

	// "fmt"
	"os"
	"strconv"

	"github.com/t1d333/refal5-tree/pkg/runtime"
)

// rmcc.ref
// rmcc1.ref: ERROR Function Get is not defined
// rmcc1.ref: ERROR Function Put is not defined

// random.ref
// random.ref: ERROR Function Br is not defined
// random.ref: ERROR Function Chr is not defined
// random.ref: ERROR Function Cp is not defined
// random.ref: ERROR Function Dg is not defined
// random.ref: ERROR Function Div is not defined
// random.ref: ERROR Function Explode_Ext is not defined
// random.ref: ERROR Function Implode_Ext is not defined
// random.ref: ERROR Function Mod is not defined
// random.ref: ERROR Function Ord is not defined
// random.ref: ERROR Function Prout is not defined
// random.ref: ERROR Function Putout is not defined
// random.ref: ERROR Function RandomDigit is not defined
// random.ref: ERROR Function Random is not defined
// random.ref: ERROR Function Rp is not defined
// random.ref: ERROR Function Type is not defined

const (
	MaxOpenFiles = 40
)

var openFiles [MaxOpenFiles]*os.File = [MaxOpenFiles]*os.File{os.Stdin, nil}

func strIntToRefalLong(number string) []runtime.R5Number {
	n := new(big.Int)
	_, ok := n.SetString(number, 10)
	if !ok {
		return []runtime.R5Number{}
	}

	var result []runtime.R5Number
	base := big.NewInt(0).Lsh(big.NewInt(1), 32)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for n.Cmp(zero) > 0 {
		n, mod = new(big.Int).DivMod(n, base, mod)
		result = append(result, runtime.R5Number(mod.Uint64()))
	}

	return result
}

func bigIntToRefalLong(number *big.Int) []runtime.R5Number {
	n := number

	var result []runtime.R5Number
	base := big.NewInt(0).Lsh(big.NewInt(1), 32)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for n.Cmp(zero) > 0 {
		n, mod = new(big.Int).DivMod(n, base, mod)
		result = append(result, runtime.R5Number(mod.Uint64()))
	}

	return result
}

func fromDigitsToBigInt(parts []runtime.R5Number) *big.Int {
	result := big.NewInt(0)
	base := big.NewInt(1)

	tmp := new(big.Int)

	for _, part := range parts {
		tmp.SetUint64(uint64(part))
		tmp.Mul(tmp, base)
		result.Add(result, tmp)
		base.Lsh(base, 32)
	}

	return result
}

func parseRefalLongInt(l, r int, arg *runtime.Rope) (*big.Int, error) {
	curr := l
	node := arg.Get(curr)
	sign := 1

	if charNode, ok := node.(*runtime.R5NodeChar); ok &&
		(charNode.Char == '-' || charNode.Char == '+') {
		curr += 1
		if charNode.Char == '-' {
			sign = -1
		}

	} else if ok {
		return nil, fmt.Errorf("Undefined symbol")
	}

	digits := []runtime.R5Number{}

	for curr < r {
		node = arg.Get(curr)

		number, ok := node.(*runtime.R5NodeNumber)

		if !ok {
			return nil, fmt.Errorf("Undefined symbol")
		}

		digits = append([]runtime.R5Number{number.Number}, digits...)
		curr += 1
	}

	result := fromDigitsToBigInt(digits)

	if sign < 1 {
		return result.Neg(result), nil
	}

	return result, nil
}

func parseAtithmArgs(l, r int, arg *runtime.Rope) (*big.Int, *big.Int, error) {
	if r-l <= 1 {
		return nil, nil, fmt.Errorf("Empty arg")
	}

	curr := l + 1

	node := arg.Get(curr)

	if bracketNode, ok := node.(*runtime.R5NodeOpenBracket); ok {
		lhs, err := parseRefalLongInt(curr+1, curr+bracketNode.CloseOffset, arg)
		if err != nil {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		rhs, err := parseRefalLongInt(curr+bracketNode.CloseOffset+1, r, arg)
		if err != nil {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		return lhs, rhs, nil

	} else if charNode, ok := node.(*runtime.R5NodeChar); ok {
		if charNode.Char != '+' && charNode.Char != '-' {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		if curr >= r {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		lhs, err := parseRefalLongInt(curr, curr+2, arg)
		if err != nil {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		rhs, err := parseRefalLongInt(curr+2, r, arg)
		if err != nil {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		return lhs, rhs, nil

	} else if _, ok := node.(*runtime.R5NodeNumber); ok {
		lhs, err := parseRefalLongInt(curr, curr+1, arg)
		if err != nil {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		rhs, err := parseRefalLongInt(curr+1, r, arg)
		if err != nil {
			return nil, nil, fmt.Errorf("Recognition failed")
		}

		return lhs, rhs, nil
	}

	return nil, nil, fmt.Errorf("Recognition failed")
}

func R5tRandom(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l > 2 {
		panic("Recognition Failed")
	}

	curr := l + 1
	lengthNode, ok := arg.Get(curr).(*runtime.R5NodeNumber)

	if !ok {
		panic("Recognition Failed")
	}

	length := int32(1)

	if lengthNode.Number > 0 {
		length = rand.Int32N(int32(lengthNode.Number) + 1)
	}

	result := []runtime.R5Node{}
	for i := int32(0); i < length; i++ {
		randomNum := rand.Int32()
		result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number(randomNum)})
	}

	fmt.Println("----", result, length)

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(result),
		}}, *rhsStack...)
}

func R5tRandomDigit(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l > 2 {
		panic("Recognition Failed")
	}

	curr := l + 1
	numberNode, ok := arg.Get(curr).(*runtime.R5NodeNumber)

	if !ok {
		panic("Recognition Failed")
	}

	randomNum := rand.Int32N(int32(numberNode.Number))

	result := []runtime.R5Node{&runtime.R5NodeNumber{Number: runtime.R5Number(randomNum)}}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(result),
		}}, *rhsStack...)
}

func R5tStep(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l > 2 {
		panic("Recognition failed")
	}

	step := []runtime.R5Node{&runtime.R5NodeNumber{runtime.R5Number(runtime.StepCounter)}}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(step),
		}}, *rhsStack...)
}

func R5tGet(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1

	numberNode, ok := arg.Get(curr).(*runtime.R5NodeNumber)
	curr += 1

	if !ok {
		panic("Recognition failed")
	}

	fileNo := numberNode.Number % runtime.R5Number(MaxOpenFiles)

	file := openFiles[fileNo]

	buf := make([]byte, 1)

	eof := false
	result := []runtime.R5Node{}

	for {
		n, err := file.Read(buf)
		if err != nil {
			eof = true
			break
		}
		if n == 0 {
			break
		}

		if buf[0] == '\n' {
			break
		}

		result = append(result, &runtime.R5NodeChar{Char: buf[0]})
	}

	if eof {
		result = append(result, &runtime.R5NodeNumber{Number: 0})
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(result),
		}}, *rhsStack...)
}

func R5tPut(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1

	numberNode, ok := arg.Get(curr).(*runtime.R5NodeNumber)
	curr += 1

	if !ok {
		panic("Recognition failed")
	}

	fileNo := numberNode.Number % runtime.R5Number(MaxOpenFiles)

	// TODO: check if file is open
	file := openFiles[fileNo]

	for curr < r {
		node := arg.Get(curr)
		curr += 1

		if charNode, ok := node.(*runtime.R5NodeChar); ok {
			_, err := fmt.Fprintf(file, "%c", charNode.Char)
			if err != nil {
				// TODO: hanlde
			}
		}

		if strNode, ok := node.(*runtime.R5NodeString); ok {
			_, err := fmt.Fprintf(file, "%s ", strNode.String)
			if err != nil {
				// TODO: hanlde
			}
		}

		if numNode, ok := node.(*runtime.R5NodeNumber); ok {
			_, err := fmt.Fprintf(file, "%d ", numNode.Number)
			if err != nil {
				// TODO: hanlde
			}
		}

		if _, ok := node.(*runtime.R5NodeOpenBracket); ok {
			_, err := fmt.Fprintf(file, "(")
			if err != nil {
				// TODO: hanlde
			}
		}

		if _, ok := node.(*runtime.R5NodeCloseBracket); ok {
			_, err := fmt.Fprintf(file, ")")
			if err != nil {
				// TODO: hanlde
			}
		}

	}
}

func R5tProut(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1

	for curr < r {
		node := arg.Get(curr)
		curr += 1

		if charNode, ok := node.(*runtime.R5NodeChar); ok {
			_, err := fmt.Printf("%c", charNode.Char)
			if err != nil {
				// TODO: hanlde
			}
		}

		if strNode, ok := node.(*runtime.R5NodeString); ok {
			_, err := fmt.Printf("%s ", strNode.String)
			if err != nil {
				// TODO: hanlde
			}
		}

		if numNode, ok := node.(*runtime.R5NodeNumber); ok {
			_, err := fmt.Printf("%d ", numNode.Number)
			if err != nil {
				// TODO: hanlde
			}
		}

		if _, ok := node.(*runtime.R5NodeOpenBracket); ok {
			_, err := fmt.Printf("(")
			if err != nil {
				// TODO: hanlde
			}
		}

		if _, ok := node.(*runtime.R5NodeCloseBracket); ok {
			_, err := fmt.Printf(")")
			if err != nil {
				// TODO: hanlde
			}
		}

	}

	fmt.Printf("\n")
}

/*
<Close s.FileNo> == пусто
Семантика. Закрывает открытый файл с номером s.FileNo % 40. Если файл с этим номером не был открыт, функция ничего не делает.
*/
func R5tClose(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l < 2 || r-l > 2 {
		panic("Recognition failed")
	}

	curr := l + 1
	fileNo := -1

	if numberNode, ok := arg.Get(curr).(*runtime.R5NodeNumber); ok {
		fileNo = int(numberNode.Number % 40)
		curr += 1
	} else {
		panic("Recognition failed")
	}

	file := openFiles[fileNo]

	if file != nil {
		file.Close()
	}
}

/*
<Open s.Mode s.FileNo e.FileName?> == пусто

s.Mode ::=

	  'r' | 'w' | 'a'
	|  r  |  w  |  a
	|  rb |  wb |  ab

e.FileName ::= s.CHAR+
*/
func R5tOpen(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l < 3 {
		panic("Recognition failed")
	}

	curr := l + 1
	openFlag := os.O_CREATE
	fileNo := -1

	currNode := arg.Get(curr)

	if charNode, ok := currNode.(*runtime.R5NodeChar); ok &&
		(charNode.Char == 'r' || charNode.Char == 'w' || charNode.Char == 'a') {
		switch charNode.Char {
		case 'r':
			openFlag |= os.O_RDONLY
		case 'w':
			openFlag |= os.O_WRONLY
		case 'a':
			openFlag |= os.O_APPEND
		}
		curr += 1
	} else if ok {
		panic("Recognition failed")
	} else if strNode, ok := currNode.(*runtime.R5NodeString); ok &&
		(strNode.String == "rb" || strNode.String == "wb" || strNode.String == "ab") {
		curr += 1
		// TODO: impl
	} else {
		panic("Recognition failed")
	}

	if numberNode, ok := arg.Get(curr).(*runtime.R5NodeNumber); ok {
		fileNo = int(numberNode.Number % 40)
		curr += 1
	} else {
		panic("Recognition failed")
	}

	// TODO: close if already opened
	fileName := ""

	for curr < r {
		node := arg.Get(curr)

		charNode, ok := node.(*runtime.R5NodeChar)
		if !ok {
			panic("Recognition failed")
		}

		fileName += string(charNode.Char)
		curr += 1
	}

	if fileName == "" {
		fileName = fmt.Sprintf("REFAL%d.DAT", fileNo)
	}

	file, err := os.OpenFile(fileName, openFlag, 0644)
	if err != nil {
		// TODO: handle
		panic("Failed to open file")
	}

	openFiles[fileNo] = file
}

/*
<Arg s.ArgNo> == e.Argument

s.ArgNo ::= s.NUMBER
e.Argument ::= s.CHAR*
Семантика: возвращает аргумент командной строки с указанным номером. Нулевой аргумент — имя вызываемой программы. Если запрашиваемый аргумент не существует — фактическое их число меньше, чем s.ArgNo, возвращается пустая строка.
*/

func R5tArg(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l <= 1 || l-r > 2 {
		panic("Recognition failed")
	}

	curr := l + 1

	argNumb, ok := arg.Get(curr).(*runtime.R5NodeNumber)

	if !ok {
		panic("Recognition failed")
	}

	osArg := []byte{}

	if argNumb.Number < runtime.R5Number(len(os.Args)) {
		osArg = []byte(os.Args[argNumb.Number])
	}

	result := []runtime.R5Node{}

	for _, c := range osArg {
		result = append(result, &runtime.R5NodeChar{Char: c})
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(result),
		}}, *rhsStack...)
}

func R5tCompare(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l <= 1 {
		return
	}

	lhs, rhs, err := parseAtithmArgs(l, r, arg)
	if err != nil {
		panic("Recognition failed")
	}

	compareResult := lhs.Cmp(rhs)

	result := &runtime.R5NodeChar{}

	switch compareResult {
	case 1:
		result.Char = '+'
	case -1:
		result.Char = '-'
	case 0:
		result.Char = '0'
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope([]runtime.R5Node{result}),
		}}, *rhsStack...)
}

func R5tAdd(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l <= 1 {
		return
	}

	var result *big.Int

	lhs, rhs, err := parseAtithmArgs(l, r, arg)
	if err != nil {
		panic("Recognition failed")
	}

	result = lhs.Add(lhs, rhs)
	sign := result.Sign()

	if sign < 0 {
		result = result.Neg(result)
	}

	resultDigits := bigIntToRefalLong(result)

	r5result := []runtime.R5Node{}

	if len(resultDigits) == 0 {
		resultDigits = append(resultDigits, 0)
	}

	for _, digit := range resultDigits {
		r5result = append([]runtime.R5Node{&runtime.R5NodeNumber{Number: digit}}, r5result...)
	}

	if sign < 0 {
		r5result = append([]runtime.R5Node{&runtime.R5NodeChar{Char: '-'}}, r5result...)
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(r5result),
		}}, *rhsStack...)
}

func R5tSub(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l <= 1 {
		return
	}

	var result *big.Int

	lhs, rhs, err := parseAtithmArgs(l, r, arg)
	if err != nil {
		panic("Recognition failed")
	}

	result = lhs.Sub(lhs, rhs)
	sign := result.Sign()

	if sign < 0 {
		result = result.Neg(result)
	}

	resultDigits := bigIntToRefalLong(result)

	r5result := []runtime.R5Node{}

	if len(resultDigits) == 0 {
		resultDigits = append(resultDigits, 0)
	}

	for _, digit := range resultDigits {
		r5result = append([]runtime.R5Node{&runtime.R5NodeNumber{Number: digit}}, r5result...)
	}

	if sign < 0 {
		r5result = append([]runtime.R5Node{&runtime.R5NodeChar{Char: '-'}}, r5result...)
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(r5result),
		}}, *rhsStack...)
}

func R5tMul(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l <= 1 {
		return
	}

	var result *big.Int

	lhs, rhs, err := parseAtithmArgs(l, r, arg)
	if err != nil {
		panic("Recognition failed")
	}

	result = lhs.Mul(lhs, rhs)
	sign := result.Sign()

	if sign < 0 {
		result = result.Neg(result)
	}

	resultDigits := bigIntToRefalLong(result)

	if len(resultDigits) == 0 {
		resultDigits = append(resultDigits, 0)
	}

	r5result := []runtime.R5Node{}

	for _, digit := range resultDigits {
		r5result = append([]runtime.R5Node{&runtime.R5NodeNumber{Number: digit}}, r5result...)
	}

	if sign < 0 {
		r5result = append([]runtime.R5Node{&runtime.R5NodeChar{Char: '-'}}, r5result...)
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(r5result),
		}}, *rhsStack...)
}

func R5tDiv(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	// TODO: need fix
	if r-l <= 1 {
		return
	}

	result := big.NewInt(0)

	lhs, rhs, err := parseAtithmArgs(l, r, arg)
	if err != nil {
		panic("Recognition failed")
	}

	mod := big.NewInt(0)
	mod.Mod(lhs, rhs)
	lhs.Sub(lhs, mod)
	fmt.Println("MOD: ", mod, lhs, rhs)

	result.Div(lhs, rhs)
	sign := result.Sign()

	if sign < 0 {
		result = result.Neg(result)
	}

	resultDigits := bigIntToRefalLong(result)

	if len(resultDigits) == 0 {
		resultDigits = append(resultDigits, 0)
	}

	r5result := []runtime.R5Node{}

	for _, digit := range resultDigits {
		r5result = append([]runtime.R5Node{&runtime.R5NodeNumber{Number: digit}}, r5result...)
	}

	if sign < 0 {
		r5result = append([]runtime.R5Node{&runtime.R5NodeChar{Char: '-'}}, r5result...)
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: runtime.NewRope(r5result),
		}}, *rhsStack...)
}

func R5tImplode_Ext(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	R5tImplode(l, r, arg, rhsStack)
}

/*
	<Implode e.Expr>
*/

func R5tImplode(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	if r-l <= 1 {
		return
	}

	curr := l + 1
	first := arg.Get(curr)
	charNode, ok := first.(*runtime.R5NodeChar)
	if !ok || (!(charNode.Char >= 'a' && charNode.Char <= 'z') &&
		!(charNode.Char >= 'A' && charNode.Char <= 'Z')) {
		nullResult := runtime.NewRope([]runtime.R5Node{&runtime.R5NodeNumber{Number: 0}})
		*rhsStack = append(
			[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
				Value: nullResult.Concat(arg),
			}}, *rhsStack...)

		return
	}

	ident := string(charNode.Char)
	curr += 1

	for {
		node := arg.Get(curr)

		charNode, ok := node.(*runtime.R5NodeChar)
		if !ok {
			break
		}

		if !(charNode.Char >= 'a' && charNode.Char <= 'z') &&
			!(charNode.Char >= 'A' && charNode.Char <= 'Z') &&
			!(charNode.Char >= '0' && charNode.Char <= '9') &&
			charNode.Char != '-' && charNode.Char != '_' {
			break
		}

		ident += string(charNode.Char)
		curr += 1
	}

	identResult := runtime.NewRope([]runtime.R5Node{&runtime.R5NodeString{String: ident}})
	_, other := arg.Split(curr)

	*rhsStack = append([]runtime.ViewFieldNode{
		&runtime.RopeViewFieldNode{
			Value: identResult.Concat(other),
		},
	}, *rhsStack...)
}

/*
<Explode s.Identifier> возвращает строку символов, которая составляла s.Idenitifier .
*/

func R5tExplode_Ext(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	R5tExplode(l, r, arg, rhsStack)
}

func R5tExplode(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	begin := l + 1

	if r-l > 2 || r-l <= 1 {
		panic("Recognition failed")
	}

	curr := arg.Get(begin)

	if curr.Type() != runtime.R5DatatagString {
		panic("Recognition failed")
	}

	identNode := curr.(*runtime.R5NodeString)
	identChars := []byte(identNode.String)

	ident := []runtime.R5Node{}

	for _, c := range identChars {
		ident = append(ident, &runtime.R5NodeChar{Char: c})
	}

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope(
			ident,
		),
	}}, *rhsStack...)
}

func R5tUpper(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	// NOTE: return flat rope, maybe need concat
	curr := l + 1

	result := []runtime.R5Node{}

	for curr < r {
		node := arg.Get(curr)
		curr += 1

		charNode, ok := node.(*runtime.R5NodeChar)

		if !ok {
			result = append(result, node)
			continue
		}

		upper := []byte(strings.ToUpper(string(charNode.Char)))[0]

		result = append(result, &runtime.R5NodeChar{Char: upper})

	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{Value: runtime.NewRope(result)}},
		*rhsStack...)
}

func R5tLower(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	// NOTE: return flat rope, maybe need concat
	curr := l + 1

	result := []runtime.R5Node{}

	for curr < r {
		node := arg.Get(curr)
		curr += 1

		charNode, ok := node.(*runtime.R5NodeChar)

		if !ok {
			result = append(result, node)
			continue
		}

		upper := []byte(strings.ToLower(string(charNode.Char)))[0]

		result = append(result, &runtime.R5NodeChar{Char: upper})

	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{Value: runtime.NewRope(result)}},
		*rhsStack...)
}

/*
<Numb e.Digit-string> возвращает макроцифру, представленную строкой e.Digit-string
*/

func R5tNumb(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1
	strResult := "0"

	// strIntToRefalLong

	first := arg.Get(curr)

	sign := byte('+')

	if charNode, ok := first.(*runtime.R5NodeChar); ok &&
		(charNode.Char == '-' || charNode.Char == '+') {
		sign = charNode.Char
		curr += 1
	}

	for curr < r {
		currNode := arg.Get(curr)

		if currNode.Type() != runtime.R5DatatagChar {
			break
		}

		charNode := currNode.(*runtime.R5NodeChar)

		if !(charNode.Char >= '0' && charNode.Char <= '9') {
			break
		}

		strResult += string(charNode.Char)
		curr += 1
	}

	result := []runtime.R5Node{}

	number := strIntToRefalLong(strResult)

	for _, n := range number {
		result = append([]runtime.R5Node{&runtime.R5NodeNumber{Number: n}}, result...)
	}

	if sign != '+' {
		result = append([]runtime.R5Node{&runtime.R5NodeChar{Char: sign}}, result...)
	}

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{Value: runtime.NewRope(result)}},
		*rhsStack...)
}

/*
<Symb s.Macrodigit>

является обратной к функции Numb . Она возвращает строку десятичных цифр, представляющую s.Macrodigit .

*/

func R5tSymb(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1

	// if r-l > 2 {
	// fmt.Println("ARG: ", arg.String())
	// panic("Recognition failed")
	// }

	first := arg.Get(curr)

	sign := byte(0)
	if charNode, ok := first.(*runtime.R5NodeChar); ok &&
		(charNode.Char == '-' || charNode.Char == '+') {
		sign = charNode.Char
		curr += 1
		first = arg.Get(curr)
	} else if ok {
		panic("Recognition failed")
	}

	if first.Type() != runtime.R5DatatagNumber {
		panic("Recognition failed")
	}

	numberNode := first.(*runtime.R5NodeNumber)
	numberChars := []byte(strconv.Itoa(int(numberNode.Number)))
	number := []runtime.R5Node{}

	for _, c := range numberChars {
		number = append(number, &runtime.R5NodeChar{Char: c})
	}

	if sign != 0 {
		number = append([]runtime.R5Node{&runtime.R5NodeChar{Char: sign}}, number...)
	}

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope(
			number,
		),
	}}, *rhsStack...)
}

/*
<Type e.Expr>

возвращает s.Type e.Expr , где e.Expr является неизменным, а s.Type зависит от типа первого элемента выражения e.Expr .

  s.Type   e.Expr начинается с:
  'L'      буквы
  'D'      цифры
  'F'      идентификатора или имени функции
  'N'      макроцифры
  'R'      действительного числа
  'O'      любого другого символа
  'B'      левой скобки
  '*'      e.Expr  является пустым
*/

func R5tType(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	resultType := &runtime.R5NodeChar{}
	resultSubType := &runtime.R5NodeChar{}

	result := []runtime.R5Node{resultType, resultSubType}

	if r-l < 2 {
		resultType.Char = '*'
		resultSubType.Char = '0'
		return
	} else {
		curr := l + 1

		first := arg.Get(curr)

		switch first.Type() {
		case runtime.R5DatatagChar:
			charNode := first.(*runtime.R5NodeChar)
			if charNode.Char >= 'a' && charNode.Char <= 'z' {
				resultType.Char = 'L'
				resultSubType.Char = 'l'
			} else if charNode.Char >= 'Z' && charNode.Char <= 'Z' {
				resultType.Char = 'L'
				resultSubType.Char = 'u'
			} else if unicode.IsPrint(rune(charNode.Char)) && unicode.IsUpper(rune(charNode.Char)) {
				resultType.Char = 'P'
				resultSubType.Char = 'u'
			} else if unicode.IsPrint(rune(charNode.Char)) && !unicode.IsUpper(rune(charNode.Char)) {
				resultType.Char = 'P'
				resultSubType.Char = 'l'
			} else if unicode.IsUpper(rune(charNode.Char)) {
				resultType.Char = 'P'
				resultSubType.Char = 'u'
			} else {
				resultType.Char = 'P'
				resultSubType.Char = 'l'
			}
		case runtime.R5DatatagFunction:
			resultType.Char = 'W'
			resultSubType.Char = 'q'
		case runtime.R5DatatagNumber:
			resultType.Char = 'N'
			resultSubType.Char = '0'
		case runtime.R5DatatagOpenBracket:
			resultType.Char = 'B'
			resultSubType.Char = '0'
		case runtime.R5DatatagString:
			resultType.Char = 'W'
			strNode := first.(*runtime.R5NodeString)
			if !unicode.IsDigit(rune(strNode.String[0])) && !unicode.IsLetter(rune(strNode.String[0])) {
				resultSubType.Char = 'q'
			} else {
				resultSubType.Char = 'i'
				for _, c := range strNode.String {
					if c != '-' && c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
						resultSubType.Char = 'q'
						break
					}
				}
			}
		}
	}

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope(result).Concat(arg),
	}}, *rhsStack...)
}

func R5tLenw(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	currIdx := l + 1
	argLen := 0

	for currIdx < r {
		argLen += 1
		curr := arg.Get(currIdx)
		if curr.Type() == runtime.R5DatatagOpenBracket {
			bracketNode := curr.(*runtime.R5NodeOpenBracket)
			currIdx += bracketNode.CloseOffset
		}
		currIdx += 1
	}

	charNode := &runtime.R5NodeNumber{Number: runtime.R5Number(argLen)}
	tmpRope := runtime.NewRope([]runtime.R5Node{charNode}).Concat(arg)

	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{Value: tmpRope}}, *rhsStack...)
}

func R5tMod(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tChr(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1
	result := []runtime.R5Node{}

	for curr < r {
		node := arg.Get(curr)

		curr += 1
		if numberNode, ok := node.(*runtime.R5NodeNumber); !ok {
			result = append(result, node)
			continue
		} else {
			result = append(result, &runtime.R5NodeChar{Char: byte(numberNode.Number % 256)})
		}
	}

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope(result),
	}}, *rhsStack...)
}

func R5tOrd(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1
	result := []runtime.R5Node{}

	for curr < r {
		node := arg.Get(curr)

		curr += 1
		if charNode, ok := node.(*runtime.R5NodeChar); !ok {
			result = append(result, node)
			continue
		} else {
			if charNode.Char != '\\' {
				result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number(charNode.Char)})
				continue
			}

			charNode = arg.Get(curr).(*runtime.R5NodeChar)
			switch charNode.Char {
			case '\\':
				result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number('\\')})
			case 't':
				result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number('\t')})
			case 'r':
				result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number('\r')})
			case 'n':
				result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number('\n')})
			case '"':
				result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number('"')})
			}

			// result = append(result, &runtime.R5NodeNumber{Number: runtime.R5Number()})
		}

	}

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope(result),
	}}, *rhsStack...)
}

// TODO: implement
func R5tCard(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text() // Возвращает строку без символа конца строки
		lineBytes := []byte(line)
		rope := runtime.NewRope([]runtime.R5Node{})
		for _, b := range lineBytes {
			charNode := &runtime.R5NodeChar{Char: b}
			rope = rope.Insert(rope.Len(), []runtime.R5Node{charNode})

		}

		*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
			Value: rope,
		}}, *rhsStack...)
		return
	}

	endCode := &runtime.R5NodeNumber{Number: 0}
	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope([]runtime.R5Node{endCode}),
	}}, *rhsStack...)
}

func R5tExistsFile(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	filename := ""
	for i := l; i < r; i++ {
		curr := arg.Get(i)
		if curr.Type() != runtime.R5DatatagChar {
			// Panic?
			return
		}

		charNode := curr.(*runtime.R5NodeChar)
		filename += string(charNode.Char)
	}

	_, err := os.Stat(filename)

	var result *runtime.Rope
	if os.IsNotExist(err) {
		result = runtime.NewRope([]runtime.R5Node{&runtime.R5NodeFunction{
			Function: &runtime.R5Function{
				Name: "False",
			},
		}})
	} else {
		result = runtime.NewRope([]runtime.R5Node{&runtime.R5NodeFunction{
			Function: &runtime.R5Function{
				Name: "True",
			},
		}})
	}
	*rhsStack = append(
		[]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{Value: result}},
		*rhsStack...)
}

func R5tPrint(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}
