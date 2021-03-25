package datastore

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/valocode/bubbly/agent/component"
	"github.com/valocode/bubbly/api/core"
	"github.com/valocode/bubbly/env"
)

// PostResourceHandler receives a core.Data representation of the data and attempts
// to load it into the store. Publishes a reply containing a nil error on failure
func (d *DataStore) PostResourceHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	var data core.Data

	err := json.Unmarshal(m.Data, &data)

	if err != nil {
		return fmt.Errorf("failed to unmarshal into core.Data: %w", err)
	}

	if err := d.Store.Save(bCtx, core.DataBlocks{data}); err != nil {
		return fmt.Errorf("failed to save resource into store: %w", err)
	}

	pub := component.Publication{
		Subject: component.Subject(m.Reply),
		Encoder: nats.DEFAULT_ENCODER,
		Error:   nil,
	}

	if err := d.Publish(bCtx, pub); err != nil {
		return fmt.Errorf(
			`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}

func (d *DataStore) GetResourcesByKindHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	result := d.Store.Query(string(m.Data))

	if result.HasErrors() {
		return fmt.Errorf("unable to fetch resources from the data store: %v", result.Errors)
	}

	pub := component.Publication{}
	if result.Data.(map[string]interface{})[core.ResourceTableName] == nil {
		pub = component.Publication{
			Subject: component.Subject(m.Reply),
			Encoder: nats.DEFAULT_ENCODER,
			Data:    []byte{},
		}
	} else {
		var (
			resourceBlocksJSON []core.ResourceBlockJSON
			inputMap           = result.Data.(map[string]interface{})[core.ResourceTableName].([]interface{})
		)

		// loop over the return store's return and create a core.
		// ResourceBlockJSON for each index.
		for i := 0; i < len(inputMap); i++ {
			var resBlockJSON core.ResourceBlockJSON
			b, err := json.Marshal(inputMap[i])
			if err != nil {
				return fmt.Errorf("failed to marshal resources: %w", err)
			}
			err = json.Unmarshal(b, &resBlockJSON)

			resourceBlocksJSON = append(resourceBlocksJSON, resBlockJSON)
		}

		// marshal the list of core.ResourceBlockJSON ready to be returned via NATS
		b, err := json.Marshal(resourceBlocksJSON)

		if err != nil {
			return fmt.Errorf("failed to marshal []core.ResourceBlockJSON: %w", err)
		}

		pub = component.Publication{
			Subject: component.Subject(m.Reply),
			Encoder: nats.DEFAULT_ENCODER,
			Data:    b,
		}
	}

	err := d.Publish(bCtx, pub)
	if err != nil {
		return fmt.Errorf(`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}

func (d *DataStore) PostSchemaHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	var schema core.Tables

	if err := json.Unmarshal(m.Data, &schema); err != nil {
		return fmt.Errorf("failed to decode schema into core.Tables: %w", err)
	}

	if err := d.Store.Apply(bCtx, schema); err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	pub := component.Publication{
		Subject: component.Subject(m.Reply),
		Encoder: nats.DEFAULT_ENCODER,
		Error:   nil,
	}

	if err := d.Publish(bCtx, pub); err != nil {
		return fmt.Errorf(
			`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}

func (d *DataStore) QueryHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	result := d.Store.Query(string(m.Data))
	resultBytes, err := json.Marshal(result)

	if err != nil {
		return fmt.Errorf("failed to marshal Query response: %w", err)
	}

	pub := component.Publication{
		Subject: component.Subject(m.Reply),
		Data:    resultBytes,
		Encoder: nats.DEFAULT_ENCODER,
		Error:   nil,
	}

	if err := d.Publish(bCtx, pub); err != nil {
		return fmt.Errorf(
			`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}

func (d *DataStore) UploadHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("component", string(d.Type)).
		Msg("processing message")

	var dbs core.DataBlocks

	if err := json.Unmarshal(m.Data, &dbs); err != nil {
		return fmt.Errorf("failed to decode data into core.DataBlocks: %w", err)
	}

	if err := d.Store.Save(bCtx, dbs); err != nil {
		return fmt.Errorf("failed to save data to data store: %w", err)
	}

	pub := component.Publication{
		Subject: component.Subject(m.Reply),
		Encoder: nats.DEFAULT_ENCODER,
		Error:   nil,
	}

	if err := d.Publish(bCtx, pub); err != nil {
		return fmt.Errorf(
			`unable to publish message (subject "%s", value "%v") over encoded channel: %w`,
			pub.Subject,
			pub.Data,
			err,
		)
	}

	return nil
}
