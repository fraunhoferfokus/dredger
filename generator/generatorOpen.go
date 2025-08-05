package generator

import (
	"bufio"
	"embed"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	fs "dredger/fileUtils"
	oasparser "dredger/parser" // Alias für internes OpenAPI-Parser-Paket

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

const (
	AsyncPkg          = "async"
	Cmd               = "cmd"
	Public            = "public"
	CorePkg           = "core"
	LogPkg            = "core/log"
	LoggerPkg         = "core/log/logger"
	LokiPkg           = "core/log/loki"
	TracingPkg        = "core/tracing"
	PagesPkg          = "web/pages"
	RestPkg           = "rest"
	DatabasePkg       = "db"
	EntitiesPkg       = "entities"
	UsecasesPkg       = "usecases"
	MiddlewarePackage = "rest/middleware"
	DefaultPort       = 8080
)

var (
	// Einbettung der Template-Ordner common, openapi und asyncapi
	// Template files are embedded via main.go and assigned to TmplFS
	// This variable holds those embedded templates.
	TmplFS embed.FS

	Config ProjectConfig
)

// GenerateServer ist der Entry-Point für das OpenAPI-Scaffolding.
func GenerateServer(conf GeneratorConfig) error {
	spec := &openapi3.T{}
	var err error
	if conf.OpenAPIPath != "" {
		spec, err = oasparser.ParseOpenAPISpecFile(conf.OpenAPIPath)
		if err != nil || spec == nil {
			log.Error().Err(err).Msg("Failed to load OpenAPI spec file")
			return err
		}
	}

	// Initialisiere Projekt-Konfiguration
	Config.Name = conf.ModuleName
	Config.Path = conf.OutputPath

	// API‐Key‐Security erkennen
	if spec.Components != nil {
		for key, scheme := range spec.Components.SecuritySchemes {
			if scheme != nil && scheme.Value != nil && scheme.Value.Type == "apiKey" {
				conf.AddAuth = true
				conf.ApiKeyHeaderName = scheme.Value.Name
				conf.ApiKeySecurityName = key
				break
			}
		}
	}

	createProjectPathDirectory(conf)

	if conf.AddFrontend {
		generateFrontend(spec, conf)
	} else {
		generateEmptyFrontend(spec, conf)
	}

	serverConf := generateServerTemplate(spec, conf)

	generateConfigFiles(serverConf)
	generateInfoFiles(spec, serverConf)
	generateLogger(conf)
	generateTracing(conf)
	generateLifecycleFiles(spec, conf)
	generateHandlerFuncs(spec, conf)
	GenerateTypes(spec, Config)

	if conf.AddDatabase {
		generateDatabaseFiles(conf)
	}

	generateValidation(conf)
	generatePolicy(conf)
	generateJustfile(conf, serverConf)
	generateReadme(conf, serverConf)
	generateDockerfile(conf, serverConf)

	log.Info().Msg("Created all files successfully.")
	return nil
}

// createProjectPathDirectory legt die Grundordner an.
func createProjectPathDirectory(conf GeneratorConfig) {
	fs.GenerateFolder(Config.Path)
	fs.GenerateFolder(filepath.Join(Config.Path, CorePkg))
	fs.GenerateFolder(filepath.Join(Config.Path, RestPkg))
	fs.GenerateFolder(filepath.Join(Config.Path, EntitiesPkg))
	fs.GenerateFolder(filepath.Join(Config.Path, UsecasesPkg))
	if conf.AddDatabase {
		fs.GenerateFolder(filepath.Join(Config.Path, DatabasePkg))
	}
	fs.GenerateFolder(filepath.Join(Config.Path, MiddlewarePackage))
	log.Info().Msg("Created project directory.")
}

// generateServerTemplate gets all info for ServerConfig to be used by other functions
func generateServerTemplate(spec *openapi3.T, generatorConf GeneratorConfig) (serverConf ServerConfig) {
	openAPIName := fs.GetFileNameWithEnding(generatorConf.OpenAPIPath)
	conf := ServerConfig{
		Port:        DefaultPort,
		ModuleName:  generatorConf.ModuleName,
		Flags:       generatorConf.Flags,
		OpenAPIName: openAPIName,
	}

	strDefaultPort := strconv.Itoa(DefaultPort)
	if spec.Servers != nil {
		if portSpec := spec.Servers[0].Variables["port"]; portSpec != nil {
			portStr := portSpec.Default
			if portSpec.Enum != nil {
				portStr = portSpec.Enum[0]
			}
			if p, err := strconv.ParseInt(portStr, 10, 16); err == nil {
				conf.Port = int16(p)
			} else {
				log.Warn().Msg("Invalid port, using default " + strDefaultPort)
			}
		}
	} else {
		log.Warn().Msg("No servers field found, using default port " + strDefaultPort)
	}

	log.Info().Msg("Adding logging middleware.")
	return conf
}

// TODO: Make a GenerateMain-Function that is either making an async, open, both or none
// Problem: Certain information comes from the actual specs, so need to give that over as well
//
//	needs OpenAPIName if openapi
//
// Generates main and mainSvc files only once
func GenerateMain(openAPINames []OpenAPIConfig, outputPath string, moduleName string, openapi bool, asyncapi bool, dataBase bool, frontend bool) {
	if len(openAPINames) == 0 && !asyncapi && !openapi {
		log.Info().Msg("No specification was given. Default project will be made.")
	}
	mainConf := MainConfig{
		AllOpenAPINames: openAPINames,
		ModuleName:      moduleName,
		Flags: Flags{
			AddDatabase: dataBase,
			AddFrontend: frontend,
			OpenAPI:     openapi,
			AsyncAPI:    asyncapi,
		},
	}

	// main.go aus openapi/
	mainPath := filepath.Join(outputPath, "main.go")
	createFileFromTemplate(mainPath, "templates/common/main.go.tmpl", mainConf)

	// mainSvc.go aus openapi/
	svcPath := filepath.Join(outputPath, "mainSvc.go")
	if _, err := os.Stat(svcPath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(svcPath, "templates/common/mainSvc.go.tmpl", mainConf)
	}
	log.Info().Msg("Main was sucessfully generated.")
}

// ----------------------------
// parseSteps + Hilfsfunktionen

func ignore(input string) bool {
	return input == "When" || input == "And" || input == "Given" || input == "Then"
}

func retrieveRegex(input string) string {
	var regex string
	for _, ch := range input {
		if ch == '{' {
			regex += `\\{`
		} else {
			regex += string(ch)
		}
	}
	return regex
}

func parseSteps(path string) []Step {
	m := make(map[string]int)
	var listOfSteps []Step

	file, err := os.Open(path)
	if err != nil {
		log.Fatal().Msg("Could not open feature file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var prevStepConf Step
	var regexedPath string

	for scanner.Scan() {
		var stepConf Step
		stepConf.Mapping = make(map[string]int)
		stepConf.RegexAndCode = make(map[string]int)

		line := scanner.Text()
		words := strings.Fields(line)
		stringRegex := "\"([^\"]*)\""

		for i, word := range words {
			if word == "Scenario:" && i == 0 {
				for _, tok := range words[i:] {
					if ok, _ := regexp.MatchString(stringRegex, tok); ok {
						regexedPath = retrieveRegex(tok)
					}
				}
				break
			}
			if word == "Feature:" && i == 0 {
				break
			}

			switch {
			case word == http.MethodGet || word == http.MethodPost || word == http.MethodPut || word == http.MethodDelete:
				stepConf.Name += word
				lower := strings.ToLower(word)
				r := []rune(lower)
				stepConf.Method = string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))

			case i >= 1 && words[i-1] == "to":
				for _, tok := range words[i:] {
					if ok, _ := regexp.MatchString(stringRegex, tok); ok {
						stepConf.Endpoint += tok
						break
					}
				}

			case isStatusCode(word):
				stepConf.StatusCode = word

			case i >= 1 && (words[i-1] == "payload" || words[i-1] == "Payload" || words[i-1] == "PAYLOAD"):
				for _, tok := range words[i:] {
					stepConf.Payload += tok
				}
				if i == len(words)-1 && stepConf.StatusCode == "" {
					prevStepConf = stepConf
				}

			default:
				if !ignore(word) {
					if stepConf.Name == "" {
						stepConf.Name += strings.ToLower(word)
					} else {
						if ok, _ := regexp.MatchString(stringRegex, word); !ok {
							r := []rune(word)
							stepConf.Name += string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
						}
					}
				}
			}

			// Dann-Logik
			if strings.HasPrefix(stepConf.Name, "the") {
				if code, err := strconv.Atoi(stepConf.StatusCode); err == nil && code != 0 && regexedPath != "" {
					prevStepConf.RegexPaths = append(prevStepConf.RegexPaths, regexedPath)
					prevStepConf.StatusCode = strconv.Itoa(code)
					prevStepConf.Mapping[prevStepConf.Endpoint] = code
					prevStepConf.RegexAndCode[regexedPath] = code
				}
			}

			// Ende Zeile → Liste füllen
			if i == len(words)-1 {
				if code, err := strconv.Atoi(stepConf.StatusCode); err == nil {
					m[stepConf.Endpoint] = code
				}
				if strings.HasPrefix(stepConf.Name, "the") && prevStepConf.Name != "" {
					found := false
					for idx, existing := range listOfSteps {
						if existing.Name == prevStepConf.Name {
							found = true
							listOfSteps[idx].Mapping[prevStepConf.Endpoint], _ = strconv.Atoi(prevStepConf.StatusCode)
							listOfSteps[idx].RegexAndCode[regexedPath] = m[prevStepConf.Endpoint]
							listOfSteps[idx].RegexPaths = append(listOfSteps[idx].RegexPaths, prevStepConf.RegexPaths[0])
							break
						}
					}
					if !found {
						listOfSteps = append(listOfSteps, prevStepConf)
					}
				}
			}

			if i == len(words)-1 && stepConf.StatusCode == "" {
				prevStepConf = stepConf
			}
		}
	}

	// Abschließendes Aufbereiten
	for i, s := range listOfSteps {
		s.RealName = s.RealName + createName(s.Name)
		localhost := "http://localhost:8080"
		for _, p := range s.RegexPaths {
			s.PathsWithHost = append(s.PathsWithHost, strings.ReplaceAll(localhost+p, "\"", ""))
		}
		listOfSteps[i] = s
	}

	return listOfSteps
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// Adds a space, if it´s snake case
func AddedSpace(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	return matchAllCap.ReplaceAllString(snake, "${1} ${2}")
}

func createName(str string) string {
	words := strings.Fields(AddedSpace(str))
	result := ""
	for i, w := range words {
		switch {
		case w == "i" && i == 0:
			result += "^" + strings.ToUpper(w) + " "
		case w == "to":
			result += `to "([^"]*)" `
		case w == "payload":
			result += `payload "([^"]*)"`
		case i == len(words)-1:
			result += "$"
		case w == "put" || w == "get" || w == "post" || w == "delete":
			result += strings.ToUpper(w) + " "
		default:
			result += w + " "
		}
	}
	return result
}

func contains(elem string, arr []string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func getAllEndpoints(listing Listing) []string {
	var eps []string
	for _, s := range listing.Steps {
		for _, p := range s.RegexPaths {
			if !contains(p, eps) {
				eps = append(eps, p)
			}
		}
	}
	return eps
}

// isStatusCode checks if the provided string represents an HTTP status code
// between 200 and 600.
func isStatusCode(s string) bool {
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return n >= http.StatusOK && n <= 600
}
