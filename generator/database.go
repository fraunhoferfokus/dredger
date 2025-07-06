package generator

import (
	"path/filepath"

	fs "dredger/fileUtils"

	"github.com/rs/zerolog/log"
)

// Kann für OpenAPI als auch für AsyncAPI verwendet werden
func generateDatabaseFiles(conf GeneratorConfig) {
	log.Info().Msg("Adding SQLite database.")

	// erzeugt <DatabaseName>.db und <DatabaseName>.go
	fileName := conf.DatabaseName + ".db"
	filePath := filepath.Join(conf.OutputPath, DatabasePkg, fileName)
	fs.GenerateFile(filePath)

	goFile := filepath.Join(conf.OutputPath, DatabasePkg, conf.DatabaseName+".go")
	templateFile := "templates/openapi/db/database.go.tmpl"
	createFileFromTemplate(goFile, templateFile, conf)
}
