package generator

import (
	fs "dredger/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateDatabaseFiles(conf GeneratorConfig) {
	log.Info().Msg("Adding SQLite database.")

	fileName := conf.DatabaseName
	filePath := filepath.Join(config.Path, DatabasePkg, fileName)
	templateFile := "templates/db/database.go.tmpl"

	fs.GenerateFile(filePath + ".db")
	createFileFromTemplate(filePath+".go", templateFile, conf)
}
