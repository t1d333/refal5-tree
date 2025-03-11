package runtime

func InitViewField(gofunc *R5Function) *Rope {
	rope := NewRope([]R5Node{
		&R5NodeOpenCall{
			CloseOffset: 3,
		},
		&R5NodeFunction{
			Function: gofunc,
		},

		&R5NodeChar{
			Char: '1',
		},
		&R5NodeCloseCall{
			OpenOffset: 3,
		},
	})
	return &rope
}
