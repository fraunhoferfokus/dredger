package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
)

func ParseComponentSpecFile(path string) (*openapi3.T, error) {
	// read file and modify content (add minimal openapi)
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}
	specContent, err := wrapComponentFile(contents)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap component spec: %v", err)
	}

	// parse openapi
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true // TODO make configurable by flag
	spec, err := loadSpec(loader, path, specContent)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Component Spec file loaded successfully.")
	// validate
	err = spec.Validate(loader.Context)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Component Spec file validated successfully.")

	return spec, err
}

// wrapComponentsFile adds minimal OpenAPI wrapper to a components-only file.
// Then it can be parsed by the openapi3 package.
func wrapComponentFile(data []byte) ([]byte, error) {
	var partial map[string]any
	if err := yaml.Unmarshal(data, &partial); err != nil {
		return nil, err
	}

	// Wrap in minimal OpenAPI structure
	wrapped := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   "Partial Spec",
			"version": "1.0.0",
		},
		"paths": map[string]any{},
	}

	// Move components to right place
	if components, ok := partial["components"]; ok {
		wrapped["components"] = components
	}

	return yaml.Marshal(wrapped)
}

// loadSpec loads the given spec either from fileContent or (if nil) from specPath.
// If an error due to reference occurs, it retries with changed paths (using loadSpecWithDir).
func loadSpec(loader *openapi3.Loader, specPath string, fileContent []byte) (*openapi3.T, error) {
	// load spec without changing dir for ref resolve
	var spec *openapi3.T
	var err error
	if fileContent == nil {
		spec, err = loader.LoadFromFile(specPath)
	} else {
		spec, err = loader.LoadFromData(fileContent)
	}

	// return success
	if err == nil {
		return spec, nil
	}
	// if error due to refs, retry with changed dir
	if strings.HasPrefix(err.Error(), "error resolving reference") {
		return loadSpecWithDir(loader, specPath, fileContent)
	}
	return spec, err
}

// loadSpecWithDir loads the given openapi spec file (from fileContent or spec path)
// while resolving refs in the same directory as the spec file.
// When fileContent is nil, the specPath will be used to read the file.
func loadSpecWithDir(loader *openapi3.Loader, specPath string, fileContent []byte) (*openapi3.T, error) {
	// Remember original directory
	origDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Change to spec's directory
	specDir := filepath.Dir(specPath)
	if err := os.Chdir(specDir); err != nil {
		return nil, err
	}
	// Restore original directory
	defer os.Chdir(origDir)

	if fileContent == nil {
		filename := filepath.Base(specPath)
		return loader.LoadFromFile(filename)
	}
	return loader.LoadFromData(fileContent)
}
