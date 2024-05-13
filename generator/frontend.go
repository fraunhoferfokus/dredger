package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func generateEmptyFrontend(_ *openapi3.T, conf GeneratorConfig) {
	frontendPath := filepath.Join(conf.OutputPath, "web")
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/web/README.md.tmpl", nil)
}

func generateFrontend(_ *openapi3.T, conf GeneratorConfig) {
	generateOpenAPIDoc(conf)

	// create folders
	restPath := filepath.Join(conf.OutputPath, "rest")
	frontendPath := filepath.Join(conf.OutputPath, "web")
	javascriptPath := filepath.Join(frontendPath, "js")
	stylesheetPath := filepath.Join(frontendPath, "css")
	pagesPath := filepath.Join(frontendPath, "pages")
	localesPath := filepath.Join(pagesPath, "locales")
	publicPath := filepath.Join(frontendPath, "public")

	fs.GenerateFolder(frontendPath)
	fs.GenerateFolder(javascriptPath)
	fs.GenerateFolder(stylesheetPath)
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

	// files in stylesheet directory
	fs.CopyWebFile("web/css", stylesheetPath, "bootstrap-icons.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "bootstrap.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "pico.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "pico.colors.min.css", true)
	fs.CopyWebFile("web/css", stylesheetPath, "simple.min.css", true)

	// files in web directory
	fs.CopyWebFile("web", frontendPath, "web.go", true)

	// files in pages directory
	fs.CopyWebFile("web/pages", restPath, "render.go", true)
	createFileFromTemplate(filepath.Join(pagesPath, "localize.go"), "templates/web/pages/localize.go.tmpl", conf)
	if _, err := os.Stat(filepath.Join(pagesPath, "languages.templ")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(pagesPath, "languages.templ"), "templates/web/pages/languages.templ.tmpl", conf)
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
