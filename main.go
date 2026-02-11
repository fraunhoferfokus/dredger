package main

import (
	"embed"
	"os"

	"dredger/cli"
	// cli "dredger/cli"  from fokus version

	// Wir brauchen hier die Generator-Pakete,
	// damit wir ihnen die eingebetteten Templates geben können:
	genOpenAPI "dredger/generator"
	//genAsyncAPI "dredger/generator/asyncapi"
	//async "dredger/parser"

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
	//genAsyncAPI.TmplFS = tmplFS

	// Jetzt startet die CLI wie gewohnt:
	cli.Execute()
	//spec, err := async.ParseAsyncAPISpecFile("examples/simple/weather-example.json")
}
