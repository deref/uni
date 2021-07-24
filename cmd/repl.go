package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

func init() {
	rootCmd.AddCommand(replCmd)
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start a Read Evaluate Print Loop.",
	Long:  "Start a Read Evaluate Print Loop.",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}

		node := exec.Command("node", "--interactive")
		node.Stdin = os.Stdin
		node.Stdout = os.Stdout
		node.Stderr = os.Stderr

		err := node.Run()

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	},
}
