package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func channelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "channel",
		Short: "Manage channels",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("channel command placeholder")
			return nil
		},
	}
}
