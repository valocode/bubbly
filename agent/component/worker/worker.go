package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/verifa/bubbly/api"

	"github.com/verifa/bubbly/agent/component"
	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/env"
	"github.com/verifa/bubbly/interval"
)

const (
	defaultPollTimeout = 60
)

func New(bCtx *env.BubblyContext) *Worker {
	w := &Worker{
		ComponentCore: &component.ComponentCore{
			Type: component.WorkerComponent,
			NATSServer: &component.NATS{
				Config: bCtx.AgentConfig.NATSServerConfig,
			},
			DesiredSubscriptions: nil,
		},
		ResourceWorker: &interval.ResourceWorker{},
	}

	w.DesiredSubscriptions = w.defaultSubscriptions()

	bCtx.Logger.Debug().Msg("successfully initialised a worker")

	return w
}

// TODO: describe more about the Worker
type Worker struct {
	*component.ComponentCore
	ResourceWorker *interval.ResourceWorker
}

// pollResources attempts to poll any available data store
func (w *Worker) pollResources(bCtx *env.BubblyContext) (*component.Publication, error) {
	// We want to fetch all resource of type pipeline run from the data
	// store. So form a graphql query representing such
	resQuery := fmt.Sprintf(`
		{
			%s(kind: "%s") {
				name
				kind
				api_version
				metadata
				spec
			}
		}
	`, core.ResourceTableName, core.RunResourceKind)

	// embed the query into a Publication
	pub := component.Publication{
		Subject: component.StoreGetResourcesByKind,
		Encoder: nats.DEFAULT_ENCODER,
		Data:    []byte(resQuery),
	}

	for {
		// request the resource(s) from any available data store.
		reply, err := w.Request(bCtx, pub)

		// if there is no error,
		// then we've at least been sent a Publication from a data store
		// which might contain some Run resources
		if err == nil {
			resBlockJson := []core.ResourceBlockJSON{}
			err = json.Unmarshal(reply.Data, &resBlockJson)

			// if nil, then there are no resources in the _resource table of
			// the data store matching the required constraint (
			// Run type)
			if resBlockJson == nil {
				// just log
				bCtx.Logger.Debug().Err(err).Msg(
					"worker failed to request run resources from data store")
			} else if err != nil {
				// we fail to unmarshal correctly. Just log,
				// but it might be better to actually error here as a failure
				// to unmarshal may indicate a corrupt _resource table format?
				bCtx.Logger.Debug().Err(err).Msg(
					"worker failed to request run resources from data store")
			} else if reflect.DeepEqual(resBlockJson, []core.ResourceBlockJSON{}) {
				// handle the case where the response is non-nil but doesn't
				// contain any resources
				bCtx.Logger.Debug().Err(err).Str("required_kind", string(core.RunResourceKind)).Msg("no resources of required kind")
			} else {
				return reply, nil
			}
		}

		// if there is an error,
		// then a data store is either unavailable or not subscribed the the
		// necessary subject. Log this...
		bCtx.Logger.Debug().
			Int("timeout", defaultPollTimeout).
			Str("component", string(w.Type)).
			Err(err).
			Msg("waiting for interval resource(s) from a data store in order to start")

		// and wait to try again
		time.Sleep(defaultPollTimeout * time.Second)
	}

	return &pub, nil
}

// Run runs the interval.ResourceWorker
func (w *Worker) Run(bCtx *env.BubblyContext, agentContext context.Context) error {
	bCtx.Logger.Debug().
		Str(
			"component",
			string(w.Type)).
		Msg("running component")

	nSubs, err := w.BulkSubscribe(bCtx)

	if err != nil {
		return fmt.Errorf("error during bulk subscription: %w", err)
	}

	w.Subscriptions = nSubs
	bCtx.Logger.Debug().
		Str("component", string(w.Type)).
		Interface("subscriptions", w.Subscriptions).
		Msg("component is listening for subscriptions")

	return w.Listen(agentContext)
}

// run is a goroutine invoked from public Run method
func (w *Worker) run(bCtx *env.BubblyContext, ch chan error) {
	// poll for Run resources from the data store
	reply, err := w.pollResources(bCtx)

	if err != nil {
		ch <- fmt.Errorf("worker failed while polling for resources: %w", err)
	}

	resourcesBlockJSON := []core.ResourceBlockJSON{}
	err = json.Unmarshal(reply.Data, &resourcesBlockJSON)
	if err != nil {
		ch <- fmt.Errorf("failed to unmarshal run resources from data store"+
			": %w", err)
	}

	var resources []core.Resource

	for _, resBlockJSON := range resourcesBlockJSON {
		resBlock, err := resBlockJSON.ResourceBlock()

		if err != nil {
			ch <- fmt.Errorf("failed to form resourceBlock: %w", err)
		}
		res, err := api.NewResource(&resBlock)

		if err != nil {
			ch <- fmt.Errorf("failed to form resource: %w", err)
		}

		resources = append(resources, res)
	}

	// worker now has access to resources, so can "do" the work of running them
	// over their intervals
	w.ResourceWorker.ParseResources(bCtx, resources)

	err = w.ResourceWorker.Run(bCtx)
	if err != nil {
		ch <- fmt.Errorf("interval worker failure: %w", err)
	}
}

// a list of DesiredSubscriptions that the data store attempts to subscribe to
func (w *Worker) defaultSubscriptions() component.DesiredSubscriptions {
	return component.DesiredSubscriptions{
		component.DesiredSubscription{
			Subject: component.WorkerPostRunResource,
			Queue:   component.WorkerQueue,
			Handler: w.PostRunResourceHandler,
			Encoder: nats.DEFAULT_ENCODER,
		},
	}
}
