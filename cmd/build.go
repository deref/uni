package cmd

import (
	"fmt"

	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [package]",
	Short: "Builds packages targeting Node.",
	Long: `Builds packages targeting Node.
Given no arguments, builds all packages. Otherwise, builds only the specified package.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		switch len(args) {
		case 0:
			// TODO: Parallelism.
			for pkgName := range repo.Packages {
				fmt.Println("building", pkgName)
				if err := internal.Build(repo, pkgName); err != nil {
					return err
				}
			}
			return nil
		case 1:
			pkgName := args[0]
			return internal.Build(repo, pkgName)
		default:
			panic("unreachable")
		}
	},
}
