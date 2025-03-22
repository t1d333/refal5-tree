package runtime

// (('3' '5')) 5 ('2' '5') '2' '4'
func InitViewField(gofunc *R5Function) []ViewFieldNode {
	viewField := []ViewFieldNode{
		&OpenCallViewFieldNode{
			Function: *gofunc,
		},
		&RopeViewFieldNode{Value: NewRope([]R5Node{
			&R5NodeOpenBracket{
				CloseOffset: 8,
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
			&R5NodeCloseBracket{
				OpenOffset: 8,
			},
		})},
		&CloseCallViewFieldNode{},
	}
	return viewField
}
