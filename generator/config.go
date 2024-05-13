package generator

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateConfigFiles(serverConf ServerConfig) {
	// create app.env file if not exist
	fileName := ".env"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/core/app.env.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

	// create config.go file
	fileName = "config.go"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/core/config.go.tmpl"
	createFileFromTemplate(filePath, templateFile, serverConf)

	// create configSvc.go extension file if not exist
	fileName = "configSvc.go"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/core/configSvc.go.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

	// create version file if not exist
	fileName = "version"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/core/version"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
		if err = os.Symlink(filePath, fileName); err != nil {
			log.Error().Err(err).Str("source", filePath).Str("target", fileName).Msg("Could not create symbolic Link")
		}
		// if runtime.GOOS == "windows" {
		// 	extCmd.RunCommand("mklink /h "+fileName+" "+filePath, config.Path)
		// } else {
		// 	extCmd.RunCommand("ln -s "+filePath+" "+fileName, config.Path)
		// }
	}
}
