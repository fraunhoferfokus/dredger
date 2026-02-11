package generator

import (
	"errors"
	"os"
	"path/filepath"

	fs "dredger/fileUtils"

	"github.com/rs/zerolog/log"
)

func generateReadme(conf GeneratorConfig, serverConf ServerConfig) {
	type readmeConfig struct {
		ModuleName string
		Port       int16
	}
	var readmeConf readmeConfig
	readmeConf.ModuleName = conf.ModuleName
	readmeConf.Port = serverConf.Port

	// ENVIRONMENT.md
	fileName := "ENVIRONMENT.md"
	filePath := filepath.Join(Config.Path, fileName)
	templateFile := "templates/common/ENVIRONMENT.md.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Info().Msg("CREATE ENVIRONMENT.md")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, readmeConf)
	}
}
