package main

import (
	"fmt"
	"log"
	"os"

	"github.com/t1d333/refal5-tree/internal/compiler"
	"github.com/urfave/cli/v2"
)

func main() {
	refalCompiler := compiler.NewRefal5Compiler()
	app := &cli.App{
		Name:  "refal5-tree",
		Usage: "Refal5 compiler with tree strings representation",
		Action: func(ctx *cli.Context) error {
			fmt.Println(ctx.NArg())
			fmt.Println(ctx.Args())
			files := ctx.Args()
			fmt.Println(files.Get(0))
			for i := 0; i < files.Len(); i++ {
				refalCompiler.Compile([]string{files.Get(i)}, compiler.CompilerOptions{})
			}
			return nil
		},
		Authors: []*cli.Author{{
			Name:  "Kirill Kiselev",
			Email: "kiselevka2003@yandex.ru",
		}},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
