package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func updateOAPIOperation(op *openapi3.Operation, opID string, opSummary string, opDefault string) {
	op.OperationID = opID
	op.Summary = opSummary
	if opDefaultResp := op.Responses.Map()[opDefault]; opDefaultResp != nil {
		defaultRef := &openapi3.ResponseRef{
			Value: opDefaultResp.Value,
		}
		op.Responses.Set("default", defaultRef)
	}
}

func createOAPIResponse(rDesc string) *openapi3.Response {
	r := openapi3.NewResponse()
	r.Description = &rDesc
	return r
}
