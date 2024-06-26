package tests

import (
	"dredger/parser"
	"testing"
)

func TestFileDoesNotExist(t *testing.T) {
	errorMessage := "file not found"
	_, err := parser.ParseOpenAPISpecFile("../../examples/apiwithexamples.json")
	if err.Error() != errorMessage {
		t.Errorf("Actual error %v, and expected %v", err, errorMessage)
	}
}

func TestFileExists(t *testing.T) {
	spec, err := parser.ParseOpenAPISpecFile("../../examples/stores.yaml")
	if err != nil {
		t.Errorf("Error was not expected")
	}
	if spec == nil {
		t.Errorf("Spec should not be null")
	}
}
