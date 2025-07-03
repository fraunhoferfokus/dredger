package generator

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generatePolicy(conf GeneratorConfig) {
	log.Info().Msg("Adding policy middleware.")

	// policy.go
	fileName := "policy.go"
	filePath := filepath.Join(config.Path, MiddlewarePackage, fileName)
	templateFile := "templates/openapi/middleware/policy.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)

	// authz.rego
	fileName = "authz.rego"
	filePath = filepath.Join(config.Path, MiddlewarePackage, fileName)
	templateFile = "templates/openapi/middleware/authz.rego.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, conf)
	}
}
