package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

func init() {
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish [package]",
	Short: "Publishes pre-packed tgz files.",
	Long: `Publishes pre-packed tgz files.
The package must already be packed. Use the pack command.`,
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
				fmt.Println("publishing", pkgName)
				if err := internal.Publish(repo, pkg); err != nil {
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
			return internal.Publish(repo, pkg)
		default:
			panic("unreachable")
		}
	},
}
