package cli

import (
	extCmd "dredger/cmd"
	"dredger/core"
	gen "dredger/generator"
	"errors"
	"os"
	"path/filepath"

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

	asyncPath    string // NEU
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dredger",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and RapidDoc-Clientcode for your application by providing an OpenAPI Specification",
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
	Short: "Create BDD test file from the feature file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gen.GenerateBdd(args[0])
	},
}

var generateCmd = &cobra.Command{
	Use:     "generate <path to OpenAPI Specification>",
	Short:   "Create server and client API code from OpenApi Spec",
	Long:    "Generate Go-Server code and RapidDoc-Clientcode for your application by providing an OpenAPI Specification",
	Example: "generate ./stores.yaml -o ./outputPath -n StoresAPI",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		openAPIPath := args[0]

		// output project path
		if projectPath == "" {
			projectPath = "src"
		}

		// output project name
		if projectName == "" {
			projectName = "build"
		}

		projectDestination := filepath.Join(projectPath)
		config := gen.GeneratorConfig{
			OpenAPIPath:  openAPIPath,

			AsyncAPIPath: asyncPath,// NEU

			OutputPath:   projectDestination,
			ModuleName:   projectName,
			DatabaseName: "database",
			Flags: gen.Flags{
				AddDatabase: databaseFlag,
				AddFrontend: frontendFlag,
			},
		}

		log.Debug().Msg("Generating project...")
		err := gen.GenerateServer(config)

		if err != nil {
			log.Error().Msg("Aborting...")
			return
		}

		log.Info().Msg("Running external commands...")

		// Create go.mod if not exist
		fileName := "go.mod"
		filePath := filepath.Join(config.OutputPath, fileName)
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

		log.Info().Msg("DONE project created at: " + projectDestination)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// add generate flags
	generateCmd.Flags().StringVarP(&projectPath, "output", "o", "src", "path where generated code gets stored")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "default", "module name of generated code")
	// optional code generation
	generateCmd.Flags().BoolVarP(&databaseFlag, "database", "D", false, "add sqlite3 database in generated code")
	generateCmd.Flags().BoolVarP(&frontendFlag, "frontend", "f", false, "add frontend code")

	generateCmd.Flags().StringVarP(
		&asyncPath,        // Zeiger auf unsere neue Variable
		"async",           // Langform: --async
		"a",               // Kurzform : -a
		"",                // Default  : leer = kein AsyncAPI-File
		"path to AsyncAPI spec file",
    )

	// add generate command
	rootCmd.AddCommand(generateCmd)

	//add bdd command
	rootCmd.AddCommand(generateBdd)

	// add version command
	rootCmd.AddCommand(showVersion)
}
