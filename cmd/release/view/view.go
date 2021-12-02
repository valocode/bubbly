package view

import (
	"errors"
	"fmt"

	"github.com/ryanuber/columnize"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/config"
	"github.com/valocode/bubbly/release"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		View a bubbly release
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# View a bubbly release

		bubbly release view
		`,
	)
)

func New(bCtx *config.BubblyConfig) *cobra.Command {
	var (
		commit   string
		project  string
		repo     string
		policies bool
		log      bool
	)
	cmd := &cobra.Command{
		Use:     "view [flags]",
		Short:   "View a bubbly release",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if commit == "" {
				// If no release was provided, get one from the current bubbly release
				var err error
				commit, err = release.Commit(bCtx)
				if err != nil {
					return err
				}
			}
			resp, err := client.GetReleases(bCtx, &api.ReleaseGetRequest{
				Commit:   commit,
				Projects: project,
				Repos:    repo,
				Policies: policies,
				Log:      log,
			})
			if err != nil {
				return err
			}
			switch len(resp.Releases) {
			case 0:
				fmt.Printf("No release for current commit %s. Please create one.\n", commit)
			case 1:
				printRelease(resp.Releases[0], policies, log)
			default:
				return errors.New("received more than one release")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&commit, "commit", "", "Specify the git commit for which to view a release")
	cmd.Flags().StringVar(&project, "project", "", "Specify the project name for which to view a release")
	cmd.Flags().StringVar(&repo, "repo", "", "Specify the repo name for which to view a release")
	cmd.Flags().BoolVar(&policies, "policies", false, "Whether to also fetch the policies that apply to this release")
	cmd.Flags().BoolVar(&log, "log", false, "Whether to also fetch the release log")

	return cmd
}

func printRelease(rel *api.Release, policies bool, log bool) {
	fmt.Println("Project: " + *rel.Project.Name)
	fmt.Println("Repo: " + *rel.Repo.Name)
	fmt.Println("Commit: " + *rel.Commit.Hash)
	fmt.Println("Branch: " + *rel.Commit.Branch)
	fmt.Println("Name: " + *rel.Release.Name)
	fmt.Println("Version: " + *rel.Release.Version)
	fmt.Println("")

	if policies {
		fmt.Println("=== Policies ===")
		if len(rel.Policies) == 0 {
			fmt.Println("No policies")
		}
		for _, p := range rel.Policies {
			fmt.Println("  - " + *p.Name)
		}
	}

	fmt.Println("")
	fmt.Println("=== Violations ===")
	var violationLines []string
	violationLines = append(violationLines, "Message | Severity | Type")
	for _, v := range rel.Violations {
		violationLines = append(violationLines, fmt.Sprintf(
			"%s | %s | %s", *v.Message, *v.Severity, *v.Type,
		))
	}
	if len(rel.Violations) == 0 {
		fmt.Println("No violations")
	} else {
		fmt.Println(columnize.SimpleFormat(violationLines))
	}

	if log {
		fmt.Println("")
		fmt.Println("=== Log ===")
		if len(rel.Entries) == 0 {
			fmt.Println("No entries")
		}
		for _, e := range rel.Entries {
			fmt.Println("  - " + e.Type.String() + " at " + e.Time.String())
		}
	}

	fmt.Println("")
}
