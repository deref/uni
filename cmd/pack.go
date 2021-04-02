package cmd

import (
	"fmt"

	"github.com/deref/uni/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(packCmd)
}

var packCmd = &cobra.Command{
	Use:   "pack [package]",
	Short: "Packs pre-built packages into a tgz files.",
	Long: `Packs packages into tgz files.
Prints the absolute path of each created package.
The package must already be built. Use the build command.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()

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
			res, err := internal.Pack(repo, pkg)
			if err != nil {
				return err
			}
			fmt.Println(res.PackagePath)
		}
		return nil
	},
}
