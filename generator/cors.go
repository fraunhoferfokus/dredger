package generator

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateCorsConfig(conf GeneratorConfig) {
	log.Info().Msg("Adding cors config middleware.")

	// cors.go
	fileName := "cors.go"
	filePath := filepath.Join(Config.Path, MiddlewarePackage, fileName)
	templateFile := "templates/openapi/middleware/cors.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)
}
