package runtime

// (('3' '5')) 5 ('2' '5') '2' '4'
func InitViewField(gofunc *R5Function) []ViewFieldNode {
	viewField := []ViewFieldNode{
		&OpenCallViewFieldNode{
			Function: *gofunc,
		},
		&RopeViewFieldNode{Value: NewRope([]R5Node{
			&R5NodeOpenBracket{CloseOffset: 4},
			&R5NodeChar{
				Char: '8',
			},
			&R5NodeChar{
				Char: '3',
			},
			&R5NodeChar{
				Char: '8',
			},
			&R5NodeCloseBracket{OpenOffset: 4},

			&R5NodeOpenBracket{CloseOffset: 4},
			&R5NodeChar{
				Char: '3',
			},
			&R5NodeChar{
				Char: '8',
			},
			&R5NodeChar{
				Char: '3',
			},
			&R5NodeOpenBracket{CloseOffset: 4},
			&R5NodeChar{
				Char: '3',
			},
		})},
		&CloseCallViewFieldNode{},
	}
	return viewField
}
