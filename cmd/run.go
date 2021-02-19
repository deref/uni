package cmd

import (
	"github.com/brandonbloom/unirepo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <script> [args...]",
	Short: "Build and run an entrypoint.",
	Long: `Builds and runs the given entrypoint file.
The script must export a function called "main", which will receive the
given string args.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Run(args[0], args[1:])
	},
}
