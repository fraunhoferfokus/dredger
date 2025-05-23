// Edit this file, as it is a specific handler function for your service
package rest

import (
	"build/core/log"
	"build/core/tracing"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Deletes a store with the given store id.
func DeleteStoreByID(c echo.Context) error {
	// trace span
	ctx := c.Request().Context()
	ctx, span := tracing.Tracer.Start(ctx, "logMessage")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	log.Info().Str("traceId", traceId).Str("spanId", spanId).Str("path", "/").Msg("DeleteStoreByID")

	// session, err := getSession(c)
	// if err != nil {
	// 	log.Error().Err(err).Msg("DeleteStoreByID failed")
	// 	return c.NoContent(http.StatusInternalServerError)
	// }

	id := c.Param("id")

	// implement your functionality best using a function from a separate file, e.g. usecases/DeleteStoreByIDDo.go

	// 400 => Invalid store id.
	// 404 => Store with given id wasn't found.
	// 200 => A store was deleted successfully.
	return c.String(http.StatusNotImplemented, "Temporary handler stub.")
}
