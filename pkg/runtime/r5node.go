package runtime

import (
	"fmt"
	"strconv"
)

type R5Datatag int

const (
	R5DatatagIllegal R5Datatag = iota
	R5DatatagChar
	R5DatatagFunction
	R5DatatagNumber
	R5DatatagString
	R5DatatagOpenBracket
	R5DatatagCloseBracket
	R5DatatagOpenCall
	R5DatatagCloseCall
)

type R5Node interface {
	Type() R5Datatag
	String() string
}

type R5NodeIllegal struct{}

func (n *R5NodeIllegal) Type() R5Datatag {
	return R5DatatagIllegal
}

func (n *R5NodeIllegal) String() string {
	return "illegal"
}

type R5NodeChar struct {
	Char byte
}

func (n *R5NodeChar) String() string {
	tmp := strconv.Quote(string(n.Char))
	return fmt.Sprintf("'%s'", tmp)
}

func (n *R5NodeChar) Type() R5Datatag {
	return R5DatatagChar
}

type R5Number uint32

type R5NodeNumber struct {
	Number R5Number
}

func (n *R5NodeNumber) Type() R5Datatag {
	return R5DatatagNumber
}

func (n *R5NodeNumber) String() string {
	return fmt.Sprintf("%d", n.Number)
}

type R5NodeString struct {
	Value string
}

func (n *R5NodeString) Type() R5Datatag {
	return R5DatatagString
}

func (n *R5NodeString) String() string {
	return fmt.Sprintf("\"%s\"", n.Value)
}

type R5NodeFunction struct {
	Function *R5Function
}

func (n *R5NodeFunction) Type() R5Datatag {
	return R5DatatagFunction
}

func (n *R5NodeFunction) String() string {
	return fmt.Sprintf("%s", n.Function.Name)
}

type R5NodeOpenBracket struct {
	CloseOffset int
}

func (n *R5NodeOpenBracket) String() string {
	return "(" 
}

func (n *R5NodeOpenBracket) Type() R5Datatag {
	return R5DatatagOpenBracket
}

type R5NodeCloseBracket struct {
	OpenOffset int
}

func (n *R5NodeCloseBracket) String() string {
	return "(" 
}

func (n *R5NodeCloseBracket) Type() R5Datatag {
	return R5DatatagCloseBracket
}

type R5NodeOpenCall struct {
	CloseOffset int
}

func (n *R5NodeOpenCall) String() string {
	return "<" 
}

func (n *R5NodeOpenCall) Type() R5Datatag {
	return R5DatatagOpenCall
}


type R5NodeCloseCall struct {
	OpenOffset int
}

func (n *R5NodeCloseCall) String() string {
	return ">" 
}

func (n *R5NodeCloseCall) Type() R5Datatag {
	return R5DatatagCloseCall
}
