package get

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/client"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/util/normalise"
)

var (
	_       cmdutil.Options = (*QueryOptions)(nil)
	getLong                 = normalise.LongDesc(`
		Display one or many bubbly resources

		    $ bubbly get (KIND | ID | all) [flags]

		`)

	getExample = normalise.Examples(`
		# Display all bubbly resources stored on the bubbly server
		bubbly get all

		# Display all bubbly resources of kind extract
		bubbly get extract

		# Display a specific bubbly resource
		bubbly get default/extract/sonarqube

		# Display a specific bubbly resource and associated events
		bubbly get default/extract/sonarqube --events

		`)
)

// QueryOptions holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type QueryOptions struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string
}

// NewCmdQuery creates a new cobra.Command representing "bubbly get"
func NewCmdQuery(bCtx *env.BubblyContext) *cobra.Command {
	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "query string",
		Short:   "Run a GraphQL query against the bubbly server",
		Long:    getLong + "\n\n",
		Example: getExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// We set a condition that exactly one arg must be received
			var query = args[0]

			client, err := client.New(bCtx)
			if err != nil {
				return fmt.Errorf("error creating bubbly client: %w", err)
			}
			// TODO: add authentication
			bytes, err := client.Query(bCtx, nil, query)
			if err != nil {
				return fmt.Errorf("error making GraphQL query: %w", err)
			}

			var result graphql.Result
			if err := json.Unmarshal(bytes, &result); err != nil {
				return fmt.Errorf("error decoding GraphQL respose into result: %w", err)
			}
			if result.HasErrors() {
				var errStr []string
				for _, err := range result.Errors {
					errStr = append(errStr, err.Message)
				}
				return fmt.Errorf("query returned %d errors: %s", len(result.Errors), strings.Join(errStr, "\n"))
			}

			pretty, err := json.MarshalIndent(result.Data, "", "  ")
			if err != nil {
				return fmt.Errorf("error pretty printing GraphQL result: %w", err)
			}

			fmt.Println(string(pretty))

			return nil
		},
	}

	return cmd
}
