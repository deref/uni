package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
			Watch: watch,
		}

		var err error
		opts.Entrypoint, err = filepath.Abs(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		opts.Args = args[1:]

		return internal.Run(repo, opts)
	},
}
