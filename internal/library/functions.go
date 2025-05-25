package library

var LibraryFunctions = map[string]interface{}{
	"Add":         struct{}{},
	"Sub":         struct{}{},
	"Mul":         struct{}{},
	"Div":         struct{}{},
	"Mod":         struct{}{},
	"Lenw":        struct{}{},
	"Numb":        struct{}{},
	"Symb":        struct{}{},
	"Explode":     struct{}{},
	"Explode_Ext": struct{}{},
	"Implode":     struct{}{},
	"Implode_Ext": struct{}{},
	"Arg":         struct{}{},
	"Compare":     struct{}{},
	"Open":        struct{}{},
	"Close":       struct{}{},
	"Upper":       struct{}{},
	"Lower":       struct{}{},
	"Prout":       struct{}{},
	"Put":         struct{}{},
	"Putout":      struct{}{},
	"Card":        struct{}{},
	"Get":         struct{}{},
	"Step":        struct{}{},
	"Random":      struct{}{},
	"RandomDigit": struct{}{},
	"Ord":         struct{}{},
	"Chr":         struct{}{},
	"Type":        struct{}{},
	"Exit":        struct{}{},
}

var LibraryFuncionAliases = map[string]string{
	"+": "Add",
	"-": "Sub",
	"*": "Mul",
	"/": "Div",
	"%": "Mod",
}

var LibraryFuncionOriginToAlias = map[string]string{
	"Add": "+",
	"Sub": "-",
	"Mul": "*",
	"Div": "/",
	"Mod": "%",
}
