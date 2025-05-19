package runtime

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
}

type R5NodeIllegal struct{}

func (n *R5NodeIllegal) Type() R5Datatag {
	return R5DatatagIllegal
}

type R5NodeChar struct {
	Char byte
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

type R5NodeString struct {
	String string
}

func (n *R5NodeString) Type() R5Datatag {
	return R5DatatagString
}

type R5NodeFunction struct {
	Function *R5Function
}

func (n *R5NodeFunction) Type() R5Datatag {
	return R5DatatagFunction
}

type R5NodeOpenBracket struct {
	CloseOffset int
}

func (n *R5NodeOpenBracket) Type() R5Datatag {
	return R5DatatagOpenBracket
}

type R5NodeCloseBracket struct {
	OpenOffset int
}

func (n *R5NodeCloseBracket) Type() R5Datatag {
	return R5DatatagCloseBracket
}

type R5NodeOpenCall struct {
	CloseOffset int
}

func (n *R5NodeOpenCall) Type() R5Datatag {
	return R5DatatagOpenCall
}

type R5NodeCloseCall struct {
	OpenOffset int
}

func (n *R5NodeCloseCall) Type() R5Datatag {
	return R5DatatagCloseCall
}
