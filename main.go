package main

import (
	//cli "dredger/cli"
	//generator "dredger/generator/openapi"

	"embed"
	"strings"

	//"os"

	"fmt"

	async "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/parser"
	v3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"

	//
	//"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates
var tmplFS embed.FS

func main() {
	// Set up zerolog time format
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// Export embed template filesystem to generator package
	//generator.TmplFS = tmplFS

	//cli.Execute()

	fileParams := async.FromFileParams{
		Path:         "./examples/simple/asyncapiv3.json",
		MajorVersion: 3,
	}
	spec, err := async.FromFile(fileParams)
	if err != nil {
		log.Err(err).Msg("Could not get spec from file")
	}

	specV3, ok := spec.(*v3.Specification)
	if !ok {
		log.Error().Msg("Returned spec is not of type v3.Specification")
		return
	}
	specV3.Process() //processes full spec

	fmt.Println("AsyncAPI Version:		", specV3.Version)
	fmt.Println("Channels:			", specV3.Channels)
	fmt.Println("Operations:			", specV3.Operations)
	fmt.Println("\nMessage payloads in components:")
	msg, ok := specV3.Components.Messages["UserSignedUp"]

	if !ok {
		log.Fatal().Msg("Message 'UserSignedUp' not found.")
	}

	schema := msg.Payload

	if schema.Type == "object" {
		for propName, prop := range schema.Properties {
			fmt.Println("Property:			", propName)
			fmt.Println("  Type:				", prop.Type)
			fmt.Println("  Description:			", prop.Description)
		}
	}

	fmt.Println("Operations:			", specV3.Operations)
	ops, ok := specV3.Operations["onUserSignUp"]
	if !ok {
		log.Fatal().Msg("Message 'UserSignedUp' not found.")
	}

	if ops.Action == "receive" {
		//whatever für nats subscribe erhalten von Event
		fmt.Println("This is a subscriber")
	}

	//Case wir wollen an die Refs
	//TODO ordentliches json/yaml finden, damit wir das probieren können
	if ops.Channel.Reference != "" {
		// Ref sieht so aus: "#/channels/userSignedUp"
		refParts := strings.Split(ops.Channel.Reference, "/")
		if len(refParts) == 3 && refParts[1] == "channels" {
			channelName := refParts[2]
			channel, ok := specV3.Channels[channelName]
			if !ok {
				log.Fatal().Msg("Channel not found: " + channelName)
			}
			fmt.Println("Channel name:			", channel.Name)
		}
	}
	//ops.Channel.Follow()
}
