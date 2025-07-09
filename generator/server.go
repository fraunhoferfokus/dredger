package generator

import (
	"errors"
	"path"

	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/rs/zerolog/log"
)

// Needes Info for Internal-file
type ChannelInfo struct {
	OperationName string
	ChannelName   string
}
type InternalConfig struct {
	Channels []ChannelInfo
}

// Needes Info for subscribers-file
type SubscribeOperation struct {
	ModuleName string
	Operations []Operation
}

// Generate the internal.go File in templates, needs: Channels, OperationName, ChannelName,
func GenerateInternalFile(spec *asyncapiv3.Specification, genConf GeneratorConfig) error {
	log.Debug().Msg("In GenerateInternalFile before checking spec")
	if spec == nil {
		err := errors.New("Could not generate internal-file. Specification not available.")
		log.Error().Err(err).Msg("")
		return err
	}
	log.Debug().Msg("After spec check before extracting SubInfos")
	channelInfos := extractSubInfo(spec)
	conf := InternalConfig{
		Channels: channelInfos,
	}
	log.Debug().Msg("After extracting SubInfos")
	filePath := path.Join(genConf.OutputPath, AsyncPkg, "server", "internal_"+spec.Info.Title+".go")
	tmplPath := path.Join("templates", "openapi", AsyncPkg, "server", "internal.go.tmpl")
	//filepath und tmplpath bestimmen und daraus dann die createFileFromTemplate(filepath, tmplPath und das c füllen)
	log.Debug().Msgf("Extracted %d channels", len(channelInfos))
	for i, ch := range channelInfos {
		log.Debug().Msgf("Channel %d: OperationName=%s, ChannelName=%s", i, ch.OperationName, ch.ChannelName)
	}

	createFileFromTemplate(filePath, tmplPath, conf)

	return nil
}

func GenerateSubscriberFile(spec *asyncapiv3.Specification, genConf GeneratorConfig) error {
	if spec == nil {
		err := errors.New("Could not generate subscribers-file. Specification not available.")
		log.Error().Err(err).Msg("")
		return err
	}

	operations := extractOperations(spec)
	conf := SubscribeOperation{
		ModuleName: genConf.ModuleName,
		Operations: operations,
	}

	filePath := path.Join(genConf.OutputPath, AsyncPkg, "server", "subscribers_"+spec.Info.Title+".go")
	tmplPath := path.Join("templates", "openapi", AsyncPkg, "server", "subscribers.go.tmpl")
	//filepath und tmplpath bestimmen und daraus dann die createFileFromTemplate(filepath, tmplPath und das c füllen)
	createFileFromTemplate(filePath, tmplPath, conf)
	return nil
}

func extractOperations(spec *asyncapiv3.Specification) []Operation {
	allOperations := []Operation{}
	for opName, op := range spec.Operations {
		if len(op.Messages) == 0 { // Default if no messages noted
			allOperations = append(allOperations, Operation{
				OperationName: "OperationName",
				ChannelName:   "ChannelName",
				MessageName:   "MessageName",
			})
		} else {
			if op.Action.IsReceive() {
				for _, msg := range op.Messages {
					allOperations = append(allOperations, Operation{
						OperationName: opName,
						ChannelName:   checkChannel(op.Channel.Reference),
						MessageName:   checkMessage(msg),
					})
				}
			}
		}
	}
	return allOperations
}

func checkMessage(message *asyncapiv3.Message) string {
	if message.Reference != "" {
		return path.Base(message.Reference)
	} else {
		log.Error().Msg("Message has to be a written as a ref when noted in operations messages.")
		return "IncorrectNotation"
	}
}

// Filters all Subscribe Operations and gets their OperationName and the ChannelName they belong to
// returns InternalConfig to be used for template
func extractSubInfo(spec *asyncapiv3.Specification) []ChannelInfo {
	log.Debug().Msg("In extractSubInfo before iterating over Operations in spec: " + spec.Info.Title)
	allSubChans := []ChannelInfo{}
	for opName, op := range spec.Operations {
		if op.Action.IsReceive() && op.Channel != nil {
			allSubChans = append(allSubChans,
				ChannelInfo{
					OperationName: opName,
					ChannelName:   checkChannel(op.Channel.Reference), //TODO get ChannelName not Channelstruct
				})
			log.Debug().Msg("Appended OperationName: " + opName + " and ChannelName: " + checkChannel(op.Channel.Reference))
		} else {
			log.Warn().Str("operation", opName).Msg("Missing channel reference")
		}
	}
	return allSubChans
}

// Might need for publish operations
func extractPubInfo(spec *asyncapiv3.Specification) []ChannelInfo {
	allPubChans := []ChannelInfo{}
	for opName, op := range spec.Operations {
		if op.Action.IsSend() && op.Channel != nil {
			allPubChans = append(allPubChans,
				ChannelInfo{
					OperationName: opName,
					ChannelName:   checkChannel(op.Channel.Reference), //TODO get ChannelName not Channelstruct
				})
		} else {
			log.Warn().Str("operation", opName).Msg("Missing channel reference")
		}
	}
	return allPubChans
}

// Either returns the ChannelName or logs the problem of having noted the channel in the operation in operations incorrectly
// and returning the DefaultName hinting at the incorrect notation so that Code can still be made from the template
func checkChannel(channelRef string) string {
	if channelRef != "" { // is actually a $ref
		return path.Base(channelRef)
	} else { // not ref, so wrong notation
		log.Error().Msg("Channel has to be a written as a ref when noted in operations channel.")
		return "IncorrectNotation"
	}
}
