package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var verbose bool

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&verbose, "verbose", false, "Prints version information for all dependencies too.")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of monoclean.",
	Long:  "Print the version of monoclean.",
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			panic("debug.ReadBuildInfo() failed")
		}
		printInfo := func(mod debug.Module) {
			if verbose {
				fmt.Println(mod.Path, mod.Version)
			} else {
				fmt.Println(mod.Version)
			}
		}
		printInfo(buildInfo.Main)
		if verbose {
			for _, dep := range buildInfo.Deps {
				fmt.Println(dep.Path, dep.Version)
			}
		}
	},
}
