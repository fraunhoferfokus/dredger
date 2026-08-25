package generator

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateCSRF(conf GeneratorConfig) {
	log.Info().Msg("Adding csrf middleware.")

	// Binder-Middleware
	fileName := "csrf.go"
	filePath := filepath.Join(conf.OutputPath, MiddlewarePackage, fileName)
	templateFile := "templates/openapi/middleware/csrf.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)
}
