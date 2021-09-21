package policy

import (
	"github.com/valocode/bubbly/cmd/policy/list"
	"github.com/valocode/bubbly/cmd/policy/save"
	"github.com/valocode/bubbly/cmd/policy/set"
	"github.com/valocode/bubbly/cmd/policy/view"
	"github.com/valocode/bubbly/env"

	"github.com/spf13/cobra"
)

func New(bCtx *env.BubblyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy <command>",
		Short: "Manage bubbly policies",
		Long:  `Manage bubbly policies`,
	}

	cmd.AddCommand(list.New(bCtx))
	cmd.AddCommand(save.New(bCtx))
	cmd.AddCommand(set.New(bCtx))
	cmd.AddCommand(view.New(bCtx))
	return cmd
}
