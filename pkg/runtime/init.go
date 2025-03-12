package runtime

func InitViewField(gofunc *R5Function) *Rope {
	rope := NewRope([]R5Node{
		&R5NodeOpenCall{
			CloseOffset: 3,
		},
		&R5NodeFunction{
			Function: gofunc,
		},
		&R5NodeOpenBracket{
			CloseOffset: 2,
		},
		&R5NodeChar{
			Char: '1',
		},
		&R5NodeCloseBracket{
			OpenOffset: 2,
		},
		&R5NodeNumber{
			Number: 5,
		},
		&R5NodeChar{
			Char: '1',
		},
		&R5NodeChar{
			Char: '2',
		},
		&R5NodeCloseCall{
			OpenOffset: 3,
		},
	})
	return &rope
}
