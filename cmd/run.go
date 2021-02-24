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

Entrypoint files are modules that export a "main" function, which will be called
with the given 'args' as positional parameters.

Main functions may return an integer status code, and the return value will
will be awaited.  If no status code is returned, the default is 0.

After awaiting a return value, the process will be terminated immediately.  Any
pending events will not be executed; main is responsible for graceful shutdown.

Example:

export const main = async (...args: string[]) => {
  console.log("see uni run");
  return 0; // Return an exit code (optional).
}
`,
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
