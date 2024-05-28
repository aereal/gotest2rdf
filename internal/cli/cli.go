package cli

import (
	"errors"
	"flag"
	"io"
	"log/slog"
	"os"

	"github.com/aereal/gotest2rdf"
)

type App struct {
	Input       io.Reader
	Out, ErrOut io.Writer
	Args        []string
}

func (a *App) Run() int {
	a.configureLogger()
	fs := flag.NewFlagSet(a.Args[0], flag.ContinueOnError)
	var backlogSize int
	var inputPath string
	var outputPath string
	fs.IntVar(&backlogSize, "backlog-size", 3, "")
	fs.StringVar(&inputPath, "input", "", "input file path (default: stdin)")
	fs.StringVar(&outputPath, "output", "", "output file path (default: stdout)")
	err := fs.Parse(a.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		slog.Error("flag parse failed", slog.String("error", err.Error()))
		return 1
	}

	var input io.Reader = a.Input
	if inputPath != "" {
		f, err := os.Open(inputPath)
		if err != nil {
			slog.Error("failed to open input", slog.String("path", inputPath), slog.String("error", err.Error()))
			return 1
		}
		defer f.Close()
		input = f
	}
	var output io.Writer = a.Out
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			slog.Error("failed to open output", slog.String("path", outputPath), slog.String("error", err.Error()))
			return 1
		}
		defer f.Close()
		output = f
	}
	if err := gotest2rdf.Transform(input, output, gotest2rdf.WithBacklogSize(backlogSize)); err != nil {
		slog.Error("transform failed", slog.String("error", err.Error()))
		return 1
	}
	return 0
}

func (a *App) configureLogger() {
	h := slog.NewTextHandler(a.ErrOut, nil)
	slog.SetDefault(slog.New(h))
}
