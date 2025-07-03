package generator

import (
	//	rest "dredger/templates/openapi/web/pages"
	//	"encoding/json"
	//	"strings"
	//	"errors"
	"path"
	"path/filepath"

	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/rs/zerolog/log"
)

type Operation struct {
	OperationName string
	ChannelName   string
	MessageName   string
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
	tmplPath := path.Join("templates", "openapi", AsyncPkg, "publishers", "channel.go.tmpl")
	absPath, _ := filepath.Abs(tmplPath)
	log.Info().Msgf("Loading template from: %s", absPath)
	for _, c := range configs {
		destPath := path.Join(fpath, c.ChannelName+".go")
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

//Thought i needed them, but didn't  , but might come in handy at a later time

func ResolveMessageRef(ref string, spec *asyncapiv3.Specification) *asyncapiv3.Message {
	// Extrahiere den Namen aus dem $ref
	refName := path.Base(ref)

	// Schaue im Components-Block nach der Message
	msg, ok := spec.Components.Messages[refName]
	if !ok {
		log.Info().Str("ref", ref).Msg("Message not found in components")
		log.Info().Msg("message '" + refName + "' not found in components")
		return nil
	}

	return msg
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
