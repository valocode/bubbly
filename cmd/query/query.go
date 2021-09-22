package query

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/config"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Run a GraphQL query against the Bubbly server

			$ bubbly query
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Run a GraphQL query against the Bubbly server

		bubbly events
		`,
	)
)

func New(bCtx *config.BubblyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query [flags]",
		Short:   "Run a GraphQL query against the Bubbly server",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var query string
			if len(args) == 1 {
				query = args[0]
			} else {
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("reading query from stdin: %w", err)
				}
				query = string(b)
			}
			if query == "" {
				return fmt.Errorf("cannot execute empty query")
			}
			resp, err := client.RunQuery(bCtx, query)
			if err != nil {
				return fmt.Errorf("getting events: %w", err)
			}

			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("pretty printing response: %w", err)
			}
			fmt.Println(string(b))
			return nil
		},
	}

	return cmd
}
