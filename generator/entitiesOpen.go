package generator

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gobeam/stringy"
	"github.com/rs/zerolog/log"
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
	Kind        string              // Composite kind: "allof", "anyof", "oneof", "alias" — empty for regular schemas
	Variants    []VariantDefinition // Variants for composite schemas (used when Kind is set)
}

// VariantDefinition represents a variant of a composite schema (allOf, anyOf, oneOf).
// For allOf, variants are embedded types in the resulting struct.
// For anyOf/oneOf, variants represent alternative types to be chosen from.
type VariantDefinition struct {
	Name      string              // Type name for inline variants, or ref target name for $ref variants
	IsRef     bool                // True if this variant is a $ref to an already-generated type
	RefTarget string              // Go type name of the referenced type (e.g., "User")
	Props     []TypeDefinition    // Property definitions for inline object variants
	Kind      string              // "" for a ref or inline object; "allof"/"anyof"/"oneof" for a nested composite variant
	Variants  []VariantDefinition // Sub-variants for nested composite variants (used when Kind is set)
}

type ImportDefinition struct {
	Name string
	URL  string
}

type ImportsConfig struct {
	ImportDefs []ImportDefinition
}

// schemaQueue is a helper for queuing schemas for type definition generation.
type schemaQueue struct {
	Name   string
	Schema *openapi3.SchemaRef
}

// Aus den Schemas in Components die Typdefinitionen und generiert entities,imports,structs und validate files
func GenerateTypes(spec *openapi3.T, pConf ProjectConfig) {
	if spec == nil || spec.Components == nil {
		return
	}
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
			"templates/common/entities/entity_allof.tmpl",
			"templates/common/entities/entity_anyof.tmpl",
			"templates/common/entities/entity_oneof.tmpl",
			"templates/common/entities/entity_regular.tmpl",
			"templates/common/entities/variant_types.tmpl",
			"templates/common/entities/property_types.tmpl",
		}
		createFileFromTemplates(filePath, templateFiles, conf)
	}
}

