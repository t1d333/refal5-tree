package main

import (
	// "log"
	// "os"

	// "github.com/t1d333/refal5-tree/internal/compiler"
	// "github.com/urfave/cli/v2"
	"fmt"

	"github.com/t1d333/refal5-tree/pkg/runtime"
)

func main() {
	n := &runtime.R5NodeOpenCall{}
	n1 := &runtime.R5NodeOpenBracket{}
	n2 := &runtime.R5NodeNumber{}
	rope := runtime.NewRope([]runtime.R5Node{n, n1, n2})
	rope2 := runtime.NewRope(
		[]runtime.R5Node{
			&runtime.R5NodeCloseCall{},
			&runtime.R5NodeCloseBracket{},
			&runtime.R5NodeChar{},
		},
	)
	tmp := rope.Concat(rope2)
	fmt.Println(tmp.Get(4).Type() == tmp.Get(3).Type())
	fmt.Println(tmp.Len())
	// refalCompiler := compiler.NewRefal5Compiler()
	// app := &cli.App{
	// 	Name:  "refal5-tree",
	// 	Usage: "Refal5 compiler with tree strings representation",
	// 	Action: func(ctx *cli.Context) error {
	// 		files := ctx.Args()
	// 		for i := 0; i < files.Len(); i++ {
	// 			refalCompiler.Compile([]string{files.Get(i)}, compiler.CompilerOptions{})
	// 		}
	// 		return nil
	// 	},
	// 	Authors: []*cli.Author{{
	// 		Name:  "Kirill Kiselev",
	// 		Email: "kiselevka2003@yandex.ru",
	// 	}},
	// }
	//
	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }
}
