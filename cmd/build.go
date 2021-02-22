package cmd

import (
	"fmt"

	"github.com/brandonbloom/uni/internal"
	"github.com/spf13/cobra"
)

var buildOpts internal.BuildOptions

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&buildOpts.Version, "version", "", "version to put in package.json")
}

var buildCmd = &cobra.Command{
	Use:   "build [package]",
	Short: "Builds packages targeting Node.",
	Long: `Builds packages targeting Node.
Given no arguments, builds all packages. Otherwise, builds only the specified package.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}
		switch len(args) {
		case 0:
			// TODO: Parallelism.
			for pkgName, pkg := range repo.Packages {
				buildOpts.Package = pkg
				fmt.Println("building", pkgName)
				if err := internal.Build(repo, buildOpts); err != nil {
					return err
				}
			}
			return nil
		case 1:
			pkgName := args[0]
			pkg, ok := repo.Packages[pkgName]
			if !ok {
				return fmt.Errorf("no such package: %q", pkgName)
			}
			buildOpts.Package = pkg
			return internal.Build(repo, buildOpts)
		default:
			panic("unreachable")
		}
	},
}
