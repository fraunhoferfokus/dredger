// Package handler implementiert Subscriber-Callbacks für Channels
package handler

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"multi-service/async/schemas"
	"context"
)

// HandlerFunc erzeugt eine nats.Msg-Callback-Funktion für den angegebenen Channel
func HandlerFunc(channel string, schemaDir string, fn func(context.Context, []byte) error) nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx := context.Background()
		// Validierung
		if errs, err := schemas.ValidatePayload(schemaDir, channel, msg.Data); err != nil {
			log.Error().Err(err).Str("channel", channel).Msg("Schema read error")
			return
		} else if len(errs) > 0 {
			for _, e := range errs {
				log.Warn().Str("channel", channel).Str("error", e.String()).Msg("Schema violation")
			}
			return
		}
		// Business-Logic
		if err := fn(ctx, msg.Data); err != nil {
			log.Error().Err(err).Str("channel", channel).Msg("Handler failed")
		}
	}
}
