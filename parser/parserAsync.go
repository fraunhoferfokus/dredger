package parser

import (
	fs "dredger/fileUtils"
	"errors"

	//"embed"

	//"github.com/rs/zerolog/log"

	//"github.com/rs/zerolog/internal/json"
	async "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/parser"
	v3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/rs/zerolog/log"
)

const schemaFile = "./examples/schemas/asyncapiv3Schema.json"

// ParseAsyncAPISpecFile liest eine AsyncAPI-Datei (YAML oder JSON) ein,
// prüft Basisfelder und gibt sie als Struct zurück.
func ParseAsyncAPISpecFile(path string) (*v3.Specification, error) {
	if !fs.CheckIfFileExists(path) {
		return nil, errors.New("file not found")
	}

	fileParams := async.FromFileParams{
		Path:         path,
		MajorVersion: 3,
	}
	spec, err := async.FromFile(fileParams)
	if err != nil {
		log.Err(err).Msg("Could not get spec from file")
	}

	specV3, ok := spec.(*v3.Specification)
	if !ok {
		log.Error().Msg("Returned spec is not of type v3.Specification")
	}
	err = specV3.Process()
	if err != nil {
		log.Err(err).Msg("Failed to process asyncapi spec")
		return nil, err
	}
	return specV3, nil
}
