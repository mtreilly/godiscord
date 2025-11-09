package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "discord",
		Short: "Discord SDK CLI",
	}

	rootCmd.AddCommand(webhookCmd())
	rootCmd.AddCommand(messageCmd())
	rootCmd.AddCommand(channelCmd())
	rootCmd.AddCommand(guildCmd())
	rootCmd.AddCommand(interactionCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
