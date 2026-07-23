package parser

import (
	fs "dredger/fileUtils"
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func ParseOpenAPISpecFile(path string) (*openapi3.T, error) {
	if !fs.CheckIfFileExists(path) {
		return nil, errors.New("file not found")
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true // TODO make configurable by flag

	// spec, err := loader.LoadFromFile(path)
	spec, err := loadSpec(loader, path, nil)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("OpenAPI Spec file loaded successfully.")

	err = spec.Validate(loader.Context)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("OpenAPI Spec file validated successfully.")

	return spec, err
}
