package generator

import (
	"bufio"
	"embed"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	//"dredger/generator"
	"path/filepath"
	"strconv"

	fs "dredger/fileUtils"
	"dredger/parser"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

const (
	Cmd               = "cmd"
	Public            = "public"
	CorePkg           = "core"
	LogPkg            = "core/log"
	LoggerPkg         = "core/log/logger"
	TracingPkg        = "core/tracing"
	PagesPkg          = "web/pages"
	RestPkg           = "rest"
	DatabasePkg       = "db"
	EntitiesPkg       = "entities"
	UsecasesPkg       = "usecases"
	AuthzPkg          = "rest/middleware"
	MiddlewarePackage = "rest/middleware"
	DefaultPort       = 8080
)

var (
	config ProjectConfig
	TmplFS embed.FS
)

func GenerateServer(conf GeneratorConfig) error {
	spec, err := parser.ParseOpenAPISpecFile(conf.OpenAPIPath)
	if err != nil || spec == nil {
		log.Error().Err(err).Msg("Failed to load OpenAPI spec file")
		return err
	}

	// Init project config
	config.Name = conf.ModuleName
	config.Path = conf.OutputPath

	// config apikey
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

	GenerateTypes(spec, config)

	if conf.AddDatabase {
		generateDatabaseFiles(conf)
	}

	if conf.AddAuth {
		generateAuthzFile(conf)
	}

	generateValidation(conf)
	generatePolicy(conf)

	generateJustfile(conf, serverConf)

	generateDockerfile(conf, serverConf)

	log.Info().Msg("Created all files successfully.")

	return nil
}

func createProjectPathDirectory(conf GeneratorConfig) {
	// Generates basic folder structure
	fs.GenerateFolder(config.Path)
	fs.GenerateFolder(filepath.Join(config.Path, CorePkg))
	fs.GenerateFolder(filepath.Join(config.Path, RestPkg))
	fs.GenerateFolder(filepath.Join(config.Path, EntitiesPkg))
	fs.GenerateFolder(filepath.Join(config.Path, UsecasesPkg))
	if conf.AddDatabase {
		fs.GenerateFolder(filepath.Join(config.Path, DatabasePkg))
	}
	if conf.AddAuth {
		fs.GenerateFolder(filepath.Join(config.Path, AuthzPkg))
	}
	fs.GenerateFolder(filepath.Join(config.Path, MiddlewarePackage))

	log.Info().Msg("Created project directory.")
}

func generateServerTemplate(spec *openapi3.T, generatorConf GeneratorConfig) (serverConf ServerConfig) {
	openAPIName := fs.GetFileName(generatorConf.OpenAPIPath)
	conf := ServerConfig{
		Port:        DefaultPort,
		ModuleName:  generatorConf.ModuleName,
		Flags:       generatorConf.Flags,
		OpenAPIName: openAPIName,
	}

	strDefaultPort := strconv.Itoa(DefaultPort)

	if spec.Servers != nil {
		serverSpec := spec.Servers[0]
		if portSpec := serverSpec.Variables["port"]; portSpec != nil {
			portStr := portSpec.Default
			if portSpec.Enum != nil {
				portStr = portSpec.Enum[0]
			}

			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Warn().Msg("Failed to convert port, using" + strDefaultPort + "instead.")
			} else {
				conf.Port = int16(port)
			}
		} else {
			log.Warn().Msg("No port field was found, using" + strDefaultPort + "instead.")
		}
	} else {
		log.Warn().Msg("No servers field was found, using port " + strDefaultPort + " instead.")
	}

	log.Info().Msg("Adding logging middleware.")

	fileName := "main.go"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/main.go.tmpl"

	log.Debug().Msg("Creating server at port " + strconv.Itoa(int(conf.Port)) + "...")
	createFileFromTemplate(filePath, templateFile, conf)

	return conf
}

// ----------------------------
func ignore(input string) bool {
	return input == "When" || input == "And" || input == "Given" || input == "Then"
}

func retrieveRegex(input string) string {
	regex := ""
	for _, j := range input {
		if string(j) == "{" {
			regex = regex + "\\\\" + string(j)
		} else {
			regex = regex + string(j)
		}
	}
	return regex
}

