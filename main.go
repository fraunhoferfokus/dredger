package main

import (
	"embed"

	"dredger/cli"

	// Wir brauchen hier das OpenAPI-Generator‐Package,
	// damit wir ihm die eingebetteten Templates geben können:
	genOpenAPI "dredger/generator/openapi"
)

//go:embed templates
var tmplFS embed.FS

func main() {
	// Hier übergeben wir die eingebetteten Dateien an das OpenAPI‐Generator‐Package:
	genOpenAPI.TmplFS = tmplFS

	// Jetzt startet die CLI wie gewohnt:
	cli.Execute()
}
