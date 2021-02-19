package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(depsCmd)
}

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install dependencies",
	Long:  `Generates a package.json file and runs your package manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deps")
	},
}