func parseSteps(path string) []Step {
	//We use this map to connect each Step struct which represents a step in godog, to the answer that it requires
	m := make(map[string]int)
	var listOfSteps []Step
	file, err := os.Open(path)
	if err != nil {
		log.Fatal()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var prevStepConf Step
	regexedPath := ""

	for scanner.Scan() {
		var stepConf Step
		stepConf.Mapping = make(map[string]int)
		stepConf.RegexAndCode = make(map[string]int)
		line := scanner.Text()
		words := strings.Fields(line)
		stringRegex := "\"([^\"]*)\""
		for i, word := range words {
			//We skip this since we don't need to create any method for them
			if word == "Scenario:" && i == 0 {
				for _, j := range words[i:] {
					if ok, _ := regexp.MatchString(stringRegex, j); ok {
						regexedPath = retrieveRegex(j)
					}
				}
				break
			}
			if word == "Feature:" && i == 0 {
				break
			} else {
				//Retrieve the Method being used
				if word == http.MethodPut || word == http.MethodGet || word == http.MethodPost || word == http.MethodDelete {
					stepConf.Name = stepConf.Name + word
					word = strings.ToLower(word)
					r := []rune(word)
					stepConf.Method = string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
				} else if i >= 1 && (words[i-1] == "to") {
					//After "to" in the predefined structure we always receive the endpoint
					for _, word := range words[i:] {
						if value, _ := regexp.MatchString(stringRegex, word); value {
							stepConf.Endpoint = stepConf.Endpoint + word
							break
						}
					}
				} else if n, err := strconv.Atoi(word); err == nil && 200 <= n && n <= 600 {
					//Retrieve the status code
					stepConf.StatusCode = word
				} else if i >= 1 && (words[i-1] == "payload" || words[i-1] == "Payload" || words[i-1] == "PAYLOAD") {
					//After the word "payload" in the predefined structure comes the payload of the request
					for _, j := range words[i:] {
						stepConf.Payload = stepConf.Payload + j
					}
					//If we are in a "Then" part of the scenario, the next word is "the", there receive the status code
					//We also change the values in the previous step. We do this because the previous step, the actual
					//method is where we send the requests to the server whereas the "Then" part of the scenario
					//is just used to check whether the response was correct based on the request's method, url and payload
					if i == len(words)-1 && len(stepConf.StatusCode) == 0 {
						prevStepConf = stepConf
					}
					break
				} else {
					//Fill the name of the step. ignore() is used to ignore keywords such as When, Then, Given, And, since
					//they are not part of the name
					ignore := ignore(word)
					if len(stepConf.Name) == 0 && !ignore {
						//if the name is empty we add the word in it but firstly it has to be lower cased since the name uses
						//camel case convention
						stepConf.Name = stepConf.Name + strings.ToLower(word)
					} else if len(stepConf.Name) != 0 && !ignore {
						//if the name is not empty we capitalize the first letter of the word and add the word to the name
						//we also want to make sure that the name does not contain any substrings in quotes
						value, _ := regexp.MatchString(stringRegex, word)
						if !value {
							r := []rune(word)
							stepConf.Name = stepConf.Name + string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
						}
					}
				}
				//we receive the status code of the step in the "Then" part of the scenario. That's why we use prevStepConf.
				if strings.HasPrefix(stepConf.Name, "the") {
					code, _ := strconv.Atoi(stepConf.StatusCode)
					if code != 0 && len(regexedPath) != 0 {
						m[prevStepConf.Endpoint] = code
						prevStepConf.RegexPaths = append(prevStepConf.RegexPaths, regexedPath)
						prevStepConf.StatusCode = strconv.Itoa(code)
						value, _ := strconv.Atoi(prevStepConf.StatusCode)
						prevStepConf.Mapping[prevStepConf.Endpoint] = value
						prevStepConf.RegexAndCode[regexedPath] = value
					}
				}
				//Finally if we have arrived at the end of the string we are looking to add prevStepConf in our list of steps
				if i == len(words)-1 {
					flag := false
					code, _ := strconv.Atoi(stepConf.StatusCode)
					//Because the function handlers are not implemented, the payloads do not play an important role. Endpoints do.
					m[stepConf.Endpoint] = code
					if strings.HasPrefix(stepConf.Name, "the") && len(prevStepConf.Name) != 0 {
						//if the list of steps is empty then we just add the step
						if len(listOfSteps) == 0 {
							listOfSteps = append(listOfSteps, prevStepConf)
						} else {
							//otherwise we check whether there is a step with the same name or not
							for _, j := range listOfSteps {
								//if there is, then we want to add this step to the mapping of the step with the same name
								//thus we don't add redundant steps in the list, and also
								//it is exactly like the schema that godog requires in order to work
								if prevStepConf.Name == j.Name {
									flag = true
									value, _ := strconv.Atoi(prevStepConf.StatusCode)
									j.Mapping[prevStepConf.Endpoint] = value
									j.RegexAndCode[regexedPath] = value
									j.RegexPaths = append(j.RegexPaths, prevStepConf.RegexPaths[0])
								}
							}
							//in the case where the name is not there and also the list is not empty we use the flag to check
							//and we just append the step to the list
							if !flag {
								listOfSteps = append(listOfSteps, prevStepConf)
							}
						}
					}
				}
			}
			//after we're done with the initialization of the stepConf, we initialize prevStepConf and move on
			if i == len(words)-1 && len(stepConf.StatusCode) == 0 {
				prevStepConf = stepConf
			}
		}
	}
	//add the string that we want to use in the godog file
	for i, v := range listOfSteps {
		v.RealName = v.RealName + createName(v.Name)
		localhost := "http://localhost:8080"
		for _, k := range v.RegexPaths {
			in := strings.ReplaceAll(localhost+k, "\"", "")
			v.PathsWithHost = append(v.PathsWithHost, in)
		}
		listOfSteps[i] = v
	}
	return listOfSteps
}

// add space between words in camelcase
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func AddedSpace(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1} ${2}")
	return strings.ToLower(snake)
}

// we use this to create the string for the function which will be used in InitializeScenario() function
func createName(str string) string {
	nameToReturn := ""
	transformed := AddedSpace(str)
	for i, word := range strings.Fields(transformed) {
		if word == "i" && i == 0 {
			nameToReturn = nameToReturn + "^" + strings.ToUpper(word) + " "
		} else if word == "to" {
			nameToReturn = nameToReturn + "to \"([^\"]*)\" "
		} else if word == "payload" {
			nameToReturn = nameToReturn + "payload \"([^\"]*)\""
		} else if i == len(transformed)-1 {
			nameToReturn = nameToReturn + "$"
		} else if word == "put" || word == "get" || word == "post" || word == "delete" {
			nameToReturn = nameToReturn + strings.ToUpper(word) + " "
		} else {
			nameToReturn = nameToReturn + word + " "
		}
	}
	return nameToReturn
}

func contains(element string, arr []string) bool {
	for _, j := range arr {
		if element == j {
			return true
		}
	}
	return false
}

func getAllEndpoints(listing Listing) []string {
	var slice []string
	for _, k := range listing.Steps {
		for _, j := range k.RegexPaths {
			if !contains(j, slice) {
				slice = append(slice, j)
			}
		}
	}
	return slice
}
