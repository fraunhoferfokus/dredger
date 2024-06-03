package core

import (
	_ "embed"
	"strings"
)

var (
	// Version labels the release
	Version string = strings.TrimSpace(version)
	//go:embed version
	version string
)
