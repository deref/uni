package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const rootDescription = "Unirepo is a tool for managing uniform TypeScript monorepos."

var rootCmd = &cobra.Command{
	Use:   "uni",
	Short: rootDescription,
	Long:  rootDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
