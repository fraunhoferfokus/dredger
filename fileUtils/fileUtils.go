package fileUtils

import (
	//"fmt"
	"fmt"
	"github.com/ghodss/yaml"
	"path/filepath"
	"strings"
)

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

func ConvertYAMLIntoJSON(file string) string {
	json, err := yaml.YAMLToJSON([]byte(file))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println(string(json))
	return string(json)
}
