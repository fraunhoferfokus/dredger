package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateLogger(conf GeneratorConfig) {
	fs.GenerateFolder(filepath.Join(conf.OutputPath, LogPkg))
	fs.GenerateFolder(filepath.Join(conf.OutputPath, LoggerPkg))

	// create log.go.tmpl file
	filePath := filepath.Join(conf.OutputPath, LogPkg, "log.go")
	templateFile := "templates/core/log/log.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)

	createFileFromTemplate(filepath.Join(conf.OutputPath, LoggerPkg, "event.go"), "templates/core/log/logger/event.go", conf)
	createFileFromTemplate(filepath.Join(conf.OutputPath, LoggerPkg, "zerolog.go"), "templates/core/log/logger/zerolog.go", conf)

	log.Info().Msg("Created logger successfully.")
}
