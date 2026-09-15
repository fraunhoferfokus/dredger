package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateCSRF(conf GeneratorConfig) {
	// Binder-Middleware
	fileName := "csrf.go"
	path := filepath.Join(conf.OutputPath, MiddlewarePackage, fileName)
	templateFile := "templates/openapi/middleware/csrf.go.tmpl"
	if fs.CheckIfFileExists(path) { // dont overwrite as CSRF policy is service specific
		log.Info().Msg("Keeping existing csrf middleware.")
	}
	log.Info().Msg("Adding csrf middleware.")
	createFileFromTemplate(path, templateFile, conf)
}
