package runtime

// (('3' '5')) 5 ('2' '5') '2' '4'
func InitViewField(gofunc *R5Function) *Rope {
	rope := NewRope([]R5Node{
		&R5NodeOpenCall{
			CloseOffset: 15,
		},
		&R5NodeFunction{
			Function: gofunc,
		},
		&R5NodeOpenBracket{
			CloseOffset: 5,
		},
		&R5NodeOpenBracket{
			CloseOffset: 3,
		},
		&R5NodeChar{
			Char: '3',
		},
		&R5NodeChar{
			Char: '3',
		},
		&R5NodeCloseBracket{
			OpenOffset: 3,
		},
		&R5NodeCloseBracket{
			OpenOffset: 5,
		},
		&R5NodeNumber{
			Number: 5,
		},
		&R5NodeOpenBracket{
			CloseOffset: 3,
		},
		&R5NodeChar{
			Char: '2',
		},
		&R5NodeChar{
			Char: '5',
		},
		&R5NodeCloseBracket{
			OpenOffset: 3,
		},
		&R5NodeChar{
			Char: '2',
		},
		&R5NodeChar{
			Char: '4',
		},

		&R5NodeCloseCall{
			OpenOffset: 15,
		},
	})
	return &rope
}
