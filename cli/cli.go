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
	asyncPath    string
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
	Use:     "generate <path to Spec>",
	Short:   "Create server code from OpenAPI or AsyncAPI Spec",
	Long:    "Je nach übergebener Spec (OpenAPI bzw. AsyncAPI) wird der passende Generator aufgerufen.",
	Example: "  dredger generate ./stores.yaml -o ./outputPath -n StoresAPI",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		specPath := args[0]
		if projectPath == "" {
			projectPath = "src"
		}
		if projectName == "" {
			projectName = "build"
		}
		projectDestination := filepath.Join(projectPath)

		// Falls --async gesetzt, erzwinge AsyncAPI
		if asyncPath != "" {
			log.Info().Msg("AsyncAPI via --async übergeben.")
			spec, err := parser.ParseAsyncAPISpecFile(asyncPath)
			if err != nil {
				log.Error().Err(err).Msg("AsyncAPI: Fehler beim Parsen")
				return
			}
			if err := genAsyncAPI.GenerateService(spec, projectDestination, projectName); err != nil {
				log.Error().Err(err).Msg("AsyncAPI: Fehler beim Generieren")
			}
			return
		}

		// Spectype automatisch erkennen
		isAsync, isOpen, err := detectSpecType(specPath)
		if err != nil {
			log.Error().Err(err).Msg("Konnte Spec-Datei nicht öffnen oder lesen")
			return
		}
		switch {
		case isAsync:
			log.Info().Msg("Erkannt: AsyncAPI-Spec – wir parsen & generieren mit dem AsyncAPI-Generator")
			spec, err := parser.ParseAsyncAPISpecFile(specPath)
			if err != nil {
				log.Error().Err(err).Msg("AsyncAPI: Fehler beim Parsen")
				return
			}
			if err := genAsyncAPI.GenerateService(spec, projectDestination, projectName); err != nil {
				log.Error().Err(err).Msg("AsyncAPI: Fehler beim Generieren")
			}

		case isOpen:
			log.Info().Msg("Erkannt: OpenAPI-Spec – wir parsen & generieren mit dem OpenAPI-Generator")
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
			log.Error().Msg("Datei ist weder gültige AsyncAPI- noch gültige OpenAPI-Spec.")
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

	generateCmd.Flags().StringVarP(
		&asyncPath, "async", "a", "",
		"Pfad zur AsyncAPI-Spec (für AsyncAPI-Generator)",
	)
}

// detectSpecType liest bis 1 MiB und sucht nach asyncapi/openapi/swagger
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
	if strings.Contains(text, "\"swagger\"") || strings.HasPrefix(text, "swagger:") {
		return false, true, nil
	}
	return false, false, nil
}
