package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func generateEmptyFrontend(_ *openapi3.T, conf GeneratorConfig) {
	frontendPath := filepath.Join(conf.OutputPath, "web")
	fs.GenerateFolder(frontendPath)
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/web/README.md.tmpl", conf)
}

func generateFrontend(spec *openapi3.T, conf GeneratorConfig) {
	generateOpenAPIDoc(conf)

	if spec.Paths.Find("/events") == nil || (spec.Paths.Find("/events").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/events").Operations()[http.MethodGet].Tags, "builtin")) {
		log.Debug().Msg("Generating default /events endpoint.")

		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service support sse"))
		updateOAPIOperation(op, "HandleEvents", "support for sse", "200")
		spec.AddOperation("/events", http.MethodGet, op)
		spec.AddOperation("/events", http.MethodPost, op)
	}

	// create folders
	restPath := filepath.Join(conf.OutputPath, "rest")
	frontendPath := filepath.Join(conf.OutputPath, "web")
	javascriptPath := filepath.Join(frontendPath, "js")
	stylesheetPath := filepath.Join(frontendPath, "css")
	imagesPath := filepath.Join(frontendPath, "images")
	fontsPath := filepath.Join(stylesheetPath, "fonts")
	pagesPath := filepath.Join(frontendPath, "pages")
	localesPath := filepath.Join(pagesPath, "locales")
	publicPath := filepath.Join(frontendPath, "public")

	fs.GenerateFolder(frontendPath)
	fs.GenerateFolder(javascriptPath)
	fs.GenerateFolder(stylesheetPath)
	fs.GenerateFolder(imagesPath)
	fs.GenerateFolder(fontsPath)
	fs.GenerateFolder(pagesPath)
	fs.GenerateFolder(localesPath)
	fs.GenerateFolder(publicPath)

	// files in root directory
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/web/README.md.tmpl", conf)

	// files in javascript directory
	fs.CopyWebFile("web/js", javascriptPath, "bootstrap.bundle.min.js", true)
	fs.CopyWebFile("web/js", javascriptPath, "htmx-sse.js", true)
	fs.CopyWebFile("web/js", javascriptPath, "htmx.min.js", true)
	fs.CopyWebFile("web/js", javascriptPath, "hyperscript.js", true)
	fs.CopyWebFile("web/js", javascriptPath, "rapidoc-min.js", true)
	fs.CopyWebFile("web/js", javascriptPath, "sse.js", true)

	// files in stylesheet directory
	fs.CopyWebFile("web/css", stylesheetPath, "bootstrap-icons.min.css", true)
	fs.CopyWebFile("web/css/fonts", fontsPath, "bootstrap-icons.woff", true)
	fs.CopyWebFile("web/css/fonts", fontsPath, "bootstrap-icons.woff2", true)
	fs.CopyWebFile("web/css", stylesheetPath, "bootstrap.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "pico.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "pico.colors.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "simple.min.css", true)

	// files in web directory
	fs.CopyWebFile("web", frontendPath, "web.go", true)

	// files in pages directory
	fs.CopyWebFile("web/pages", restPath, "render.go", true)
	fs.CopyWebFile("web/pages", restPath, "progress.go", true)
	createFileFromTemplate(filepath.Join(pagesPath, "localize.go"), "templates/web/pages/localize.go.tmpl", conf)
	if _, err := os.Stat(filepath.Join(pagesPath, "languages.templ")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(pagesPath, "languages.templ"), "templates/web/pages/languages.templ.tmpl", conf)
	}
	if spec.Paths.Find("/index.html") != nil && spec.Paths.Find("/index.html").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/index.html").Operations()[http.MethodGet].Tags, "builtin") {
		if _, err := os.Stat(filepath.Join(pagesPath, "index.templ")); errors.Is(err, os.ErrNotExist) {
			createFileFromTemplate(filepath.Join(pagesPath, "index.templ"), "templates/web/pages/index.templ.tmpl", conf)
		}
		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service delivers index page"))
		updateOAPIOperation(op, "GetIndex", "successfully deliver index page", "200")
		spec.AddOperation("/index.html", http.MethodGet, op)
	}
	if spec.Paths.Find("/") != nil && spec.Paths.Find("/").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/").Operations()[http.MethodGet].Tags, "builtin") {
		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service delivers index page"))
		updateOAPIOperation(op, "GetRoot", "successfully deliver index page", "200")
		spec.AddOperation("/", http.MethodGet, op)
	}
	if spec.Paths.Find("/content.html") != nil && spec.Paths.Find("/content.html").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/content.html").Operations()[http.MethodGet].Tags, "builtin") {
		if _, err := os.Stat(filepath.Join(pagesPath, "content.templ")); errors.Is(err, os.ErrNotExist) {
			createFileFromTemplate(filepath.Join(pagesPath, "content.templ"), "templates/web/pages/content.templ.tmpl", conf)
		}
		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service delivers content page"))
		updateOAPIOperation(op, "GetContent", "successfully deliver content page", "200")
		spec.AddOperation("/content.html", http.MethodGet, op)
	}
	if _, err := os.Stat(filepath.Join(localesPath, "locale.de.toml")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(localesPath, "locale.de.toml"), "templates/web/pages/locales/locale.de.toml", conf)
		createFileFromTemplate(filepath.Join(localesPath, "locale.en.toml"), "templates/web/pages/locales/locale.en.toml", conf)
	}

	// files in public directory
	fs.CopyWebFile("web", publicPath, "README-public.md", false)

	log.Info().Msg("Created Frontend successfully.")
}

