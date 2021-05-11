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
)

type BubblyFileParser struct {
	Release        *ReleaseSpec        `hcl:"release,block"`
	ResourceBlocks core.ResourceBlocks `hcl:"resource,block"`
}

func CreateRelease(bCtx *env.BubblyContext) (*ReleaseSpec, error) {
	var fileParser BubblyFileParser
	err := parser.ParseConfig(bCtx, &fileParser)
	if err != nil {
		return nil, fmt.Errorf("error parsing bubbly configs: %w", err)
	}

	d, err := fileParser.Release.Data()
	if err != nil {
		return nil, fmt.Errorf("error creating data blocks for release: %w", err)
	}

	client, err := client.New(bCtx)
	if err != nil {
		return nil, fmt.Errorf("error creating bubbly client: %w", err)
	}

	dBytes, err := json.Marshal(d)
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
