package cli

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"dredger/core"
	genAsyncAPI "dredger/generator/asyncapi"
	genOpenAPI "dredger/generator/openapi"
	"dredger/parser"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Variablen für die Flags
var (
	projectPath  string
	projectName  string
	databaseFlag bool
	frontendFlag bool
)

// rootCmd repräsentiert den Basis-Befehl
var rootCmd = &cobra.Command{
	Use:   "dredger",
	Short: "Create server and client code from OpenAPI/AsyncAPI Spec",
	Long:  "Generate Go‐Server‐Code (für OpenAPI) oder AsyncAPI‐Code, je nachdem welche Spec man übergibt.",
}

var showVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the dredger tool",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		println("dredger v" + core.Version)
	},
}

var generateCmd = &cobra.Command{
	Use:     "generate <path to Spec> [more specs...]",
	Short:   "Create server code from OpenAPI or AsyncAPI Spec",
	Long:    "Je nach übergebener Spec (OpenAPI bzw. AsyncAPI) wird der passende Generator aufgerufen.",
	Example: "  dredger generate api.yaml async.yaml moreasync.yaml -o ./out -n multi",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if projectPath == "" {
			projectPath = "src"
		}
		if projectName == "" {
			projectName = "build"
		}
		projectDestination := filepath.Join(projectPath)

		specPaths := args

		for _, specPath := range specPaths {
			specPath = strings.TrimSpace(specPath)
			if specPath == "" || specPath == "\\" {
				// Ignore stray arguments from malformed line breaks
				continue
			}
			isAsync, isOpen, err := detectSpecType(specPath)
			if err != nil {
				log.Error().Err(err).Msg("Konnte Spec-Datei nicht öffnen oder lesen")
				continue
			}

			switch {
			case isAsync:
				log.Info().Msgf("Erkannt: AsyncAPI-Spec %s – wir parsen & generieren", specPath)
				spec, err := parser.ParseAsyncAPISpecFile(specPath)
				if err != nil {
					log.Error().Err(err).Msg("AsyncAPI: Fehler beim Parsen")
					continue
				}
				if err := genAsyncAPI.GenerateService(spec, projectDestination, projectName); err != nil {
					log.Error().Err(err).Msg("AsyncAPI: Fehler beim Generieren")
				}

			case isOpen:
				log.Info().Msgf("Erkannt: OpenAPI-Spec %s – wir parsen & generieren", specPath)
				config := genOpenAPI.GeneratorConfig{
					OpenAPIPath:  specPath,
					OutputPath:   projectDestination,
					ModuleName:   projectName,
					DatabaseName: "database",
					Flags: genOpenAPI.Flags{
						AddDatabase: databaseFlag,
						AddFrontend: frontendFlag,
					},
				}
				if err := genOpenAPI.GenerateServer(config); err != nil {
					log.Error().Err(err).Msg("OpenAPI: Fehler beim Generieren")
				}

			default:
				log.Error().Msgf("Datei %s ist weder gültige AsyncAPI- noch gültige OpenAPI-Spec.", specPath)
			}
		}
	},
}

func Execute() {
	rootCmd.AddCommand(generateCmd, showVersion)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	generateCmd.Flags().StringVarP(&projectPath, "output", "o", "src", "Pfad, in dem der Code erzeugt wird")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "default", "Modulname des erzeugten Codes")
	generateCmd.Flags().BoolVarP(&databaseFlag, "database", "D", false, "füge SQLite3-Datenbank in den generierten Code ein")
	generateCmd.Flags().BoolVarP(&frontendFlag, "frontend", "f", false, "füge Frontend-Code hinzu")

}

// automatische spec erkennung 💃
func detectSpecType(specPath string) (isAsync bool, isOpenAPI bool, err error) {
	f, err := os.Open(specPath)
	if err != nil {
		return false, false, err
	}
	defer f.Close()

	buf := make([]byte, 1024*1024)
	n, err := io.ReadFull(f, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return false, false, err
	}
	text := strings.ToLower(string(buf[:n]))

	if strings.Contains(text, "\"asyncapi\"") || strings.HasPrefix(text, "asyncapi:") {
		return true, false, nil
	}
	if strings.Contains(text, "\"openapi\"") || strings.HasPrefix(text, "openapi:") {
		return false, true, nil
	}
	//veraltete "schreibweise" jetzt openapi
	if strings.Contains(text, "\"swagger\"") || strings.HasPrefix(text, "swagger:") {
		return false, true, nil
	}
	return false, false, nil
}
