package main

import (
	"embed"

	"dredger/cli"

	// Wir brauchen hier die Generator-Pakete,
	// damit wir ihnen die eingebetteten Templates geben können:
	genAsyncAPI "dredger/generator/asyncapi"
	genOpenAPI "dredger/generator/openapi"
)

//go:embed templates
var tmplFS embed.FS

func main() {
	// Hier übergeben wir die eingebetteten Dateien an die Generator-Packages:
	genOpenAPI.TmplFS = tmplFS
	genAsyncAPI.TmplFS = tmplFS

	// Jetzt startet die CLI wie gewohnt:
	cli.Execute()
}
