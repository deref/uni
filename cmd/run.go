package cmd

import (
	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [flags] <script> [args...]",
	Short: "Build and run an entrypoint.",
	Long: `Builds and runs the given entrypoint file.
The script must export a function called "main", which will receive the
given string args.`,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		return internal.Run(repo, args[0], args[1:])
	},
}
