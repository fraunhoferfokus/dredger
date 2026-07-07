package generator

import (
	parser "dredger/parser"

	"github.com/rs/zerolog/log"
)

func GenerateComponents(conf GeneratorConfig) error {
	spec, err := parser.ParseComponentSpecFile(conf.OpenAPIPath)
	if err != nil {
		return err
	}

	Config.Name = conf.ModuleName
	Config.Path = conf.OutputPath

	createProjectPathDirectory(conf)

	GenerateTypes(spec, Config)

	log.Info().Msg("Created all files successfully.")
	return nil
}
