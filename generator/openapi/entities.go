package generator

import (
	"math"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gobeam/stringy"
)

var IMPORT_UUID bool
var IMPORT_TIME bool

type ModelConfig struct {
	Imports     ImportsConfig
	SchemaDefs  map[string][]TypeDefinition
	ProjectName string
}

type TypeDefinition struct {
	Name        string
	Type        string
	MinLength   uint64
	MaxLength   uint64
	Pattern     string
	Minimum     float64
	Maximum     float64
	MarshalName string
	NestedTypes []TypeDefinition
}

type ImportDefinition struct {
	Name string
	URL  string
}

type ImportsConfig struct {
	ImportDefs []ImportDefinition
}

func GenerateTypes(spec *openapi3.T, pConf ProjectConfig) {
	if spec != nil && spec.Components != nil {
		schemaDefs := generateTypeDefs(&spec.Components.Schemas)
		imports := generateImports()
		var conf ModelConfig
		conf.Imports = imports
		conf.ProjectName = pConf.Name

		for schema, defs := range schemaDefs {
			conf.SchemaDefs = map[string][]TypeDefinition{schema: defs}
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

func generateTypeDefs(schemas *openapi3.Schemas) map[string][]TypeDefinition {
	schemaDefs := make(map[string][]TypeDefinition, len(*schemas))
	for name, ref := range *schemas {
		var goType string
		switch {
		case ref.Value.Type.Includes("number"):
			switch ref.Value.Format {
			case "float":
				goType = "float32"
			case "double":
				goType = "float64"
			default:
				goType = "float"
			}
		case ref.Value.Type.Includes("integer"):
			goType = "int"
			if ref.Value.Format != "" {
				goType = ref.Value.Format
			}
		case ref.Value.Type.Includes("boolean"):
			goType = "bool"
		case ref.Value.Type.Includes("string"):
			switch ref.Value.Format {
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
		case ref.Value.Type.Includes("array"):
			items, _ := toGoType(ref.Value.Items)
			goType = "[]" + items
		case ref.Value.Type.Includes("object"):
			goType, _ = toGoType(ref.Value.AdditionalProperties.Schema)
		default:
			types := ref.Value.Type.Slice()
			if len(types) > 0 {
				goType = types[0]
			}
		}

		schemaDefs[name] = []TypeDefinition{{
			Name:        name,
			Type:        goType,
			MinLength:   ref.Value.MinLength,
			MaxLength:   uintOrMax(ref.Value.MaxLength),
			Pattern:     ref.Value.Pattern,
			Minimum:     floatOrMin(ref.Value.Min),
			Maximum:     floatOrMax(ref.Value.Max),
			MarshalName: stringy.New(name).LcFirst(),
			NestedTypes: nil,
		}}
	}
	return schemaDefs
}

func uintOrMax(x *uint64) uint64 {
	if x != nil {
		return *x
	}
	return math.MaxUint64
}

func floatOrMin(x *float64) float64 {
	if x != nil {
		return *x
	}
	return -math.MaxFloat64
}

func floatOrMax(x *float64) float64 {
	if x != nil {
		return *x
	}
	return math.MaxFloat64
}

func toGoType(sRef *openapi3.SchemaRef) (string, bool) {
	if sRef == nil || sRef.Value == nil {
		return "interface{}", false
	}
	switch {
	case sRef.Value.Type.Includes("number"):
		switch sRef.Value.Format {
		case "float":
			return "float32", false
		case "double":
			return "float64", false
		default:
			return "float64", false
		}
	case sRef.Value.Type.Includes("integer"):
		if sRef.Value.Format != "" {
			return sRef.Value.Format, false
		}
		return "int", false
	case sRef.Value.Type.Includes("boolean"):
		return "bool", false
	case sRef.Value.Type.Includes("string"):
		switch sRef.Value.Format {
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
	case sRef.Value.Type.Includes("array"):
		items, _ := toGoType(sRef.Value.Items)
		return "[]" + items, false
	case sRef.Value.Type.Includes("object"):
		if sRef.Value.AdditionalProperties.Schema != nil {
			t, _ := toGoType(sRef.Value.AdditionalProperties.Schema)
			return "map[string]" + t, false
		}
		if sRef.Ref != "" {
			parts := strings.Split(sRef.Ref, "/")
			return parts[len(parts)-1], false
		}
		return "map[string]interface{}", false
	default:
		types := sRef.Value.Type.Slice()
		if len(types) > 0 {
			return types[0], false
		}
		return "interface{}", false
	}
}

func generateImports() ImportsConfig {
	var defs []ImportDefinition
	if IMPORT_UUID {
		defs = append(defs, ImportDefinition{"uuid", "\"github.com/google/uuid\""})
	}
	if IMPORT_TIME {
		defs = append(defs, ImportDefinition{"time", ""})
	}
	return ImportsConfig{ImportDefs: defs}
}
