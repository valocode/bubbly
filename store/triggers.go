package store

import (
	"encoding/json"
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/valocode/bubbly/api"
	"github.com/valocode/bubbly/api/common"
	"github.com/valocode/bubbly/api/core"
	v1 "github.com/valocode/bubbly/api/v1"
	"github.com/valocode/bubbly/client"
	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/events"
	"github.com/valocode/bubbly/server"
)

type Kind int

const (
	Active Kind = iota
	Passive
)

func (k Kind) String() string {
	return [...]string{
		"Active",
		"Passive"}[k]
}

// A trigger is a wrapper around a visitFn, which defines what should happen
// at each node of the data tree.
// A trigger can be seen as a mechanism for introducing store-level
// "cause-and-effect", in that it will typically evaluate changes
// to subsections of the store (cause) and react with some effect.
// Typically, active triggers generate core.DataBlocks
// representing resultant desired changes (effect) whilst passive triggers
// have effects external to the store (for example, making a NATS publication
// to notify a microservice about something)
type trigger struct {
	id          string
	description string
	visitFn     visitFn
	Kind        Kind
}

// HandleTriggers evaluates changes to the dataTree against a set of triggers.
// Each trigger contains a visitFn that should "trigger" some action.
// Active triggers typically trigger the injection of
// new dataBlocks into the store on certain conditions.
// Once all triggers have been evaluated, the dataBlocks are converted
// to a dataTree and returned to be saved to the store
func HandleTriggers(bCtx *env.BubblyContext, tree dataTree, triggers []*trigger, kind Kind) (dataTree, error) {
	// First make sure we reset the tree so that we can traverse it again
	tree.reset()
	triggerBlocks := core.DataBlocks{}
	for _, t := range triggers {
		if t.Kind != kind {
			continue
		}
		blocks, err := tree.traverse(bCtx, t.visitFn)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to use trigger to traverse data tree: %w",
				err)
		}
		triggerBlocks = append(triggerBlocks, blocks...)
	}

	// if the triggers have generated dataBlocks (and are therefore active triggers),
	// they should be returned as a dataTree for saving to the store
	if len(triggerBlocks) > 0 && kind == Active {
		triggerTree, err := createDataTree(triggerBlocks)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to create tree of data blocks while handling triggers: %w", err)
		}
		return triggerTree, nil
	}

	return nil, nil
}

var internalTriggers = []*trigger{eventStoreTrigger, remoteRunTrigger}

var eventStoreTrigger = &trigger{
	id:          "default/trigger/event_store_trigger",
	description: "update event store upon new/updated entry to resource store",
	Kind:        Active,
	visitFn: func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error {
		switch node.Data.TableName {
		// if the table is _resource, then we know a resource has been loaded to the store
		case core.ResourceTableName:
			// if the _resource is provided only to provide a FK for another
			// _event entry (and therefore the only field provided is the "id"),
			// then we do not want to update the event store as that will take
			// place _anyway_.
			if len(node.Data.Fields) != 1 && node.Data.IsValidResource() {
				fields := node.Data.Fields

				// get the id of the resource
				id := fields["id"]

				if id.IsNull() {
					return fmt.Errorf(`DataBlock missing required "%s" field`, "id")
				}

				newEventBlocks := core.DataBlocks{
					core.Data{
						TableName: core.ResourceTableName,
						Fields: core.DataFields{
							"id": id,
						},
					},
					core.Data{
						TableName: core.EventTableName,
						Fields: map[string]cty.Value{
							"status": cty.StringVal(events.ResourceCreatedUpdated.String()),
							"time":   cty.StringVal(events.TimeNow()),
							"error":  cty.StringVal(""),
						},
						Joins: []string{core.ResourceTableName},
					},
				}

				*blocks = append(*blocks, newEventBlocks...)
			}

		case core.SchemaTableName:
			// TODO
		default:
			// fmt.Printf("trigger has no action for table: %s\n", node.Data.TableName)
		}
		return nil
	},
}

var remoteRunTrigger = &trigger{
	id:          "default/trigger/remote_run_trigger",
	description: "make a NATS publication upon new/updated entry to a run resource",
	Kind:        Passive,
	visitFn: func(bCtx *env.BubblyContext, node *dataNode, blocks *core.DataBlocks) error {
		switch node.Data.TableName {
		// if the table is _resource, then we know a resource has been loaded to the store
		case core.ResourceTableName:
			// if the _resource is provided only to provide a FK for another
			// _event entry (and therefore the only field provided is the "id"),
			// then we do not want to update the event store as that will take
			// place _anyway_.

			if len(node.Data.Fields) != 1 && node.Data.IsValidResource() {
				fields := node.Data.Fields

				// get the kind of the resource
				kind := fields["kind"]

				var resKind core.ResourceKind

				if err := gocty.FromCtyValue(kind, &resKind); err != nil {
					return fmt.Errorf(`failed to convert kind "%s" to core.ResourceKind`, kind)
				}

				// make sure the resource is of kind run
				if resKind == core.RunResourceKind {
					// make sure that the run resource is of type remote
					resJSON, _ := node.Data.ToResourceBlockJSON()

					resBlock, err := resJSON.ResourceBlock()

					if err != nil {
						return fmt.Errorf("failed to form resource from block: %w", err)
					}

					res, err := api.NewResource(&resBlock)
					if err != nil {
						return fmt.Errorf("failed to create resource from block: %w", err)
					}

					r := res.(*v1.Run)
					// TODO: need to pass RequestAuth instead of nil
					runCtx := core.NewResourceContext(cty.NilVal, api.NewResource, nil)
					if err := common.DecodeBody(bCtx, r.SpecHCL.Body, &r.Spec, runCtx); err != nil {
						return fmt.Errorf("failed to form resource from block: %w", err)
					}

					if r.Spec.Remote == nil {
						bCtx.Logger.Debug().Str("resource", r.String()).Msg("run is of type local and therefore should not be run by a bubbly worker")
						return nil
					}

					// resource validated as a remote run resource.
					// Now ship it to an available worker
					wr := server.WorkerRun{
						Name: r.Name(),
					}

					nc, err := client.New(bCtx)
					if err != nil {
						return fmt.Errorf("failed to connect to the NATS server: %w", err)
					}
					defer nc.Close()

					rBytes, err := json.Marshal(wr)
					if err != nil {
						return fmt.Errorf("failed to marshal ID into WorkerRun: %w", err)
					}

					// TODO: need to pass RequestAuth instead of nil
					if err := nc.PostResourceToWorker(bCtx, nil, rBytes); err != nil {
						return fmt.Errorf("failed to post resource to worker: %w", err)
					}

				}
			}
		}
		return nil
	},
}
