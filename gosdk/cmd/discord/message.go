package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func messageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "message",
		Short: "Send or edit messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("message command placeholder")
			return nil
		},
	}
}
