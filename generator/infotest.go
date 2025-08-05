package generator

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

func generateInfoFiles(spec *openapi3.T, serverConf ServerConfig) {
	// info.go
	fileName := "info.go"
	filePath := filepath.Join(Config.Path, CorePkg, fileName)
	templateFile := "templates/common/core/info.go.tmpl"
	createFileFromTemplate(filePath, templateFile, serverConf)

	// infoSvc.go
	fileName = "infoSvc.go"
	filePath = filepath.Join(Config.Path, CorePkg, fileName)
	templateFile = "templates/common/core/infoSvc.go.tmpl"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

	// Zusatz‐Endpoints (/infoz etc.) wie bisher …
	if spec.Paths.Find("/infoz") != nil &&
		spec.Paths.Find("/infoz").Operations()[http.MethodGet] != nil &&
		slices.Contains(spec.Paths.Find("/infoz").Operations()[http.MethodGet].Tags, "builtin") {
		//log.Debug().Msg("Generating default /infoz endpoint.")
		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service is ready"))
		updateOAPIOperation(op, "GetInfo", "Returns infos about the service", "200")
		spec.AddOperation("/infoz", http.MethodGet, op)
	}
}

func generateInfoFilesAsync(spec *asyncapiv3.Specification, serverConf ServerConfig) {
	// info.go
	fileName := "info.go"
	filePath := filepath.Join(Config.Path, CorePkg, fileName)
<<<<<<< HEAD:generator/infotest.go
<<<<<<< HEAD
	templateFile := "templates/common/core/info.go.tmpl"
=======
	templateFile := "templates/openapi/core/info.go.tmpl"
>>>>>>> 8c640d3 (Added function for async)
=======
	templateFile := "templates/openapi/core/info.go.tmpl"
>>>>>>> 8c640d3 (Added function for async):generator/info.go
	createFileFromTemplate(filePath, templateFile, serverConf)

	// infoSvc.go
	fileName = "infoSvc.go"
	filePath = filepath.Join(Config.Path, CorePkg, fileName)
<<<<<<< HEAD:generator/infotest.go
<<<<<<< HEAD
	templateFile = "templates/common/core/infoSvc.go.tmpl"
=======
	templateFile = "templates/openapi/core/infoSvc.go.tmpl"
>>>>>>> 8c640d3 (Added function for async)
=======
	templateFile = "templates/openapi/core/infoSvc.go.tmpl"
>>>>>>> 8c640d3 (Added function for async):generator/info.go
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		createFileFromTemplate(filePath, templateFile, serverConf)
	}

}
