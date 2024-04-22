package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateAuthzFile(conf GeneratorConfig) {
	log.Info().Msg("Adding auth middleware.")

	fileName := "authz.go"
	filePath := filepath.Join(config.Path, AuthzPkg, fileName)
	templateFile := "templates/middleware/authz.go.tmpl"

	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)
}
