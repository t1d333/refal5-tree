package runtime

// (('3' '5')) 5 ('2' '5') '2' '4'
func InitViewField(gofunc *R5Function) []ViewFieldNode {
	viewField := []ViewFieldNode{
		&OpenCallViewFieldNode{
			Function: *gofunc,
		},
		&RopeViewFieldNode{Value: NewRope([]R5Node{
			&R5NodeChar{Char: '1'},
			&R5NodeNumber{
				Number: 5,
			},
			&R5NodeFunction{
				Function: &R5Function{
					Name: "No",
				},
			},
			&R5NodeNumber{
				Number: 5,
			},
		})},
		&CloseCallViewFieldNode{},
	}
	return viewField
}
