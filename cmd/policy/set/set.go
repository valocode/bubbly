package set

import (
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Set a Bubbly policy to a project or repository
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Set a Bubbly policy to a project or repository

		bubbly policy set
		`,
	)
)

func New(bCtx *env.BubblyContext) *cobra.Command {

	var (
		name           string
		setProjects    []string
		notSetProjects []string
		setRepos       []string
		notSetRepos    []string
	)
	cmd := &cobra.Command{
		Use:     "set <policy> [flags]",
		Short:   "Set a Bubbly policy",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := client.SetPolicy(bCtx, &api.ReleasePolicySetRequest{
				Policy: &name,
				Affects: &api.ReleasePolicyAffects{
					Projects:    setProjects,
					NotProjects: notSetProjects,
					Repos:       setRepos,
					NotRepos:    notSetRepos,
				},
			}); err != nil {
				return err
			}
			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.StringVar(
		&name,
		"name",
		"",
		`Provide the name of the policy (default is filename without extension)`,
	)
	f.StringSliceVar(&setProjects, "projects", nil, "List of project (names) to associate the policy with")
	f.StringSliceVar(&notSetProjects, "not-projects", nil, "List of project (names) to disassociate the policy with")
	f.StringSliceVar(&setRepos, "repos", nil, "List of repo (names) to associate the policy with")
	f.StringSliceVar(&notSetRepos, "not-repos", nil, "List of repo (names) to disassociate the policy with")
	return cmd
}
