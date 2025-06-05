// Edit this file, as it is a specific handler function for your service
package rest

import (
	"simpleApi/core/log"
	"simpleApi/core/tracing"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Show API version details
func GetVersionDetailsv2(c echo.Context) error {
	// trace span
	ctx := c.Request().Context()
	ctx, span := tracing.Tracer.Start(ctx, "logMessage")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	log.Info().Str("traceId", traceId).Str("spanId", spanId).Str("path", "/").Msg("GetVersionDetailsv2")

	// session, err := getSession(c)
	// if err != nil {
	// 	log.Error().Err(err).Msg("GetVersionDetailsv2 failed")
	// 	return c.NoContent(http.StatusInternalServerError)
	// }

	// implement your functionality best using a function from a separate file, e.g. usecases/GetVersionDetailsv2Do.go

	// 200 => 200 response
	// 203 => 203 response
	return c.String(http.StatusNotImplemented, "Temporary handler stub.")
}
