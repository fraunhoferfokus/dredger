package generator

import (
	"os"
	"path"
	"strings"
	"text/template"

	fs "dredger/fileUtils"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobeam/stringy"
	"github.com/rs/zerolog/log"
)

func snakecase(s string) string {
	return stringy.New(s).SnakeCase("?", "").Get()
}

func camelcase(s string) string {
	// return stringy.New(s).CamelCase("?", "", "#", "").Get()  // from fokus
	if strings.ContainsAny(s, "-_.") {
		return stringy.New(s).CamelCase().Get()
	}
	return s
}

func lcFirst(s string) string {
	return stringy.New(s).LcFirst()
}

func ucFirst(s string) string {
	return stringy.New(s).UcFirst()

func createFileFromTemplate(filePath string, tmplPath string, config interface{}) {
	templateName := path.Base(tmplPath)
	funcmap := sprig.FuncMap()
	funcmap["camelcase"] = camelcase
	funcmap["snakecase"] = snakecase
	funcmap["lcfirst"] = lcFirst
	funcmap["ucfirst"] = ucFirst

	// Create file and open it
	fs.GenerateFile(filePath)
	file, fErr := os.OpenFile(filePath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		log.Fatal().Err(fErr).Msg("Failed creating file.")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into file
	tmpl := template.Must(template.New(templateName).Funcs(funcmap).ParseFS(TmplFS, tmplPath))
	tmplErr := tmpl.Execute(file, config)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template.")
		panic(tmplErr)
	}

	log.Info().Str("template", templateName).Msg("CREATE " + filePath)
}

func createFileFromTemplates(filePath string, tmplPaths []string, config interface{}) {
	templateName := path.Base(tmplPaths[0])

	// Create file and open it
	fs.GenerateFile(filePath)
	file, fErr := os.OpenFile(filePath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		log.Fatal().Err(fErr).Msg("Failed creating file.")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into file
	tmpl := template.Must(template.New(templateName).Funcs(sprig.FuncMap()).ParseFS(TmplFS, tmplPaths...))
	tmplErr := tmpl.Execute(file, config)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template.")
		panic(tmplErr)
	}

	log.Info().Msg("CREATE " + filePath)
}
