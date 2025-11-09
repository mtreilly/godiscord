package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func guildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "guild",
		Short: "Query guild metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("guild command placeholder")
			return nil
		},
	}
}
