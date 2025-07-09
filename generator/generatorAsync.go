package generator

import (
	//	"embed"
	"path/filepath"
	//"strings"

	fs "dredger/fileUtils"
	parser "dredger/parser"

	//"dredger/parser"

	//	open "dredger/generator/openapi"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/rs/zerolog/log"
)

// TmplFS holds the embedded async templates. It is assigned by main.go.
//var TmplFS embed.FS

//var config AsyncProjectConfig

type AsyncProjectConfig struct {
	Name string
	Path string
}

// mainConfig contains values passed to the templates.
type mainConfig struct {
	ModuleName string
	Title      string
}

// GenerateService creates a minimal async service based on the given spec.
func GenerateService(spec *asyncapiv3.Specification, outputPath, moduleName string) error {

	conf := mainConfig{
		ModuleName: moduleName,
		Title:      spec.Info.Title,
	}

	fs.GenerateFolder(outputPath)
	base := filepath.Join(outputPath, "src")
	subdirs := []string{"cmd", "internal", "logger", "tracing", "web"}
	fs.GenerateFolder(base)
	for _, d := range subdirs {
		fs.GenerateFolder(filepath.Join(base, d))
	}
	fs.GenerateFolder(filepath.Join(base, "cmd", "publisher"))
	fs.GenerateFolder(filepath.Join(base, "cmd", "server"))
	fs.GenerateFolder(filepath.Join(base, "internal", "config"))
	fs.GenerateFolder(filepath.Join(base, "internal", "server"))
	fs.GenerateFolder(filepath.Join(base, "server", "subscribers"))

	createFileFromTemplate(filepath.Join(outputPath, "README.md"), "templates/asyncapi/README.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(outputPath, "ENVIRONMENT.md"), "templates/asyncapi/ENVIRONMENT.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(outputPath, "go.mod"), "templates/asyncapi/go.mod.tmpl", conf)

	createFileFromTemplate(filepath.Join(base, "cmd", "publisher", "channel.go"), "templates/asyncapi/src/cmd/publisher/channel.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "cmd", "server", "main.go"), "templates/asyncapi/src/cmd/server/main.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "internal", "structs", "envelope.go"), "templates/asyncapi/src/internal/structs/envelope.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "internal", "structs", "message.go"), "templates/asyncapi/src/internal/structs/message.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "internal", "server", "subscribers", "channel.go"), "templates/asyncapi/src/internal/server/subscribers/channel.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "internal", "server", "mainSvc.go"), "templates/asyncapi/src/internal/server/mainSvc.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "logger", "logger.go"), "templates/asyncapi/src/logger/logger.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "tracing", "tracing.go"), "templates/asyncapi/src/tracing/tracing.go.tmpl", conf)
	createFileFromTemplate(filepath.Join(base, "web", "index.html"), "templates/asyncapi/src/web/index.html.tmpl", conf)

	log.Info().Msg("Created AsyncAPI service files successfully")
	return nil
}

func GenerateAsyncService(conf GeneratorConfig) error {
	spec := &asyncapiv3.Specification{}
	var err error
	//log.Debug().Msg("AsyncAPIPath: " + conf.AsyncAPIPath)
	if conf.AsyncAPIPath != "" {
		spec, err = parser.ParseAsyncAPISpecFile(conf.AsyncAPIPath)
		if err != nil || spec == nil {
			log.Error().Err(err).Msg("Failed to load AsyncAPI spec file")
			return err
		}
	}

	//Config.Path = conf.OutputPath
	//Config.Name = conf.ModuleName

	//log.Debug().Str("Spec info title to check if actually processed", spec.Info.Version).Msg("Check spec actually processed or not")
	createProjectPathDirectoryAsync(conf)

	if conf.AddFrontend {
		generateFrontendAsync(spec, conf)
	} else {
		generateEmptyFrontendAsync(spec, conf)
	}

	// Hauptdateien
	createFileFromTemplate(filepath.Join(conf.OutputPath, "README.md"), "templates/asyncapi/README.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(conf.OutputPath, "ENVIRONMENT.md"), "templates/asyncapi/ENVIRONMENT.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(conf.OutputPath, "go.mod"), "templates/asyncapi/go.mod.tmpl", conf)

	// Publisher Channels generieren
	GenerateChannelFile(spec, conf)
	// Internal Server files generieren
	GenerateInternalFile(spec, conf)
	// Subscriber files generieren
	GenerateSubscriberFile(spec, conf)
	// Generating logging files
	generateLogger(conf)
	// Generating data base files
	if conf.AddDatabase {
		generateDatabaseFiles(conf)
	}
	//Generating entities and structs
	GenerateAsyncTypes(spec, ProjectConfig{
		Name: conf.ModuleName,
		Path: conf.OutputPath,
	})

	log.Info().Msg("Created AsyncAPI NATS service files successfully")
	return nil
}

//TODO: Funktion that creates the files and then a function that checks if the info was actually there
// example Messages: get from Components all refs and path.Base the name from those
// Check with all Channels if those messages exist in a channel, otherwise error, log message unused or something

// createProjectPathDirectory legt die Grundordner bei Async an
func createProjectPathDirectoryAsync(conf GeneratorConfig) {
	fs.GenerateFolder(conf.OutputPath)
	//log.Debug().Msg("Created " + conf.OutputPath)
	fs.GenerateFolder(filepath.Join(conf.OutputPath, CorePkg))
	//log.Debug().Msg("Created " + CorePkg)
	fs.GenerateFolder(filepath.Join(conf.OutputPath, AsyncPkg))
	fs.GenerateFolder(filepath.Join(conf.OutputPath, AsyncPkg, "publishers"))
	fs.GenerateFolder(filepath.Join(conf.OutputPath, AsyncPkg, "server"))
	//log.Debug().Msg("Created " + AsyncPkg)
	fs.GenerateFolder(filepath.Join(conf.OutputPath, RestPkg))
	//log.Debug().Msg("Created " + RestPkg)
	fs.GenerateFolder(filepath.Join(conf.OutputPath, EntitiesPkg))
	//log.Debug().Msg("Created " + EntitiesPkg)
	fs.GenerateFolder(filepath.Join(conf.OutputPath, UsecasesPkg))
	//log.Debug().Msg("Created " + UsecasesPkg)
	if conf.AddDatabase {
		fs.GenerateFolder(filepath.Join(conf.OutputPath, DatabasePkg))
		//log.Debug().Msg("Created " + DatabasePkg)
	}
	fs.GenerateFolder(filepath.Join(conf.OutputPath, MiddlewarePackage))
	//log.Info().Msg("Created project directory.")
}
