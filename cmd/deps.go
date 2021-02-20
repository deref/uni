package cmd

import (
	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(depsCmd)
}

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install dependencies",
	Long:  `Generates a root package.json file and runs your package manager.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		return internal.Deps(repo)
	},
}
