package view

import (
	"errors"
	"fmt"

	"github.com/ryanuber/columnize"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
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

func New(bCtx *env.BubblyContext) *cobra.Command {

	var (
		commit   string
		repo     string
		policies bool
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
			rel, err := client.GetRelease(bCtx, &api.ReleaseGetRequest{
				Commit:   &commit,
				Repo:     &repo,
				Policies: policies,
			})
			if err != nil {
				return err
			}
			switch len(rel.Releases) {
			case 0:
				fmt.Printf("No release for current commit %s. Please create one.\n", commit)
			case 1:
				printRelease(rel.Releases[0], policies)
			default:
				return errors.New("received more than one release")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&commit, "commit", "", "Specify the git commit for which to view a release")
	cmd.Flags().StringVar(&repo, "repo", "", "Specify the repo name for which to view a release. Requires also the --commit flag")
	cmd.Flags().BoolVar(&policies, "policies", false, "Whether to also fetch the policies that apply to this release")

	return cmd
}

func printRelease(rel *api.ReleaseRead, policies bool) {
	fmt.Println("Project: " + *rel.Project.Name)
	fmt.Println("Repo: " + *rel.Project.Name)
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
			"%s | %s | %s", *v.Message, *v.Severity, "TODO",
		))
	}
	if len(rel.Violations) == 0 {
		fmt.Println("No violations")
	} else {
		fmt.Println(columnize.SimpleFormat(violationLines))
	}

	fmt.Println("")
}
