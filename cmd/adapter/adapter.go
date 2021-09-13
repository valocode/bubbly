package adapter

import (
	"github.com/valocode/bubbly/cmd/adapter/run"
	"github.com/valocode/bubbly/cmd/adapter/save"
	"github.com/valocode/bubbly/cmd/adapter/view"
	"github.com/valocode/bubbly/env"

	"github.com/spf13/cobra"
)

func New(bCtx *env.BubblyConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adapter <command>",
		Short: "Manage bubbly adapters",
		Long:  `Manage bubbly adapters`,
	}

	cmd.AddCommand(save.New(bCtx))
	cmd.AddCommand(run.New(bCtx))
	cmd.AddCommand(view.New(bCtx))
	return cmd
}
