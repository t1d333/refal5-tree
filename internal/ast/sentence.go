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

func (s *SentenceNode) String() string {
	res := ""
	for _, p := range s.Lhs {
		res += " " + p.String()
	}

	for _, c := range s.Condtitions {
		res += ", "
		for _, p := range c.Pattern {
			res += " " + p.String()
		}
		res += " : "
		for _, r := range c.Result {
			res += " " + r.String()
		}
	}

	if s.Rhs.GetSentenceRhsType() == SentenceRhsResultType {
		res += " = "
		rhsRes := s.Rhs.(*SentenceRhsResultNode)
		for _, r := range rhsRes.Result {
			res += " " + r.String()
		}
	}

	return res
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
