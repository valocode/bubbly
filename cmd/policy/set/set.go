package set

import (
	"fmt"

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
			resp, err := client.GetPolicies(bCtx, &api.ReleasePolicyGetRequest{
				Name: name,
			})
			if err != nil {
				return fmt.Errorf("getting policies: %w", err)
			}
			if len(resp.Policies) == 0 {
				fmt.Println("No policies with name: " + name)
				return nil
			}
			policy := resp.Policies[0]
			if err := client.SetPolicy(bCtx, &api.ReleasePolicyUpdateRequest{
				ID: policy.ID,
				Policy: &api.ReleasePolicyUpdate{
					Affects: &api.ReleasePolicyAffectsSet{
						Projects:    setProjects,
						NotProjects: notSetProjects,
						Repos:       setRepos,
						NotRepos:    notSetRepos,
					},
				},
			}); err != nil {
				return err
			}
			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.StringSliceVar(&setProjects, "projects", nil, "List of project (names) to associate the policy with")
	f.StringSliceVar(&notSetProjects, "not-projects", nil, "List of project (names) to disassociate the policy with")
	f.StringSliceVar(&setRepos, "repos", nil, "List of repo (names) to associate the policy with")
	f.StringSliceVar(&notSetRepos, "not-repos", nil, "List of repo (names) to disassociate the policy with")
	return cmd
}
