package get

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/client"
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
	arg     string

	// flags
	showEvents bool

	// resolved
	resources []builtin.Resource
	events    []builtin.Event
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
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args
			o.arg = args[0]

			if err := o.Validate(cmd); err != nil {
				return err
			}

			if err := o.Resolve(); err != nil {
				return err
			}

			if runError := o.Run(); runError != nil {
				return runError
			}
			o.Print()
			return nil
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&o.showEvents,
		"events",
		"e",
		false,
		"specify whether to display resource events")

	return cmd, o
}

// Validate checks the GetOptions to see if there is sufficient information run the command.
func (o *GetOptions) Validate(cmd *cobra.Command) error {
	return nil
}

// Resolve resolves various GetOptions attributes from the provided arguments to cmd
func (o *GetOptions) Resolve() error {
	return nil
}

// Run runs the get command over the validated GetOptions configuration
func (o *GetOptions) Run() error {

	var (
		resourceFilter string
		resourceQuery  string
		eventQuery     string
		resWrap        builtin.Resource_Wrap
		eventWrap      builtin.Event_Wrap
	)

	client, err := client.New(o.bCtx)
	if err != nil {
		return fmt.Errorf("error creating bubbly client: %w", err)
	}

	switch {
	case o.arg == "all":
		resourceFilter = "last: 20"
	case strings.ContainsAny(o.arg, "/"):
		resourceFilter = "(id: \"" + o.arg + "\")"
	default:
		resourceFilter = "(kind: \"" + o.arg + "\")"
	}
	resourceQuery = fmt.Sprintf(`
	{
		_resource(%s) {
			id
			_event(last: 1) {
				status
				time
			}
		}
	}
	`, resourceFilter)

	o.bCtx.Logger.Debug().
		Str("query", resourceQuery).
		Msg("getting resources matching query")

	if err := client.QueryType(o.bCtx, nil, resourceQuery, &resWrap); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	o.resources = resWrap.Resource

	// Handle events

	eventQuery = `
	{
		_event(last: 20, order_by: {time: desc}) {
			status
			time
			error
			_resource {
				id
			}
		}
	}
	`

	if o.showEvents {
		o.bCtx.Logger.Debug().
			Str("query", eventQuery).
			Msg("getting events matching query")
		if err := client.QueryType(o.bCtx, nil, eventQuery, &eventWrap); err != nil {
			return fmt.Errorf("error executing query: %w", err)
		}
	}
	o.events = eventWrap.Event

	// // If we should show events, make sure we get the events
	// if o.events {
	// 	eventFilter = "last: 20"
	// }

	// if o.arg == "all" {
	// 	resourceFilter = ""
	// } else {
	// 	if strings.ContainsAny(o.arg, "/") {
	// 		resourceFilter = "(id: \"" + o.arg + "\")"
	// 	} else {
	// 		resourceFilter = "(kind: \"" + o.arg + "\")"
	// 	}
	// }

	// query = fmt.Sprintf(`
	// {
	// 	_event(%s) {
	// 		status
	// 		time
	// 		error
	// 		_resource%s {
	// 			id
	// 		}
	// 	}
	// }`, eventFilter, resourceFilter)

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

	if o.showEvents {
		if o.bCtx.CLIConfig.Color {
			color.Blue("\nEvents")
		} else {
			fmt.Println("\nEvents")
		}
		for _, e := range o.events {
			eventLines = append(eventLines, eventPrintLine(e.Resource.Id, e))
		}
	}

	fmt.Println(columnize.SimpleFormat(eventLines))

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
