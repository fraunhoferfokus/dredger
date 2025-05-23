package main

import (
	cli "dredger/cli"
	generator "dredger/generator/openapi"

	"embed"
	"os"

	async "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/parser"
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
	// Export embed template filesystem to generator package
	generator.TmplFS = tmplFS

	cli.Execute()

	fileParams := async.FromFileParams{
		Path:         "./examples/simple/asyncapiv3.json",
		MajorVersion: 3,
	}
	spec, err := async.FromFile(fileParams)
	if err != nil {
		log.Err(err).Msg("Could not get spec from file")
	}

	spec.Process() //processes full spec

}
