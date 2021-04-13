package datastore

import (
	"encoding/json"
	"fmt"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

func (d *DataStore) getResourcesByKindHandler(bCtx *env.BubblyContext, subject string, reply string, data component.MessageData) (interface{}, error) {
	bCtx.Logger.Debug().
		Str("subject", subject).
		Str("component", string(d.Type)).
		Msg("processing message")

	return d.Store.Query(string(data.Data)), nil
}

func (d *DataStore) postSchemaHandler(bCtx *env.BubblyContext, subject string, reply string, data component.MessageData) (interface{}, error) {
	bCtx.Logger.Debug().
		Str("subject", subject).
		Str("component", string(d.Type)).
		Msg("processing message")

	var schema core.Tables
	if err := json.Unmarshal(data.Data, &schema); err != nil {
		return nil, fmt.Errorf("failed to decode schema into core.Tables: %w", err)
	}

	if err := d.Store.Apply(schema); err != nil {
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}
	return nil, nil
}

func (d *DataStore) queryHandler(bCtx *env.BubblyContext, subject string, reply string, data component.MessageData) (interface{}, error) {
	bCtx.Logger.Debug().
		Str("subject", subject).
		Str("component", string(d.Type)).
		Msg("processing message")

	result := d.Store.Query(string(data.Data))
	if result.HasErrors() {
		return nil, fmt.Errorf("error while querying the data store: %v", result.Errors)
	}
	return result, nil
}

func (d *DataStore) uploadHandler(bCtx *env.BubblyContext, subject string, reply string, data component.MessageData) (interface{}, error) {
	bCtx.Logger.Debug().
		Str("subject", subject).
		Str("component", string(d.Type)).
		Msg("processing message")

	var dbs core.DataBlocks
	if err := json.Unmarshal(data.Data, &dbs); err != nil {
		return nil, fmt.Errorf("failed to decode data into core.DataBlocks: %w", err)
	}
	if err := d.Store.Save(dbs); err != nil {
		return nil, fmt.Errorf("failed to save data to data store: %w", err)
	}

	return nil, nil
}
