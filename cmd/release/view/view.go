package view

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/bubbly/builtin"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = cmdutil.LongDesc(`
		View a release

		    $ bubbly release view

		`)

	cmdExample = cmdutil.Examples(`
		# View a release
		bubbly release view
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string
	Release *builtin.Release

	// flags
	filename string
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command: "view",
		bCtx:    bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "view a release",
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
			if err := o.run(); err != nil {
				return err
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
	release, err := bubbly.GetRelease(o.bCtx, o.filename)
	if err != nil {
		return err
	}
	o.Release = release
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {
	status := cmdutil.ReleaseStatusColor(builtin.ReleaseStatusByStages(*o.Release))

	fmt.Println("Project: " + o.Release.Project.Id)
	fmt.Println("Name: " + o.Release.Name)
	fmt.Println("Version: " + o.Release.Version)
	fmt.Println("Status: " + status)
	fmt.Println("")
	fmt.Println("Inputs:")
	for _, input := range o.Release.ReleaseInput {
		fmt.Println("  - " + input.Type)
	}
	fmt.Println("")
	fmt.Println("Stages:")
	for _, stage := range o.Release.ReleaseStage {
		stageStatus := cmdutil.ReleaseStatusColor(builtin.ReleaseStageStatus(stage))

		fmt.Println("")
		fmt.Println(stage.Name + ": " + stageStatus)
		for _, criteria := range stage.ReleaseCriteria {
			criteriaStatus := cmdutil.ReleaseStatusColor(builtin.ReleaseCriteriaStatus(criteria))
			fmt.Println("    - Criteria: " + criteria.EntryName + " (" + criteriaStatus + ")")
		}
	}
}
