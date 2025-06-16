// Edit this file, as it is a specific handler function for your service
package rest

import (
	"multiService/core/log"
	"multiService/core/tracing"

	"net/http"

	"github.com/labstack/echo/v4"
)

// List API versions
func ListVersionsv2(c echo.Context) error {
	// trace span
	ctx := c.Request().Context()
	ctx, span := tracing.Tracer.Start(ctx, "logMessage")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	log.Info().Str("traceId", traceId).Str("spanId", spanId).Str("path", "/").Msg("ListVersionsv2")

	// session, err := getSession(c)
	// if err != nil {
	// 	log.Error().Err(err).Msg("ListVersionsv2 failed")
	// 	return c.NoContent(http.StatusInternalServerError)
	// }

	// implement your functionality best using a function from a separate file, e.g. usecases/ListVersionsv2Do.go

	// 300 => 300 response
	// 200 => 200 response
	return c.String(http.StatusNotImplemented, "Temporary handler stub.")
}
