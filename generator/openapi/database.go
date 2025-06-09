package generator

import (
	"path/filepath"

	fs "dredger/fileUtils"

	"github.com/rs/zerolog/log"
)

func generateDatabaseFiles(conf GeneratorConfig) {
	log.Info().Msg("Adding SQLite database.")

	// erzeugt <DatabaseName>.db und <DatabaseName>.go
	fileName := conf.DatabaseName + ".db"
	filePath := filepath.Join(config.Path, DatabasePkg, fileName)
	fs.GenerateFile(filePath)

	goFile := filepath.Join(config.Path, DatabasePkg, conf.DatabaseName+".go")
	templateFile := "templates/openapi/db/database.go.tmpl"
	createFileFromTemplate(goFile, templateFile, conf)
}
