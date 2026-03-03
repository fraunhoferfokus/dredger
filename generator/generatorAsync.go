package generator

import (
	//	"embed"
	"path/filepath"
	"strconv"
	"strings"

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

	Config.Path = conf.OutputPath
	Config.Name = conf.ModuleName

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
	// Generating CORS configuration
	generateCorsConfig(conf)
	// Generating Tracing files
	generateTracing(conf)
	// Generating data base files
	if conf.AddDatabase {
		generateDatabaseFiles(conf)
	}
	serverConf := generateServerTemplateAsync(spec, conf)
	// Generating the config files
	generateConfigFiles(serverConf)
	generateJustfile(conf, serverConf)
	//generateReadme(conf, serverConf)
	//generateDockerfile(conf, serverConf)
	generateInfoFilesAsync(spec, serverConf)
	//Generating entities and structs
	GenerateAsyncTypes(spec, ProjectConfig{
		Name: conf.ModuleName,
		Path: conf.OutputPath,
	})

	log.Info().Msg("Created AsyncAPI NATS service files successfully")
	return nil
}

// generateServerTemplate gets all info for ServerConfig to be used by other functions
func generateServerTemplateAsync(spec *asyncapiv3.Specification, generatorConf GeneratorConfig) (serverConf ServerConfig) {
	asyncAPIName := fs.GetFileNameWithEnding(generatorConf.AsyncAPIPath)
	conf := ServerConfig{
		Port:        DefaultPort,
		ModuleName:  generatorConf.ModuleName,
		Flags:       generatorConf.Flags,
		OpenAPIName: asyncAPIName,
	}

	strDefaultPort := strconv.Itoa(DefaultPort)
	if spec.Servers != nil {
		for _, server := range spec.Servers {
			if server.Host != "" && server.Protocol == "nats" {
				portStr := server.Host
				parts := strings.Split(portStr, ":")
				if len(parts) == 2 {
					portStr = strings.TrimSpace(parts[1])
					if p, err := strconv.ParseInt(portStr, 10, 16); err == nil {
						conf.Port = int16(p)
					} else {
						log.Warn().Msg("Invalid port, using default " + strDefaultPort)
					}
				}
			} // if you extend for more protocols you could add more if cases: server.Protocol== "kafka"
		}

	} else {
		log.Warn().Msg("No servers field found, using default port " + strDefaultPort)
	}

	log.Info().Msg("Adding logging middleware.")
	return conf
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
	fs.GenerateFolder(filepath.Join(conf.OutputPath, AsyncPkg, "subscribers"))
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
