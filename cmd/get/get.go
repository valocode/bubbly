package get

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/cmd/util"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*GetOptions)(nil)
	getLong                 = util.LongDesc(`
		Display one or many bubbly resources

		    $ bubbly get (KIND | ID | all) [flags]

		`)

	getExample = util.Examples(`
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

// GetOptions holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type GetOptions struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string

	// flags
	events bool

	// resolved
	query     string
	resources []builtin.Resource
}

// NewCmdGet creates a new cobra.Command representing "bubbly get"
func NewCmdGet(bCtx *env.BubblyContext) (*cobra.Command, *GetOptions) {
	o := &GetOptions{
		Command: "get",
		bCtx:    bCtx,
	}

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "get (KIND | ID | all) [flags]",
		Short:   "Display one or many bubbly resources",
		Long:    getLong + "\n\n",
		Example: getExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args

			validationError := o.Validate(cmd)

			if validationError != nil {
				return validationError
			}

			resolveError := o.Resolve()

			if resolveError != nil {
				return resolveError
			}

			bCtx.Logger.Debug().
				Str("query", o.query).
				Msg("getting resources matching query")

			runError := o.Run()

			if runError != nil {
				return runError
			}

			o.Print()

			return nil
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&o.events,
		"events",
		"e",
		false,
		"specify whether to display resource events")

	return cmd, o
}

// Validate checks the GetOptions to see if there is sufficient information run the command.
func (o *GetOptions) Validate(cmd *cobra.Command) error {
	// if no arguments are provided, error, as bubbly cannot reasonably conclude
	// the desired collection of resources to query for
	if len(o.Args) != 1 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", o.Args)
	}

	return nil
}

// Resolve resolves various GetOptions attributes from the provided arguments to cmd
func (o *GetOptions) Resolve() error {

	o.query = formGetQuery(o.Args[0])

	return nil
}

// Run runs the get command over the validated GetOptions configuration
func (o *GetOptions) Run() error {
	resources, err := bubbly.QueryResources(o.bCtx, o.query)

	if err != nil {
		switch err {
		case bubbly.ErrNoResourcesFound:
			return fmt.Errorf("no resources found")
		default:
			return fmt.Errorf("failed to query for resources: %w", err)
		}
	}
	o.resources = resources
	return nil
}

// Print formats and prints the resources and events returned from the query
func (o *GetOptions) Print() {
	var resourceLines []string

	if o.bCtx.CLIConfig.Color {
		color.Blue("Resources")
	} else {
		fmt.Println("Resources")
	}
	for _, r := range o.resources {
		resourceLines = append(resourceLines, resourcePrintLine(r))
	}

	fmt.Println(columnize.SimpleFormat(resourceLines))

	var eventLines []string

	if o.events {
		if o.bCtx.CLIConfig.Color {
			color.Blue("\nEvents")
		} else {
			fmt.Println("\nEvents")
		}
		for _, r := range o.resources {
			for _, e := range r.Event {
				eventLines = append(eventLines, eventPrintLine(r.Id, e))
			}
		}
	}

	fmt.Println(columnize.SimpleFormat(eventLines))

}

// formGetQuery construct the graphql query string to be sent to the bubbly
// store.
// TODO: consider usage of a query builder library
func formGetQuery(input string) string {
	var (
		getArgument struct {
			name  string
			value string
		}

		getQuery string
	)

	if input == "all" {
		getQuery = fmt.Sprintf(`
			{
				%s {
					id
					%s {
						status
						time
						error
					}
				}
			}
		`, core.ResourceTableName, core.EventTableName)
	} else {
		if strings.ContainsAny(input, "/") {
			getArgument.name = "id"
		} else {
			getArgument.name = "kind"
		}
		getArgument.value = input

		getQuery = fmt.Sprintf(`
			{
				%s(%s: "%s", last: 10) {
					id
					%s {
						status
						time
						error
					}
				}
			}
		`, core.ResourceTableName, getArgument.name, getArgument.value,
			core.EventTableName)
	}

	return getQuery
}

// resourcePrintLine returns a formatted string representing the printout of
// a resource
func resourcePrintLine(r builtin.Resource) string {
	latestEvent := r.Event[len(r.Event)-1]
	return fmt.Sprintf("%s | %s | %s\n", r.Id, latestEvent.Status, latestEvent.Time)
}

// eventPrintLine returns a formatted string representing the printout of
// an event
func eventPrintLine(resID string, e builtin.Event) string {
	return fmt.Sprintf("%s | %s | %s | %s\n", resID, e.Status, e.Time, e.Error)
}
