// multi-service/async/warp_server — minimaler HTTP-Server für Health-Checks
package warp_server

import (
    "fmt"
    "net/http"

    "github.com/rs/zerolog/log"
)

func Start(port int) {
    http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("healthy"))
    })
    addr := fmt.Sprintf(":%d", port)
    log.Info().Str("addr", addr).Msg("Starting HTTP health server")
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal().Err(err).Msg("HTTP server failed")
    }
}
