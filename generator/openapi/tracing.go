package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateTracing(conf GeneratorConfig) {
	fs.GenerateFolder(filepath.Join(conf.OutputPath, TracingPkg))

	// create tracing.go.tmpl file
	filePath := filepath.Join(conf.OutputPath, TracingPkg, "tracing.go")
	templateFile := "templates/core/tracing/tracing.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)

	log.Info().Msg("Created tracing successfully.")
}
