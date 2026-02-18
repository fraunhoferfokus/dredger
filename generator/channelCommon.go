package generator

import (
	//	rest "dredger/templates/openapi/web/pages"
	//	"encoding/json"
	//	"strings"
	//	"errors"

	"path"
	"path/filepath"
	"strings"

	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/rs/zerolog/log"
)

type Operation struct {
	OperationName string
	ChannelName   string
	Messages      []Message
}

type Message struct {
	MessageName       string
	MessageStructName string
}

type GenConfig struct {
	ModuleName    string
	OperationName string
	ChannelName   string
	MessageName   []string
}

//IDEE: Nur die SendOperations, also braucht man erstmal eine Map von allen Send-Operations zu erstellen
// Nachdem alle Send-Operations in einer Map sind, kann man durch diese iterieren und die Info aus diesen
// SendOperations extrahieren in Form von "Send"Operation{ModuleName, OperationName, ChannelName}, welches
// dann an createFileFromTemplate ins Template gelangt

func GenerateChannelFile(spec *asyncapiv3.Specification, conf GeneratorConfig) {
	var sendOps = GetPublishChannelOperations(spec)
	configs := []GenConfig{}
	for _, op := range sendOps {
		configs = append(configs, GenConfig{
			ModuleName:    conf.ModuleName,
			OperationName: op.OperationName,
			ChannelName:   op.ChannelName,
		})
	}

	fpath := path.Join(conf.OutputPath, AsyncPkg, "publishers")
	tmplPath := path.Join("templates", "asyncapi", AsyncPkg, "publishers", "channel.go.tmpl")
	absPath, _ := filepath.Abs(tmplPath)
	log.Info().Msgf("Loading template from: %s", absPath)
	for _, c := range configs {
		destPath := path.Join(fpath, lcFirst(c.OperationName)+".go")
		createFileFromTemplate(destPath, tmplPath, c)
	}
	log.Info().Msg("Finished generating all files for Publishers folder")
}

// FIXME: immer diese Fehlermeldung: panic: template: pattern matches no files: `templates\openapi\async\publishers\channel.go.tmpl`

// Returns an Array of Operations from spec, that are only Send-Operations (from spec)
func GetPublishChannelOperations(spec *asyncapiv3.Specification) []Operation {
	var result []Operation
	for opName, op := range spec.Operations {
		if op.Action == "send" {
			result = append(result, Operation{
				OperationName: opName,
				ChannelName:   path.Base(op.Channel.Reference),
			})
		}
	}
	log.Info().Msg("Getting Send-Operations")
	return result
}

// Extracts message from spec to a given ref
func ResolveMessageRef(spec *asyncapiv3.Specification, ref string) *asyncapiv3.Message {
	const prefix = "#/"

	// Verify ref structure
	if !strings.HasPrefix(ref, prefix) {
		log.Error().Str("ref format", ref).Str("prefix", prefix).Msg("unsupported ref format - must be in prefix")
		return nil
	}

	// Allow messages from Components and Channels
	const channelPrefix = prefix + "channels/"
	const componentPrefix = prefix + "components/messages"
	refName := path.Base(ref)
	if strings.HasPrefix(ref, channelPrefix) {
		channelName := strings.Split(strings.Replace(ref, channelPrefix, "", 1), "/")[0]
		channel, ok := spec.Channels[channelName]
		if !ok {
			log.Error().Str("ref", ref).Str("channel name", channelName).Any("channels", spec.Channels).Msg("could not find reference in channels")
			return nil
		}
		msg, ok := channel.Messages[refName]
		if !ok {
			log.Error().Str("ref", ref).Msg("could not find reference in the channels messages")
			log.Debug().Str("ref", ref).Any("messages", channel.Messages).Any("result", msg).Msg("could not find reference in the channels messages")
		}
		return msg
	} else if strings.HasPrefix(ref, componentPrefix) {
		msg, ok := spec.Components.Messages[refName]
		if !ok {
			log.Error().Str("ref", ref).Msg("could not find reference in components messages")
			return nil
		}
		return msg
	}

	log.Error().Str("ref format", ref).Str("prefix", prefix).Msg("unsupported ref format - must be in prefix and either in either in components or channels")
	return nil
}

func ResolveChannelRef(ref string, spec *asyncapiv3.Specification) *asyncapiv3.Channel {
	refName := path.Base(ref)
	if refName != "" {
		ch, ok := spec.Channels[refName]
		if !ok {
			log.Info().Str("ref", ref).Msg("Channels not found in spec channels")
			log.Info().Msg("channel '" + refName + "' not found in spec channels")
			return nil
		}

		return ch
	}
	return nil
}
