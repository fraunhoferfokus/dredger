package generator

import (
	"embed"
	"path/filepath"

	fs "dredger/fileUtils"
	"dredger/parser"

	"github.com/rs/zerolog/log"
)

// TmplFS holds the embedded async templates. It is assigned by main.go.
var TmplFS embed.FS

// mainConfig contains values passed to the templates.
type mainConfig struct {
	ModuleName string
	Title      string
}

// GenerateService creates a minimal async service based on the given spec.
func GenerateService(specPath, outputPath, moduleName string) error {
	spec, err := parser.ParseAsyncAPISpecFile(specPath)
	if err != nil {
		return err
	}

	conf := mainConfig{
		ModuleName: moduleName,
		Title:      spec.Info.Title,
	}

	fs.GenerateFolder(outputPath)
	base := filepath.Join(outputPath, "src")
	subdirs := []string{"cli", "config", "handler", "logger", "model", "schemas", "utils", "tracing", "warp_server", "policy"}
	fs.GenerateFolder(base)
	for _, d := range subdirs {
		fs.GenerateFolder(filepath.Join(base, d))
	}

	createFileFromTemplate(filepath.Join(outputPath, "README.md"), "templates/asyncapi/README.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(outputPath, "ENVIRONMENT.md"), "templates/asyncapi/ENVIRONMENT.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(outputPath, "go.mod"), "templates/asyncapi/go.mod.tmpl", conf)

	createFileFromTemplate(filepath.Join(base, "main.go"), "templates/asyncapi/src/main.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "cli", "cli.go"), "templates/asyncapi/src/cli/cli.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "config", "config.go"), "templates/asyncapi/src/config/config.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "handler", "handler.go"), "templates/asyncapi/src/handler/handler.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "logger", "logger.go"), "templates/asyncapi/src/logger/logger.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "model", "model.go"), "templates/asyncapi/src/model/model.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "schemas", "schema.go"), "templates/asyncapi/src/schemas/schema.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "utils", "utils.go"), "templates/asyncapi/src/utils/utils.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "utils", "common.go"), "templates/asyncapi/src/utils/common.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "tracing", "tracing.go"), "templates/asyncapi/src/tracing/tracing.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "warp_server", "server.go"), "templates/asyncapi/src/warp_server/server.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "policy", "policy.go"), "templates/asyncapi/src/policy/policy.go.tmpl", conf)

	log.Info().Msg("Created AsyncAPI service files successfully")
	return nil
}
