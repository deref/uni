package cmd

import (
	"fmt"

	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(packCmd)
}

var packCmd = &cobra.Command{
	Use:   "pack [package]",
	Short: "Packs a pre-built package into a tgz file.",
	Long: `Packs a package into a tgz file.
The package must already be built. Use the build command.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		switch len(args) {
		case 0:
			// TODO: Parallelism.
			for pkgName, pkg := range repo.Packages {
				fmt.Println("packing", pkgName)
				if err := internal.Pack(repo, pkg); err != nil {
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
			return internal.Pack(repo, pkg)
		default:
			panic("unreachable")
		}
	},
}
