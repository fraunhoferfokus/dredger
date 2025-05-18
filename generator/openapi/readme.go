package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"os"
	"path/filepath"

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

	fileName := "ENVIRONMENT.md"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/ENVIRONMENT.md.tmpl"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Info().Msg("CREATE ENVIRONMENT.md")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, readmeConf)
	}
}
