package create

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ryanuber/columnize"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/release"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Create a bubbly release
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Create a bubbly release

		bubbly release create
		`,
	)
)

func New(bCtx *env.BubblyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [flags]",
		Short:   "Create a bubbly release",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			req, err := release.CreateRelease(bCtx)
			if err != nil {
				return err
			}

			if err := client.CreateRelease(bCtx, req); err != nil {
				return err
			}
			var tag = ""
			if req.Commit.Tag != nil {
				tag = *req.Commit.Tag
			}
			color.Green("Release Created!")
			// Print the release info nicely formatted
			releaseInfo := []string{
				"Repo: | " + *req.Repo.Name,
				"Commit: | " + *req.Commit.Hash,
				"Tag: | " + tag,
				"Branch: | " + *req.Commit.Branch,
				"Name: | " + *req.Release.Name,
				"Version: | " + *req.Release.Version,
			}
			fmt.Println(columnize.SimpleFormat(releaseInfo))
			return nil
		},
	}
	return cmd
}
