package util

import (
	"github.com/fatih/color"
	"github.com/valocode/bubbly/bubbly/builtin"
)

func ReleaseStatusColor(status builtin.ReleaseStatus) string {
	// TODO: handle no color mode
	var (
		green  = color.New(color.FgYellow).SprintFunc()
		yellow = color.New(color.FgYellow).SprintFunc()
		red    = color.New(color.FgRed).SprintFunc()
	)

	switch status {
	case builtin.BlockedReleaseStatus:
		return red(status)
	case builtin.PendingReleaseStatus:
		return yellow(status)
	case builtin.ReadyReleaseStatus:
		return green(status)
	default:
		return string(status)
	}
}
