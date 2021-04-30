package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/deref/uni/internal"
	"github.com/spf13/cobra"
)

var runOpts = internal.RunOptions{}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&runOpts.Watch, "watch", false, "re-runs command when source files change")
	runCmd.Flags().BoolVar(&runOpts.BuildOnly, "build-only", false, "(internal) exit before running, skip temporary file cleanup, and print path to build output")
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

Unhandled exceptions and promise rejections will be logged to stderr and the
process will immediately exit with status code 1.

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

		var err error
		runOpts.Entrypoint, err = filepath.Abs(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		runOpts.Args = args[1:]

		ctx, cancel := context.WithCancel(cmd.Context())
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT)
		go func() {
			_ = <-shutdown
			fmt.Fprintf(os.Stderr, "Received shutdown signal. Waiting for process to finish.\n")
			cancel()
		}()

		err = internal.Run(ctx, repo, runOpts)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	},
}
