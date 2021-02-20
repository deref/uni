package cmd

import (
	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [...packages]",
	Short: "Builds packages targeting Node.",
	Long:  `Builds packages targeting Node. TODO: Target browsers. TODO: Say more.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		pkgName := "@unirepo/example-util" // XXX get from args, have --all flag.
		return internal.Build(repo, pkgName)
	},
}
