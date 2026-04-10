package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"os"
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

	// locales (localize for en+de)
	targetCoreDir := filepath.Join(Config.Path, "core")
	localesPath := filepath.Join(targetCoreDir, "locales")
	fs.GenerateFolder(localesPath)
	createFileFromTemplate(filepath.Join(targetCoreDir, "localize.go"), "templates/common/core/localize.go.tmpl", conf)
	if _, err := os.Stat(filepath.Join(localesPath, "locale.de.toml")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(localesPath, "locale.de.toml"), "templates/common/core/locales/locale.de.toml", conf)
		createFileFromTemplate(filepath.Join(localesPath, "locale.en.toml"), "templates/common/core/locales/locale.en.toml", conf)
	}
}
