package datastore

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

func (d *DataStore) getResourcesByKindHandler(bCtx *env.BubblyContext, m *nats.Msg) (interface{}, error) {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	return d.Store.Query(string(m.Data)), nil
}

func (d *DataStore) postSchemaHandler(bCtx *env.BubblyContext, m *nats.Msg) (interface{}, error) {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	var schema core.Tables
	if err := json.Unmarshal(m.Data, &schema); err != nil {
		return nil, fmt.Errorf("failed to decode schema into core.Tables: %w", err)
	}

	if err := d.Store.Apply(bCtx, schema); err != nil {
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}
	return nil, nil
}

func (d *DataStore) queryHandler(bCtx *env.BubblyContext, m *nats.Msg) (interface{}, error) {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	result := d.Store.Query(string(m.Data))
	if result.HasErrors() {
		return nil, fmt.Errorf("error while querying the data store: %v", result.Errors)
	}
	return result, nil
}

func (d *DataStore) uploadHandler(bCtx *env.BubblyContext, m *nats.Msg) (interface{}, error) {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	var dbs core.DataBlocks
	if err := json.Unmarshal(m.Data, &dbs); err != nil {
		return nil, fmt.Errorf("failed to decode data into core.DataBlocks: %w", err)
	}
	if err := d.Store.Save(bCtx, dbs); err != nil {
		return nil, fmt.Errorf("failed to save data to data store: %w", err)
	}

	return nil, nil
}
