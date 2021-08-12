package run

import (
	"errors"
	"fmt"
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
		opts        adapter.RunOptions
		adapterArgs adapter.RunArgs
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
				if opts.Filename == "" {
					return errors.New("provide either an adapter name[:tag] or a --from-file to read the adapter from")
				}
			}
			if len(args) == 1 {
				// TODO: warn about the --from-file overriding the adapter name?
				// if opts.Filename != "" {
				// }
				var err error
				opts.Name, opts.Tag, err = adapter.ParseAdpaterID(args[0])
				if err != nil {
					return err
				}
			}
			// First attempt is to try to read the file locally.
			// This depends on the user provided options - if the user provided options
			// to read locally and it fails, then fail the run. Otherwise continue
			// to fetching the adapter remotely
			var (
				// a           *adapter.Adapter
				adapterPath string
				failOnError bool
			)
			if opts.Filename != "" {
				adapterPath = opts.Filename
				failOnError = true
			}
			if opts.Filename == "" && opts.BaseDir != "" {
				adapterPath = path.Join(opts.BaseDir, opts.Name+".adapter.hcl")
				failOnError = true
			}
			if adapterPath == "" {
				// Set the default adapter path
				adapterPath = path.Join(".bubbly", opts.Name+".adapter.hcl")
			}
			a, err := adapter.FromFile(adapterPath)
			if err != nil {
				if failOnError {
					return err
				}
				// Else, try to fetch the adapter remotely
			}
			if a == nil {
				a, err = client.GetAdapter(bCtx, &api.AdapterGetRequest{
					Name: &opts.Name,
					Tag:  &opts.Tag,
				})
				if err != nil {
					return err
				}
			}

			output, err := a.Run(adapterArgs)
			if err != nil {
				return err
			}

			commit := "asdasd"
			if output.HasCodeScan() {
				err := client.SaveCodeScan(bCtx, &api.CodeScanRequest{
					Commit:   &commit,
					CodeScan: output.CodeScan,
				})
				if err != nil {
					return fmt.Errorf("error saving code scan: %w", err)
				}
			}
			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.StringVarP(
		&opts.Filename,
		"from-file", "f",
		"",
		`Get the adapter to run from a specific file containing an adapter.`,
	)
	f.StringVar(
		&adapterArgs.Filename,
		"arg-file",
		"",
		`An argument for the adapter providing an input file.
Only valid for adapter types that read an input file (json, yaml, csv, etc).
Default is <name>.adapter.<type>, e.g. snyk.adapter.json`,
	)

	return cmd
}
