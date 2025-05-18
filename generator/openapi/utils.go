package generator

import (
	"regexp"
	"strings"
)

func validateStatusCode(code string) bool {
	return regexp.MustCompile(`[1-5](\d\d|XX)`).MatchString(code)
}

func convertPathParams(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}
