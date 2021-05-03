package list

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"

	"github.com/valocode/bubbly/client"
	cmdutil "github.com/valocode/bubbly/cmd/util"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/util/normalise"
)

var (
	_       cmdutil.Options = (*options)(nil)
	cmdLong                 = normalise.LongDesc(`
		List bubbly releases

		    $ bubbly release list

		`)

	cmdExample = normalise.Examples(`
		# List bubbly releases
		bubbly release list
		`)
)

type releases struct {
	Releases []struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		ReleaseItems []struct {
			Type   string `json:"type"`
			Commit struct {
				Repo struct {
					Name string
				} `json:"repo"`
			} `json:"commit"`
		} `json:"release_item"`
		ReleaseStages []struct {
			Name             string `json:"name"`
			ReleaseCriterias []struct {
				EntryName string `json:"entry_name"`
			} `json:"release_criteria"`
		}
		ReleaseEntries []struct {
			Name   string `json:"name"`
			Result bool   `json:"result"`
		}
	} `json:"release"`
}

// options holds everything necessary to run the command.
// Flag values received to the command are loaded into this struct
type options struct {
	cmdutil.Options
	BubblyContext *env.BubblyContext
	Command       string
	Args          []string

	// flags
	releases *releases
}

// New creates a new cobra command
func New(bCtx *env.BubblyContext) *cobra.Command {
	o := &options{
		Command:       "list",
		BubblyContext: bCtx,
	}

	// cmd represents the apply command
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list bubbly releases",
		Long:    cmdLong + "\n\n",
		Example: cmdExample,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args

			if err := o.validate(cmd); err != nil {
				return err
			}
			if err := o.resolve(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			o.Print()
			return nil
		},
	}

	return cmd
}

// validate checks the cmd options
func (o *options) validate(cmd *cobra.Command) error {
	// Nothing to do
	return nil
}

// resolve resolves args for the command
func (o *options) resolve() error {
	return nil
}

// run runs the command over the validated options
func (o *options) run() error {
	releaseQuery := `
{
	release {
		name
		version
		release_item(filter_on: true) {
			type
			commit {
				repo {
					name
				}
			}
		}
		release_stage(filter_on: true) {
			name
			release_criteria(filter_on: true) {
				entry_name
			}
		}
	}
}
	`

	releases := &releases{}

	client, err := client.New(o.BubblyContext)
	if err != nil {
		return fmt.Errorf("error creating bubbly client: %w", err)
	}
	bytes, err := client.Query(o.BubblyContext, nil, releaseQuery)
	if err != nil {
		return fmt.Errorf("error making GraphQL query: %w", err)
	}

	var results graphql.Result
	results.Data = releases
	if err := json.Unmarshal(bytes, &results); err != nil {
		return fmt.Errorf("error unmarshalling GraphQL results: %w", err)
	}
	if results.HasErrors() {
		var msgs []string
		for _, err := range results.Errors {
			msgs = append(msgs, err.Message)
		}
		return fmt.Errorf("GraphQL query returned errors: %s", strings.Join(msgs, "\n\n"))
	}

	o.releases = releases
	return nil
}

// Print prints the successful outcome of the cmd
func (o *options) Print() {

	var releaseLines []string
	releaseLines = append(releaseLines, "Name | Version | Type | Status")
	for _, rel := range o.releases.Releases {
		var (
			relType   string
			relStatus bool

			criterion []string
		)
		for _, item := range rel.ReleaseItems {
			relType = item.Type
		}
		for _, stage := range rel.ReleaseStages {
			for _, criteria := range stage.ReleaseCriterias {
				criterion = append(criterion, criteria.EntryName)
			}
		}
		relStatus = true
		for _, entry := range rel.ReleaseEntries {
			var entryOK bool
			for _, criteria := range criterion {
				if criteria == entry.Name {
					entryOK = true
				}
			}
			if !entryOK {
				relStatus = false
				break
			}
		}
		relStatusStr := "READY"
		if relStatus {
			relStatusStr = "BLOCKED"
		}
		releaseLines = append(releaseLines, fmt.Sprintf(
			"%s | %s | %s | %s ", rel.Name, rel.Version, relType, relStatusStr,
		))
	}
	fmt.Println("Releases")
	fmt.Println(columnize.SimpleFormat(releaseLines))
}
