package ast

type SentenceRhsType int

const (
	SentenceRhsResultType = iota
	SentenceRhsBlockType  = iota
)

type SentenceNode struct {
	Lhs []*PatternNode
	// NOTE: need conditions?
	Condtitions []*ConditionNode
	Rhs         *SentenceRhsNode
}

type SentenceRhsNode interface {
	GetSentenceRhsType() SentenceRhsType
}

type SentenceRhsResultNode struct{}

func (*SentenceRhsResultNode) GetRhsType() SentenceRhsType {
	return SentenceRhsResultType
}

type SentenceRhsBlockNode struct{}

func (*SentenceRhsBlockNode) GetRhsType() SentenceRhsType {
	return SentenceRhsBlockType
}
