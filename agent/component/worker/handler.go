package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/cenkalti/backoff/v4"
	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/api"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
)

const (
	defaultInitialStoreDelay = 0
)

func (w *Worker) PostRunResourceHandler(bCtx *env.BubblyContext, m *nats.Msg) error {
	var resJSON core.ResourceBlockJSON

	err := json.Unmarshal(m.Data, &resJSON)

	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Str("res", resJSON.SpecRaw).
		Str("component", string(w.Type)).
		Msg("processing resource")

	if err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	resBlock, err := resJSON.ResourceBlock()

	if err != nil {
		return fmt.Errorf("failed to form resource from block: %w", err)
	}

	res, err := api.NewResource(&resBlock)

	w.ResourceWorker.ParseResources(bCtx, []core.Resource{res})

	err = w.ResourceWorker.Run(bCtx)
	if err != nil {
		return fmt.Errorf("interval worker failure: %w", err)
	}

	return nil
}

// PostRunResourceHandler receives a notification from the store triggers
// that a new run
// resource will soon be added to the store. The ID is passed from the store,
// which the worker uses to query and pull the full resource from the _resource
// table. This allows us to be sure that run resources trig
// to run it. It valides that the received resource is of kind `run` and type `OneOffRun`.
func (w *Worker) PostRunResourceHandler2(bCtx *env.BubblyContext, m *nats.Msg) error {
	bCtx.Logger.Debug().
		Interface("subscription", m.Sub).
		Interface("id", string(m.Data)).
		Str("component", string(w.Type)).
		Msg("processing message")

	resQuery := fmt.Sprintf(`
		{
			%s(id: "%s") {
				name
				kind
				api_version
				metadata
				spec
			}
		}
	`, core.ResourceTableName, m.Data)

	bCtx.Logger.Debug().Str("query", resQuery).Msg("worker querying for run resource")

	req := component.Publication{
		Subject: component.StoreGetResource,
		Data:    []byte(resQuery),
		Encoder: nats.DEFAULT_ENCODER,
	}

	var replyPub *component.Publication

	// operation will continue to request the run resource from the store
	operation := func() error {
		reply, err := w.Request(bCtx, req)
		if err != nil {
			bCtx.Logger.Debug().Str("resource", string(m.Data)).Err(err).Msg("worker failed to get run resource from store. Retrying after timeout")
			return err
		} else if reply.Error != nil {
			bCtx.Logger.Debug().Str("resource", string(m.Data)).Err(reply.Error).Msg("worker failed to get run resource from store. Retrying after timeout")
			return reply.Error
		}

		replyPub = reply
		return nil // or an error
	}
	bo := backoff.NewExponentialBackOff()

	bo.InitialInterval = 10 * time.Second
	bo.MaxElapsedTime = 120 * time.Second

	time.Sleep(defaultInitialStoreDelay * time.Second)
	// continue to request the run resource from the store until 120 seconds have passed
	err := backoff.Retry(operation, bo)

	if err != nil {
		return fmt.Errorf(
			`failed to get resource from query "%s" even after retrying: %w`,
			resQuery,
			err,
		)
	}
	//
	// // reply is a Publication received from a bubbly store
	// reply, err = w.Request(bCtx, req)
	//
	// if err != nil {
	// 	bCtx.Logger.Debug().Err(err).Int("timeout", defaultWorkerGetResourceTimeout).Msg("worker failed to get run resource from store. Retrying after timeout")
	// 	time.Sleep(defaultWorkerGetResourceTimeout * time.Second)
	// 	// reply is a Publication received from a bubbly store
	// 	reply, err = w.Request(bCtx, req)
	// } else if reply.Error != nil {
	// 	bCtx.Logger.Debug().Err(reply.Error).Int("timeout", defaultWorkerGetResourceTimeout).Msg("worker failed to get run resource from store. Retrying after timeout")
	// 	time.Sleep(defaultWorkerGetResourceTimeout * time.Second)
	// 	// reply is a Publication received from a bubbly store
	// 	reply, err = w.Request(bCtx, req)
	// }
	//
	// if err != nil {
	// 	return fmt.Errorf(
	// 		`failed to get resource from query "%s" even after retrying: %w`,
	// 		resQuery,
	// 		err,
	// 	)
	// }
	//
	// if reply != nil && reply.Error != nil {
	// 	return fmt.Errorf(
	// 		`failed to get resource from query "%s" even after retrying: %w`,
	// 		resQuery,
	// 		reply.Error,
	// 	)
	// }

	var result interface{}

	err = json.Unmarshal(replyPub.Data, &result)

	if err != nil {
		return fmt.Errorf("error unmarshalling resource: %w", err)
	}

	if result == nil || result.(map[string]interface{})[core.ResourceTableName] == nil {
		return fmt.Errorf("no resource found matching resource id `%s`", string(m.Data))
	}

	var (
		resJSON  core.ResourceBlockJSON
		inputMap = result.(map[string]interface{})[core.ResourceTableName].([]interface{})
	)
	b, err := json.Marshal(inputMap[0])

	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}
	err = json.Unmarshal(b, &resJSON)
	if err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	resBlock, err := resJSON.ResourceBlock()

	if err != nil {
		return fmt.Errorf("failed to form resource from block: %w", err)
	}

	res, err := api.NewResource(&resBlock)

	w.ResourceWorker.ParseResources(bCtx, []core.Resource{res})

	err = w.ResourceWorker.Run(bCtx)
	if err != nil {
		return fmt.Errorf("interval worker failure: %w", err)
	}

	return nil

	// w.Publish()
	//
	// var d2 core.Data
	//
	// err = json.Unmarshal(m.Data, &d2)
	//
	// x2, err := d2.Fields.MarshalJSON()
	//
	// p := core.ResourceBlock{}
	//
	// if err := json.Unmarshal(x2, &p); err != nil {
	// 	fmt.Errorf("oh drats")
	// }
	//
	// if err != nil {
	// 	fmt.Errorf("failed to unmarshal into core.Data: %w", err)
	// }
	//
	// resBlockJson := []core.ResourceBlockJSON{}
	//
	// var kind, name, version, specRaw string
	// var meta core.Metadata
	// gocty.FromCtyValue(d2.Fields["kind"], &kind)
	// gocty.FromCtyValue(d2.Fields["name"], &name)
	// gocty.FromCtyValue(d2.Fields["api_version"], &version)
	// gocty.FromCtyValue(d2.Fields["metadata"], &meta)
	// gocty.FromCtyValue(d2.Fields["spec"], &specRaw)
	//
	// resJSON := core.ResourceBlockJSON{
	// 	ResourceBlockAlias: core.ResourceBlockAlias{
	// 		ResourceKind:       kind,
	// 		ResourceName:       name,
	// 		ResourceAPIVersion: core.APIVersion(version),
	// 		Metadata:           &meta,
	// 	},
	// 	SpecRaw: specRaw,
	// }
	//
	// spew.Dump(resJSON)
	//
	// d3, _ := json.Marshal(d2)
	//
	// t := core.ResourceBlockJSON{}
	// if err := json.Unmarshal(d3, &t); err != nil {
	// 	fmt.Errorf("damn")
	// }
	//
	// if err := json.Unmarshal(m.Data, &resBlockJson); err != nil {
	// 	fmt.Errorf("oh no!")
	// }
	//
	// b, err = json.Marshal(d2.Data)
	//
	// if err != nil {
	// 	fmt.Errorf("Oh no!")
	// }
	//
	// var x core.ResourceBlockJSON
	// if err := json.Unmarshal(b, &x); err != nil {
	// 	fmt.Errorf("o noes")
	// }
	//
	// // d2.UnmarshalJSON()
	// // var rbj core.ResourceBlockJSON{}
	// //
	// // rbj.ConstructFromData(data core.Data)
	//
	// // TODO: handle running of resource here
	//
	// resourcesBlockJSON := core.ResourceBlockJSON{}
	// if err := json.Unmarshal(m.Data, &resourcesBlockJSON); err != nil {
	// 	fmt.Errorf("OH NOOOOO")
	// 	return fmt.Errorf(`failed to unmarshal received data into a core.ResourceBlockJSON: %w`, err)
	// }

	// return nil
}
