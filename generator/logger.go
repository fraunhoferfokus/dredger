package generator

import (
	fs "dredger/fileUtils"
	//"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateLogger(conf GeneratorConfig) {
	fs.GenerateFolder(filepath.Join(conf.OutputPath, LogPkg))
	fs.GenerateFolder(filepath.Join(conf.OutputPath, LoggerPkg))
	fs.GenerateFolder(filepath.Join(conf.OutputPath, LokiPkg))

	// log.go
	filePath := filepath.Join(conf.OutputPath, LogPkg, "log.go")
	templateFile := "templates/openapi/core/log/log.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)

	// event.go & zerolog.go & loki.go
	createFileFromTemplate(
		filepath.Join(conf.OutputPath, LoggerPkg, "event.go"),
		"templates/openapi/core/log/logger/event.go",
		conf,
	)
	createFileFromTemplate(
		filepath.Join(conf.OutputPath, LoggerPkg, "zerolog.go"),
		"templates/openapi/core/log/logger/zerolog.go",
		conf,
	)
	createFileFromTemplate(
		filepath.Join(conf.OutputPath, LokiPkg, "loki.go"),
		"templates/openapi/core/log/loki/loki.go",
		conf,
	)

	// Logger-Middleware
	fileName := "logger.go"
	filePath = filepath.Join(conf.OutputPath, MiddlewarePackage, fileName)
	templateFile = "templates/openapi/middleware/logger.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)

	log.Info().Msg("Created logger successfully.")
}
