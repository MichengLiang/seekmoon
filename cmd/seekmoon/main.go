// Package main starts the SeekMoon CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yumiaura/seekmoon/internal/app"
	"github.com/yumiaura/seekmoon/internal/cli"
)

func main() {
	rt, err := app.NewRuntime()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(cli.ExecuteWithCode(context.Background(), rt, cli.Options{Out: os.Stdout, Err: os.Stderr}, os.Args[1:]...))
}
