package generator

import (
	"errors"
	"os"
	"path/filepath"

	fs "dredger/fileUtils"
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

	// 1) Dockerfile
	fileName := "Dockerfile"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/common/Dockerfile.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE Dockerfile")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}

	// 2) .gitlab-ci.yml
	fileName = ".gitlab-ci.yml"
	filePath = filepath.Join(config.Path, fileName)
	templateFile = "templates/common/gitlab-ci.yml.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE .gitlab-ci.yml")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}

	// 3) .gitignore
	fileName = ".gitignore"
	filePath = filepath.Join(config.Path, fileName)
	templateFile = "templates/common/gitignore.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE .gitignore")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}

	// 4) .gitleaksignore
	fileName = ".gitleaksignore"
	filePath = filepath.Join(config.Path, fileName)
	templateFile = "templates/common/gitleaksignore.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE .gitleaksignore")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}

	// 5) image-manifest.yml
	fileName = "image-manifest.yml"
	filePath = filepath.Join(config.Path, fileName)
	templateFile = "templates/common/image-manifest.yml.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("CREATE image-manifest.yml")
		fs.GenerateFile(filePath)
		createFileFromTemplate(filePath, templateFile, dockerfileConf)
	}
}