// function to get the port specified in the OpenAPI Spec
func getServerPort(spec *openapi3.T) (port int16) {
	if spec.Servers != nil {
		serverSpec := spec.Servers[0]
		if portSpec := serverSpec.Variables["port"]; portSpec != nil {
			portStr := portSpec.Default
			if portSpec.Enum != nil {
				portStr = portSpec.Enum[0]
			}

			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Warn().Msg("Failed to convert port, using 8080 instead.")
				return 8080
			} else {
				return int16(port)
			}
		} else {
			log.Warn().Msg("Failed to convert port, using 8080 instead.")
			return 8080
		}
	} else {
		log.Warn().Msg("Failed to convert port, using 8080 instead.")
		return 8080
	}
}

func createSchemas(spec *openapi3.T) (schemas Schemas) {
	schemas.List = make([]SchemaConf, 0)
	schemas.IsNotEmpty = false

	if spec != nil && spec.Components != nil && spec.Components.Schemas != nil {
		schemaStrings := toString(reflect.ValueOf(spec.Components.Schemas).MapKeys())

		for i := range schemaStrings {
			tmpSchemaName := schemaStrings[i]

			// check if schema has x-label == "form" -> if yes add schema to list
			schemaInformation, _ := spec.Components.Schemas[tmpSchemaName].Value.MarshalJSON()
			if strings.Contains(string(schemaInformation[:]), "\"x-label\":\"form\"") {
				var schema SchemaConf

				// add names
				schema.Name = strings.ReplaceAll(strings.ToLower(tmpSchemaName), " ", "")
				schema.H1Name = strings.Title(tmpSchemaName)
				schema.ComponentName = strings.ReplaceAll(schema.H1Name, " ", "")

				// add properties
				schema.Properties = make([]PropertyConf, 0)
				tmpSchemaPropertyNames := reflect.ValueOf(spec.Components.Schemas[tmpSchemaName].Value.Properties).MapKeys()
				for j := range tmpSchemaPropertyNames {
					tmpSchemaPropertyName := tmpSchemaPropertyNames[j].Interface().(string)
					var tmpPropertyConf PropertyConf
					tmpPropertyConf.Name = tmpSchemaPropertyName
					tmpPropertyConf.LabelName = strings.Title(tmpSchemaPropertyName)
					types := spec.Components.Schemas[tmpSchemaName].Value.Properties[tmpSchemaPropertyName].Value.Type.Slice()
					if len(types) > 0 {
						tmpPropertyConf.Type = types[0]
					}
					schema.Properties = append(schema.Properties, tmpPropertyConf)
				}

				schemas.List = append(schemas.List, schema)
				schemas.IsNotEmpty = true
			}

		}
	}

	return schemas

}

// function to convert an []reflect.Value to []string
func toString(inputArray []reflect.Value) (resultArray []string) {
	for i := range inputArray {
		resultArray = append(resultArray, inputArray[i].Interface().(string))
	}

	return resultArray
}

func generateOpenAPIDoc(conf GeneratorConfig) {
	// create folder
	type templateConfig struct {
		GeneratorConfig
		OpenAPIFile string
	}
	path := filepath.Join(conf.OutputPath, "web", "doc")
	fs.GenerateFolder(path)

	template := templateConfig{
		GeneratorConfig: conf,
		OpenAPIFile:     fs.GetFileNameWithEnding(conf.OpenAPIPath),
	}

	// create static html files
	createFileFromTemplate(filepath.Join(path, "index.html"), "templates/openapidoc/index.html.tmpl", template)

	// copy OpenAPI Specification in this directory
	fs.CopyFile(conf.OpenAPIPath, path, template.OpenAPIFile)

	log.Info().Msg("Created OpenAPI Documentation successfully.")
}
