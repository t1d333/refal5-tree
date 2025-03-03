package ast

type SentenceRhsType int

const (
	SentenceRhsResultType = iota
	SentenceRhsBlockType  = iota
)

type SentenceNode struct {
	Lhs []PatternNode
	// NOTE: need conditions?
	Condtitions []*ConditionNode
	Rhs         SentenceRhsNode
}

type SentenceRhsNode interface {
	GetSentenceRhsType() SentenceRhsType
}

type SentenceRhsResultNode struct {
	Result []ResultNode
}

func (*SentenceRhsResultNode) GetSentenceRhsType() SentenceRhsType {
	return SentenceRhsResultType
}

type SentenceRhsBlockNode struct {
	Result []ResultNode
	Body   []*SentenceNode
}

func (*SentenceRhsBlockNode) GetSentenceRhsType() SentenceRhsType {
	return SentenceRhsBlockType
}
