package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateDockerfile(conf GeneratorConfig, serverConf ServerConfig) {
	type dockerfileConfig struct {
		ModuleName string
		Port       int16
	}

	var dockerfileConf dockerfileConfig
	dockerfileConf.ModuleName = conf.ModuleName
	dockerfileConf.Port = serverConf.Port

	fileName := "Dockerfile"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/Dockerfile.tmpl"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE Dockerfile")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}

	fileName = ".gitlab-ci.yml"
	filePath = filepath.Join(config.Path, fileName)
	templateFile = "templates/gitlab-ci.yml.tmpl"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE .gitlab-ci")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}

	fileName = "image-manifest.yml"
	filePath = filepath.Join(config.Path, fileName)
	templateFile = "templates/image-manifest.yml.tmpl"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE image-manifest")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}
}
