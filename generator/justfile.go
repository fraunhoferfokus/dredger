package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateJustfile(conf GeneratorConfig, serverConf ServerConfig) {
	type justfileConfig struct {
		ModuleName string
		Port       int16
	}

	var justfileConf justfileConfig
	justfileConf.ModuleName = conf.ModuleName
	justfileConf.Port = serverConf.Port

	fileName := "Justfile"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/Justfile.tmpl"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Info().Msg("CREATE Justfile")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, justfileConf)
	}
}
