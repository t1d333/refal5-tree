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

	r1 := runtime.NewRope([]runtime.R5Node{
		&runtime.R5NodeNumber{
			Number: 0,
		},
		&runtime.R5NodeNumber{
			Number: 1,
		},
		&runtime.R5NodeNumber{
			Number: 2,
		},
		&runtime.R5NodeNumber{
			Number: 3,
		},
	})

	r2 := runtime.NewRope([]runtime.R5Node{
		&runtime.R5NodeNumber{
			Number: 4,
		},
		&runtime.R5NodeNumber{
			Number: 5,
		},
		&runtime.R5NodeNumber{
			Number: 6,
		},
		&runtime.R5NodeNumber{
			Number: 7,
		},
	})

	r3 := r1.Concat(r2)

	// r2, r3 := r1.Split(3)
	// fmt.Println(r2.Len(), r3.Len())
	// r := r2.Concat(*r3)
	// fmt.Println(r.Len())
	// r3.Insert(1, []runtime.R5Node{&runtime.R5NodeChar{Char: 'A'}})
	// r3.Delete(1)
	
	r4, r5 := r3.Split(5)
	fmt.Println(r4.Len(), r5.Len())
	r5.Delete(2)
	for i := 0; i < r5.Len(); i++ {
		fmt.Println(i, r5.Get(i))
	}
}
