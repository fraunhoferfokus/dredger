package generator

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

// GenerateBdd erzeugt ein Godog-Testfile aus einer Feature-Datei.
func GenerateBdd(path string) {
	stepListing := parseSteps(path)
	_, err := os.Create("server_godog_test.go")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create BDD test file")
	}

	f, err := os.OpenFile("server_godog_test.go", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open BDD test file")
	}
	defer f.Close()

	// hol die Vorlage aus common
	content, _ := os.ReadFile("templates/common/bdd.go.tmpl")
	t := template.Must(template.New("bdd-tmpl").
		Funcs(sprig.FuncMap()).
		Parse(string(content)))

	if err := t.Execute(f, stepListing); err != nil {
		log.Fatal().Err(err).Msg("Failed to render BDD template")
	}
}
