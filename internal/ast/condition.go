package ast

type ConditionNode struct {
	Pattern []*PatternNode
	Result  []*ResultNode
}
