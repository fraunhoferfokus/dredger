// Edit this file, as it is a specific handler function for your service
package rest

import (
	"build/core/log"
	"build/core/tracing"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Returns a list of all the stores.
func GetAllStores(c echo.Context) error {
	// trace span
	ctx := c.Request().Context()
	ctx, span := tracing.Tracer.Start(ctx, "logMessage")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	log.Info().Str("traceId", traceId).Str("spanId", spanId).Str("path", "/").Msg("GetAllStores")

	// session, err := getSession(c)
	// if err != nil {
	// 	log.Error().Err(err).Msg("GetAllStores failed")
	// 	return c.NoContent(http.StatusInternalServerError)
	// }

	// implement your functionality best using a function from a separate file, e.g. usecases/GetAllStoresDo.go

	// 200 => An array of Store objects.
	return c.String(http.StatusNotImplemented, "Temporary handler stub.")
}
