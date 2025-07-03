package generator

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func generateInfoFiles(spec *openapi3.T, serverConf ServerConfig) {
	// info.go
	fileName := "info.go"
	filePath := filepath.Join(config.Path, CorePkg, fileName)
	templateFile := "templates/openapi/core/info.go.tmpl"
	createFileFromTemplate(filePath, templateFile, serverConf)

	// infoSvc.go
	fileName = "infoSvc.go"
	filePath = filepath.Join(config.Path, CorePkg, fileName)
	templateFile = "templates/openapi/core/infoSvc.go.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

	// Zusatz‐Endpoints (/infoz etc.) wie bisher …
	if spec.Paths.Find("/infoz") != nil &&
		spec.Paths.Find("/infoz").Operations()[http.MethodGet] != nil &&
		slices.Contains(spec.Paths.Find("/infoz").Operations()[http.MethodGet].Tags, "builtin") {
		log.Debug().Msg("Generating default /infoz endpoint.")
		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service is ready"))
		updateOAPIOperation(op, "GetInfo", "Returns infos about the service", "200")
		spec.AddOperation("/infoz", http.MethodGet, op)
	}
}
