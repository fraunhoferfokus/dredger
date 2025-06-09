// Package schemas enthält JSON-Schema Validierung für Payloads
package schemas

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"path/filepath"
)

// ValidatePayload lädt das Schema und validiert das JSON-Payload.
func ValidatePayload(schemaDir, channel string, payload []byte) ([]gojsonschema.ResultError, error) {
	schemaPath := filepath.Join(schemaDir, channel+".json")
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}
	loader := gojsonschema.NewBytesLoader(schemaBytes)
	document := gojsonschema.NewBytesLoader(payload)
	result, err := gojsonschema.Validate(loader, document)
	if err != nil {
		return nil, err
	}
	if result.Valid() {
		return nil, nil
	}
	return result.Errors(), nil
}
