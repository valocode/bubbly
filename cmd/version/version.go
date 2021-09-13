package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valocode/bubbly/env"
)

// New creates a new Cobra command
func New(bCtx *env.BubblyConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the Bubbly version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bubbly version", bCtx.Version.Version)
			fmt.Println("Commit:", bCtx.Version.Commit)
			fmt.Println("Date:", bCtx.Version.Date)
		},
	}
}
