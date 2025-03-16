package runtime

// (('3' '5')) 5 ('2' '5') '2' '4'
func InitViewField(gofunc *R5Function) *Rope {
	rope := NewRope([]R5Node{
		&R5NodeOpenCall{
			CloseOffset: 8,
		},
		&R5NodeFunction{
			Function: gofunc,
		},
		&R5NodeChar{
			Char: '8',
		},
		&R5NodeChar{
			Char: '3',
		},
		&R5NodeNumber{
			Number: 5,
		},
		&R5NodeChar{
			Char: '3',
		},
		&R5NodeChar{
			Char: '8',
		},
		&R5NodeChar{
			Char: '8',
		},
		&R5NodeChar{
			Char: '3',
		},
		&R5NodeCloseCall{
			OpenOffset: 8,
		},
	})
	return &rope
}
