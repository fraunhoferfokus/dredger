package main

import (
	"embed"
	"os"

	"dredger/cli"

	// Wir brauchen hier die Generator-Pakete,
	// damit wir ihnen die eingebetteten Templates geben können:
	genAsyncAPI "dredger/generator/asyncapi"
	genOpenAPI "dredger/generator/openapi"

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
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// Hier übergeben wir die eingebetteten Dateien an die Generator-Packages:
	genOpenAPI.TmplFS = tmplFS
	genAsyncAPI.TmplFS = tmplFS

	// Jetzt startet die CLI wie gewohnt:
	cli.Execute()
}
