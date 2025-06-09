package parser

import (
	//fs "dredger/fileUtils"
	"path/filepath"
	"strings"

	//"embed"
	gen "dredger/generator/asyncapi"
	"fmt"

	//"github.com/rs/zerolog/log"
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"

	//"github.com/rs/zerolog/internal/json"
	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

const schemaFile = "./examples/schemas/asyncapiv3Schema.json"

// ParseAsyncAPISpecFile liest eine AsyncAPI-Datei (YAML oder JSON) ein,
// prüft Basisfelder und gibt sie als Struct zurück.
func ParseAsyncAPISpecFile(file string) (gen.Specification, error) {
	//Case YAML
	if Check_YAML(file) {
		jsonBytes, err := ConvertYAMLIntoJSON(file)
		if err != nil {
			return gen.Specification{}, err
		}

		var spec gen.Specification
		if err := json.Unmarshal(jsonBytes, &spec); err != nil {
			return gen.Specification{}, fmt.Errorf("JSON decode to struct failed: %w", err)
		}

		fmt.Printf("Succesfully parsed yaml to json example version: %v", spec.Version)
		return spec, nil
	}

	//Case JSON
	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Msgf("schema file %s not found, skipping validation", schemaFile)
			sch = nil
		} else {
			return gen.Specification{}, fmt.Errorf("Schema compile error: %w", err)
		}
	}

	f, err := os.Open(file)
	if err != nil {
		return gen.Specification{}, err
	}
	defer f.Close()

	// Unmarshal for validation
	var raw any
	json.NewDecoder(f).Decode(&raw)

	if sch != nil {
		err = sch.Validate(raw)
		if err != nil {
			return gen.Specification{}, fmt.Errorf("Validation failed: %w", err)
		}
	}

	// Reset reader to start
	f.Seek(0, 0)

	// Ins Struct decodieren
	var spec gen.Specification
	err = json.NewDecoder(f).Decode(&spec)
	if err != nil {
		return gen.Specification{}, fmt.Errorf("JSON decode error: %w", err)
	}
	fmt.Printf("Succesfully parsed example version: %v", spec.Version)
	return spec, nil
}

func Check_YAML(file string) bool {
	if isYMLFile(file) {
		return true
	}
	return false
}

func isJSONFile(file string) bool {
	end := strings.ToLower(filepath.Ext(file))
	return end == ".json"
}

func isYMLFile(file string) bool {
	end := strings.ToLower(filepath.Ext(file))
	return end == ".yaml"
}

func ConvertYAMLIntoJSON(filePath string) ([]byte, error) {
	yamlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Warn().Err(err).Msg("Could not read YAML file.")
		return nil, err
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		log.Warn().Err(err).Msg("Could not convert YAML to JSON.")
		return nil, err
	}

	return jsonBytes, nil
}
