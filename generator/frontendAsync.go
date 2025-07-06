package generator

import (
	fs "dredger/fileUtils"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/rs/zerolog/log"
)

func generateEmptyFrontendAsync(_ *asyncapiv3.Specification, conf GeneratorConfig) {
	frontendPath := filepath.Join(conf.OutputPath, "web")
	fs.GenerateFolder(frontendPath)
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/openapi/web/README.md.tmpl", conf)
}

func generateFrontendAsync(spec *asyncapiv3.Specification, conf GeneratorConfig) {
	// create folders
	corePath := filepath.Join(conf.OutputPath, "core")
	asyncPath := filepath.Join(conf.OutputPath, "async")
	restPath := filepath.Join(conf.OutputPath, "rest")
	frontendPath := filepath.Join(conf.OutputPath, "web")
	javascriptPath := filepath.Join(frontendPath, "js")
	stylesheetPath := filepath.Join(frontendPath, "css")
	imagesPath := filepath.Join(frontendPath, "images")
	fontsPath := filepath.Join(stylesheetPath, "fonts")
	pagesPath := filepath.Join(frontendPath, "pages")
	localesPath := filepath.Join(corePath, "locales")
	publicPath := filepath.Join(frontendPath, "public")
	docPath := filepath.Join(frontendPath, "doc")

	fs.GenerateFolder(frontendPath)
	fs.GenerateFolder(javascriptPath)
	fs.GenerateFolder(stylesheetPath)
	fs.GenerateFolder(imagesPath)
	fs.GenerateFolder(fontsPath)
	fs.GenerateFolder(pagesPath)
	fs.GenerateFolder(localesPath)
	fs.GenerateFolder(publicPath)
	fs.GenerateFolder(docPath)
	fs.GenerateFolder(asyncPath)

	// files in root directory
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/openapi/web/README.md.tmpl", conf)

	// files in javascript directory
	fs.CopyWebFile("openapi/web/js", javascriptPath, "bootstrap.bundle.min.js", true)
	fs.CopyWebFile("openapi/web/js", javascriptPath, "htmx.min.js", true)
	fs.CopyWebFile("openapi/web/js", javascriptPath, "hyperscript.js", true)
	fs.CopyWebFile("openapi/web/js", javascriptPath, "sse.js", true)
	fs.CopyWebFile("openapi/web/js", javascriptPath, "rapidoc-min.js", true)
	fs.CopyWebFile("openapi/web/js", javascriptPath, "elements.min.js", true)

	// files in stylesheet directory
	fs.CopyWebFile("openapi/web/css", stylesheetPath, "bootstrap-icons.min.css", true)
	fs.CopyWebFile("openapi/web/css/fonts", fontsPath, "bootstrap-icons.woff", true)
	fs.CopyWebFile("openapi/web/css/fonts", fontsPath, "bootstrap-icons.woff2", true)
	fs.CopyWebFile("openapi/web/css", stylesheetPath, "bootstrap.min.css", true)
	fs.CopyWebFile("openapi/web/css", stylesheetPath, "pico.min.css", true)
	fs.CopyWebFile("openapi/web/css", stylesheetPath, "pico.colors.min.css", true)
	fs.CopyWebFile("openapi/web/css", stylesheetPath, "elements.min.css", true)

	// files in images directory
	fs.CopyWebFile("openapi/web/images", imagesPath, "favicon.ico", false)

	// files in web directory
	fs.CopyWebFile("openapi/web", frontendPath, "web.go", true)

	// files in core directory
	createFileFromTemplate(filepath.Join(corePath, "localize.go"), "templates/openapi/core/localize.go.tmpl", conf)
	if _, err := os.Stat(filepath.Join(localesPath, "locale.de.toml")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(localesPath, "locale.de.toml"), "templates/openapi/core/locales/locale.de.toml", conf)
		createFileFromTemplate(filepath.Join(localesPath, "locale.en.toml"), "templates/openapi/core/locales/locale.en.toml", conf)
	}

	// files in pages directory
	fs.CopyWebFile("openapi/web/pages", restPath, "render.go", true)
	if _, err := os.Stat(filepath.Join(pagesPath, "languages.templ")); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filepath.Join(pagesPath, "languages.templ"), "templates/openapi/web/pages/languages.templ.tmpl", conf)
	}

	// files in public directory

	tmplData := frontendTemplateConfig{
		Title:    spec.Info.Title,
		Version:  spec.Info.Version,
		Channels: extractChannels(spec),
	}

	createFileFromTemplate(
		filepath.Join(publicPath, "index.html"),
		"templates/openapi/web/public/index.html.tmpl",
		tmplData,
	)
	fs.CopyWebFile(path.Join("openapi/web", "public"), publicPath, "README.md", false)

	// files in doc directory
	fs.CopyWebFile(path.Join("openapi/web", "doc"), docPath, "README.md", false)

	log.Info().Msg("Created Frontend successfully.")
}

type frontendTemplateConfig struct {
	Title    string
	Version  string
	Channels []channelInfo
}

type channelInfo struct {
	Name   string
	Title  string
	Fields []fieldInfo
}

type fieldInfo struct {
	GoName   string
	JSONName string
	Label    string
}

func extractChannels(spec *asyncapiv3.Specification) []channelInfo {
	var channels []channelInfo
	for name, ch := range spec.Channels {
		c := channelInfo{
			Name:  name,
			Title: ch.Description,
		}

		channels = append(channels, c)
	}
	return channels
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
