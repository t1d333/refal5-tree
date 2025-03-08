package runtime

func InitViewField(gofunc *R5Function) *Rope {
	rope := NewRope([]R5Node{
		&R5NodeOpenCall{},
		&R5NodeFunction{
			Function: gofunc,
		},
		&R5NodeCloseCall{},
	})
	return &rope
}
