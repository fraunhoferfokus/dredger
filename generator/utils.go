package generator

import (
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func validateStatusCode(code string) bool {
	return regexp.MustCompile(`[1-5](\d\d|XX)`).MatchString(code)
}

func convertPathParams(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}

// sortedResponseCodes returns the keys of an OpenAPI responses map in a stable,
// readable order: numeric/wildcard status codes sorted lexicographically
// ("200" < "404" < "5XX"), with "default" always last.
func sortedResponseCodes(responses map[string]*openapi3.ResponseRef) []string {
	codes := slices.Sorted(maps.Keys(responses))
	slices.SortStableFunc(codes, func(a, b string) int {
		if (a == "default") != (b == "default") {
			if a == "default" {
				return 1
			}
			return -1
		}
		return 0
	})
	return codes
}
