package generator

import (
	"math"
	"path"
	"path/filepath"
	"strings"

	//	open "dredger/generator"
	"github.com/gobeam/stringy"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

type TypeDefinitionAsync struct {
	Name        string
	Type        string
	MinLength   uint
	MaxLength   uint
	Pattern     string
	Minimum     float64
	Maximum     float64
	MarshalName string
	NestedTypes []TypeDefinitionAsync
}

type ModelConfigAsync struct {
	Imports     ImportsConfig
	SchemaDefs  map[string][]TypeDefinitionAsync
	ProjectName string
}

func GenerateAsyncTypes(spec *asyncapiv3.Specification, pConf ProjectConfig) {
	if spec != nil && len(spec.Components.Messages) != 0 {
		schemaDefs := generateAsyncTypeDefs(spec.Components.Schemas)
		imports := generateAsyncImports()
		var conf ModelConfigAsync
		conf.Imports = imports
		conf.ProjectName = pConf.Name

		for schema, defs := range schemaDefs {
			conf.SchemaDefs = map[string][]TypeDefinitionAsync{schema: defs}
			fileName := strings.ToLower(schema) + ".go"
			filePath := filepath.Join(pConf.Path, EntitiesPkg, fileName)
			templateFiles := []string{
				"templates/openapi/entities/entities.go.tmpl",
				"templates/openapi/entities/imports.tmpl",
				"templates/openapi/entities/structs.tmpl",
				"templates/openapi/entities/validate.tmpl",
			}
			createFileFromTemplates(filePath, templateFiles, conf)
		}
	}
}

func generateAsyncTypeDefs(schemas map[string]*asyncapiv3.Schema) map[string][]TypeDefinitionAsync {
	schemaDefs := make(map[string][]TypeDefinitionAsync, len(schemas))
	for name, ref := range schemas {
		var goType string
		switch {

		case ref.Type == "number":
			switch ref.Format {
			case "float":
				goType = "float32"
			case "double":
				goType = "float64"
			default:
				goType = "float"
			}
		case ref.Type == "integer":
			goType = "int"
			if ref.Format != "" {
				goType = ref.Format
			}
		case ref.Type == "boolean":
			goType = "bool"
		case ref.Type == "string":
			switch ref.Format {
			case "binary":
				goType = "[]byte"
			case "date":
				IMPORT_TIME = true
				goType = "time.Time"
			case "uuid":
				IMPORT_UUID = true
				goType = "uuid.UUID"
			default:
				goType = "string"
			}
		case ref.Type == "array":
			items, _ := toGoTypeAsync(ref.Items)
			goType = "[]" + items
		case ref.Type == "object":
			if ref.AdditionalProperties != nil {
				goType, _ = toGoTypeAsync(ref.AdditionalProperties)
			} else {
				goType = "map[string]interface{}"
			}
		default:
			types := ref.Type
			if len(types) > 0 {
				goType = types
			}
		}

		schemaDefs[name] = []TypeDefinitionAsync{{
			Name:        name,
			Type:        goType,
			MinLength:   ref.MinLength,
			MaxLength:   uintOrMaxA(&ref.MaxLength),
			Pattern:     ref.Pattern,
			Minimum:     floatOrMinA(&ref.Minimum),
			Maximum:     floatOrMaxA(&ref.Maximum),
			MarshalName: stringy.New(name).LcFirst(),
			NestedTypes: nil,
		}}
	}
	return schemaDefs
}

func uintOrMaxA(x *uint) uint {
	if x != nil {
		return *x
	}
	return math.MaxUint
}

func floatOrMinA(x *float64) float64 {
	if x != nil {
		return *x
	}
	return -math.MaxFloat64
}

func floatOrMaxA(x *float64) float64 {
	if x != nil {
		return *x
	}
	return math.MaxFloat64
}

func toGoTypeAsync(s *asyncapiv3.Schema) (string, bool) {
	if s.Reference == "" || s.ReferenceTo == nil {
		return "interface{}", false
	}
	switch {
	case s.Type == "number":
		switch s.Format {
		case "float":
			return "float32", false
		case "double":
			return "float64", false
		default:
			return "float", false
		}
	case s.Type == "integer":
		if s.Format != "" {
			return s.Format, false
		}
		return "int", false
	case s.Type == "boolean":
		return "bool", false
	case s.Type == "string":
		switch s.Format {
		case "binary":
			return "[]byte", false
		case "date":
			IMPORT_TIME = true
			return "time.Time", false
		case "uuid":
			IMPORT_UUID = true
			return "uuid.UUID", false
		default:
			return "string", false
		}
	case s.Type == "array":
		items, _ := toGoTypeAsync(s.Items)
		return "[]" + items, false
	case s.Type == "object":
		if s.AdditionalProperties != nil {
			t, _ := toGoTypeAsync(s.AdditionalProperties)
			return "map[string]" + t, false
		}
		if s.Reference != "" {
			return path.Base(s.Reference), false
		}
		return "map[string]interface{}", false
	default:
		types := s.Type
		if len(types) > 0 {
			return types, false
		}
		return "interface{}", false
	}
}

func generateAsyncImports() ImportsConfig {
	var defs []ImportDefinition
	if IMPORT_UUID {
		defs = append(defs, ImportDefinition{"uuid", "\"github.com/google/uuid\""})
	}
	if IMPORT_TIME {
		defs = append(defs, ImportDefinition{"time", ""})
	}
	return ImportsConfig{ImportDefs: defs}
}
