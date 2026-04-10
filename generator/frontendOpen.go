package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"net/http"
	"os"
	"path"
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
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/common/web/README.md.tmpl", conf)
}

func generateFrontend(spec *openapi3.T, conf GeneratorConfig) {
	generateOpenAPIDoc(conf)

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
	docPath := filepath.Join(frontendPath, "doc")

	fs.GenerateFolder(frontendPath)
	fs.GenerateFolder(javascriptPath)
	fs.GenerateFolder(stylesheetPath)
	fs.GenerateFolder(imagesPath)
	fs.GenerateFolder(fontsPath)
	fs.GenerateFolder(pagesPath)
	fs.GenerateFolder(localesPath)
	fs.GenerateFolder(publicPath)
	fs.GenerateFolder(docPath)

	// files in root directory
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/common/web/README.md.tmpl", conf)

	// files in javascript directory
	fs.CopyWebFile("common/web/js", javascriptPath, "bootstrap.bundle.min.js", true)
	fs.CopyWebFile("common/web/js", javascriptPath, "htmx.min.js", true)
	fs.CopyWebFile("common/web/js", javascriptPath, "hyperscript.js", true)
	fs.CopyWebFile("common/web/js", javascriptPath, "sse.js", true)
	fs.CopyWebFile("common/web/js", javascriptPath, "rapidoc-min.js", true)
	fs.CopyWebFile("common/web/js", javascriptPath, "elements.min.js", true)

	// files in stylesheet directory
	fs.CopyWebFile("common/web/css", stylesheetPath, "bootstrap-icons.min.css", true)
	fs.CopyWebFile("common/web/css/fonts", fontsPath, "bootstrap-icons.woff", true)
	fs.CopyWebFile("common/web/css/fonts", fontsPath, "bootstrap-icons.woff2", true)
	fs.CopyWebFile("common/web/css", stylesheetPath, "bootstrap.min.css", true)
	fs.CopyWebFile("common/web/css", stylesheetPath, "pico.min.css", true)
	fs.CopyWebFile("common/web/css", stylesheetPath, "pico.colors.min.css", true)
	fs.CopyWebFile("common/web/css", stylesheetPath, "elements.min.css", true)

	// files in images directory
	fs.CopyWebFile("common/web/images", imagesPath, "favicon.ico", false)

	// files in web directory
	fs.CopyWebFile("common/web", frontendPath, "web.go", true)

	// files in pages directory
	fs.CopyWebFile("common/web/pages", restPath, "render.go", true)
	if _, err := os.Stat(filepath.Join(pagesPath, "languages.templ")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(pagesPath, "languages.templ"), "templates/common/web/pages/languages.templ.tmpl", conf)
	}
	if spec.Paths.Find("/index.html") != nil && spec.Paths.Find("/index.html").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/index.html").Operations()[http.MethodGet].Tags, "builtin") {
		if _, err := os.Stat(filepath.Join(pagesPath, "index.templ")); errors.Is(err, os.ErrNotExist) {
			createFileFromTemplate(filepath.Join(pagesPath, "index.templ"), "templates/common/web/pages/index.templ.tmpl", conf)
			createFileFromTemplate(filepath.Join(pagesPath, "content.templ"), "templates/common/web/pages/content.templ.tmpl", conf)
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
			createFileFromTemplate(filepath.Join(pagesPath, "content.templ"), "templates/common/web/pages/content.templ.tmpl", conf)
		}
		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service delivers content page"))
		updateOAPIOperation(op, "GetContent", "successfully deliver content page", "200")
		spec.AddOperation("/content.html", http.MethodGet, op)
	}

	// files in public directory
	fs.CopyWebFile(path.Join("common/web", "public"), publicPath, "README.md", false)

	// files in doc directory
	fs.CopyWebFile(path.Join("common/web", "doc"), docPath, "README.md", false)

	// support for events
	if spec.Paths.Find("/events") != nil && spec.Paths.Find("/events").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/events").Operations()[http.MethodGet].Tags, "builtin") {
		//log.Debug().Msg("Generating default /events endpoint.")
		createFileFromTemplate(filepath.Join(restPath, "progress.go"), "templates/common/web/pages/progress.go.tmpl", conf)
		createFileFromTemplate(filepath.Join(restPath, "notice.go"), "templates/common/web/pages/notice.go.tmpl", conf)

		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service support sse"))
		updateOAPIOperation(op, "HandleEvents", "support for sse", "200")
		spec.AddOperation("/events", http.MethodGet, op)
		spec.AddOperation("/events", http.MethodPost, op)
	}

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

			port, err := strconv.ParseInt(portStr, 10, 16)
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
	docPath := filepath.Join(conf.OutputPath, "web", "doc")
	fs.GenerateFolder(docPath)

	template := templateConfig{
		GeneratorConfig: conf,
		OpenAPIFile:     fs.GetFileNameWithEnding(conf.OpenAPIPath),
	}

	// create static html files
	createFileFromTemplate(filepath.Join(docPath, "rapidoc.html"), "templates/openapi/rapidoc/index.html.tmpl", template)
	createFileFromTemplate(filepath.Join(docPath, "elements.html"), "templates/openapi/elements/index.html.tmpl", template)

	// copy OpenAPI Specification in this directory
	fs.CopyFile(conf.OpenAPIPath, docPath, template.OpenAPIFile)
	// add symlink to project root
	specPath := filepath.Join(docPath, template.OpenAPIFile)
	linkFilename := "OpenAPI" + path.Ext(template.OpenAPIFile) // static filename for project root
	linkPath := filepath.Join(Config.Path, linkFilename)
	if !fs.CheckIfFileExists(linkPath) { // skip it file (symlink) already exists
		if err := os.Symlink(specPath, linkPath); err != nil {
			log.Warn().Err(err).Str("source", specPath).Str("target", linkPath).Msg("Failed to create Symlink for OpenAPI specification file")
		}
	}

	log.Info().Msg("Created OpenAPI Documentation successfully.")
}
