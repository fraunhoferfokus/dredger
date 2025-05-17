// Edit this file, as it is a specific handler function for your service
package rest

import (
	"build/core/log"
	"build/core/tracing"
	"build/entities"

	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

// Creates a new store.
func CreateStore(c echo.Context) error {
	// trace span
	ctx := c.Request().Context()
	ctx, span := tracing.Tracer.Start(ctx, "logMessage")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()
	log.Info().Str("traceId", traceId).Str("spanId", spanId).Str("path", "/").Msg("CreateStore")

	// session, err := getSession(c)
	// if err != nil {
	// 	log.Error().Err(err).Msg("CreateStore failed")
	// 	return c.NoContent(http.StatusInternalServerError)
	// }

	// Parse body
	contentType := c.Request().Header.Get("Content-Type")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Error().Err(err).Str("traceId", traceId).Str("spanId", spanId).Msg("Can't read body")
		return c.String(http.StatusBadRequest, "Can't read body")
	}

	var content entities.Store

	if contentType == "application/json" {
		err := json.Unmarshal(body, &content)
		if err != nil {
			log.Error().Err(err).Str("traceId", traceId).Str("spanId", spanId).Msg("Could not unmarshall JSON input")
			return c.String(http.StatusBadRequest, "Could not unmarshall JSON input")
		}
	} else if contentType == "application/yaml" {
		err := yaml.Unmarshal(body, &content)
		if err != nil {
			log.Error().Err(err).Str("traceId", traceId).Str("spanId", spanId).Msg("Could not unmarshall YAML input")
			return c.String(http.StatusBadRequest, "Could not unmarshall YAML input")
		}
	} else if contentType == "application/xml" {
		err := xml.Unmarshal(body, &content)
		if err != nil {
			log.Error().Err(err).Str("traceId", traceId).Str("spanId", spanId).Msg("Could not unmarshall XML input")
			return c.String(http.StatusBadRequest, "Could not unmarshall XML input")
		}
	} else {
		log.Error().Str("content type", contentType).Str("traceId", traceId).Str("spanId", spanId).Msg("Wrong content type")
		return c.String(http.StatusUnsupportedMediaType, "Wrong content type")
	}

	// validate content
	if err := content.Validate(); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// var content contains the payload from the request body

	// implement your functionality best using a function from a separate file, e.g. usecases/CreateStoreDo.go

	// 200 => A store was created successfully.
	// 400 => Invalid store properties.
	return c.String(http.StatusNotImplemented, "Temporary handler stub.")
}
