package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

var depsOpts internal.InstallDependenciesOptions

func init() {
	rootCmd.AddCommand(depsCmd)
	depsCmd.Flags().BoolVar(&depsOpts.Frozen, "frozen", false, "(unstable) prevents modification of dependency lock file")
}

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install dependencies",
	Long:  `Generates a root package.json, tsconfig.json and other config files and runs your package manager.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}
		return internal.InstallDependencies(repo, depsOpts)
	},
}
