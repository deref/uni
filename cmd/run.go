package cmd

import (
	"github.com/brandonbloom/uni/internal"
	"github.com/spf13/cobra"
)

var watch bool

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&watch, "watch", false, "re-runs command when source files change")
}

var runCmd = &cobra.Command{
	Use:   "run [flags] <script> [args...]",
	Short: "Build and run an entrypoint.",
	Long: `Builds and runs the given entrypoint file.
The script must export a function called "main", which will receive the
given string args.`,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}
		opts := internal.RunOptions{
			Watch:      watch,
			Entrypoint: args[0],
			Args:       args[1:],
		}
		return internal.Run(repo, opts)
	},
}
