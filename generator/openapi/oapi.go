package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func updateOAPIOperation(op *openapi3.Operation, opID string, opSummary string, opDefault string) {
	op.OperationID = opID
	op.Summary = opSummary
	op.Responses.Default().Value = op.Responses.Map()[opDefault].Value
}

func createOAPIResponse(rDesc string) *openapi3.Response {
	r := openapi3.NewResponse()
	r.Description = &rDesc
	return r
}
