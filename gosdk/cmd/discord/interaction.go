package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func interactionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interaction",
		Short: "Respond to interactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("interaction command placeholder")
			return nil
		},
	}
}
