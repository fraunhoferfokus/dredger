package generator

import (
	"math"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
)

var IMPORT_UUID bool
var IMPORT_TIME bool

type ModelCOnfig struct {
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
	// only if Type is struct
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
		schemaDefs := generateStructDefs(&spec.Components.Schemas)
		imports := generateImports()
		var conf ModelCOnfig
		conf.Imports = imports
		conf.ProjectName = pConf.Name

		for schema, t := range schemaDefs {
			conf.SchemaDefs = map[string][]TypeDefinition{schema: t}
			fileName := strings.ToLower(schema) + ".go"
			filePath := filepath.Join(pConf.Path, EntitiesPkg, fileName)
			templateFiles := []string{"templates/entities/entities.go.tmpl", "templates/entities/imports.tmpl", "templates/entities/structs.tmpl", "templates/entities/validate.tmpl"}
			createFileFromTemplates(filePath, templateFiles, conf)
		}
	}
}

func generateStructDefs(schemas *openapi3.Schemas) map[string][]TypeDefinition {
	schemaDefs := make(map[string][]TypeDefinition, len(*schemas))

	for schemaName, ref := range *schemas {
		schemaDefs[schemaName] = generateTypeDefs(&ref.Value.Properties)
	}
	return schemaDefs
}

func uintOrMax(x *uint64) uint64 {
	if x != nil {
		return *x
	}
	return math.MaxInt64
}

func floatOrMin(x *float64) float64 {
	if x != nil {
		return *x
	}
	return math.MaxFloat64 * -1
}

func floatOrMax(x *float64) float64 {
	if x != nil {
		return *x
	}
	return math.MaxFloat64
}

func generateTypeDefs(properties *openapi3.Schemas) []TypeDefinition {
	typeDefs := make([]TypeDefinition, len(*properties))
	i := 0
	for name, property := range *properties {
		goType, nested := toGoType(property)
		var nestedGoTypes []TypeDefinition
		if nested {
			nestedGoTypes = generateTypeDefs(&property.Value.Properties)
		}

		// first letter to lower case
		marshalName := []rune(name)
		marshalName[0] = unicode.ToLower(marshalName[0])
		propertyDef := TypeDefinition{
			name,
			goType,
			property.Value.MinLength,
			uintOrMax(property.Value.MaxLength),
			property.Value.Pattern,
			floatOrMin(property.Value.Min),
			floatOrMax(property.Value.Max),
			string(marshalName),
			nestedGoTypes,
		}
		typeDefs[i], i = propertyDef, i+1
	}

	return typeDefs
}

// schema type to generated go type
func toGoType(sRef *openapi3.SchemaRef) (goType string, nested bool) {

	if sRef.Value.Type.Includes("number") {
		switch sRef.Value.Format {
		case "float":
			goType = "float32"
		case "double":
			goType = "float64"
		default:
			goType = "float"
		}
	} else if sRef.Value.Type.Includes("integer") {
		goType = "int"
		if sRef.Value.Format != "" {
			goType = sRef.Value.Format
		}
	} else if sRef.Value.Type.Includes("boolean") {
		goType = "bool"
	} else if sRef.Value.Type.Includes("string") {
		switch sRef.Value.Format {
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
	} else if sRef.Value.Type.Includes("array") {
		items, _ := toGoType(sRef.Value.Items)
		goType = "[]" + items
	} else if sRef.Value.Type.Includes("object") {
		if sRef.Value.AdditionalProperties.Schema != nil {
			if sRef.Value.AdditionalProperties.Schema.Ref != "" {
				splitRef := strings.Split(sRef.Value.AdditionalProperties.Schema.Ref, "/")
				goType = "map[string]" + splitRef[len(splitRef)-1]
			} else {
				goType = "map[string]??"
			}
		} else if sRef.Ref != "" {
			// checks if object type is defined by reference elsewhere in the schema
			splitRef := strings.Split(sRef.Ref, "/")
			goType = splitRef[len(splitRef)-1]
		} else {
			goType = "struct"
			nested = true
		}
	} else {
		types := sRef.Value.Type.Slice()
		if len(types) > 0 {
			goType = types[0]
		}
	}
	return goType, nested
}

func generateImports() ImportsConfig {
	var importDefs []ImportDefinition
	if IMPORT_UUID {
		importDefs = append(importDefs, ImportDefinition{"", "\"github.com/google/uuid\""})
	}
	if IMPORT_TIME {
		importDefs = append(importDefs, ImportDefinition{"time", ""})
	}

	conf := ImportsConfig{
		importDefs,
	}

	return conf
}
