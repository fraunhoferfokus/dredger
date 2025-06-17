// Package config liest Umgebungsvariablen ein
package config

import (
	"os"
)

// GetEnv liest eine Umgebungsvariable oder liefert Default.
func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
