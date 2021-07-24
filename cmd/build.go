package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

var buildOpts internal.BuildOptions

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&buildOpts.Version, "version", "", "version to put in package.json")
	buildCmd.Flags().BoolVar(&buildOpts.Watch, "watch", false, "rebuilds each time source files change")
	buildCmd.Flags().BoolVar(&buildOpts.Types, "types", false, "also build a .d.ts file")
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

		var packages map[string]*internal.Package
		switch len(args) {
		case 0:
			packages = repo.Packages
		case 1:
			pkgName := args[0]
			pkg, ok := repo.Packages[pkgName]
			if !ok {
				return fmt.Errorf("no such package: %q", pkgName)
			}
			packages = map[string]*internal.Package{
				pkgName: pkg,
			}
		default:
			panic("unreachable")
		}

		// TODO: Parallelism.
		for _, pkg := range packages {
			buildOpts.Package = pkg
			if err := internal.Build(repo, buildOpts); err != nil {
				return err
			}
		}
		return nil
	},
}
