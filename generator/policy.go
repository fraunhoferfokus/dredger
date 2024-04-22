package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generatePolicy(conf GeneratorConfig) {
	log.Info().Msg("Adding policy middleware.")

	fileName := "policy.go"
	filePath := filepath.Join(config.Path, MiddlewarePackage, fileName)
	templateFile := "templates/middleware/policy.go.tmpl"
	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)

	fileName = "authz.go"
	filePath = filepath.Join(config.Path, MiddlewarePackage, fileName)
	templateFile = "templates/middleware/authz.go.tmpl"
	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)

	fileName = "authz.rego"
	filePath = filepath.Join(config.Path, MiddlewarePackage, fileName)
	templateFile = "templates/middleware/authz.rego.tmpl"
	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)
}
