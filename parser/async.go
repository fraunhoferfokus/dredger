package parser

import (
	"errors"
	"os"

	fs "dredger/fileUtils"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// AsyncAPIDoc – minimaler Ausschnitt; für Schritt 1 völlig ausreichend
type AsyncAPIDoc struct {
	Asyncapi string `yaml:"asyncapi"`
	Info     struct {
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
	} `yaml:"info"`
	Channels map[string]interface{} `yaml:"channels"`
}

// ParseAsyncAPISpecFile liest eine AsyncAPI-Datei (YAML oder JSON) ein,
// prüft Basisfelder und gibt sie als Struct zurück.
func ParseAsyncAPISpecFile(path string) (*AsyncAPIDoc, error) {
	if !fs.CheckIfFileExists(path) {
		return nil, errors.New("file not found")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc AsyncAPIDoc
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	// Mini-Validierung
	if doc.Asyncapi == "" || len(doc.Channels) == 0 {
		return nil, errors.New("invalid AsyncAPI: missing asyncapi or channels")
	}

	log.Info().
		Str("title", doc.Info.Title).
		Str("version", doc.Info.Version).
		Int("channels", len(doc.Channels)).
		Msg("AsyncAPI spec loaded successfully")

	return &doc, nil
}
