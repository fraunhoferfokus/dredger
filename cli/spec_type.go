package cli

import (
	"bufio"
	"os"
	"strings"
)

// Spec type enum
type specType int

const (
	openApiSpec specType = iota
	asyncApiSpec
	componentSpec
)

// automatische spec erkennung 💃
func detectSpecType(specPath string) (specType, error) {
	firstLine, err := firstNonCommentLine(specPath)
	if err != nil {
		return -1, err
	}

	if strings.Contains(firstLine, "\"asyncapi\"") || strings.HasPrefix(firstLine, "asyncapi:") {
		return asyncApiSpec, nil
	}
	if strings.Contains(firstLine, "\"openapi\"") || strings.HasPrefix(firstLine, "openapi:") {
		return openApiSpec, nil
	}
	//veraltete "schreibweise" jetzt openapi
	if strings.Contains(firstLine, "\"swagger\"") || strings.HasPrefix(firstLine, "swagger:") {
		return openApiSpec, nil
	}
	if strings.Contains(firstLine, "\"components\"") || strings.HasPrefix(firstLine, "components:") {
		return componentSpec, nil
	}
	return -1, nil
}

// firstNonCommentLine returns the first line of the file which is not empty and does not have a comment with `#`
func firstNonCommentLine(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	// go line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip empty lines
		if line == "" {
			continue
		}
		// skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		// return first line
		return strings.ToLower(line), nil
	}
	// check for error scanning file
	if err := scanner.Err(); err != nil {
		return "", err
	}
	// first (non comment line) is empty
	return "", nil
}
