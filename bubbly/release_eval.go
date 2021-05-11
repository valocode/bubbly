package bubbly

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
	"github.com/valocode/bubbly/store"
	"github.com/zclconf/go-cty/cty"
)

func EvalReleaseCriteria(bCtx *env.BubblyContext, criteriaName string) (*ReleaseSpec, error) {
	var fileParser BubblyFileParser
	err := parser.ParseConfig(bCtx, &fileParser)
	if err != nil {
		return nil, fmt.Errorf("error parsing bubbly configs: %w", err)
	}

	// TODO: hacked/hardcoded -- get data from fileParser.Relase
	releaseData := core.Data{
		TableName: "release",
		Fields: core.DataFields{
			"name":    cty.StringVal("github.com/valocode/bubbly"),
			"version": cty.StringVal("e97d7ddb0599c8e3839f7bf377568657ccc562c0"),
		},
		Policy: core.ReferencePolicy,
	}

	criteria, err := criteriaByName(fileParser.Release, criteriaName)
	if err != nil {
		return nil, err
	}
	dEntry, err := criteria.EntryLog(bCtx, releaseData)
	if err != nil {
		return nil, err
	}

	var data core.DataBlocks
	data = append(data, releaseData)
	data = append(data, dEntry...)

	// TODO: DELETE DUMMY DATA!
	// data = append(data, core.Data{
	// 	TableName: "release_criteria",
	// 	Fields: core.DataFields{
	// 		"entry_name": cty.StringVal("gofmt"),
	// 	},
	// 	Joins:  []string{"release"},
	// 	Policy: core.ReferenceIfExistsPolicy,
	// })

	// data = append(data, core.Data{
	// 	TableName: "release_entry",
	// 	Fields: core.DataFields{
	// 		"name": cty.StringVal("gosec"),
	// 	},
	// 	Joins:  []string{"release_criteria"},
	// 	Policy: core.CreatePolicy,
	// })

	// data = append(data, core.Data{
	// 	TableName: "release_criteria",
	// 	Fields: core.DataFields{
	// 		"entry_name": cty.StringVal("performance_test"),
	// 	},
	// 	Joins:  []string{"release"},
	// 	Policy: core.ReferenceIfExistsPolicy,
	// })
	// data = append(data, core.Data{
	// 	TableName: "release_entry",
	// 	Fields: core.DataFields{
	// 		"name": cty.StringVal("performance_test"),
	// 	},
	// 	Joins:  []string{"release_criteria"},
	// 	Policy: core.CreatePolicy,
	// })

	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("error creating bubbly client: %w", err)
	}

	dBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling release data block: %w", err)
	}
	// TODO: auth
	if err := client.Load(bCtx, nil, dBytes); err != nil {

		// TODO: this doesn't work when errors are sent over HTTP/NATS...
		if errors.Is(err, store.ErrDataCreateExists) {
			return nil, fmt.Errorf("release already exists")
		}
		return nil, fmt.Errorf("error saving release data block: %w", err)
	}
	return fileParser.Release, nil
}

func criteriaByName(release *ReleaseSpec, criteriaName string) (*releaseCriteria, error) {
	for _, stages := range release.Stages {
		for _, criteria := range stages.Criterion {
			if criteria.Name == criteriaName {
				return &criteria, nil
			}
		}
	}
	return nil, fmt.Errorf("no criteria found in release with name: %s", criteriaName)
}
