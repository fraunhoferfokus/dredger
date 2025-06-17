package generator

import (
	"embed"
	"path/filepath"
	"strings"

	fs "dredger/fileUtils"
	parser "dredger/parser"

	//"dredger/parser"

	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
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
func GenerateService(spec Specification, outputPath, moduleName string) error {

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

func GenerateAsyncService(conf AsyncAPIConfig) error {
	spec, err := parser.ParseAsyncAPISpecFile(conf.AsyncAPIPath)
	if err != nil || spec == nil {
		log.Error().Err(err).Msg("Failed to load AsyncAPI spec file")
		return err
	}

	base := filepath.Join(conf.OutputPath, "src")
	subdirs := []string{"cmd", "internal", "logger", "tracing", "web"}
	fs.GenerateFolder(base)
	for _, d := range subdirs {
		fs.GenerateFolder(filepath.Join(base, d))
	}

	fs.GenerateFolder(filepath.Join(base, "cmd", "publisher"))
	fs.GenerateFolder(filepath.Join(base, "cmd", "server"))
	fs.GenerateFolder(filepath.Join(base, "internal", "config"))
	fs.GenerateFolder(filepath.Join(base, "internal", "server"))
	fs.GenerateFolder(filepath.Join(base, "internal", "server", "subscribers"))

	// Hauptdateien
	createFileFromTemplate(filepath.Join(conf.OutputPath, "README.md"), "templates/asyncapi/README.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(conf.OutputPath, "ENVIRONMENT.md"), "templates/asyncapi/ENVIRONMENT.md.tmpl", conf)
	createFileFromTemplate(filepath.Join(conf.OutputPath, "go.mod"), "templates/asyncapi/go.mod.tmpl", conf)

	// Publisher & Subscriber Channels generieren
	config := &AsyncGenConfig{}
	for _, ch := range spec.Channels {
		//var msgs = &ch.Messages
		for _, op := range spec.Operations {
			if op.Action == "send" {
				//				opName := GetMessageNameFromRef(op)
				config = &AsyncGenConfig{
					ModuleName:  conf.ModuleName,
					ChannelName: ch.Name,
					Description: spec.Info.Description,
					//					MessageName: opName,
					MessageName: GetAMessage(spec.Components.Messages),
					Action:      string(op.Action),
					Channels:    GetChannels(spec.Channels),
				}
			}

			if op.Action == "receive" {
				//opName := GetMessageNameFromRef(op)
				config = &AsyncGenConfig{
					ModuleName:  conf.ModuleName,
					ChannelName: ch.Name,
					Description: spec.Info.Description,
					MessageName: GetAMessage(spec.Components.Messages),
					Action:      string(op.Action),
					Channels:    GetChannels(spec.Channels),
				}
			}
		}
	}

	createFileFromTemplate(filepath.Join(base, "cmd", "publisher", "channel.go"), "templates/asyncapi/src/cmd/publisher/channel.go.tmpl", config)
	createFileFromTemplate(filepath.Join(base, "cmd", "server", "main.go"), "templates/asyncapi/src/cmd/server/main.go.tmpl", config)
	//	createFileFromTemplate(filepath.Join(base, "internal", "structs", "envelope.go"), "templates/asyncapi/src/internal/structs/envelope.go.tmpl", config)
	//	createFileFromTemplate(filepath.Join(base, "internal", "structs", "message.go"), "templates/asyncapi/src/internal/structs/message.go.tmpl", config)
	createFileFromTemplate(filepath.Join(base, "internal", "server", "subscribers", "channel.go"), "templates/asyncapi/src/internal/server/subscribers/channel.go.tmpl", config)
	createFileFromTemplate(filepath.Join(base, "internal", "server", "mainSvc.go"), "templates/asyncapi/src/internal/server/mainSvc.go.tmpl", config)
	createFileFromTemplate(filepath.Join(base, "logger", "logger.go"), "templates/asyncapi/src/logger/logger.go.tmpl", config)
	createFileFromTemplate(filepath.Join(base, "tracing", "tracing.go"), "templates/asyncapi/src/tracing/tracing.go.tmpl", config)
	createFileFromTemplate(filepath.Join(base, "web", "index.html"), "templates/asyncapi/src/web/index.html.tmpl", config)

	log.Info().Msg("Created AsyncAPI NATS service files successfully")
	return nil
}

func GetMessages(message map[string]*asyncapiv3.Message) []string {
	messages := make([]string, 0, len(message))
	for name := range message {
		messages = append(messages, name)
	}
	return messages
}

func GetAMessage(message map[string]*asyncapiv3.Message) string {
	for key := range message {
		return key
	}
	return ""
}

func GetChannels(channel map[string]*asyncapiv3.Channel) []string {
	channels := make([]string, 0, len(channel))
	for name := range channel {
		channels = append(channels, name)
	}
	return channels
}

func GetMessageFromChannel(msgs map[string]*asyncapiv3.Message) string {
	for _, msg := range msgs {
		return msg.Name
	}
	return ""
}

func GetMessage(op *asyncapiv3.Operation) string {
	if len(op.Messages) == 1 {
		return op.Messages[0].Name // TODO: change
	}

	return op.Channel.Name
}

func GetMessageNameFromRef(op *asyncapiv3.Operation) string {
	if op.Reference == "" {
		return "Default"
	}
	paths := strings.Split(op.Reference, "/")
	return paths[len(paths)-1]
}
