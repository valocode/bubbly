package util

import "github.com/fatih/color"

func ReleaseStatusColor(status string) string {
	// TODO: handle no color mode
	var (
		green  = color.New(color.FgYellow).SprintFunc()
		yellow = color.New(color.FgYellow).SprintFunc()
		red    = color.New(color.FgRed).SprintFunc()
	)

	switch status {
	case "BLOCKED":
		return red(status)
	case "PENDING":
		return yellow(status)
	case "READY":
		return green(status)
	default:
		return status
	}
}
