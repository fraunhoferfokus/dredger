# ENVIRONMENT

| Variable                   | Beschreibung                                            | Default     |
|----------------------------|----------------------------------------------------------|-------------|
| SERVICE_PORT               | Port, auf dem der NATS-Listener Ihr Service bindet      | `4222`      |
| SERVER_URL                 | NATS-Server URL                                         | `nats://localhost:4222` |
| LOG_LEVEL                  | Log-Level (`DEBUG`, `INFO`, `WARN`, `ERROR`)            | `INFO`      |
| TRACING_ENABLED            | OpenTelemetry aktivieren (`true`/`false`)               | `false`     |

Alle weiteren Channel-Settings werden automatisch aus der Spec generiert.
