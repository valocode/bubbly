package bubbly

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
)

func EvalReleaseCriteria(bCtx *env.BubblyContext, filename string, criteriaName string) (*ReleaseSpec, error) {

	release, err := createReleaseSpec(bCtx, filename)
	if err != nil {
		return nil, fmt.Errorf("error creating release spec: %w", err)
	}

	// Get the reference to the release data block
	releaseRef, err := release.DataRef()
	if err != nil {
		return nil, fmt.Errorf("unable to process release definition: %w", err)
	}

	criteria, err := criteriaByName(release, criteriaName)
	if err != nil {
		return nil, err
	}
	// Evaluate the release criteria and create the release entry data blocks
	dEntry, err := criteria.Evaluate(bCtx, releaseRef, release.BaseDir)
	if err != nil {
		return nil, err
	}

	var data core.DataBlocks
	data = append(data, releaseRef...)
	data = append(data, dEntry...)

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
		return nil, fmt.Errorf("error saving release data block: %w", err)
	}
	return release, nil
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
