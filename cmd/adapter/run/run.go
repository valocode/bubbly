package run

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
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
		basedir  string
		filename string
		name     string
		tag      string
		// adapterArgs adapter.RunArgs
	)

	cmd := &cobra.Command{
		Use:     "run [name[:tag]] [flags]",
		Short:   "Run a Bubbly adapter",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// Then user must provide --from-file
				if filename == "" {
					return errors.New("provide either an adapter name[:tag] or a --from-file")
				}
			}
			if len(args) == 1 {
				// TODO: warn about the --from-file overriding the adapter name?
				// if opts.Filename != "" {
				// }
				var err error
				name, tag, err = adapter.ParseAdpaterID(args[0])
				if err != nil {
					return err
				}
			}
			// First attempt is to try to read the file locally.
			// This depends on the user provided options - if the user provided options
			// to read locally and it fails, then fail the run. Otherwise continue
			// to fetching the adapter remotely
			var (
				adapterPath string
				failOnError bool
			)
			if filename != "" {
				adapterPath = filename
				failOnError = true
			} else {
				adapterPath = path.Join(basedir, "adapters", name+".rego")
			}
			if basedir != ".bubbly" {
				failOnError = true
			}

			var result *adapter.AdapterResult
			if _, err := os.Stat(adapterPath); err == nil {
				result, err = adapter.RunFromFile(adapterPath)
				if err != nil {
					return fmt.Errorf("error running adapter: %w", err)
				}
			} else if os.IsNotExist(err) {
				if failOnError {
					return fmt.Errorf("adapter does not exist: %s", adapterPath)
				}
			} else {
				return fmt.Errorf("error checking for adapter existence: %w", err)
			}

			if result == nil {
				a, err := client.GetAdapter(bCtx, &api.AdapterGetRequest{
					Name: &name,
					Tag:  &tag,
				})
				if err != nil {
					return err
				}
				result, err = adapter.Run(*a.Module)
				if err != nil {
					return fmt.Errorf("error running adapter: %w", err)
				}
			}

			commit := "asdasd"
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

	f := cmd.PersistentFlags()
	f.StringVarP(
		&filename,
		"from-file", "f",
		"",
		`Get the adapter to run from a specific file containing an adapter.`,
	)
	f.StringVarP(
		&basedir,
		"directory", "d",
		".bubbly",
		`Use the following directory as the base instead of the default ".bubbly".`,
	)
	// 	f.StringVar(
	// 		&adapterArgs.Filename,
	// 		"arg-file",
	// 		"",
	// 		`An argument for the adapter providing an input file.
	// Only valid for adapter types that read an input file (json, yaml, csv, etc).
	// Default is <name>.adapter.<type>, e.g. snyk.adapter.json`,
	// 	)

	return cmd
}
