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

// Aus den Schemas in Components die Typdefinitionen und generiert entities,imports,structs und validate files
func GenerateTypes(spec *openapi3.T, pConf ProjectConfig) {
	if spec != nil && spec.Components != nil {
		schemaDefs := generateTypeDefs(&spec.Components.Schemas)
		imports := generateImports()
		var conf ModelConfig
		conf.Imports = imports
		conf.ProjectName = pConf.Name

		for schema, defs := range schemaDefs {
			//log.Debug().Str("Operationname", schema).Msg("SchemaDefs")
			conf.SchemaDefs = map[string][]TypeDefinition{schema: defs}
			fileName := strings.ToLower(schema) + ".go"
			filePath := filepath.Join(pConf.Path, EntitiesPkg, fileName)
			templateFiles := []string{
				"templates/common/entities/entities.go.tmpl",
				"templates/common/entities/imports.tmpl",
				"templates/common/entities/structs.tmpl",
				"templates/common/entities/validate.tmpl",
			}
			createFileFromTemplates(filePath, templateFiles, conf)
		}
	}
}

func generateTypeDefs(schemas *openapi3.Schemas) map[string][]TypeDefinition {
	schemaDefs := make(map[string][]TypeDefinition, len(*schemas))
	for schemaName, ref := range *schemas {
		//fmt.Printf("%s: %#v\n\n", schemaName, ref.Value.Type)
		var goType string
		if ref.Value.Type.Includes("number") {
			switch ref.Value.Format {
			case "float":
				goType = "float32"
			case "double":
				goType = "float64"
			default:
				goType = "float"
			}
			schemaDefs[schemaName] = []TypeDefinition{{
				schemaName,
				goType,
				ref.Value.MinLength,
				uintOrMax(ref.Value.MaxLength),
				ref.Value.Pattern,
				floatOrMin(ref.Value.Min),
				floatOrMax(ref.Value.Max),
				stringy.New(schemaName).LcFirst(),
				[]TypeDefinition{},
			}}
		} else if ref.Value.Type.Includes("integer") {
			goType = "int"
			if ref.Value.Format != "" {
				goType = ref.Value.Format
			}
			schemaDefs[schemaName] = []TypeDefinition{{
				schemaName,
				goType,
				ref.Value.MinLength,
				uintOrMax(ref.Value.MaxLength),
				ref.Value.Pattern,
				floatOrMin(ref.Value.Min),
				floatOrMax(ref.Value.Max),
				stringy.New(schemaName).LcFirst(),
				[]TypeDefinition{},
			}}
		} else if ref.Value.Type.Includes("boolean") {
			goType = "bool"
			schemaDefs[schemaName] = []TypeDefinition{{
				schemaName,
				goType,
				ref.Value.MinLength,
				uintOrMax(ref.Value.MaxLength),
				ref.Value.Pattern,
				floatOrMin(ref.Value.Min),
				floatOrMax(ref.Value.Max),
				stringy.New(schemaName).LcFirst(),
				[]TypeDefinition{},
			}}
		} else if ref.Value.Type.Includes("string") {
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
			schemaDefs[schemaName] = []TypeDefinition{{
				schemaName,
				goType,
				ref.Value.MinLength,
				uintOrMax(ref.Value.MaxLength),
				ref.Value.Pattern,
				floatOrMin(ref.Value.Min),
				floatOrMax(ref.Value.Max),
				stringy.New(schemaName).LcFirst(),
				[]TypeDefinition{},
			}}
		} else if ref.Value.Type.Includes("array") {
			items, _ := toGoType(ref.Value.Items)
			goType = "[]" + items
			schemaDefs[schemaName] = []TypeDefinition{{
				schemaName,
				goType,
				ref.Value.MinLength,
				uintOrMax(ref.Value.MaxLength),
				ref.Value.Pattern,
				floatOrMin(ref.Value.Min),
				floatOrMax(ref.Value.Max),
				stringy.New(schemaName).LcFirst(),
				[]TypeDefinition{},
			}}
		} else if ref.Value.Type.Includes("object") {
			schemaDefs[schemaName] = generatePropertyDefs(&ref.Value.Properties)
		}
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

func generatePropertyDefs(properties *openapi3.Schemas) []TypeDefinition {
	typeDefs := make([]TypeDefinition, len(*properties))
	i := 0
	for name, property := range *properties {
		goType, nested := toGoType(property)
		var nestedGoTypes []TypeDefinition
		if nested {
			nestedGoTypes = generatePropertyDefs(&property.Value.Properties)
		}

		propertyDef := TypeDefinition{
			name,
			goType,
			property.Value.MinLength,
			uintOrMax(property.Value.MaxLength),
			property.Value.Pattern,
			floatOrMin(property.Value.Min),
			floatOrMax(property.Value.Max),
			stringy.New(name).LcFirst(),
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
