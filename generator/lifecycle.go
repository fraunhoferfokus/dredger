package generator

import (
	fs "dredger/fileUtils"
	"net/http"
	"path/filepath"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func generateLifecycleFiles(spec *openapi3.T, conf GeneratorConfig) {
	if spec.Paths.Find("/livez") == nil || (spec.Paths.Find("/livez").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/livez").Operations()[http.MethodGet].Tags, "builtin")) {
		log.Debug().Msg("Generating default /livez endpoint.")

		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The service is alive"))
		updateOAPIOperation(op, "GetLive", "Returns live-state of the service", "200")
		spec.AddOperation("/livez", http.MethodGet, op)
	}

	if spec.Paths.Find("/readyz") == nil || (spec.Paths.Find("/readyz").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/readyz").Operations()[http.MethodGet].Tags, "builtin")) {
		log.Debug().Msg("Generating default /readyz endpoint.")

		op := openapi3.NewOperation()
		op.AddResponse(200, createOAPIResponse("The service is ready"))
		op.AddResponse(http.StatusOK, createOAPIResponse("The service is not ready"))
		updateOAPIOperation(op, "GetReady", "Returns ready-state of the service", "200")
		spec.AddOperation("/readyz", http.MethodGet, op)
	}

	if spec.Paths.Find("/robots.txt") == nil || (spec.Paths.Find("/robots.txt").Operations()[http.MethodGet] != nil && slices.Contains(spec.Paths.Find("/robots.txt").Operations()[http.MethodGet].Tags, "builtin")) {
		log.Debug().Msg("Generating default /robots.txt endpoint.")

		op := openapi3.NewOperation()
		op.AddResponse(http.StatusOK, createOAPIResponse("The robots.txt"))
		updateOAPIOperation(op, "GetRobots", "Returns robots.txt", "200")
		spec.AddOperation("/robots.txt", http.MethodGet, op)
	}

	restPath := filepath.Join(conf.OutputPath, "rest")
	fs.CopyWebFile("web", restPath, "robots.txt", false)
}
