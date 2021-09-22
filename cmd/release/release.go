package release

import (
	"github.com/valocode/bubbly/cmd/release/create"
	"github.com/valocode/bubbly/cmd/release/view"
	"github.com/valocode/bubbly/config"

	"github.com/spf13/cobra"
)

func New(bCtx *config.BubblyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <command>",
		Short: "Manage bubbly releases",
		Long:  `Manage bubbly releases`,
	}

	cmd.AddCommand(create.New(bCtx))
	cmd.AddCommand(view.New(bCtx))
	return cmd
}
