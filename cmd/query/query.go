package query

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/graphql-go/graphql"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/cmd/util"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = util.LongDesc(`
		Perform a GraphQL query

		    $ bubbly query QUERY_STRING

		`)

	cmdExample = util.Examples(`
		# Perform a GraphQL query
		bubbly query QUERY_STRING
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string

	query  string
	result string
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command: "query",
		bCtx:    bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "query",
		Short:   "perform a graphql query",
		Long:    cmdLong + "\n\n",
		Example: cmdExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.query = args[0]

			if err := o.validate(cmd); err != nil {
				return err
			}
			if err := o.resolve(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			o.Print()
			return nil
		},
	}

	return cmd
}

// validate checks the cmd options
func (o *options) validate(cmd *cobra.Command) error {
	// Nothing to do
	return nil
}

// resolve resolves args for the command
func (o *options) resolve() error {
	return nil
}

// run runs the command over the validated options
func (o *options) run() error {
	client, err := client.New(o.bCtx)
	if err != nil {
		return fmt.Errorf("error creating bubbly client: %w", err)
	}
	// TODO: add authentication
	bytes, err := client.Query(o.bCtx, nil, o.query)
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

	o.result = string(pretty)
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {
	fmt.Printf("\nResult:\n%s\n\n", o.result)
	color.Green("Query successfully handled!")
}
