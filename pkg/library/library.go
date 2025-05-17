package library

import (
	"bufio"
	// "fmt"
	"os"
	"strconv"

	"github.com/t1d333/refal5-tree/pkg/runtime"
)

// Refal-05 lib usage

// Add 1
// Arg 1
// Chr 3
// Div 2
// First 1
// Get 2
// Implode 1
// Mod 2
// Open 2
// Ord 4
// Prout 7
// Putout 2
// Sub 2
// Type 19
// Upper 1
// GetEnv 4
// System 1
// Exit 3
// Close 2
// ExistFile 1
// Implode_Ext 1
// Compare 6
// ListOfBuiltin 2

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
<Explode s.Identifier>

возвращает строку символов, которая составляла s.Idenitifier .
*/

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

/*
<Numb e.Digit-string>

возвращает макроцифру, представленную строкой e.Digit-string
*/

func R5tNumb(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1
	strResult := "0"

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

	result, _ := strconv.Atoi(strResult)

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope(
			[]runtime.R5Node{&runtime.R5NodeNumber{Number: runtime.R5Number(result)}},
		),
	}}, *rhsStack...)
}

/*
<Symb s.Macrodigit>

является обратной к функции Numb . Она возвращает строку десятичных цифр, представляющую s.Macrodigit .

*/

func R5tSymb(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
	curr := l + 1

	if r-l > 2 {
		panic("Recognition failed")
	}

	first := arg.Get(curr)

	if first.Type() != runtime.R5DatatagNumber {
		panic("Recognition failed")
	}

	numberNode := first.(*runtime.R5NodeNumber)
	numberChars := []byte(strconv.Itoa(int(numberNode.Number)))
	number := []runtime.R5Node{}

	for _, c := range numberChars {
		number = append(number, &runtime.R5NodeChar{Char: c})
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
	if arg.Len() == 0 {
		panic("Recognition failed")
	}

	first := arg.Get(0)
	result := &runtime.R5NodeChar{}

	switch first.Type() {
	case runtime.R5DatatagChar:
		charNode := first.(*runtime.R5NodeChar)
		if (charNode.Char >= 'a' && charNode.Char <= 'z') ||
			(charNode.Char >= 'Z' && charNode.Char <= 'Z') {
			result.Char = 'L'
		} else if charNode.Char >= '0' && charNode.Char <= '9' {
			result.Char = 'D'
		} else {
			result.Char = 'O'
		}
	case runtime.R5DatatagFunction:
		result.Char = 'F'
	case runtime.R5DatatagNumber:
		result.Char = 'N'
	case runtime.R5DatatagOpenBracket:
		result.Char = 'B'
	case runtime.R5DatatagOpenCall:
		result.Char = 'B'
	case runtime.R5DatatagString:
		result.Char = 'O'
	default:
		panic("Recognition failed")

	}

	*rhsStack = append([]runtime.ViewFieldNode{&runtime.RopeViewFieldNode{
		Value: runtime.NewRope([]runtime.R5Node{result}),
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

func R5tAdd(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tSub(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tDiv(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tMod(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tArg(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}

func R5tChr(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
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

func R5tProut(l, r int, arg *runtime.Rope, rhsStack *[]runtime.ViewFieldNode) {
}
