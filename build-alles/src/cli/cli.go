// multi-service/async/cli — Kommandozeilen-Interface mit cobra
package cli

import (
    "github.com/spf13/cobra"
    "os"
)

func NewAsyncCommand(run func(channel, message string) error) *cobra.Command {
    var channel string
    var message string

    cmd := &cobra.Command{
        Use:   "async",
        Short: "Sendet Nachrichten gemäß AsyncAPI-Spec über NATS",
        RunE: func(_ *cobra.Command, _ []string) error {
            return run(channel, message)
        },
    }

    cmd.Flags().StringVarP(&channel, "channel", "c", "", "Ziel-Channel (Subject)")
    cmd.Flags().StringVarP(&message, "message", "m", "", "Nachrichten-Payload (JSON)")
    _ = cmd.MarkFlagRequired("channel")

    return cmd
}

// In deinem main.go bindest du diesen Befehl dann z.B. so:
// rootCmd.AddCommand(cli.NewAsyncCommand(asyncRun))
