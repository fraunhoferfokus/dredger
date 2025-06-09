package generator

import (
	"errors"
	"os"
	"path/filepath"

	fs "dredger/fileUtils"

	"github.com/rs/zerolog/log"
)

// generateConfigFiles legt .env, config.go, configSvc.go und version an.
func generateConfigFiles(serverConf ServerConfig) {
	// 1) .env
	fileName := ".env"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/common/ENVIRONMENT.md.tmpl" // falls du eine andere Vorlage willst, passe hier an
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

	// 2) config.go
	fileName = "config.go"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/openapi/core/config.go.tmpl"
	createFileFromTemplate(filePath, templateFile, serverConf)

	// 3) configSvc.go
	fileName = "configSvc.go"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/openapi/core/configSvc.go.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

	// 4) version (und Symlink)
	fileName = "version"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/openapi/core/version"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
		if err := os.Symlink(filePath, fileName); err != nil {
			log.Warn().
				Err(err).
				Str("source", filePath).
				Str("target", fileName).
				Msg("Could not create symbolic link, bitte manuell anlegen")
		}
	}
}
