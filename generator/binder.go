package generator

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateBinder(conf GeneratorConfig) {
	log.Info().Msg("Adding binder middleware.")

	// Binder-Middleware
	fileName := "binder.go"
	filePath := filepath.Join(conf.OutputPath, MiddlewarePackage, fileName)
	templateFile := "templates/openapi/middleware/binder.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)
}
