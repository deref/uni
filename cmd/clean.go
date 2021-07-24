package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes build output.",
	Long:  `Removes build output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		return internal.Clean(repo)
	},
}
