package cmd

import (
	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes build output.",
	Long:  `Removes build output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Clean()
	},
}
