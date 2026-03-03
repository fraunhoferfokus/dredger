package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateValidation(conf GeneratorConfig) {
	log.Info().Msg("Adding validation middleware.")

	fileName := "validation.go"

	filePath := filepath.Join(Config.Path, MiddlewarePackage, fileName)
	templateFile := "templates/openapi/middleware/validation.go.tmpl"
	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)
}
