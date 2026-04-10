package middleware

import (
	"fmt"
	"io"
	"strings"

	"github.com/labstack/echo/v5"
	"gopkg.in/yaml.v3"
)

// Bind binds path params, query params and the request body into provided type `dest`.
// The binder binds body based on Content-Type header and adds support for YAML to the default echo binder.
func Bind(c *echo.Context, dest any) error {
	req := (*c).Request()

	// Normalize content type (strip charset etc.)
	ct := req.Header.Get(echo.HeaderContentType)
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}
	ct = strings.TrimSpace(ct)

	// Handle YAML
	if ct == "application/x-yaml" || ct == "application/yaml" || ct == "text/yaml" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("invalid request body (cannot process yaml)")
		}

		if err := yaml.Unmarshal(body, dest); err != nil {
			return err
		}

		return nil
	}

	// Fallback to Echo's default binder
	return (*c).Bind(dest)
}
