package cli

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"dredger/core"

	genOpenAPI "dredger/generator/openapi"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Variablen für die Flags
var (
	projectPath  string
	projectName  string
	databaseFlag bool
	frontendFlag bool

	// (Optional) Flag, falls man explizit einen anderen AsyncAPI-Pfad angeben will.
	asyncPath string
)

// rootCmd repräsentiert den Basis-Befehl
var rootCmd = &cobra.Command{
	Use:   "dredger",
	Short: "Create server and client code from OpenAPI/AsyncAPI Spec",
	Long:  "Generate Go‐Server-Code (für OpenAPI) oder (demnächst) AsyncAPI-Code, je nachdem welche Spec man übergibt.",
}

var showVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the dredger tool",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		println("dredger v" + core.Version)
	},
}

var generateBdd = &cobra.Command{
	Use:   "generate-bdd <path to feature file>",
	Short: "Create BDD test file from the feature file (OpenAPI)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		genOpenAPI.GenerateBdd(args[0])
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

		// 1. Standardwerte für Output / Modul
		if projectPath == "" {
			projectPath = "src"
		}
		if projectName == "" {
			projectName = "build"
		}
		projectDestination := filepath.Join(projectPath)

		// 2. Wenn der Benutzer explizit --async gesetzt hat, ignorieren wir alles andere und zeigen nur eine Info.
		if asyncPath != "" {
			log.Info().Msg("AsyncAPI-Spec wurde via --async übergeben. AsyncAPI-Generator ist noch nicht implementiert.")
			return
		}

		// 3. Sonst ermitteln wir automatisch, ob specPath eine AsyncAPI oder OpenAPI ist:
		isAsync, isOpenAPI, err := detectSpecType(specPath)
		if err != nil {
			log.Error().Err(err).Msg("Konnte Spec-Datei nicht öffnen oder lesen")
			return
		}

		// 4. Je nach Ergebnis rufen wir den richtigen Zweig auf
		switch {
		case isAsync:
			log.Info().Msg("Erkannt: AsyncAPI-Spec – noch nicht implementiert")
			return

		case isOpenAPI:
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

// Execute hängt alle Unterbefehle an rootCmd und führt aus
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags für „generate“
	generateCmd.Flags().StringVarP(&projectPath, "output", "o", "src", "Pfad, in dem der Code erzeugt wird")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "default", "Modulname des erzeugten Codes")
	generateCmd.Flags().BoolVarP(&databaseFlag, "database", "D", false, "füge SQLite3-Datenbank in den generierten Code ein")
	generateCmd.Flags().BoolVarP(&frontendFlag, "frontend", "f", false, "füge Frontend-Code hinzu")

	// Optional: falls man wirklich eine AsyncAPI übergeben möchte
	generateCmd.Flags().StringVarP(
		&asyncPath,
		"async", "a", "",
		"Pfad zur AsyncAPI-Spec (falls man gezielt AsyncAPI generieren will)",
	)

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(generateBdd)
	rootCmd.AddCommand(showVersion)
}

// detectSpecType liest maximal ~1 MB der Datei ein und sucht nach den Schlüsselwörtern
// „asyncapi“, „openapi“ oder „swagger“, um zu entscheiden, um welchen Spec-Typ es sich handelt.
func detectSpecType(specPath string) (isAsync bool, isOpenAPI bool, err error) {
	f, err := os.Open(specPath)
	if err != nil {
		return false, false, err
	}
	defer f.Close()

	// Lies bis zu 1 MiB (falls Datei groß)
	buf := make([]byte, 1024*1024)
	n, err := io.ReadFull(f, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return false, false, err
	}
	buf = buf[:n]

	text := string(buf)
	lower := strings.ToLower(text)

	// AsyncAPI: prüfe auf „\"asyncapi\"“ (JSON) oder „asyncapi:“ (YAML)
	if strings.Contains(lower, "\"asyncapi\"") || strings.HasPrefix(lower, "asyncapi:") {
		return true, false, nil
	}
	// OpenAPI v3: prüfe auf „\"openapi\"“ (JSON) oder „openapi:“ (YAML)
	if strings.Contains(lower, "\"openapi\"") || strings.HasPrefix(lower, "openapi:") {
		return false, true, nil
	}
	// Swagger v2: prüfe auf „\"swagger\"“ (JSON) oder „swagger:“ (YAML)
	if strings.Contains(lower, "\"swagger\"") || strings.HasPrefix(lower, "swagger:") {
		return false, true, nil
	}

	return false, false, nil
}
