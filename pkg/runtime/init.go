package runtime

func InitViewField(gofunc *R5Function) []ViewFieldNode {
	viewField := []ViewFieldNode{
		&OpenCallViewFieldNode{
			Function: *gofunc,
		},
		&CloseCallViewFieldNode{},
	}
	return viewField
}
