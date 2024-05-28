package main

import (
	"os"

	"github.com/aereal/gotest2rdf/internal/cli"
)

func main() {
	os.Exit((&cli.App{Input: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr, Args: os.Args}).Run())
}
