package run

import (
	"errors"
	"fmt"

	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/release"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Run a Bubbly adapter
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Run a Bubbly adapter

		bubbly adapter run
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	var (
		// runner runner.AdapterRun
		trace      bool
		inputFiles string
		// adapterArgs adapter.RunArgs
	)

	cmd := &cobra.Command{
		Use:     "run <name[:tag] | adapter-file> [flags]",
		Short:   "Run a Bubbly adapter",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			adapterID := args[0]
			//
			// Get the bubbly release, through the git commit
			//
			commit, err := release.Commit(bCtx)
			if err != nil {
				return fmt.Errorf("error getting commit from release: %w", err)
			}
			result, err := adapter.RunFromID(bCtx, adapterID,
				adapter.WithInputFiles(inputFiles),
				adapter.WithTracing(trace),
			)
			if err != nil {
				return err
			}

			if result == nil {
				return errors.New("did not receive any results from the adapter")
			}
			if result.CodeScan != nil {
				err := client.SaveCodeScan(bCtx, &api.CodeScanRequest{
					Commit:   &commit,
					CodeScan: result.CodeScan,
				})
				if err != nil {
					return err
				}
			}
			if result.TestRun != nil {
				err := client.SaveTestRun(bCtx, &api.TestRunRequest{
					Commit:  &commit,
					TestRun: result.TestRun,
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&trace, "trace", false, "Enable more verbose trace output for Rego queries")
	f.StringVar(&inputFiles, "input", "", "Provide a comma-separated list of files to parse and provide as input")

	return cmd
}
