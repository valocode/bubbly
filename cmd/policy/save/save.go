package save

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/policy"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Save a Bubbly policy to the bubbly server
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Save a Bubbly policy to the bubbly server

		bubbly policy save
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
		Use:     "save <policy-file> [flags]",
		Short:   "Save a Bubbly policy",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			b, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			if name == "" {
				fname := filepath.Base(filename)
				name = strings.TrimSuffix(fname, filepath.Ext(fname))
			}
			module := string(b)
			if err := policy.Validate(module, policy.WithResolver(&policy.EntResolver{})); err != nil {
				return fmt.Errorf("validating policy module: %w", err)
			}
			if err := client.SavePolicy(bCtx, &api.ReleasePolicySaveRequest{
				Policy: &api.ReleasePolicyCreate{
					ReleasePolicyModelCreate: *ent.NewReleasePolicyModelCreate().
						SetName(name).
						SetModule(module),
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
	f.StringVar(
		&name,
		"name",
		"",
		`Provide the name of the policy (default is filename without extension)`,
	)
	f.StringSliceVar(&setProjects, "set-projects", nil, "List of project (names) to associate the policy with")
	f.StringSliceVar(&notSetProjects, "not-set-projects", nil, "List of project (names) to disassociate the policy with")
	f.StringSliceVar(&setRepos, "set-repos", nil, "List of repo (names) to associate the policy with")
	f.StringSliceVar(&notSetRepos, "not-set-repos", nil, "List of repo (names) to disassociate the policy with")
	return cmd
}
