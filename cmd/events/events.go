package events

import (
	"fmt"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/ent"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/store/api"

	"github.com/spf13/cobra"

	cmdutil "github.com/valocode/bubbly/cmd/util"
)

var (
	cmdLong = cmdutil.LongDesc(
		`
		Get events from the Bubbly server

			$ bubbly events
		`,
	)

	cmdExamples = cmdutil.Examples(
		`
		# Get events from the Bubbly server

		bubbly events
		`,
	)
)

func New(bCtx *env.BubblyConfig) *cobra.Command {
	var req api.EventGetRequest

	cmd := &cobra.Command{
		Use:     "events [flags]",
		Short:   "Get events from the Bubbly server",
		Long:    cmdLong + "\n\n",
		Example: cmdExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := client.GetEvents(bCtx, &req)
			if err != nil {
				return fmt.Errorf("getting events: %w", err)
			}

			printEvents(resp.Events)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&req.Project, "project", "", "The project to filter events by")
	f.StringVar(&req.Repo, "repo", "", "The repo to filter events by")
	f.StringVar(&req.Commit, "commit", "", "The commit to filter events by")
	f.StringVar(&req.ReleaseName, "release-name", "", "The release name to filter events by")
	f.StringVar(&req.ReleaseVersion, "release-version", "", "The release version to filter events by")
	f.StringVar(&req.Last, "last", "20", "The last number of events to show (default: 20)")

	return cmd
}

func printEvents(dbEvents []*ent.Event) {
	var eventLines []string
	eventLines = append(eventLines, "Message | Type | Status | Time")
	for _, e := range dbEvents {
		// The message might contain new lines, therefore split by those
		// and only add other columns on the first index so that formatting is
		// not messed up
		for idx, msg := range strings.Split(e.Message, "\n") {
			if idx == 0 {
				eventLines = append(eventLines, fmt.Sprintf(
					"%s | %s | %s | %s", msg, e.Type, e.Status, e.Time,
				))
			} else {
				eventLines = append(eventLines, fmt.Sprintf(
					"%s | %s | %s | %s", msg, "", "", "",
				))
			}
		}
	}
	if len(dbEvents) == 0 {
		fmt.Println("No events")
	} else {
		fmt.Println(columnize.SimpleFormat(eventLines))
	}
}
