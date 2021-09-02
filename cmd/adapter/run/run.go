package run

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/labstack/gommon/color"
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
		dryrun     bool
		inputFiles string
		output     string
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
			if dryrun {
				fmt.Println("")
				fmt.Println("Running in dry run mode so not sending results to bubbly server")
				fmt.Println("")
			}
			if result.CodeScan != nil && !dryrun {
				err := client.SaveCodeScan(bCtx, &api.CodeScanRequest{
					Commit:   &commit,
					CodeScan: result.CodeScan,
				})
				if err != nil {
					return err
				}
			}
			if result.TestRun != nil && !dryrun {
				err := client.SaveTestRun(bCtx, &api.TestRunRequest{
					Commit:  &commit,
					TestRun: result.TestRun,
				})
				if err != nil {
					return err
				}
			}

			switch output {
			case "json":
				printJSON(result)
			// case "stdout":
			// 	printStdout(result)
			case "":
				// Do nothing
			default:
				return fmt.Errorf("unknown output format: %s", output)

			}
			fmt.Printf("\n\nAdapter %s ran successfully!\n", adapterID)
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&trace, "trace", false, "Enable more verbose trace output for Rego queries")
	f.BoolVar(&dryrun, "dry", false, "Enable dry run mode where the results are not sent to the bubbly server")
	f.StringVar(&inputFiles, "input", "", "Provide a comma-separated list of files to parse and provide as input")
	f.StringVarP(&output, "output", "o", "", "The output format of the results (only json supported right now)")

	return cmd
}

func printJSON(result *adapter.AdapterResult) {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		color.Red(fmt.Sprintf("error JSON encoding output: %v", err))
		return
	}
	fmt.Println(string(b))
}

// func printStdout(result *adapter.AdapterResult) {
// 	if scan := result.CodeScan; scan != nil {
// 		fmt.Println("=== Code Scan ===")
// 		fmt.Println("Tool: ", *scan.Tool)
// 		fmt.Println("Metadata: ", scan.Metadata)
// 		fmt.Println("Time: ", scan.Time)

// 		fmt.Println("=== Code Issues ===")
// 		var codeIssueLines []string
// 		codeIssueLines = append(codeIssueLines, "RuleID | Message | Severity | Type")
// 		for _, v := range scan.Issues {
// 			codeIssueLines = append(codeIssueLines, fmt.Sprintf(
// 				"%s | %s | %s | %s", *v.RuleID, *v.Message, *v.Severity, *v.Type,
// 			))
// 		}
// 		if len(scan.Issues) == 0 {
// 			fmt.Println("No code issues")
// 		} else {
// 			fmt.Println(columnize.SimpleFormat(codeIssueLines))
// 		}
// 	}
// }