func generateTypeDefs(schemas *openapi3.Schemas) map[string][]TypeDefinition {
	schemaDefs := make(map[string][]TypeDefinition)
	// build schemas into queue
	queue := []schemaQueue{}
	for schemaName, ref := range *schemas {
		queue = append(queue, schemaQueue{schemaName, ref})
	}

	for len(queue) > 0 {
		// pop
		first := queue[0]
		queue = queue[1:]
		schemaName := first.Name
		ref := first.Schema
		// process schema
		var goType string
		if ref.Ref != "" { // refs are just aliases
			originalType := extractRefType(ref.Ref)
			schemaDefs[schemaName] = []TypeDefinition{{
				schemaName,
				originalType,
				0, 0, "", 0, 0, // ignore length (etc) requirements
				stringy.New(schemaName).LcFirst(),
				[]TypeDefinition{},
				"alias", nil,
			}}
			log.Info().Str("ref", ref.Ref).Msg("Processing ref as type")
		} else if ref.Value.Type.Includes("number") {
			switch ref.Value.Format {
			case "float":
				goType = "float32"
			case "double":
				goType = "float64"
			default:
				goType = "float64" // as float does not exist
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
				"",
				nil,
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
				"",
				nil,
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
				"",
				nil,
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
				"",
				nil,
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
				"",
				nil,
			}}
		} else if ref.Value.Type.Includes("object") {
			schemaDefs[schemaName] = generatePropertyDefs(&ref.Value.Properties, PascalCase(schemaName))
		} else if ref.Value.AllOf != nil {
			schemaDefs[schemaName] = generateComposedType(ref, &ref.Value.AllOf, schemaName, "allof")
		} else if ref.Value.AnyOf != nil {
			schemaDefs[schemaName] = generateComposedType(ref, &ref.Value.AnyOf, schemaName, "anyof")
		} else if ref.Value.OneOf != nil {
			schemaDefs[schemaName] = generateComposedType(ref, &ref.Value.OneOf, schemaName, "oneof")
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

// PascalCase converts a name to UpperCamelCase (PascalCase) for exported type names
func PascalCase(name string) string {
	// Use the existing camelcase function and capitalize first letter
	camel := camelcase(name)
	return stringy.New(camel).UcFirst()
}

// generatePropertyDefs builds the property definitions for an object's properties.
// prefix is the PascalCase name of the enclosing type, used to derive unique helper
// type names for properties whose schema is itself a composite or nested object.
func generatePropertyDefs(properties *openapi3.Schemas, prefix string) []TypeDefinition {
	typeDefs := make([]TypeDefinition, len(*properties))
	i := 0
	for name, property := range *properties {
		propertyDef := TypeDefinition{
			Name:        name,
			MinLength:   property.Value.MinLength,
			MaxLength:   uintOrMax(property.Value.MaxLength),
			Pattern:     property.Value.Pattern,
			Minimum:     floatOrMin(property.Value.Min),
			Maximum:     floatOrMax(property.Value.Max),
			MarshalName: stringy.New(name).LcFirst(),
		}

		goType, nested := toGoType(property)

		// A property whose schema is a composite (oneOf/anyOf, or an allOf that
		// toGoType can't reduce to a single $ref) has no direct Go type. Generate a
		// named helper composite type in the same file and reference it by name.
		if goType == "" && !nested && property.Ref == "" {
			if kind := compositeKind(property); kind != "" {
				helperName := PascalCase(prefix) + PascalCase(name)
				propertyDef.Type = helperName
				propertyDef.Kind = kind
				propertyDef.Variants = buildVariants(compositeSchemas(property, kind), helperName, kind)
				typeDefs[i], i = propertyDef, i+1
				continue
			}
		}

		if nested {
			nestedGoTypes := generatePropertyDefs(&property.Value.Properties, PascalCase(prefix)+PascalCase(name))
			if len(nestedGoTypes) == 0 { // allow empty structs
				goType += "{}"
			}
			propertyDef.NestedTypes = nestedGoTypes
		}
		propertyDef.Type = goType
		typeDefs[i], i = propertyDef, i+1
	}

	return typeDefs
}

// compositeKind returns "allof"/"anyof"/"oneof" if the schema is a composite, else "".
func compositeKind(sRef *openapi3.SchemaRef) string {
	if sRef.Value == nil {
		return ""
	}
	switch {
	case sRef.Value.AllOf != nil:
		return "allof"
	case sRef.Value.AnyOf != nil:
		return "anyof"
	case sRef.Value.OneOf != nil:
		return "oneof"
	}
	return ""
}

// compositeSchemas returns the subschemas of a composite for the given kind.
func compositeSchemas(sRef *openapi3.SchemaRef, kind string) *openapi3.SchemaRefs {
	switch kind {
	case "allof":
		return &sRef.Value.AllOf
	case "anyof":
		return &sRef.Value.AnyOf
	case "oneof":
		return &sRef.Value.OneOf
	}
	return nil
}

// schema type to generated go type
func toGoType(sRef *openapi3.SchemaRef) (goType string, nested bool) {
	// Resolve $ref first — applies to any type (including object types without explicit "type:")
	if sRef.Ref != "" {
		goType := extractRefType(sRef.Ref)
		return goType, false
	}

	// Handle allOf — resolve single $ref allOf to the ref target type before type checks
	if sRef.Value != nil && sRef.Value.AllOf != nil {
		for _, variant := range sRef.Value.AllOf {
			if variant.Ref != "" {
				return extractRefType(variant.Ref), false
			}
		}
	}

	if sRef.Value.Type.Includes("number") {
		switch sRef.Value.Format {
		case "float":
			goType = "float32"
		case "double":
			goType = "float64"
		default:
			goType = "float64" // as float does not exists
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
				goType = "map[string]" + extractRefType(sRef.Value.AdditionalProperties.Schema.Ref)
			} else {
				goType = "map[string]??"
			}
		} else {
			goType = "struct"
			nested = true
		}
	}
	return goType, nested
}

// generateComposedType returns the TypeDefinition for a composed type so with `kind` = "allof", "anyof", "oneof".
// ref is the original openapi schema.
// schema is the type to be composed of.
// schemaName is the name of the new schema (type).
func generateComposedType(ref *openapi3.SchemaRef, schema *openapi3.SchemaRefs, schemaName, kind string) []TypeDefinition {
	variants := buildVariants(schema, PascalCase(schemaName), kind)
	return []TypeDefinition{{
		Name:        schemaName,
		Type:        "struct",
		MinLength:   ref.Value.MinLength,
		MaxLength:   uintOrMax(ref.Value.MaxLength),
		Pattern:     ref.Value.Pattern,
		Minimum:     floatOrMin(ref.Value.Min),
		Maximum:     floatOrMax(ref.Value.Max),
		MarshalName: stringy.New(schemaName).LcFirst(),
		NestedTypes: nil,
		Kind:        kind,
		Variants:    variants,
	}}
}

// buildVariants builds the variant definitions for the subschemas of a composite
// (allOf/anyOf/oneOf). baseName is the PascalCase name used to derive helper type
// names for inline variants; kind is the parent composite kind, used only for logging.
func buildVariants(schema *openapi3.SchemaRefs, baseName, kind string) []VariantDefinition {
	variants := make([]VariantDefinition, 0, len(*schema))
	for i, variant := range *schema {
		partName := fmt.Sprintf("%sPart%d", baseName, i)
		switch {
		case variant.Ref != "":
			// $ref variant — already has a generated type, embed it directly
			targetType := extractRefType(variant.Ref)
			variants = append(variants, VariantDefinition{
				Name:      targetType,
				IsRef:     true,
				RefTarget: targetType,
			})
		case variant.Value == nil:
			log.Warn().
				Str("schema", baseName).
				Int("variantIndex", i).
				Msgf("Ignoring nil %s variant", kind)
		case variant.Value.AllOf != nil:
			variants = append(variants, VariantDefinition{
				Name:     partName,
				Kind:     "allof",
				Variants: buildVariants(&variant.Value.AllOf, partName, "allof"),
			})
		case variant.Value.AnyOf != nil:
			variants = append(variants, VariantDefinition{
				Name:     partName,
				Kind:     "anyof",
				Variants: buildVariants(&variant.Value.AnyOf, partName, "anyof"),
			})
		case variant.Value.OneOf != nil:
			variants = append(variants, VariantDefinition{
				Name:     partName,
				Kind:     "oneof",
				Variants: buildVariants(&variant.Value.OneOf, partName, "oneof"),
			})
		case variant.Value.Type.Includes("object"):
			// Inline object variant — generate an intermediate type with its properties
			variants = append(variants, VariantDefinition{
				Name:  partName,
				IsRef: false,
				Props: generatePropertyDefs(&variant.Value.Properties, partName),
			})
		default:
			// Non-object, non-ref, non-composite variant (e.g. a bare primitive constraint)
			log.Warn().
				Str("schema", baseName).
				Int("variantIndex", i).
				Str("type", strings.Join(variant.Value.Type.Slice(), ",")).
				Msgf("Ignoring non-object %s variant", kind)
		}
	}
	return variants
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

// extractRefType returns the corresponding go type for a given ref.
func extractRefType(ref string) string {
	splited := strings.Split(ref, "/")
	if len(splited) == 0 {
		log.Error().Str("ref", ref).Msg("Failed to extract ref type")
		return "NotImplemented" // return invalid type
	}
	refObject := splited[len(splited)-1]
	return PascalCase(refObject)
}
