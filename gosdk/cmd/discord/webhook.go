package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func webhookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Interact with configured webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("webhook command placeholder (send/list)")
			return nil
		},
	}
	return cmd
}
