package util

import (
	"github.com/spf13/cobra"
)

type Options interface {
	Validate(cmd *cobra.Command) error
	Resolve(cmd *cobra.Command) error
	Run() error
	Print(cmd *cobra.Command)
}
