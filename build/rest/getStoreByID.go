// Edit this file, as it is a specific handler function for your service
package rest

import (
	"build/core/log"
	"build/core/tracing"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Returns a store with the given store id.
func GetStoreByID(c echo.Context) error {
	// trace span
	ctx := c.Request().Context()
	ctx, span := tracing.Tracer.Start(ctx, "logMessage")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	log.Info().Str("traceId", traceId).Str("spanId", spanId).Str("path", "/").Msg("GetStoreByID")

	// session, err := getSession(c)
	// if err != nil {
	// 	log.Error().Err(err).Msg("GetStoreByID failed")
	// 	return c.NoContent(http.StatusInternalServerError)
	// }

	id := c.Param("id")

	// implement your functionality best using a function from a separate file, e.g. usecases/GetStoreByIDDo.go

	// 200 => A store was returned successfully.
	// 400 => Invalid store id.
	// 404 => Store with given id wasn't found.
	return c.String(http.StatusNotImplemented, "Temporary handler stub.")
}
