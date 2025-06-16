// multi-service/async/utils — Helfer für NATS-Subscribe/Publish
package utils

import (
    "github.com/nats-io/nats.go"
    "github.com/rs/zerolog/log"
)

// Subscribe abonniert ein Subject und startet Handler.
func Subscribe(nc *nats.Conn, subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
    sub, err := nc.Subscribe(subject, handler)
    if err != nil {
        log.Error().Err(err).Str("subject", subject).Msg("Subscribe failed")
    }
    return sub, err
}

// QueueSubscribe abonniert mit Queue-Gruppe.
func QueueSubscribe(nc *nats.Conn, subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
    sub, err := nc.QueueSubscribe(subject, queue, handler)
    if err != nil {
        log.Error().Err(err).Str("subject", subject).Str("queue", queue).Msg("QueueSubscribe failed")
    }
    return sub, err
}

// Publish sendet eine Nachricht.
func Publish(nc *nats.Conn, subject string, data []byte) error {
    err := nc.Publish(subject, data)
    if err != nil {
        log.Error().Err(err).Str("subject", subject).Msg("Publish failed")
    }
    return err
}
