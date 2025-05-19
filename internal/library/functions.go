package library

var LibraryFunctions = map[string]interface{}{
	"Add":     struct{}{},
	"Sub":     struct{}{},
	"Mul":     struct{}{},
	"Div":     struct{}{},
	"Lenw":    struct{}{},
	"Numb":    struct{}{},
	"Symb":    struct{}{},
	"Explode": struct{}{},
	"Implode": struct{}{},
	"Arg":     struct{}{},
	"Compare": struct{}{},
	"Open":    struct{}{},
	"Close":   struct{}{},
	"Upper":   struct{}{},
	"Lower":   struct{}{},
	"Prout":   struct{}{},
	"Put":     struct{}{},
	"Get":     struct{}{},
}

var LibraryFuncionAliases = map[string]string{
	"+": "Add",
	"-": "Sub",
	"*": "Mul",
	"/": "Div",
}

var LibraryFuncionOriginToAlias = map[string]string{
	"Add": "+",
	"Sub": "-",
	"Mul": "*",
	"Div": "/",
}
