package list

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/cmd/util"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = util.LongDesc(`
		Create a release

		    $ bubbly release create

		`)

	cmdExample = util.Examples(`
		# Create a release
		bubbly release create
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string
	Release *bubbly.ReleaseSpec
	success bool

	// flags
	filename string
	noFail   bool
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command: "create",
		bCtx:    bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create a release",
		Long:    cmdLong + "\n\n",
		Example: cmdExample,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args

			if err := o.validate(cmd); err != nil {
				return err
			}
			if err := o.resolve(); err != nil {
				return err
			}

			// Let's be optimistic and set to false in case of failure
			o.success = true
			if err := o.run(); err != nil {
				o.success = false
				if !o.noFail {
					return err
				}
			}

			o.Print()
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&o.filename,
		"filename",
		"f",
		".",
		"filename or directory that contains the bubbly release definition")
	f.BoolVar(&o.noFail,
		"no-fail",
		false,
		"do not fail if create fails (e.g. if the release already exists)",
	)

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
	release, err := bubbly.CreateRelease(o.bCtx, o.filename)
	if err != nil {
		return err
	}
	o.Release = release
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {
	if o.success {
		if o.bCtx.CLIConfig.Color {
			color.Green("Release successfully created!")
		} else {
			fmt.Println("Release successfully created!")
		}
		fmt.Println("\n" + o.Release.String())
		return
	}

	if o.bCtx.CLIConfig.Color {
		color.Yellow("Release creation failed.")
	} else {
		fmt.Println("Release creation failed.")
	}
}
