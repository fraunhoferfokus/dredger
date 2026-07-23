package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	extCmd "dredger/cmd"
	"dredger/core"
	gen "dredger/generator"

	//genAsyncAPI "dredger/generator/asyncapi"
	"dredger/parser"

	"github.com/huandu/xstrings"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// variables for the flags
var (
	projectPath  string
	projectName  string
	databaseFlag bool
	frontendFlag bool
)

// Variables to know what needs to be in the main
var (
	openapi  bool
	asyncapi bool
)

// Just a temporary struct, to hold information about
var allOpenAPINames []gen.OpenAPIConfig

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

		// Reset global accumulators for this run
		gen.ResetGlobalPaths()

		specPaths := args

		for _, specPath := range specPaths {
			specPath = strings.TrimSpace(specPath)
			if specPath == "" || specPath == "\\" {
				// Ignore stray arguments from malformed line breaks
				continue
			}
			specType, err := detectSpecType(specPath)
			if err != nil {
				log.Error().Err(err).Msg("Konnte Spec-Datei nicht öffnen oder lesen")
				continue
			}

			switch specType {
			case asyncApiSpec:
				log.Info().Msgf("Erkannt: AsyncAPI-Spec %s – wir parsen & generieren", specPath)
				_, err := parser.ParseAsyncAPISpecFile(specPath)
				if err != nil {
					log.Error().Err(err).Msg("AsyncAPI: Fehler beim Parsen")
					continue
				}
				config := gen.GeneratorConfig{ //auch an den GenerateService / GenerateServer übergeben
					AsyncAPIPath: specPath,
					OutputPath:   projectDestination,
					ModuleName:   projectName,
					DatabaseName: "database",
					Flags: gen.Flags{
						AddDatabase: databaseFlag,
						AddFrontend: frontendFlag,
					},
				}
				if err := gen.GenerateAsyncService(config); err != nil {
					log.Error().Err(err).Msg("AsyncAPI: Fehler beim Generieren")
				}
				asyncapi = true
			case openApiSpec:
				log.Info().Msgf("Erkannt: OpenAPI-Spec %s – wir parsen & generieren", specPath)
				config := gen.GeneratorConfig{
					OpenAPIPath:  specPath,
					OutputPath:   projectDestination,
					ModuleName:   projectName,
					DatabaseName: "database",
					Flags: gen.Flags{
						AddDatabase: databaseFlag,
						AddFrontend: frontendFlag,
					},
				}
				if err := gen.GenerateServer(config); err != nil {
					log.Error().Err(err).Msg("OpenAPI: Fehler beim Generieren")
				}
				openapi = true
				allOpenAPINames = append(allOpenAPINames, gen.OpenAPIConfig{
					OpenAPIPath: specPath,
				})
			case componentSpec:
				log.Info().Msgf("Detected component spec %s - we parse and generate", specPath)
				_, err := parser.ParseComponentSpecFile(specPath)
				if err != nil {
					log.Error().Err(err).Msg("Component-Spec: not a valid spec")
					continue
				}
				config := gen.GeneratorConfig{
					OpenAPIPath:  specPath, // deal as openapi spec
					OutputPath:   projectDestination,
					ModuleName:   projectName,
					DatabaseName: "database",
					Flags: gen.Flags{
						AddDatabase: false,
						AddFrontend: false,
					},
				}
				// generate components (entities) only
				if err := gen.GenerateComponents(config); err != nil {
					log.Error().Err(err).Msg("Component-Spec: Error generating")
				}
			default:
				log.Error().Msgf("Datei %s ist weder gültige AsyncAPI- noch gültige OpenAPI-Spec.", specPath)
				// TODO Needs default case code for no spec given
				// GenerateDefault
			}

		}

		// IDEE: Array mit allen specPaths, welche OpenAPI sind, da sie für den OpenAPIName gebraucht werden, wenn es OpenAPI ist
		gen.GenerateMain(allOpenAPINames, projectDestination, projectName, openapi, asyncapi, databaseFlag, frontendFlag)

		// After all OpenAPI specs are processed, generate merged rest.go
		if openapi && len(allOpenAPINames) > 0 {
			// Use the first spec's path as the OpenAPIName for the merged rest.go
			openAPIName := allOpenAPINames[0].OpenAPIPath
			if fs := filepath.Base(openAPIName); fs != openAPIName {
				openAPIName = fs
			}
			gen.GenerateRestFromAccumulatedPaths(projectName, gen.Flags{
				AddDatabase: databaseFlag,
				AddFrontend: frontendFlag,
				OpenAPI:     openapi,
				AsyncAPI:    asyncapi,
			}, openAPIName)
		}

		// Create go.mod if not exist
		fileName := "go.mod"
		filePath := filepath.Join(projectDestination, fileName)
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			log.Info().Msg("RUN `go mod init " + xstrings.FirstRuneToLower(xstrings.ToCamelCase(projectName)) + "`")
			extCmd.RunCommand("go mod init "+xstrings.FirstRuneToLower(xstrings.ToCamelCase(projectName)), projectDestination)
		}

		log.Info().Msg("RUN `goimports`")
		extCmd.RunCommand("goimports -w .", projectDestination)

		log.Info().Msg("RUN `go mod tidy`")
		extCmd.RunCommand("go mod tidy", projectDestination)

		log.Info().Msg("RUN `go fmt`")
		extCmd.RunCommand("go fmt ./...", projectDestination)

		// TODO: create go files from tmpl `templ generate web/pages/*.templ`

		log.Info().Msg("DONE project created at: " + projectDestination)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(generateCmd, showVersion)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(showVersion)
	generateCmd.Flags().StringVarP(&projectPath, "output", "o", "src", "Pfad, in dem der Code erzeugt wird")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "default", "Modulname des erzeugten Codes")
	generateCmd.Flags().BoolVarP(&databaseFlag, "database", "D", false, "füge SQLite3-Datenbank in den generierten Code ein")
	generateCmd.Flags().BoolVarP(&frontendFlag, "frontend", "f", false, "füge Frontend-Code hinzu")

}
