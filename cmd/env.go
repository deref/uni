package cmd

import (
	"fmt"
	"os"

	"github.com/brandonbloom/uni/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Dump environment information.",
	Long: `Validates and prints information about the current environment.

Returns a non-zero status code if there environment fails any checks.
The printed output of this command is intended for human consumption and
not (yet?) intended to be parsed.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := mustLoadRepository()
		env, err := internal.AnalyzeEnvironment(repo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		internal.DumpEnvironment(env)
		if env.Erred {
			os.Exit(1)
		}
	},
}
