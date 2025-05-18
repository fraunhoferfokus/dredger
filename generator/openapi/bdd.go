package generator

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

// path is the path of the feature file
func GenerateBdd(path string) {
	var step Listing
	step.Steps = parseSteps(path)
	step.UniqueEndpoints = getAllEndpoints(step)

	_, e := os.Create("server_godog_test.go")
	if e != nil {
		log.Fatal()
	}

	f, err := os.OpenFile("server_godog_test.go", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	content, _ := os.ReadFile("templates/bdd.go.tmpl")
	t := template.Must(template.New("bdd-tmpl").Funcs(sprig.FuncMap()).Parse(string(content)))
	err1 := t.Execute(f, step)
	if err1 != nil {
		panic(err1)
	}
}
