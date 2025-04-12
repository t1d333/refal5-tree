package runtime

func InitViewField(gofunc *R5Function) []ViewFieldNode {
	viewField := []ViewFieldNode{
		&OpenCallViewFieldNode{
			Function: *gofunc,
		},
		&RopeViewFieldNode{Value: NewRope([]R5Node{
			&R5NodeFunction{
				Function: &R5Function{
					Name: "Other",
				},
			},
		})},
		&CloseCallViewFieldNode{},
	}
	return viewField
}
