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
//  Version from praktikum
// 		_, err := parser.ParseOpenAPISpecFile("../../examples/nonexistent.yaml")
//	if err == nil || err.Error() != "file not found" {
// 	}
}

func TestFileExists(t *testing.T) {
	spec, err := parser.ParseOpenAPISpecFile("../../examples/stores/stores.yaml")
	if err != nil {
		t.Errorf("Error was not expected: %v", err)
	}
	if spec == nil {
		t.Errorf("Spec should not be null")
	}
}
