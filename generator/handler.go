package generator

import (
	extCmd "dredger/cmd"
	fs "dredger/fileUtils"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/huandu/xstrings"
	"github.com/rs/zerolog/log"
)

func generateHandlerFuncStub(op *openapi3.Operation, method string, path string, genConf GeneratorConfig) (OperationConfig, error) {
	var conf OperationConfig

	var methodPath = method + " " + path

	if op.Security != nil {
		for _, item := range *op.Security {
			for key := range item {
				if key == genConf.ApiKeySecurityName {
					conf.AddAuth = true
					break
				}
			}
		}
	}

	conf.ModuleName = genConf.ModuleName

	conf.Method = method

	conf.Summary = op.Summary
	if op.Summary == "" {
		log.Warn().Msg("No summary found for endpoint: " + methodPath)
	}

	conf.OperationID = xstrings.FirstRuneToUpper(xstrings.ToCamelCase(op.OperationID))
	if op.OperationID == "" {
		log.Error().Msg("No operation ID found for endpoint: " + methodPath)
		return conf, errors.New("no operation id, can't create function")
	}
	conf.Schema = ""
	if op.RequestBody != nil {
		if mt := op.RequestBody.Value.Content.Get("application/json"); mt != nil {
			x := strings.Split(mt.Schema.Ref, "/")
			if len(x) > 0 {
				conf.Schema = x[len(x)-1]
			}
		} else if mt := op.RequestBody.Value.Content.Get("application/yaml"); mt != nil {
			x := strings.Split(mt.Schema.Ref, "/")
			if len(x) > 0 {
				conf.Schema = x[len(x)-1]
			}
		} else if mt := op.RequestBody.Value.Content.Get("application/xml"); mt != nil {
			x := strings.Split(mt.Schema.Ref, "/")
			if len(x) > 0 {
				conf.Schema = x[len(x)-1]
			}
		}
	} else if t, exists := op.Extensions["x-requestType"]; exists {
		if s, ok := t.(string); ok {
			conf.Schema = s
		}
	}

	hasHtmlResponse := false
	if op.Responses != nil {
		for _, resRef := range op.Responses.Map() {
			for cKey := range resRef.Value.Content {
				if cKey == "text/html" {
					hasHtmlResponse = true
				}
			}
		}
	}
	if genConf.AddFrontend && hasHtmlResponse && slices.Contains(op.Tags, "page") {
		fileName := xstrings.FirstRuneToLower(xstrings.ToCamelCase(conf.OperationID)) + ".templ"
		filePath := filepath.Join(Config.Path, PagesPkg, fileName)
		templateFile := "templates/openapi/web/pages.templ.tmpl"
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			createFileFromTemplate(filePath, templateFile, conf)
		}
	}

	for resKey, resRef := range op.Responses.Map() {
		if !validateStatusCode(resKey) && resKey != "default" {
			log.Warn().Msg("Status code " + resKey + " for endpoint " + methodPath + " is not a valid status code.")
		}

		conf.Responses = append(conf.Responses, ResponseConfig{resKey, *resRef.Value.Description})
	}

	conf.PathParams = map[string]string{}
	for _, p := range op.Parameters {
		if p != nil && p.Value.In == "path" {
			conf.PathParams[p.Value.Name] = p.Value.Schema.Value.Format
			if p.Value.Schema.Value.Format == "" {
				types := p.Value.Schema.Value.Type.Slice()
				if len(types) > 0 {
					conf.PathParams[p.Value.Name] = types[0]
				}
			}
		}
	}

	conf.QueryParams = map[string]string{}
	for _, p := range op.Parameters {
		if p != nil && p.Value.In == "query" {
			conf.QueryParams[p.Value.Name] = p.Value.Schema.Value.Format
			if p.Value.Schema.Value.Format == "" {
				types := p.Value.Schema.Value.Type.Slice()
				if len(types) > 0 {
					conf.QueryParams[p.Value.Name] = types[0]
				}
			}
		}
	}

	canBeEdited := true
	fileName := xstrings.FirstRuneToLower(xstrings.ToCamelCase(conf.OperationID)) + ".go"
	filePath := filepath.Join(Config.Path, RestPkg, fileName)
	templateFile := "templates/openapi/rest/handlerFunc.go.tmpl"
	if hasHtmlResponse && slices.Contains(op.Tags, "page") {
		templateFile = "templates/openapi/rest/pageHandlerFunc.go.tmpl"
	}
	if conf.OperationID == "GetLive" {
		canBeEdited = false
		templateFile = "templates/openapi/rest/getLive.go.tmpl"
	}
	if conf.OperationID == "GetInfo" {
		canBeEdited = false
		templateFile = "templates/openapi/rest/getInfo.go.tmpl"
	}
	if conf.OperationID == "GetRobots" {
		canBeEdited = false
		templateFile = "templates/openapi/rest/getRobots.go.tmpl"
	}
	if conf.OperationID == "GetIndex" {
		templateFile = "templates/openapi/rest/getIndex.go.tmpl"
	}
	if conf.OperationID == "GetRoot" {
		templateFile = "templates/openapi/rest/getRoot.go.tmpl"
	}
	if conf.OperationID == "GetContent" {
		templateFile = "templates/openapi/rest/getContent.go.tmpl"
	}
	if conf.OperationID == "HandleEvents" {
		canBeEdited = false
		templateFile = "templates/openapi/rest/handleEvents.go.tmpl"
	}

	log.Debug().Str("operation", conf.OperationID).Str("template", templateFile).Msg("Generate handler")
	if _, err := os.Stat(filePath); !canBeEdited || errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, conf)
	}
	// remove unused imports
	extCmd.RunCommand("goimports -w "+fileName, filepath.Join(Config.Path, RestPkg))
	return conf, nil
}

func generateHandlerFuncs(spec *openapi3.T, genConf GeneratorConfig) {
	conf := HandlerConfig{
		ModuleName:  genConf.ModuleName,
		OpenAPIPath: fs.GetFileNameWithEnding(genConf.OpenAPIPath),
		AddAuth:     genConf.AddAuth,
		Flags:       genConf.Flags,
	}
	conf.ModuleName = genConf.ModuleName
	conf.Flags = genConf.Flags

	for _, item := range spec.Security {
		for key := range item {
			if key == genConf.ApiKeySecurityName {
				conf.AddGlobalAuth = true
				break
			}
		}
	}

	for path, pathObj := range spec.Paths.Map() {
		var newPath PathConfig
		newPath.Path = convertPathParams(path)

		for method, op := range pathObj.Operations() {
			if !slices.Contains(op.Tags, "builtin") {
				opConfig, err := generateHandlerFuncStub(op, method, newPath.Path, genConf)

				if err != nil {
					log.Err(err).Msg("Skipping generation of handler function for endpoint " + method + " " + path)
				}

				newPath.Operations = append(newPath.Operations, opConfig)
			}
		}

		conf.Paths = append(conf.Paths, newPath)
	}

	fileName := "rest.go"
	filePath := filepath.Join(Config.Path, RestPkg, fileName)
	templateFile := "templates/openapi/rest/handler.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)

	fileName = "restSvc.go"
	filePath = filepath.Join(Config.Path, RestPkg, fileName)
	templateFile = "templates/openapi/rest/restSvc.go.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, conf)
	}

}
