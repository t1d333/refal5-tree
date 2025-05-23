package main

import (
	"log"
	"os"

	"github.com/t1d333/refal5-tree/internal/compiler"
	"github.com/urfave/cli/v2"
	// "fmt"
	// "github.com/t1d333/refal5-tree/pkg/runtime"
)

func main() {
	refalCompiler := compiler.NewRefal5Compiler()
	app := &cli.App{
		Name:  "refal5-tree",
		Usage: "Refal5 compiler with tree strings representation",
		Action: func(ctx *cli.Context) error {
			filesArg := ctx.Args()
			names := make([]string, filesArg.Len())

			for i := 0; i < filesArg.Len(); i++ {
				names[i] = filesArg.Get(i)
			}
			refalCompiler.Compile(names, compiler.CompilerOptions{})
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
