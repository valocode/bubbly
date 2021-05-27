package list

import (
	"fmt"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/bubbly"
	"github.com/valocode/bubbly/bubbly/builtin"
	"github.com/valocode/bubbly/cmd/util"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = util.LongDesc(`
		List bubbly releases

		    $ bubbly release list

		`)

	cmdExample = util.Examples(`
		# List bubbly releases
		bubbly release list
		`)
)

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	bCtx    *env.BubblyContext
	Command string
	Args    []string

	// flags
	releases *builtin.Release_Wrap
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command: "list",
		bCtx:    bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list bubbly releases",
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

	releases, err := bubbly.ListReleases(o.bCtx)
	if err != nil {
		return err
	}
	o.releases = releases
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {

	var releaseLines []string
	releaseLines = append(releaseLines, "Name | Version | Type | Status")
	for _, rel := range o.releases.Release {
		var (
			relType []string
		)
		for _, item := range rel.ReleaseItem {
			relType = append(relType, item.Type)
		}
		relStatusStr := cmdutil.ReleaseStatusColor(builtin.ReleaseStatusByStages(rel))
		releaseLines = append(releaseLines, fmt.Sprintf(
			"%s | %s | %s | %s ", rel.Name, rel.Version, strings.Join(relType, ","), relStatusStr,
		))
	}
	fmt.Println("Releases")
	fmt.Println(columnize.SimpleFormat(releaseLines))
}
