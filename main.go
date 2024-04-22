package main

import (
	cli "dredger/cli"
	generator "dredger/generator"

	"embed"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates
var tmplFS embed.FS

func main() {
	// Set up zerolog time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Export embed template filesystem to generator package
	generator.TmplFS = tmplFS

	cli.Execute()
}
