package generator

import (
	"os"
	"path/filepath"

	fs "dredger/fileUtils"

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
	filePath := filepath.Join(Config.Path, fileName)
	templateFile := "templates/common/Justfile.tmpl"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Info().Msg("CREATE Justfile")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, justfileConf)
	}
}
