package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/color"
	"github.com/open-policy-agent/opa/tester"
	"github.com/valocode/bubbly/adapter"
	"github.com/valocode/bubbly/env"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Run Bubbly adapter tests
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Run Bubbly adapter tests

		bubbly adapter test
		`,
	)
)

func New(bCtx *env.BubblyConfig) *cobra.Command {
	var (
		trace  bool
		output string
	)

	cmd := &cobra.Command{
		Use:     "test <directory> [flags]",
		Short:   "Run Bubbly adapter tests in the given directory",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dirpath := args[0]
			dir, err := os.ReadDir(dirpath)
			if err != nil {
				return err
			}
			var files []string
			for _, d := range dir {
				if filepath.Ext(d.Name()) == ".rego" {
					files = append(files, filepath.Join(dirpath, d.Name()))
				}
			}

			results, err := adapter.RunTests(context.Background(), files)
			if err != nil {
				return err
			}

			switch output {
			case "json":
				printJSON(results)
			case "":
				// Do nothing
			default:
				return fmt.Errorf("unknown output format: %s", output)

			}
			var (
				testsPass int
				testsFail int
				testsSkip int
			)
			for _, result := range results {
				switch {
				case result.Pass():
					testsPass++
				case result.Fail:
					testsFail++
				case result.Skip:
					testsSkip++
				}
			}
			fmt.Printf("\n\n=== Test run results ===\n")
			fmt.Println("Tests Passed: ", testsPass)
			fmt.Println("Tests Failed: ", testsFail)
			fmt.Println("Tests Passed: ", testsSkip)
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&trace, "trace", false, "Enable more verbose trace output for Rego queries")
	f.StringVarP(&output, "output", "o", "", "The output format of the results (only json supported right now)")

	return cmd
}

func printJSON(results []*tester.Result) {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		color.Red(fmt.Sprintf("error JSON encoding output: %v", err))
		return
	}
	fmt.Println(string(b))
}
